// file: qr-service/internal/service/qr_processor.go
package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"log"
	"qr/config"
	"qr/dto"
	"qr/pkg/kafka"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/skip2/go-qrcode"
)

const (
	defaultPublicIDPrefix = "content_qr/"
	qrCodeSize            = 512
)

// QRProcessor chứa các dependency và logic xử lý
type QRProcessor struct {
	cld       *cloudinary.Cloudinary
	publisher *kafka.Publisher
	cfg       *config.Config
}

func NewQRProcessor(cld *cloudinary.Cloudinary, publisher *kafka.Publisher, cfg *config.Config) *QRProcessor {
	return &QRProcessor{
		cld:       cld,
		publisher: publisher,
		cfg:       cfg,
	}
}

// ProcessOrderRequest là hàm xử lý chính cho một message từ Kafka
func (p *QRProcessor) ProcessOrderRequest(ctx context.Context, orderRequest dto.OrderQRGenerationRequestEvent) {
	log.Printf("Processing Order QR request for Order ID %s with %d tickets", orderRequest.OrderID, len(orderRequest.Tickets))

	processedTickets, err := p.processOrderInParallel(ctx, orderRequest)
	if err != nil {
		log.Printf("Failed to process QR codes for order %s: %v", orderRequest.OrderID, err)
		// Trong hệ thống thực tế, message này sẽ được thử lại hoặc gửi vào Dead Letter Queue.
		return
	}

	err = p.publishEmailRequest(ctx, orderRequest, processedTickets)
	if err != nil {
		log.Printf("Error publishing email request for order %s: %v", orderRequest.OrderID, err)
	} else {
		log.Printf("Successfully published email request for order %s", orderRequest.OrderID)
	}
}

func (p *QRProcessor) processOrderInParallel(ctx context.Context, orderRequest dto.OrderQRGenerationRequestEvent) ([]dto.TicketDetailForEmail, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var processedTickets []dto.TicketDetailForEmail
	var errors []error

	for _, ticket := range orderRequest.Tickets {
		wg.Add(1)
		go func(t dto.TicketDetailForQR) {
			defer wg.Done()
			qrImage, err := generateQRCodeWithFixedLogo(t.QRContent)
			if err != nil {
				log.Printf("Error generating QR image for seat %d: %v", t.SeatID, err)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
			publicID := defaultPublicIDPrefix + generateContentHash(t.QRContent)
			cloudinaryURL, _, err := p.uploadToCloudinary(ctx, qrImage, publicID)
			if err != nil {
				log.Printf("Error uploading QR to Cloudinary for seat %d: %v", t.SeatID, err)
				mu.Lock()
				errors = append(errors, err)
				mu.Unlock()
				return
			}
			mu.Lock()
			processedTickets = append(processedTickets, dto.TicketDetailForEmail{
				SeatInfo:  fmt.Sprintf("Ghế số %d", t.SeatID),
				QRCodeURL: cloudinaryURL,
			})
			mu.Unlock()
		}(ticket)
	}
	wg.Wait()

	if len(errors) > 0 {
		return nil, fmt.Errorf("encountered %d errors during QR processing for order %s", len(errors), orderRequest.OrderID)
	}
	return processedTickets, nil
}

func (p *QRProcessor) publishEmailRequest(ctx context.Context, orderRequest dto.OrderQRGenerationRequestEvent, processedTickets []dto.TicketDetailForEmail) error {
	title := fmt.Sprintf("Xác nhận đơn hàng của bạn - Mã đơn %s", orderRequest.OrderID)
	emailData := dto.OrderConfirmationData{
		CustomerName: orderRequest.CustomerName,
		OrderID:      orderRequest.OrderID,
		TotalPrice:   orderRequest.TotalPrice,
		Tickets:      processedTickets,
	}
	bodyBytes, err := json.Marshal(emailData)
	if err != nil {
		return fmt.Errorf("failed to marshal order confirmation data: %w", err)
	}
	emailPayload := dto.EmailRequestEvent{
		To:    orderRequest.CustomerEmail,
		Title: title,
		Body:  string(bodyBytes),
		Type:  "order_confirmation_payment",
	}
	return p.publisher.Publish(ctx, p.cfg.Kafka.TopicEmailRequests, []byte(orderRequest.CustomerEmail), emailPayload)
}

func (p *QRProcessor) uploadToCloudinary(ctx context.Context, imageToUpload image.Image, publicID string) (string, string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, imageToUpload); err != nil {
		return "", "", fmt.Errorf("lỗi encode ảnh PNG: %v", err)
	}
	uploadParams := uploader.UploadParams{
		PublicID:   publicID,
		Format:     "png",
		Overwrite:  boolPtr(true),
		Invalidate: boolPtr(true),
	}
	uploadResult, err := p.cld.Upload.Upload(ctx, &buf, uploadParams)
	if err != nil {
		return "", "", fmt.Errorf("lỗi tải lên Cloudinary: %v", err)
	}
	return uploadResult.SecureURL, uploadResult.PublicID, nil
}

func generateQRCodeWithFixedLogo(content string) (image.Image, error) {
	qrInstance, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo QR: %v", err)
	}
	qrInstance.DisableBorder = true
	return qrInstance.Image(qrCodeSize), nil
}

func generateContentHash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

func boolPtr(b bool) *bool { return &b }

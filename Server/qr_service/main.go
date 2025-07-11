package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"os/signal"
	"qr/config"
	"qr/dto"
	"qr/internal/service"
	"qr/pkg/kafka"
	"syscall"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
	qrcode "github.com/skip2/go-qrcode"
	"github.com/twmb/franz-go/pkg/kgo"
)

const (
	defaultPublicIDPrefix = "content_qr/"
	defaultUploadsFolder  = "generic_uploads"
	defaultLogoPath       = "img.png" // Đường dẫn đến logo mặc định
)

// --- DTOs (Data Transfer Objects) ---

// OrderQRGenerationRequestEvent: DTO cho event nhận từ Kafka để yêu cầu tạo QR
type OrderQRGenerationRequestEvent struct {
	OrderID       string              `json:"orderId"`
	CustomerEmail string              `json:"customerEmail"`
	CustomerName  string              `json:"customerName"`
	TotalPrice    float64             `json:"totalPrice"`
	Tickets       []TicketDetailForQR `json:"tickets"`
}

type TicketDetailForQR struct {
	SeatID    int32  `json:"seatId"`
	QRContent string `json:"qrContent"`
}

// OrderConfirmationData: DTO để xây dựng data gửi cho Email Service
type OrderConfirmationData struct {
	CustomerName string                 `json:"customerName"`
	OrderID      string                 `json:"orderId"`
	TotalPrice   float64                `json:"totalPrice"`
	Tickets      []TicketDetailForEmail `json:"tickets"` // Danh sách vé đã có QR
}

type TicketDetailForEmail struct {
	SeatInfo  string `json:"seatInfo"`
	QRCodeURL string `json:"qrCodeUrl"`
}

// EmailRequestEvent: DTO để gửi yêu cầu đến Email Service qua Kafka
type EmailRequestEvent struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Type  string `json:"type"`
}

// QRInfoResponse: DTO cho các API response liên quan đến thông tin QR
type QRInfoResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
	CloudinaryURL string `json:"cloudinary_url,omitempty"`
	PublicID      string `json:"public_id,omitempty"`
	Content       string `json:"content,omitempty"`
	Error         string `json:"error,omitempty"`
}

// UploadImageResponse: DTO cho API response tải ảnh lên
type UploadImageResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message,omitempty"`
	CloudinaryURL string `json:"cloudinary_url,omitempty"`
	PublicID      string `json:"public_id,omitempty"`
	Error         string `json:"error,omitempty"`
}

// --- Hàm Main ---
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Không thể tải cấu hình: %v", err)
	}

	cld, err := initCloudinary(cfg)
	if err != nil {
		log.Fatalf("Không thể khởi tạo Cloudinary: %v", err)
	}

	kafkaPublisher, err := kafka.NewPublisher(cfg.Kafka)
	if err != nil {
		log.Fatalf("Không thể khởi tạo Kafka Publisher: %v", err)
	}
	defer kafkaPublisher.Close()

	qrProcessor := service.NewQRProcessor(cld, kafkaPublisher, cfg)
	consumerCtx, consumerCancel := context.WithCancel(context.Background())
	defer consumerCancel()
	go startKafkaConsumer(consumerCtx, cfg, qrProcessor)

	router := setupGinRouter(cld) // Truyền client Cloudinary vào router setup
	server := &http.Server{Addr: ":" + cfg.Server.Port, Handler: router}

	go func() {
		log.Printf("API server đang lắng nghe trên port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Đang tắt server và consumer...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server shutdown thất bại:", err)
	}
	log.Println("Server đã tắt thành công.")
}

// --- Các hàm khởi tạo và tiện ích ---

func initCloudinary(cfg *config.Config) (*cloudinary.Cloudinary, error) {
	if cfg.Cloudinary.URL == "" {
		return nil, fmt.Errorf("biến môi trường CLOUDINARY_URL chưa được thiết lập")
	}
	cld, err := cloudinary.NewFromURL(cfg.Cloudinary.URL)
	if err != nil {
		return nil, fmt.Errorf("lỗi khởi tạo Cloudinary từ URL: %w", err)
	}
	log.Println("Khởi tạo Cloudinary thành công.")
	return cld, nil
}

func generateContentHash(content string) string {
	hasher := sha256.New()
	hasher.Write([]byte(content))
	return hex.EncodeToString(hasher.Sum(nil))
}

func boolPtr(b bool) *bool { return &b }

// --- Logic Kafka Consumer ---

func startKafkaConsumer(ctx context.Context, cfg *config.Config, processor *service.QRProcessor) {
	consumerClient, err := kafka.NewConsumerClient(
		cfg.Kafka,
		cfg.Kafka.TopicOrderQRRequests,
		cfg.Kafka.GroupIDOrderQR,
	)
	if err != nil {
		log.Fatalf("Không thể tạo Kafka Consumer: %v", err)
	}
	defer consumerClient.Close()

	log.Printf("Kafka QR Order consumer đã bắt đầu. Đang lắng nghe trên topic: %s", cfg.Kafka.TopicOrderQRRequests)

	for {
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer nhận được tín hiệu dừng.")
			return
		default:
			fetches := consumerClient.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				log.Printf("Lỗi khi nhận message từ Kafka: %v", errs)
				continue
			}

			fetches.EachRecord(func(record *kgo.Record) {
				var orderRequest dto.OrderQRGenerationRequestEvent
				if err := json.Unmarshal(record.Value, &orderRequest); err != nil {
					log.Printf("Lỗi giải mã JSON từ Kafka: %v. Bỏ qua message.", err)
					if err_commit := consumerClient.CommitRecords(ctx, record); err_commit != nil {
						log.Printf("Lỗi commit message lỗi: %v", err_commit)
					}
					return
				}
				processor.ProcessOrderRequest(context.Background(), orderRequest)
				if err := consumerClient.CommitRecords(ctx, record); err != nil {
					log.Printf("Lỗi commit offset: %v", err)
				}
			})
		}
	}
}

// --- CÁC API ROUTER VÀ HANDLER ---

func setupGinRouter(cld *cloudinary.Cloudinary) *gin.Engine {
	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	{
		qrGroup := apiV1.Group("/qr")
		{
			// Truyền client Cloudinary vào các handler
			qrGroup.GET("/image", handleViewImageByContent(cld))
			qrGroup.GET("/url", handleGetQRURLByContent(cld))
			// === API MỚI ĐỂ TEST ===
			qrGroup.POST("/generate-test", handleGenerateTestQR(cld))
		}

		uploadGroup := apiV1.Group("/upload")
		{
			uploadGroup.POST("/image", handleImageUpload(cld))
		}
	}

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "timestamp": time.Now().UTC()})
	})
	return router
}

// handleImageUpload: Tải một file ảnh bất kỳ lên
func handleImageUpload(cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, UploadImageResponse{Success: false, Error: "Yêu cầu không hợp lệ, không tìm thấy file 'image'"})
			return
		}
		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, UploadImageResponse{Success: false, Error: "Không thể mở file đã tải lên"})
			return
		}
		defer src.Close()

		uploadParams := uploader.UploadParams{
			Folder:         defaultUploadsFolder,
			UseFilename:    boolPtr(false),
			UniqueFilename: boolPtr(true),
			Overwrite:      boolPtr(false),
		}

		uploadResult, err := cld.Upload.Upload(c.Request.Context(), src, uploadParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, UploadImageResponse{Success: false, Error: "Lỗi tải ảnh lên Cloudinary: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, UploadImageResponse{
			Success: true, Message: "Tải ảnh lên thành công!", CloudinaryURL: uploadResult.SecureURL, PublicID: uploadResult.PublicID,
		})
	}
}

// handleViewImageByContent: Xem ảnh QR (redirect)
func handleViewImageByContent(cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentToQuery := c.Query("content")
		if contentToQuery == "" {
			c.JSON(http.StatusBadRequest, QRInfoResponse{Success: false, Error: "Query parameter 'content' không được để trống"})
			return
		}
		publicIDToView := defaultPublicIDPrefix + generateContentHash(contentToQuery)
		imgObj, err := cld.Image(publicIDToView)
		if err != nil {
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi tạo đối tượng ảnh từ public_id"})
			return
		}
		imageURL, err := imgObj.String()
		// Cloudinary trả về URL có dạng http://res.cloudinary.com/<cloud_name>/image/upload/v<version>/<public_id>
		// Nếu ảnh không tồn tại, nó vẫn trả về URL này mà không có lỗi. Cần một cách kiểm tra khác.
		// Tuy nhiên, để đơn giản, ta tạm chấp nhận và redirect. Trình duyệt sẽ hiển thị lỗi nếu ảnh không có.
		if err != nil {
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi tạo URL ảnh: " + err.Error()})
			return
		}
		c.Redirect(http.StatusFound, imageURL)
	}
}

// handleGetQRURLByContent: Lấy URL của ảnh QR (JSON)
func handleGetQRURLByContent(cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		contentToQuery := c.Query("content")
		if contentToQuery == "" {
			c.JSON(http.StatusBadRequest, QRInfoResponse{Success: false, Error: "Query parameter 'content' không được để trống"})
			return
		}
		publicIDToView := defaultPublicIDPrefix + generateContentHash(contentToQuery)
		imgObj, err := cld.Image(publicIDToView)
		if err != nil {
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi tạo đối tượng ảnh từ public_id"})
			return
		}
		imageURL, err := imgObj.String()
		if err != nil {
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi tạo URL ảnh: " + err.Error()})
			return
		}

		// Tạm thời, chúng ta không thể kiểm tra sự tồn tại của ảnh một cách hiệu quả chỉ bằng URL.
		// Trả về URL và để client xử lý.
		c.JSON(http.StatusOK, QRInfoResponse{
			Success: true, Message: "Lấy URL ảnh thành công.", CloudinaryURL: imageURL, PublicID: publicIDToView, Content: contentToQuery,
		})
	}
}

// --- HANDLER CHO API TEST TẠO QR CÓ LOGO ---
func handleGenerateTestQR(cld *cloudinary.Cloudinary) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody struct {
			Content string `json:"content"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, QRInfoResponse{Success: false, Error: "Invalid request body. Cần có trường 'content'."})
			return
		}

		if requestBody.Content == "" {
			c.JSON(http.StatusBadRequest, QRInfoResponse{Success: false, Error: "Trường 'content' không được để trống."})
			return
		}

		// 1. Tạo ảnh QR với logo
		qrImage, err := generateQRCodeWithLogo(requestBody.Content, defaultLogoPath)
		if err != nil {
			log.Printf("Lỗi tạo QR code với logo: %v", err)
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi khi tạo ảnh QR: " + err.Error()})
			return
		}

		// 2. Chuyển ảnh thành byte stream để upload
		buf := new(bytes.Buffer)
		if err := png.Encode(buf, qrImage); err != nil {
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi khi mã hóa ảnh QR: " + err.Error()})
			return
		}
		imageBytes := buf.Bytes()

		// 3. Upload lên Cloudinary
		publicID := defaultPublicIDPrefix + generateContentHash(requestBody.Content)
		uploadParams := uploader.UploadParams{
			PublicID:       publicID,
			Folder:         defaultPublicIDPrefix,
			Overwrite:      boolPtr(true),
			UseFilename:    boolPtr(false),
			UniqueFilename: boolPtr(false),
		}

		uploadResult, err := cld.Upload.Upload(c.Request.Context(), bytes.NewReader(imageBytes), uploadParams)
		if err != nil {
			c.JSON(http.StatusInternalServerError, QRInfoResponse{Success: false, Error: "Lỗi tải ảnh QR lên Cloudinary: " + err.Error()})
			return
		}

		// 4. Trả về kết quả
		c.JSON(http.StatusOK, QRInfoResponse{
			Success:       true,
			Message:       "Tạo và tải ảnh QR có logo thành công!",
			CloudinaryURL: uploadResult.SecureURL,
			PublicID:      uploadResult.PublicID,
			Content:       requestBody.Content,
		})
	}
}

// --- HÀM TẠO QR CODE VỚI LOGO Ở GIỮA (ĐÃ KHÔI PHỤC) ---

// generateQRCodeWithLogo tạo một mã QR từ content và chèn logo vào giữa.
// logoPath là đường dẫn file cục bộ đến ảnh logo (ví dụ "img.png").
func generateQRCodeWithLogo(content string, logoPath string) (image.Image, error) {
	// Mức độ sửa lỗi cao nhất để đảm bảo QR dễ đọc dù bị che một phần bởi logo
	qr, err := qrcode.New(content, qrcode.Highest)
	if err != nil {
		return nil, fmt.Errorf("lỗi tạo đối tượng QR code: %v", err)
	}
	// Đặt màu nền và màu QR code
	qr.BackgroundColor = image.White
	qr.ForegroundColor = image.Black

	// Tạo ảnh QR code với kích thước 512x512 pixels
	qrImage := qr.Image(512)
	if qrImage == nil {
		return nil, fmt.Errorf("lỗi tạo ảnh từ đối tượng QR code")
	}

	// Tải file logo từ đường dẫn cục bộ
	logoFile, err := os.Open(logoPath)
	if err != nil {
		// Nếu không tìm thấy logo, vẫn tiếp tục nhưng ghi log cảnh báo.
		// QR code sẽ được tạo mà không có logo.
		log.Printf("Cảnh báo: Không tìm thấy file logo local '%s': %v. Sẽ tạo QR code không có logo.", logoPath, err)
	}
	var logoImage image.Image
	if logoFile != nil {
		defer logoFile.Close()
		logoImage, _, err = image.Decode(logoFile)
		if err != nil {
			return nil, fmt.Errorf("lỗi giải mã logo từ file '%s': %v", logoPath, err)
		}
	}

	// Tạo một ảnh mới để kết hợp QR và logo
	finalRect := qrImage.Bounds()
	finalImage := image.NewRGBA(finalRect)
	draw.Draw(finalImage, finalRect, qrImage, image.Point{}, draw.Src) // Vẽ QR làm nền

	// Chỉ vẽ logo nếu đã tải được
	if logoImage != nil {
		// Tỷ lệ logo so với QR code (ví dụ: 1/5 chiều rộng)
		logoProportion := 5
		logoTargetWidth := uint(qrImage.Bounds().Dx() / logoProportion)
		if logoTargetWidth == 0 {
			logoTargetWidth = 1 // Đảm bảo không chia cho 0
		}
		// Thay đổi kích thước logo
		logoResized := resize.Resize(logoTargetWidth, 0, logoImage, resize.Lanczos3)

		// Tính toán vị trí để đặt logo vào giữa
		offsetX := (finalRect.Dx() - logoResized.Bounds().Dx()) / 2
		offsetY := (finalRect.Dy() - logoResized.Bounds().Dy()) / 2
		logoDrawRect := logoResized.Bounds().Add(image.Pt(offsetX, offsetY))

		// Vẽ logo lên trên ảnh QR
		draw.Draw(finalImage, logoDrawRect, logoResized, image.Point{}, draw.Over)
	}

	return finalImage, nil
}

package worker

import (
	"context"
	"log"
	"strings"
	"time"

	"payment_service/domain/model"
	"payment_service/internal/service" // Import service interface của bạn

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type ExpirySubscriber struct {
	redisClient    *redis.Client
	invoiceService service.InvoiceServiceInterface // Dùng interface để dễ test
}

func NewExpirySubscriber(redisClient *redis.Client, invoiceService service.InvoiceServiceInterface) *ExpirySubscriber {
	return &ExpirySubscriber{
		redisClient:    redisClient,
		invoiceService: invoiceService,
	}
}

// Start bắt đầu lắng nghe sự kiện
func (s *ExpirySubscriber) Start(ctx context.Context) {
	// Kênh mà Redis sẽ publish sự kiện key hết hạn
	channel := "__keyevent@0__:expired"
	pubsub := s.redisClient.Subscribe(ctx, channel)
	defer pubsub.Close()

	log.Println("Bắt đầu lắng nghe sự kiện hết hạn từ Redis trên kênh:", channel)

	ch := pubsub.Channel()

	for msg := range ch {
		// Key có dạng "invoice_expiry:uuid-của-invoice"
		if strings.HasPrefix(msg.Payload, "invoice_expiry:") {
			invoiceIDStr := strings.TrimPrefix(msg.Payload, "invoice_expiry:")
			log.Printf("Nhận được sự kiện hết hạn cho key: %s", msg.Payload)

			// Xử lý việc hủy hóa đơn
			s.handleExpiredInvoice(ctx, invoiceIDStr)
		}
	}
}

// handleExpiredInvoice xử lý logic khi hóa đơn hết hạn
func (s *ExpirySubscriber) handleExpiredInvoice(ctx context.Context, invoiceIDStr string) {
	invoiceID, err := uuid.Parse(invoiceIDStr)
	if err != nil {
		log.Printf("LỖI: Định dạng Invoice ID không hợp lệ từ Redis event: %s", invoiceIDStr)
		return
	}

	// 1. Lấy thông tin hóa đơn từ DB để kiểm tra trạng thái hiện tại
	// Sử dụng một context mới với timeout để tránh worker bị treo
	procCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	invoice, err := s.invoiceService.GetInvoiceByID(procCtx, invoiceID)
	if err != nil {
		log.Printf("LỖI: Không thể lấy hóa đơn %s để xử lý hết hạn: %v", invoiceID, err)
		return
	}

	// 2. **QUAN TRỌNG**: Kiểm tra xem hóa đơn có còn đang chờ thanh toán không
	// Điều này tránh việc hủy một hóa đơn đã được thanh toán thành công.
	currentStatus := model.PaymentStatus(invoice.PaymentStatus.String)
	if currentStatus != model.PaymentStatusPending && currentStatus != model.PaymentStatusAwaitingConfirmation {
		log.Printf("INFO: Hóa đơn %s không còn ở trạng thái chờ (trạng thái hiện tại: %s). Bỏ qua việc hủy.", invoiceID, currentStatus)
		return
	}

	// 3. Nếu vẫn đang chờ, cập nhật trạng thái thành FAILED
	log.Printf("INFO: Hóa đơn %s đã hết hạn. Cập nhật trạng thái thành FAILED.", invoiceID)
	reason := "Giao dịch đã hết hạn thanh toán (15 phút)."

	// Sử dụng hàm đã có để đảm bảo logic nhất quán (bao gồm cả việc gửi Kafka event)
	// Giả sử payment method không quan trọng khi chỉ cần invoiceID để fail
	_, err = s.invoiceService.UpdateInvoiceStatusForPaymentFailureForUUID(procCtx, invoice.InvoiceID, model.PaymentMethod(invoice.PaymentMethod.String), reason)
	if err != nil {
		log.Printf("LỖI NGHIÊM TRỌNG: Không thể cập nhật trạng thái FAILED cho hóa đơn hết hạn %s: %v", invoiceID, err)
		// Cần có cơ chế retry hoặc alerting ở đây
	} else {
		log.Printf("THÀNH CÔNG: Đã cập nhật trạng thái FAILED cho hóa đơn hết hạn %s.", invoiceID)
	}
}

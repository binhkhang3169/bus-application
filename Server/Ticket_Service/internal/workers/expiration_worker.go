// file: internal/workers/expiration_worker.go
package workers

import (
	"context"
	"ticket-service/internal/services"
	"ticket-service/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

type ExpirationWorker struct {
	redisClient   *redis.Client
	ticketService services.ITicketService
	logger        utils.Logger
	interval      time.Duration
}

// NewExpirationWorker khởi tạo worker dọn dẹp
func NewExpirationWorker(redisClient *redis.Client, ticketService services.ITicketService, logger utils.Logger) *ExpirationWorker {
	return &ExpirationWorker{
		redisClient:   redisClient,
		ticketService: ticketService,
		logger:        logger,
		interval:      5 * time.Minute, // Chạy quét dọn mỗi 5 phút
	}
}

// Start khởi động worker
func (w *ExpirationWorker) Start(ctx context.Context) {
	w.logger.Info("Starting Expiration Worker...")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Chạy ngay lần đầu tiên khi khởi động
	w.cleanupExpiredBookings(ctx)

	for {
		select {
		case <-ticker.C:
			w.cleanupExpiredBookings(ctx)
		case <-ctx.Done():
			w.logger.Info("Stopping Expiration Worker.")
			return
		}
	}
}

// cleanupExpiredBookings quét và xử lý các booking đã hoàn thành nhưng bị bỏ rơi
func (w *ExpirationWorker) cleanupExpiredBookings(ctx context.Context) {
	w.logger.Info("Expiration Worker: Running cleanup cycle...")

	// Quét các key state, dùng SCAN để không block Redis
	var cursor uint64
	matchPattern := "booking:state:*"

	for {
		keys, nextCursor, err := w.redisClient.Scan(ctx, cursor, matchPattern, 50).Result()
		if err != nil {
			w.logger.Error("Expiration Worker: Error scanning for booking state keys: %v", err)
			return
		}

		for _, key := range keys {
			w.processKey(ctx, key)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
}

func (w *ExpirationWorker) processKey(ctx context.Context, key string) {
	stateData, err := w.redisClient.HGetAll(ctx, key).Result()
	if err != nil {
		w.logger.Error("Expiration Worker: Could not retrieve state for key %s: %v", key, err)
		return
	}

	status, ok := stateData["status"]
	if !ok {
		return // Không có trạng thái, bỏ qua
	}

	// Chỉ xử lý các vé đã hoàn thành nhưng client không bao giờ kết nối ws để "nhận"
	if status == "COMPLETED" {
		submittedAtStr, ok := stateData["submitted_at"]
		if !ok {
			return
		}

		submittedAt, err := time.Parse(time.RFC3339, submittedAtStr)
		if err != nil {
			return
		}

		// Nếu booking đã hoàn thành và tồn tại hơn 10 phút -> coi như bị bỏ rơi
		if time.Since(submittedAt) > 10*time.Minute {
			ticketID, ok := stateData["ticket_id"]
			if !ok || ticketID == "" {
				w.logger.Error("Expiration Worker: Found abandoned COMPLETED booking %s without a ticket_id.", key)
				return
			}

			w.logger.Info("Expiration Worker: Found abandoned booking. Cancelling ticket %s.", ticketID)

			// Hủy vé
			cancelCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			if err := w.ticketService.CancelTicket(cancelCtx, ticketID); err != nil {
				w.logger.Error("Expiration Worker: Failed to cancel ticket %s: %v", ticketID, err)
			} else {
				// Xóa key khỏi Redis sau khi đã xử lý xong
				w.redisClient.Del(ctx, key)
			}
		}
	}
}

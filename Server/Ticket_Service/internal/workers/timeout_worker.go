// file: internal/workers/timeout_worker.go
package workers

import (
	"context"
	"fmt"
	"strconv"
	"ticket-service/internal/services"
	"ticket-service/pkg/utils"
	"time"

	"github.com/redis/go-redis/v9"
)

const pendingWsConnectionsKey = "pending_ws_connections"

type TimeoutWorker struct {
	redisClient   *redis.Client
	ticketService services.ITicketService
	logger        utils.Logger
	interval      time.Duration
}

func NewTimeoutWorker(redisClient *redis.Client, ticketService services.ITicketService, logger utils.Logger) *TimeoutWorker {
	return &TimeoutWorker{
		redisClient:   redisClient,
		ticketService: ticketService,
		logger:        logger,
		interval:      1 * time.Second, // Quét mỗi giây để đảm bảo timeout gần như tức thì
	}
}

func (w *TimeoutWorker) Start(ctx context.Context) {
	w.logger.Info("Starting Timeout Worker...")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			w.processExpiredBookings(ctx)
		case <-ctx.Done():
			w.logger.Info("Stopping Timeout Worker.")
			return
		}
	}
}

func (w *TimeoutWorker) processExpiredBookings(ctx context.Context) {
	now := time.Now().UnixMilli()
	// Lấy tất cả các bookingId đã hết hạn (score < thời gian hiện tại)
	expiredBookingIDs, err := w.redisClient.ZRangeByScore(ctx, pendingWsConnectionsKey, &redis.ZRangeBy{
		Min: "0",
		Max: strconv.FormatInt(now, 10),
	}).Result()

	if err != nil {
		w.logger.Error("TimeoutWorker: Error fetching expired bookings from Redis: %v", err)
		return
	}

	if len(expiredBookingIDs) == 0 {
		return
	}

	w.logger.Info("TimeoutWorker: Found %d expired booking(s) that did not connect to WebSocket.", len(expiredBookingIDs))

	// Xóa các booking đã hết hạn khỏi sorted set để tránh xử lý lại
	if _, err := w.redisClient.ZRem(ctx, pendingWsConnectionsKey, expiredBookingIDs).Result(); err != nil {
		w.logger.Error("TimeoutWorker: Failed to remove expired bookings from sorted set: %v", err)
	}

	// Xử lý hủy vé cho từng bookingId
	for _, bookingID := range expiredBookingIDs {
		go w.handleCancellation(bookingID)
	}
}

func (w *TimeoutWorker) handleCancellation(bookingID string) {
	w.logger.Info("TimeoutWorker: Handling cancellation for bookingId %s", bookingID)

	// Chúng ta cần tìm ticketId tương ứng. Trạng thái COMPLETED trong Redis Hash sẽ chứa ticketId.
	redisStateKey := fmt.Sprintf("booking:state:%s", bookingID)

	// Chờ một khoảng thời gian ngắn để consumer có thể xử lý xong và ghi ticketId vào hash
	time.Sleep(5 * time.Second) // Chờ 5s phòng trường hợp consumer đang xử lý

	stateData, err := w.redisClient.HGetAll(context.Background(), redisStateKey).Result()
	if err != nil || len(stateData) == 0 {
		w.logger.Error("TimeoutWorker: Could not find state for bookingId %s after timeout. It might have been processed and cleaned up, or never existed.", bookingID)
		return
	}

	status, _ := stateData["status"]
	ticketID, ok := stateData["ticket_id"]

	if status == "COMPLETED" && ok && ticketID != "" {
		// Đây là trường hợp consumer đã tạo vé xong, nhưng client không kết nối ws.
		w.logger.Info("TimeoutWorker: Booking %s was COMPLETED. Cancelling associated ticket %s.", bookingID, ticketID)

		cancelCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := w.ticketService.CancelTicket(cancelCtx, ticketID); err != nil {
			w.logger.Error("TimeoutWorker: Failed to cancel ticket %s: %v", ticketID, err)
		} else {
			// Xóa key state sau khi đã hủy vé thành công
			w.redisClient.Del(context.Background(), redisStateKey)
		}
	} else {
		// Trường hợp QUEUED hoặc PROCESSING. Lúc này chưa có vé để hủy.
		// Chỉ cần xóa key state để consumer sau này sẽ bỏ qua.
		w.logger.Info("TimeoutWorker: Booking %s timed out with status '%s'. Deleting state key.", bookingID, status)
		w.redisClient.Del(context.Background(), redisStateKey)
	}
}

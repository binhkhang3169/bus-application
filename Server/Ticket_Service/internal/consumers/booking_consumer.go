// file: ticket-service/internal/consumers/booking_consumer.go
package consumers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"ticket-service/config"
	"ticket-service/internal/services"
	"ticket-service/pkg/kafkaclient"
	"ticket-service/pkg/utils"
	"ticket-service/pkg/websocket"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type BookingRequestConsumer struct {
	ticketService services.ITicketService
	redisClient   *redis.Client
	logger        utils.Logger
	kafkaCfg      config.KafkaConfig
	topicCfg      config.TopicConfig
}

// NewBookingRequestConsumer khởi tạo consumer mới, nhận vào toàn bộ config.
func NewBookingRequestConsumer(cfg config.Config, ticketService services.ITicketService, redisClient *redis.Client, logger utils.Logger) *BookingRequestConsumer {
	return &BookingRequestConsumer{
		ticketService: ticketService,
		redisClient:   redisClient,
		logger:        logger,
		kafkaCfg:      cfg.Kafka,
		topicCfg:      cfg.Kafka.Topics.BookingRequests,
	}
}

// Start khởi động vòng lặp consumer.
func (c *BookingRequestConsumer) Start(ctx context.Context) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(c.kafkaCfg.Seeds...),
		kgo.ConsumerGroup(c.topicCfg.GroupID),
		kgo.ConsumeTopics(c.topicCfg.Topic),
		kgo.DisableAutoCommit(),
	}
	if c.kafkaCfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}
	if c.kafkaCfg.SASLUser != "" && c.kafkaCfg.SASLPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{
			User: c.kafkaCfg.SASLUser,
			Pass: c.kafkaCfg.SASLPass,
		}.AsSha256Mechanism()))
	}
	// (Tùy chọn) Thêm logic SASL/TLS nếu cần
	// if c.kafkaCfg.EnableTLS { ... }

	client, err := kgo.NewClient(opts...)
	if err != nil {
		c.logger.Error("BookingConsumer: Failed to create Kafka client: %v", err)
		return
	}
	defer client.Close()

	c.logger.Info("Starting Kafka consumer for topic: %s", c.topicCfg.Topic)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("BookingConsumer: Shutting down.")
			return
		default:
			fetches := client.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				c.logger.Error("BookingConsumer: Fetch errors: %v", errs)
				continue
			}

			fetches.EachRecord(func(record *kgo.Record) {
				// Xử lý mỗi message trong một goroutine riêng để không block fetcher
				go c.handleMessage(context.Background(), record, client)
			})
		}
	}
}

func (c *BookingRequestConsumer) handleMessage(ctx context.Context, record *kgo.Record, client *kgo.Client) {
	var event kafkaclient.BookingRequestEvent
	if err := json.Unmarshal(record.Value, &event); err != nil {
		c.logger.Error("BookingConsumer: Failed to unmarshal booking request event: %v. Committing to skip.", err)
		if err := client.CommitRecords(ctx, record); err != nil {
			c.logger.Error("BookingConsumer: Failed to commit poison pill record: %v", err)
		}
		return
	}

	c.logger.Info("BookingConsumer: Processing booking request for bookingId: %s", event.BookingID)

	handleCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	redisStateKey := fmt.Sprintf("booking:state:%s", event.BookingID)

	stateExists, err := c.redisClient.Exists(handleCtx, redisStateKey).Result()
	if err != nil {
		c.logger.Error("BookingConsumer: Failed to check state existence for %s: %v", event.BookingID, err)
		return // Không commit để thử lại
	}
	if stateExists == 0 {
		c.logger.Info("BookingConsumer: State for bookingId %s no longer exists. Likely timed out. Skipping processing.", event.BookingID)
		if err := client.CommitRecords(ctx, record); err != nil {
			c.logger.Error("BookingConsumer: Failed to commit timed-out record: %v", err)
		}
		return
	}

	// 1. Cập nhật trạng thái trong Redis -> PROCESSING
	if err := c.redisClient.HSet(handleCtx, redisStateKey, "status", "PROCESSING").Err(); err != nil {
		c.logger.Error("BookingConsumer: Failed to update Redis state to PROCESSING for %s: %v", event.BookingID, err)
	}

	// 2. Gọi logic tạo vé cốt lõi
	ticket, err := c.ticketService.CreateTicket(handleCtx, &event.Input, event.CustomerID)

	var messageToPublish websocket.Message
	redisFields := make(map[string]interface{})

	// 3. Chuẩn bị kết quả
	if err != nil {
		c.logger.Error("[Notify] Failed to create ticket for bookingId %s: %v", event.BookingID, err)
		errorPayload := map[string]string{"error": err.Error()}
		messageToPublish = websocket.Message{Type: "error", Payload: errorPayload}
		redisFields["status"] = "FAILED"
		redisFields["error_message"] = err.Error()
	} else {
		c.logger.Info("[Notify] Ticket %s created successfully for bookingId %s.", ticket.TicketID, event.BookingID)
		ticketReturn, repoErr := c.ticketService.GetTicketByID(handleCtx, ticket.TicketID)
		if repoErr != nil {
			errorPayload := map[string]string{"error": "Không thể lấy chi tiết vé sau khi tạo."}
			messageToPublish = websocket.Message{Type: "error", Payload: errorPayload}
			redisFields["status"] = "FAILED"
			redisFields["error_message"] = "Không thể lấy chi tiết vé sau khi tạo."
		} else {
			messageToPublish = websocket.Message{Type: "result", Payload: ticketReturn}
			resultPayloadBytes, _ := json.Marshal(ticketReturn)
			redisFields["status"] = "COMPLETED"
			redisFields["ticket_id"] = ticket.TicketID
			redisFields["result_payload"] = string(resultPayloadBytes)
		}
	}

	// 4. Cập nhật trạng thái cuối cùng vào Redis Hash
	if err := c.redisClient.HSet(handleCtx, redisStateKey, redisFields).Err(); err != nil {
		c.logger.Error("BookingConsumer: Failed to update final Redis state for %s: %v", event.BookingID, err)
	}

	// 5. Publish kết quả lên kênh Pub/Sub để WebSocket client nhận
	payloadBytes, _ := json.Marshal(messageToPublish)
	redisChannel := fmt.Sprintf("booking-result:%s", event.BookingID)
	if pubErr := c.redisClient.Publish(context.Background(), redisChannel, payloadBytes).Err(); pubErr != nil {
		c.logger.Error("[Notify] Failed to publish result to Redis for bookingId %s: %v", event.BookingID, pubErr)
	} else {
		c.logger.Info("[Notify] Successfully published result to Redis channel %s", redisChannel)
	}

	// 6. Commit message Kafka sau khi đã xử lý xong
	if err := client.CommitRecords(ctx, record); err != nil {
		c.logger.Error("BookingConsumer: Failed to commit final record: %v", err)
	}
}

// file: ticket-service/internal/consumers/trip_consumer.go
package consumers

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"ticket-service/config"
	"ticket-service/internal/services"
	"ticket-service/pkg/utils"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

type TripCreatedEvent struct {
	TripID            string    `json:"tripId"`
	CreationTimestamp MilliTime `json:"creationTimestamp"`
}

type MilliTime time.Time

func (mt *MilliTime) UnmarshalJSON(b []byte) error {
	// 1. Đọc giá trị JSON vào một biến chuỗi.
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	// 2. Phân tích chuỗi theo định dạng chuẩn RFC3339Nano (ISO 8601).
	// Định dạng này rất linh hoạt, có thể xử lý cả trường hợp có và không có nano giây.
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		// Nếu thất bại, thử lại với định dạng không có nano giây để tăng tính tương thích.
		t, err = time.Parse(time.RFC3339, s)
		if err != nil {
			return err
		}
	}

	// 3. Gán giá trị đã parse thành công.
	*mt = MilliTime(t)
	return nil
}

type TripConsumer struct {
	managerTicketService services.IManagerTicketService
	logger               utils.Logger
	kafkaCfg             config.KafkaConfig
	topicCfg             config.TopicConfig
}

func NewTripConsumer(cfg config.Config, managerService services.IManagerTicketService, logger utils.Logger) *TripConsumer {
	return &TripConsumer{
		managerTicketService: managerService,
		logger:               logger,
		kafkaCfg:             cfg.Kafka,
		topicCfg:             cfg.Kafka.Topics.TripCreated,
	}
}

func (c *TripConsumer) Start(ctx context.Context) {
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

	client, err := kgo.NewClient(opts...)
	if err != nil {
		c.logger.Error("TripConsumer: Failed to create Kafka client: %v", err)
		return
	}
	defer client.Close()

	c.logger.Info("Starting Kafka consumer for topic: %s", c.topicCfg.Topic)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("TripConsumer: Shutting down.")
			return
		default:
			fetches := client.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				c.logger.Error("TripConsumer: Fetch errors: %v", errs)
				continue
			}

			fetches.EachRecord(func(record *kgo.Record) {
				var event TripCreatedEvent
				if err := json.Unmarshal(record.Value, &event); err != nil {
					c.logger.Error("TripConsumer: Failed to unmarshal event: %v. Skipping.", err)
					if err := client.CommitRecords(context.Background(), record); err != nil {
						c.logger.Error("TripConsumer: Failed to commit poison pill record: %v", err)
					}
					return
				}

				handleCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				_, err := c.managerTicketService.CreateManagerTicketsForTrip(handleCtx, event.TripID)
				if err != nil {
					c.logger.Error("Failed to create seats for TripID %s: %v. Not committing.", event.TripID, err)
					return // Không commit để thử lại
				}

				if err := client.CommitRecords(context.Background(), record); err != nil {
					c.logger.Error("TripConsumer: Failed to commit record: %v", err)
				}
			})
		}
	}
}

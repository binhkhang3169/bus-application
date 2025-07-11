// file: ticket-service/internal/consumers/ticket_status_consumer.go
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

type TicketStatusUpdateEvent struct {
	TicketID   string `json:"ticket_id"`
	StatusCode string `json:"status_code"`
}

type TicketStatusConsumer struct {
	managerTicketService services.IManagerTicketService
	logger               utils.Logger
	kafkaCfg             config.KafkaConfig
	topicCfg             config.TopicConfig
}

func NewTicketStatusConsumer(cfg config.Config, managerService services.IManagerTicketService, logger utils.Logger) *TicketStatusConsumer {
	return &TicketStatusConsumer{
		managerTicketService: managerService,
		logger:               logger,
		kafkaCfg:             cfg.Kafka,
		topicCfg:             cfg.Kafka.Topics.TicketStatusUpdates,
	}
}

func (c *TicketStatusConsumer) Start(ctx context.Context) {
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
		c.logger.Error("TicketStatusConsumer: Failed to create Kafka client: %v", err)
		return
	}
	defer client.Close()

	c.logger.Info("Starting Kafka consumer for topic: %s", c.topicCfg.Topic)

	for {
		select {
		case <-ctx.Done():
			c.logger.Info("TicketStatusConsumer: Shutting down.")
			return
		default:
			fetches := client.PollFetches(ctx)
			if errs := fetches.Errors(); len(errs) > 0 {
				c.logger.Error("TicketStatusConsumer: Fetch errors: %v", errs)
				continue
			}

			fetches.EachRecord(func(record *kgo.Record) {
				var event TicketStatusUpdateEvent
				if err := json.Unmarshal(record.Value, &event); err != nil {
					c.logger.Error("TicketStatusConsumer: Failed to unmarshal event: %v. Skipping.", err)
					if err := client.CommitRecords(context.Background(), record); err != nil {
						c.logger.Error("TicketStatusConsumer: Failed to commit poison pill record: %v", err)
					}
					return
				}

				handleCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				if err := c.managerTicketService.UpdateStatusByTicketID(handleCtx, event.TicketID, event.StatusCode); err != nil {
					c.logger.Error("Failed to process ticket status update for %s: %v. Not committing.", event.TicketID, err)
					return // Không commit để thử lại
				}

				// Commit sau khi xử lý thành công
				if err := client.CommitRecords(context.Background(), record); err != nil {
					c.logger.Error("TicketStatusConsumer: Failed to commit record: %v", err)
				}
			})
		}
	}
}

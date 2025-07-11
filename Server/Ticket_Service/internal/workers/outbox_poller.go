package workers

import (
	"context"
	"encoding/json"
	"time"

	"ticket-service/internal/db"
	"ticket-service/pkg/kafkaclient"
	"ticket-service/pkg/utils"

	"github.com/google/uuid"
)

type OutboxPoller struct {
	q         db.Querier
	publisher *kafkaclient.Publisher
	logger    utils.Logger
	interval  time.Duration
	batchSize int
}

func NewOutboxPoller(q db.Querier, publisher *kafkaclient.Publisher, logger utils.Logger) *OutboxPoller {
	return &OutboxPoller{
		q:         q,
		publisher: publisher,
		logger:    logger,
		interval:  2 * time.Second, // Poll every 2 seconds
		batchSize: 50,              // Process up to 50 events per batch
	}
}

func (p *OutboxPoller) Start(ctx context.Context) {
	p.logger.Info("Starting Outbox Poller worker...")
	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			p.logger.Info("Stopping Outbox Poller worker.")
			return
		case <-ticker.C:
			p.processBatch(ctx)
		}
	}
}

func (p *OutboxPoller) processBatch(ctx context.Context) {
	events, err := p.q.GetOutboxEvents(ctx, int32(p.batchSize))
	if err != nil {
		p.logger.Error("Failed to fetch outbox events: %v", err)
		return
	}

	if len(events) == 0 {
		return
	}

	p.logger.Info("Processing %d events from outbox...", len(events))
	var processedIDs []uuid.UUID

	for _, event := range events {
		var payloadData interface{}
		if err := json.Unmarshal(event.Payload, &payloadData); err != nil {
			p.logger.Error("Failed to unmarshal payload for event %s: %v. Marking as processed to avoid loop.", event.ID, err)
			processedIDs = append(processedIDs, event.ID) // Move to DLQ in a real scenario
			continue
		}

		// Try to publish to Kafka
		err := p.publisher.Publish(ctx, event.Topic, []byte(event.Key), payloadData)
		if err != nil {
			p.logger.Error("Failed to publish event %s to Kafka topic %s: %v. Will retry on next tick.", event.ID, event.Topic, err)
			// Don't add to processedIDs, so it will be retried
			continue
		}

		processedIDs = append(processedIDs, event.ID)
	}

	// Delete successfully processed events from the outbox
	if len(processedIDs) > 0 {
		err := p.q.DeleteOutboxEvents(ctx, processedIDs)
		if err != nil {
			p.logger.Error("CRITICAL: Failed to delete processed events from outbox: %v", err)
			// This is a critical error that requires monitoring and alerting.
		}
	}
}

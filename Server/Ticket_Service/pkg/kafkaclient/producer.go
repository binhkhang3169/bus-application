// file: ticket-service/pkg/kafkaclient/publisher.go
package kafkaclient

import (
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"fmt"
	"ticket-service/config"
	"ticket-service/domain/models"
	"ticket-service/pkg/utils"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// Publisher là một client chung để gửi message tới Kafka.
type Publisher struct {
	client *kgo.Client
	logger utils.Logger
}

// NewPublisher tạo một publisher mới với franz-go
func NewPublisher(cfg config.KafkaConfig, logger utils.Logger) (*Publisher, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Seeds...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.ProducerBatchMaxBytes(1e6),
	}

	if cfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

	if cfg.SASLUser != "" && cfg.SASLPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{
			User: cfg.SASLUser,
			Pass: cfg.SASLPass,
		}.AsSha256Mechanism()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo kafka client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("không thể ping đến kafka brokers: %w", err)
	}

	logger.Info("Kafka publisher client đã được khởi tạo thành công.")
	return &Publisher{client: client, logger: logger}, nil
}

// Publish gửi một message tới một topic được chỉ định.
func (p *Publisher) Publish(ctx context.Context, topic string, key []byte, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		p.logger.Error("KafkaPublisher: Failed to marshal payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	record := &kgo.Record{Topic: topic, Key: key, Value: payloadBytes}

	p.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			p.logger.Error("KafkaPublisher: Failed to publish message to topic %s: %v", topic, err)
		} else {
			p.logger.Info("KafkaPublisher: Published message to topic %s with key %s", topic, string(key))
		}
	})

	return nil
}

// Close đóng client Kafka.
func (p *Publisher) Close() error {
	p.client.Close()
	return nil
}

// -- Định nghĩa các DTO cho các sự kiện (giữ nguyên) --

type SeatUpdateEvent struct {
	TripID    string `json:"tripId"`
	SeatCount int    `json:"seatCount"`
}

type EmailRequestEvent struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
	Type  string `json:"type"`
}

type QRGenerationRequestEvent struct {
	TicketID      string  `json:"ticketId"`
	SeatID        int32   `json:"seatId"`
	QRContent     string  `json:"qrContent"`
	CustomerEmail string  `json:"customerEmail"`
	CustomerName  string  `json:"customerName"`
	Price         float64 `json:"price"`
}

type TicketDetailForQR struct {
	SeatID    int32  `json:"seatId"`
	QRContent string `json:"qrContent"`
}

type OrderQRGenerationRequestEvent struct {
	OrderID       string              `json:"orderId"`
	CustomerEmail string              `json:"customerEmail"`
	CustomerName  string              `json:"customerName"`
	TotalPrice    float64             `json:"totalPrice"`
	Tickets       []TicketDetailForQR `json:"tickets"`
}
type SeatsReservedEvent struct {
	TripID    string    `json:"tripId"`
	TicketID  string    `json:"ticketId"`
	SeatIDs   []int32   `json:"seatIds"`
	Timestamp time.Time `json:"timestamp"`
}

type BookingRequestEvent struct {
	BookingID  string             `json:"booking_id"`
	Input      models.TicketInput `json:"input"`
	CustomerID sql.NullInt32      `json:"customer_id"`
}

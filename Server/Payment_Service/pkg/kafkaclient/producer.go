// payment-service/pkg/kafkaclient/producer.go
package kafkaclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"payment_service/config" // THAY ĐỔI: Import config
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// Publisher là một client chung để gửi message tới Kafka sử dụng franz-go.
type Publisher struct {
	client *kgo.Client
}

// THAY ĐỔI: NewPublisher nhận cấu trúc config thay vì chuỗi URL
func NewPublisher(cfg config.KafkaConfig) (*Publisher, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Seeds...),
		kgo.RequiredAcks(kgo.AllISRAcks()), // Đảm bảo độ tin cậy cao nhất
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

	// Ping để kiểm tra kết nối ngay lúc khởi tạo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("không thể ping đến kafka brokers: %w", err)
	}

	log.Println("Kafka publisher client đã được khởi tạo thành công.")
	return &Publisher{client: client}, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, key []byte, payload interface{}) error {
	bgCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("không thể marshal payload: %w", err)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: payloadBytes,
	}

	err = p.client.ProduceSync(bgCtx, record).FirstErr()
	if err != nil {
		return fmt.Errorf("gửi message đến topic %s thất bại: %w", topic, err)
	}

	return nil
}

func (p *Publisher) Close() {
	p.client.Close()
}

// --- DTOs cho các sự kiện (giữ nguyên) ---
type NotificationEvent struct {
	UserID  *string `json:"user_id"`
	Type    string  `json:"type"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
}

type TicketStatusUpdateEvent struct {
	TicketID   string `json:"ticket_id"`
	StatusCode string `json:"status_code"`
}

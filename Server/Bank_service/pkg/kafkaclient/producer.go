// file: payment-service/pkg/kafkaclient/publisher.go
package kafkaclient

import (
	"bank/config" // Import config để lấy thông tin kafka
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// Publisher là một client chung để gửi message tới Kafka sử dụng franz-go.
type Publisher struct {
	client *kgo.Client
}

// NewPublisher tạo một publisher mới kết nối tới Kafka brokers.
// Nó sử dụng cấu hình từ `config.Config` để thiết lập kết nối, bao gồm TLS và SASL.
func NewPublisher(cfg config.Config) (*Publisher, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.KafkaSeeds...),
		// Cấu hình retry vô hạn cho producer để đảm bảo message được gửi đi.
		kgo.ProducerBatchMaxBytes(1e6),
		kgo.ProduceRequestTimeout(10 * time.Second),
		kgo.RetryTimeout(time.Second * 30),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	}

	if cfg.KafkaEnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

	// Cấu hình SASL/SCRAM-SHA-256 nếu có user/pass
	if cfg.KafkaSASLUser != "" && cfg.KafkaSASLPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{
			User: cfg.KafkaSASLUser,
			Pass: cfg.KafkaSASLPass,
		}.AsSha256Mechanism()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create kafka client: %w", err)
	}

	// Ping để kiểm tra kết nối
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping kafka brokers: %w", err)
	}

	log.Println("Kafka client initialized successfully.")
	return &Publisher{client: client}, nil
}

// Publish gửi một message tới một topic Kafka một cách bất đồng bộ.
// Lỗi (nếu có) sẽ được log lại trong hàm callback.
func (p *Publisher) Publish(ctx context.Context, topic string, key []byte, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Printf("KafkaPublisher: Failed to marshal payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: payloadBytes,
	}

	// client.Produce là lệnh bất đồng bộ (non-blocking).
	// Nó sẽ gửi message trong background.
	// Hàm callback sẽ được gọi khi message được produce thành công hoặc thất bại.
	err = p.client.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		log.Printf("CRITICAL: KafkaPublisher: Failed to publish message to topic %s: %v", topic, err)
		return err
	}

	return nil
}

// Close đóng client Kafka và giải phóng tài nguyên.
func (p *Publisher) Close() {
	// Close Flushes, then closes the client.
	p.client.Close()
}

// --- DTOs cho các sự kiện (Không thay đổi) ---

// NotificationEvent là payload cho sự kiện gửi thông báo.
type NotificationEvent struct {
	UserID  *string `json:"user_id"`
	Type    string  `json:"type"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
}

// TicketStatusUpdateEvent là payload cho sự kiện cập nhật trạng thái vé.
type TicketStatusUpdateEvent struct {
	TicketID   string `json:"ticket_id"`
	StatusCode string `json:"status_code"`
}

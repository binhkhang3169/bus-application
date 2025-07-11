package kafkaclient

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"news-management/config" // Hoặc "bank/config" tùy service
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// Publisher là một client chung để gửi message tới Kafka sử dụng franz-go.
type Publisher struct {
	client *kgo.Client
}

// NewPublisher tạo một publisher mới kết nối tới Kafka brokers.
func NewPublisher(cfg *config.Config) (*Publisher, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.KafkaSeeds...),
		kgo.ProducerBatchMaxBytes(1e6),
		kgo.RequiredAcks(kgo.AllISRAcks()), // Đảm bảo message được ghi vào tất cả replica
	}

	if cfg.KafkaEnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}

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

	// Ping để kiểm tra kết nối ngay lúc khởi tạo
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := client.Ping(ctx); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed to ping kafka brokers: %w", err)
	}

	log.Println("Kafka publisher client initialized successfully.")
	return &Publisher{client: client}, nil
}

// Publish gửi một message tới một topic Kafka một cách bất đồng bộ.
func (p *Publisher) Publish(ctx context.Context, topic string, key []byte, payload interface{}) error {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	record := &kgo.Record{
		Topic: topic,
		Key:   key,
		Value: payloadBytes,
	}

	// client.Produce là lệnh bất đồng bộ. Nó sẽ gửi message trong background.
	p.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Printf("CRITICAL: KafkaPublisher: Failed to publish message to topic %s: %v", topic, err)
		}
	})

	return nil
}

func (p *Publisher) Close() {
	p.client.Close()
}

// DTO cho sự kiện thông báo (giữ nguyên)
type NotificationEvent struct {
	UserID  *string `json:"user_id"`
	Type    string  `json:"type"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
}

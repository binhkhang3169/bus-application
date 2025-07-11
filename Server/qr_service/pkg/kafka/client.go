// file: qr-service/pkg/kafka/client.go
package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"qr/config"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

// Publisher dùng để gửi message đi
type Publisher struct {
	client *kgo.Client
}

// NewPublisher tạo một publisher mới kết nối an toàn đến Redpanda
func NewPublisher(cfg config.KafkaConfig) (*Publisher, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Seeds...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
	}
	if cfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}
	if cfg.SASLUser != "" && cfg.SASLPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{User: cfg.SASLUser, Pass: cfg.SASLPass}.AsSha256Mechanism()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo kafka publisher client: %w", err)
	}

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
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("không thể marshal payload: %w", err)
	}

	record := &kgo.Record{Topic: topic, Key: key, Value: payloadBytes}

	p.client.Produce(ctx, record, func(r *kgo.Record, err error) {
		if err != nil {
			log.Printf("Lỗi gửi message đến topic %s: %v\n", topic, err)
		}
	})
	return nil
}

func (p *Publisher) Close() {
	p.client.Close()
}

// NewConsumerClient tạo một consumer client mới
func NewConsumerClient(cfg config.KafkaConfig, topic, groupID string) (*kgo.Client, error) {
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.Seeds...),
		kgo.ConsumerGroup(groupID),
		kgo.ConsumeTopics(topic),
		kgo.DisableAutoCommit(),
	}
	if cfg.EnableTLS {
		opts = append(opts, kgo.DialTLSConfig(new(tls.Config)))
	}
	if cfg.SASLUser != "" && cfg.SASLPass != "" {
		opts = append(opts, kgo.SASL(scram.Auth{User: cfg.SASLUser, Pass: cfg.SASLPass}.AsSha256Mechanism()))
	}

	client, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo kafka consumer client: %w", err)
	}
	return client, nil
}

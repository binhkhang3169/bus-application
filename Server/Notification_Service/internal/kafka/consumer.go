// notification-service/internal/kafka/consumer.go
package kafka

import (
	"context"
	"crypto/tls" // THÊM MỚI
	"encoding/json"
	"log"
	"notification-service/config" // THÊM MỚI
	"notification-service/internal/model"
	"notification-service/internal/service"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"        // THAY ĐỔI: Thư viện mới
	"github.com/twmb/franz-go/pkg/sasl/scram" // THÊM MỚI
)

// NotificationMessage giữ nguyên cấu trúc
type NotificationMessage struct {
	UserID  *string `json:"user_id"`
	Type    string  `json:"type"`
	Title   string  `json:"title"`
	Message string  `json:"message"`
}

// StartConsumer được viết lại hoàn toàn để sử dụng franz-go
func StartConsumer(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, svc service.NotificationService) {
	defer wg.Done()

	// THAY ĐỔI: Xây dựng các tùy chọn cho client của franz-go
	opts := []kgo.Opt{
		kgo.SeedBrokers(cfg.KafkaSeeds...),
		kgo.ConsumerGroup(cfg.KafkaGroupID),
		kgo.ConsumeTopics(cfg.KafkaTopic),
		kgo.DisableAutoCommit(), // Tắt auto-commit để kiểm soát thủ công
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
		log.Fatalf("Error creating Kafka client: %v", err)
	}
	defer client.Close()

	log.Printf("Kafka Consumer Group started for topic '%s' and groupID '%s'", cfg.KafkaTopic, cfg.KafkaGroupID)

	// Vòng lặp chính để nhận message
	for {
		// Kiểm tra nếu context đã bị hủy (tín hiệu shutdown)
		select {
		case <-ctx.Done():
			log.Println("Kafka consumer shutting down...")
			return
		default:
			// Tiếp tục vòng lặp
		}

		fetches := client.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Printf("Error polling fetches: %v", errs)
			continue
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()
			log.Printf("Kafka message claimed: value = %s, topic = %s, partition = %d, offset = %d",
				string(record.Value), record.Topic, record.Partition, record.Offset)

			var kafkaMsg NotificationMessage
			if err := json.Unmarshal(record.Value, &kafkaMsg); err != nil {
				log.Printf("Error unmarshalling Kafka message: %v. Skipping (poison pill).", err)
				// Commit message lỗi để không xử lý lại
				if err := client.CommitRecords(ctx, record); err != nil {
					log.Printf("Failed to commit poison pill message: %v", err)
				}
				continue
			}

			createReq := model.CreateNotificationRequest{
				UserID:  kafkaMsg.UserID,
				Type:    kafkaMsg.Type,
				Title:   kafkaMsg.Title,
				Message: kafkaMsg.Message,
			}

			// Gọi service để xử lý nghiệp vụ
			_, err := svc.CreateNotification(context.Background(), createReq)
			if err != nil {
				log.Printf("Error processing message, will not commit and retry later: %v", err)
				// Không commit khi có lỗi, Kafka sẽ gửi lại message này
				continue
			}

			// Commit offset thủ công sau khi xử lý thành công
			log.Printf("Successfully processed message for user %v. Committing offset %d.", kafkaMsg.UserID, record.Offset)
			if err := client.CommitRecords(ctx, record); err != nil {
				log.Printf("Failed to commit message after successful processing: %v", err)
			}
		}
	}
}

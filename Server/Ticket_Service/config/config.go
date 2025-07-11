// file: ticket-service/config/config.go
package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// TopicConfig chứa cấu hình cho một topic Kafka cụ thể.
type TopicConfig struct {
	Topic   string
	GroupID string
}

// KafkaTopics định nghĩa tất cả các topic mà service tương tác.
type KafkaTopics struct {
	TripCreated         TopicConfig
	TicketStatusUpdates TopicConfig
	SeatsReserved       TopicConfig
	SeatsReleased       TopicConfig
	EmailRequests       TopicConfig
	OrderQRRequests     TopicConfig
	BookingRequests     TopicConfig
}

// THAY ĐỔI: Cấu trúc KafkaConfig được thiết kế lại cho franz-go
type KafkaConfig struct {
	Seeds     []string
	EnableTLS bool
	SASLUser  string
	SASLPass  string
	Topics    KafkaTopics
}

type Config struct {
	Server struct {
		Port         string
		ReadTimeout  time.Duration
		WriteTimeout time.Duration
		GinMode      string
	}
	Database struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	// THAY ĐỔI: Redis config giờ chỉ cần URL
	Redis struct {
		URL string
	}
	JWT struct {
		SecretKey string
	}
	URL struct {
		QRService string `mapstructure:"QR_SERVICE_URL"`
	}
	Kafka KafkaConfig
}

func LoadConfig() (Config, error) {
	_ = godotenv.Load()
	var cfg Config

	// Server configuration
	cfg.Server.Port = GetEnv("SERVER_PORT", "8084")
	readTimeout, _ := strconv.Atoi(GetEnv("SERVER_READ_TIMEOUT", "15"))
	writeTimeout, _ := strconv.Atoi(GetEnv("SERVER_WRITE_TIMEOUT", "15"))
	cfg.Server.ReadTimeout = time.Duration(readTimeout) * time.Second
	cfg.Server.WriteTimeout = time.Duration(writeTimeout) * time.Second
	cfg.Server.GinMode = GetEnv("GIN_MODE", "debug")

	// Database configuration
	cfg.Database.Host = GetEnv("DATABASE_HOST", "localhost")
	cfg.Database.Port = GetEnv("DATABASE_PORT", "5432")
	cfg.Database.User = GetEnv("DATABASE_USER", "postgres")
	cfg.Database.Password = GetEnv("DATABASE_PASSWORD", "postgres")
	cfg.Database.Name = GetEnv("DATABASE_NAME", "ticket_db")
	cfg.Database.SSLMode = GetEnv("DATABASE_SSL_MODE", "disable")

	// THAY ĐỔI: Cấu hình Redis
	cfg.Redis.URL = GetEnv("REDIS_URL", "redis://localhost:6379")

	cfg.JWT.SecretKey = GetEnv("JWT_TOKEN", "your-very-secret-key-for-jwt")
	cfg.URL.QRService = GetEnv("QR_SERVICE_URL", "http://localhost:8090/api/v1/qr/generate")

	// THAY ĐỔI: Load cấu hình Kafka mới
	kafkaEnableTLS, _ := strconv.ParseBool(GetEnv("KAFKA_ENABLE_TLS", "false"))
	cfg.Kafka.Seeds = strings.Split(GetEnv("KAFKA_SEEDS", "localhost:9092"), ",")
	cfg.Kafka.EnableTLS = kafkaEnableTLS
	cfg.Kafka.SASLUser = GetEnv("KAFKA_SASL_USER", "")
	cfg.Kafka.SASLPass = GetEnv("KAFKA_SASL_PASS", "")

	// Load topics configuration (giữ nguyên logic đọc từ env)
	cfg.Kafka.Topics.TripCreated.Topic = GetEnv("KAFKA_TOPIC_TRIP_CREATED", "trip_created")
	cfg.Kafka.Topics.TripCreated.GroupID = GetEnv("KAFKA_GROUP_ID_TRIP_CREATED", "ticket_service_trip_group")

	cfg.Kafka.Topics.TicketStatusUpdates.Topic = GetEnv("KAFKA_TOPIC_TICKET_STATUS", "ticket_status_updates")
	cfg.Kafka.Topics.TicketStatusUpdates.GroupID = GetEnv("KAFKA_GROUP_ID_TICKET_STATUS", "ticket_service_status_group")

	cfg.Kafka.Topics.BookingRequests.Topic = GetEnv("KAFKA_TOPIC_BOOKING_REQUESTS", "booking_requests")
	cfg.Kafka.Topics.BookingRequests.GroupID = GetEnv("KAFKA_GROUP_ID_BOOKING_REQUESTS", "ticket_service_booking_group")

	cfg.Kafka.Topics.SeatsReserved.Topic = GetEnv("KAFKA_TOPIC_SEATS_RESERVED", "seats_reserved")
	cfg.Kafka.Topics.SeatsReleased.Topic = GetEnv("KAFKA_TOPIC_SEATS_RELEASED", "seats_released")
	cfg.Kafka.Topics.EmailRequests.Topic = GetEnv("KAFKA_TOPIC_EMAIL_REQUESTS", "email_requests")
	cfg.Kafka.Topics.OrderQRRequests.Topic = GetEnv("KAFKA_TOPIC_QR_REQUESTS", "order_qr_requests")

	return cfg, nil
}

func GetEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// file: qr-service/config/config.go
package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type KafkaConfig struct {
	Seeds                []string
	EnableTLS            bool
	SASLUser             string
	SASLPass             string
	TopicOrderQRRequests string
	GroupIDOrderQR       string
	TopicEmailRequests   string
}

type CloudinaryConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

type Config struct {
	Server     ServerConfig
	Cloudinary CloudinaryConfig
	Kafka      KafkaConfig
}

func Load() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config

	cfg.Server.Port = getEnv("PORT", "8090")
	cfg.Cloudinary.URL = getEnv("CLOUDINARY_URL", "")

	// Load Kafka Config
	enableTLS, _ := strconv.ParseBool(getEnv("KAFKA_ENABLE_TLS", "false"))
	cfg.Kafka.EnableTLS = enableTLS
	cfg.Kafka.Seeds = strings.Split(getEnv("KAFKA_SEEDS", "localhost:9092"), ",")
	cfg.Kafka.SASLUser = getEnv("KAFKA_SASL_USER", "")
	cfg.Kafka.SASLPass = getEnv("KAFKA_SASL_PASS", "")

	// Load Kafka Topics Config
	cfg.Kafka.TopicOrderQRRequests = getEnv("KAFKA_TOPIC_ORDER_QR_REQUESTS", "order_qr_requests")
	cfg.Kafka.GroupIDOrderQR = getEnv("KAFKA_GROUP_ID_QR_ORDER", "qr_service_order_group")
	cfg.Kafka.TopicEmailRequests = getEnv("KAFKA_TOPIC_EMAIL_REQUESTS", "email_requests")

	return &cfg, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

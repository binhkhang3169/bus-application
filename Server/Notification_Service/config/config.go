// notification-service/config/config.go
package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config struct holds all configuration for the application
type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSslmode      string
	KafkaSeeds     []string
	KafkaTopic     string
	KafkaGroupID   string
	KafkaEnableTLS bool
	KafkaSASLUser  string
	KafkaSASLPass  string
	HTTPPort       string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file. Path is relative to the executable or working directory.
	// If running main.go directly from cmd/ it would be ../.env
	// Adjust the path if necessary, e.g. for tests or different execution contexts.
	err := godotenv.Load() // Tự động tìm .env ở thư mục gốc dự án khi chạy từ gốc
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	dbPortStr := getEnv("DB_PORT", "5437")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}
	kafkaEnableTLS, _ := strconv.ParseBool(getEnv("KAFKA_ENABLE_TLS", "false"))

	return &Config{
		DBHost:         getEnv("DB_HOST", "postgres_noti"),
		DBPort:         dbPort,
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "noti_service"),
		DBSslmode:      getEnv("DB_SSLMODE", "disable"),
		KafkaSeeds:     strings.Split(getEnv("KAFKA_SEEDS", "localhost:9092"), ","),
		KafkaTopic:     getEnv("KAFKA_TOPIC_NOTIFICATIONS", "notifications_topic"),
		KafkaGroupID:   getEnv("KAFKA_GROUP_ID", "notification_service_group"),
		KafkaEnableTLS: kafkaEnableTLS,
		KafkaSASLUser:  getEnv("KAFKA_SASL_USER", ""),
		KafkaSASLPass:  getEnv("KAFKA_SASL_PASS", ""),
		HTTPPort:       getEnv("HTTP_PORT", "8080"),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

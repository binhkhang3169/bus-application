package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config lưu trữ tất cả cấu hình của ứng dụng.
type Config struct {
	DBDriver         string
	DBHost           string
	DBPort           int
	DBUser           string
	DBPassword       string
	DBName           string
	DBSSLMode        string
	DBMaxConnections int
	DBMinConnections int
	ServerPort       int
	JWTSecret        string

	// Kafka configuration for franz-go
	KafkaSeeds     []string // Thay thế cho kafkaURL
	KafkaEnableTLS bool
	KafkaSASLUser  string
	KafkaSASLPass  string
}

// LoadConfig nạp cấu hình từ file .env và biến môi trường.
func LoadConfig(path string) (Config, error) {
	err := godotenv.Load(fmt.Sprintf("%s/.env", path))
	if err != nil {
		// Không tìm thấy file .env cũng không sao, có thể dùng biến môi trường đã có
	}

	dbPort, _ := strconv.Atoi(os.Getenv("DB_PORT"))
	dbMaxConns, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	dbMinConns, _ := strconv.Atoi(os.Getenv("DB_MIN_CONNECTIONS"))
	serverPort, _ := strconv.Atoi(os.Getenv("SERVER_PORT"))
	kafkaEnableTLS, _ := strconv.ParseBool(os.Getenv("KAFKA_ENABLE_TLS"))

	// Đọc KAFKA_SEEDS dưới dạng chuỗi phân tách bằng dấu phẩy
	kafkaSeeds := strings.Split(os.Getenv("KAFKA_SEEDS"), ",")

	config := Config{
		DBDriver:         os.Getenv("DB_DRIVER"),
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           dbPort,
		DBUser:           os.Getenv("DB_USER"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		DBName:           os.Getenv("DB_NAME"),
		DBSSLMode:        os.Getenv("DB_SSL_MODE"),
		DBMaxConnections: dbMaxConns,
		DBMinConnections: dbMinConns,
		ServerPort:       serverPort,
		JWTSecret:        os.Getenv("JWT_SECRET"),

		// Kafka settings
		KafkaSeeds:     kafkaSeeds,
		KafkaEnableTLS: kafkaEnableTLS,
		KafkaSASLUser:  os.Getenv("KAFKA_SASL_USER"),
		KafkaSASLPass:  os.Getenv("KAFKA_SASL_PASS"),
	}

	return config, nil
}

// DBSource trả về chuỗi kết nối PostgreSQL.
func (cfg *Config) DBSource() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName, cfg.DBSSLMode,
	)
}

// ServerAddress trả về địa chỉ để chạy HTTP server.
func (cfg *Config) ServerAddress() string {
	return fmt.Sprintf("0.0.0.0:%d", cfg.ServerPort)
}

// NOTE: KafkaURL() has been removed as franz-go uses Seed Brokers.

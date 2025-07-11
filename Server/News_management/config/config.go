package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config chứa các thông tin cấu hình ứng dụng
type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSslmode  string
	ServerPort string
	KafkaURL   string

	KafkaSeeds     []string
	KafkaEnableTLS bool
	KafkaSASLUser  string
	KafkaSASLPass  string
}

// LoadConfig tải cấu hình từ file .env
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbSslmode := os.Getenv("DB_SSLMODE")
	serverPort := os.Getenv("SERVER_PORT")
	kafkaURL := os.Getenv("KAFKA_URL")
	kafkaSeedsStr := os.Getenv("KAFKA_SEEDS")
	kafkaEnableTLS, _ := strconv.ParseBool(os.Getenv("KAFKA_ENABLE_TLS"))

	if dbHost == "" || dbPortStr == "" || dbUser == "" || dbName == "" {
		return nil, fmt.Errorf("database configuration environment variables are not fully set")
	}
	if serverPort == "" {
		serverPort = "8080" // Default port
	}

	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		return nil, fmt.Errorf("invalid DB_PORT: %w", err)
	}

	return &Config{
		DBHost:         dbHost,
		DBPort:         dbPort,
		DBUser:         dbUser,
		DBPassword:     dbPassword,
		DBName:         dbName,
		DBSslmode:      dbSslmode,
		ServerPort:     serverPort,
		KafkaURL:       kafkaURL,
		KafkaSeeds:     strings.Split(kafkaSeedsStr, ","),
		KafkaEnableTLS: kafkaEnableTLS,
		KafkaSASLUser:  os.Getenv("KAFKA_SASL_USER"),
		KafkaSASLPass:  os.Getenv("KAFKA_SASL_PASS"),
	}, nil
}

// DBSource tạo chuỗi kết nối DSN cho PostgreSQL
func (c *Config) DBSource() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, c.DBSslmode)
}

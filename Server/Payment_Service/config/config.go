package config

import (
	"os"
	"strconv"
	"strings"
)

// Config holds the application configuration
type Config struct {
	Server        ServerConfig
	Database      DatabaseConfig
	VNPay         VNPayConfig
	JWT           JWTConfig
	Stripe        StripeConfig        // <--- Thêm dòng này
	TicketService TicketServiceConfig // <--- Thêm dòng này
	KafkaConfig   KafkaConfig
	RedisConfig   RedisConfig
}

// ServerConfig holds the server configuration
type ServerConfig struct {
	Port string
}

// DatabaseConfig holds the database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Driver   string
}

// VNPayConfig holds the configuration for VNPAY integration
type VNPayConfig struct {
	TmnCode        string
	HashSecret     string
	VNPayURL       string
	ReturnURL      string
	APIUrl         string
	MerchantAPI    string
	TransactionAPI string
}

type JWTConfig struct {
	SecretKey string
}
type StripeConfig struct {
	SecretKey      string `mapstructure:"STRIPE_SECRET_KEY"`
	PublishableKey string `mapstructure:"STRIPE_PUBLISHABLE_KEY"`
	WebhookSecret  string `mapstructure:"STRIPE_WEBHOOK_SECRET"`
}

type TicketServiceConfig struct {
	URL string `mapstructure:"TICKET_SERVICE_URL"`
}

type KafkaConfig struct {
	Seeds     []string `mapstructure:"KAFKA_SEEDS"`
	Topic     string   `mapstructure:"KAFKA_TOPIC"`
	EnableTLS bool     `mapstructure:"KAFKA_ENABLE_TLS"`
	SASLUser  string   `mapstructure:"KAFKA_SASL_USER"`
	SASLPass  string   `mapstructure:"KAFKA_SASL_PASS"`
}

type RedisConfig struct {
	URL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	kafkaEnableTLS, _ := strconv.ParseBool(getEnv("KAFKA_ENABLE_TLS", "false"))

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8083"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "postgres_payment"),
			Port:     getEnv("DB_PORT", "5433"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "payment_service"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
			Driver:   getEnv("DRIVER", "postgres"),
		},
		VNPay: VNPayConfig{
			TmnCode:        getEnv("VNPAY_TMN_CODE", "NSD601AI"),
			HashSecret:     getEnv("VNPAY_HASH_SECRET", "2HEAB5G3TIDCN8R93EJQ6BP2967TZTNS"),
			VNPayURL:       getEnv("VNPAY_URL", "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"),
			ReturnURL:      getEnv("VNPAY_RETURN_URL", "http://localhost:3000/ket-qua-dat-ve"),
			APIUrl:         getEnv("VNPAY_API_URL", "http://sandbox.vnpayment.vn/merchant_webapi/merchant.html"),
			TransactionAPI: getEnv("VNPAY_TRANSACTION_API", "https://sandbox.vnpayment.vn/merchant_webapi/api/transaction"),
		},
		JWT: JWTConfig{
			SecretKey: getEnv("JWT_TOKEN", ""),
		},
		Stripe: StripeConfig{
			SecretKey:      getEnv("STRIPE_SECRET_KEY", "sk_test_YOUR_STRIPE_SECRET_KEY"),
			PublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", "pk_test_YOUR_STRIPE_PUBLISHABLE_KEY"),
			WebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", "whsec_YOUR_STRIPE_WEBHOOK_SECRET"),
		},
		TicketService: TicketServiceConfig{
			URL: getEnv("TICKET_SERVICE_URL", "TICKET_SERVICE_URL=http://ticket_service:8084/api/v1/payments"), // Ví dụ URL
		},
		KafkaConfig: KafkaConfig{
			Seeds:     strings.Split(getEnv("KAFKA_SEEDS", "localhost:9092"), ","),
			Topic:     getEnv("KAFKA_TOPIC", "ticket_status_updates"),
			EnableTLS: kafkaEnableTLS,
			SASLUser:  getEnv("KAFKA_SASL_USER", ""),
			SASLPass:  getEnv("KAFKA_SASL_PASS", ""),
		},
		// THAY ĐỔI: Gán giá trị cho Redis URL
		RedisConfig: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
	}
}

// GetDatabaseDSN returns the database connection string
func (c *DatabaseConfig) GetDatabaseDSN() string {
	return "postgres://" + c.User + ":" + c.Password + "@" + c.Host + ":" + c.Port + "/" + c.DBName + "?sslmode=" + c.SSLMode
}

// Helper function to get environment variable with a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Helper function to get integer environment variable with a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

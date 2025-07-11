package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file (if present) and system environment.
func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("Info: .env file not found or error loading. Using environment variables or defaults.")
	}
}

// GetEnv retrieves an environment variable or returns a default value.
func GetEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Env variable '%s' not found, using default value: '%s'", key, defaultValue)
	return defaultValue
}

// JWTConfig holds JWT related configuration.
type JWTConfig struct {
	SecretKey string
}

// LoadJWTConfig loads JWT configuration from environment variables.
func LoadJWTConfig() JWTConfig {
	secretKey := GetEnv("JWT_SECRET_KEY", "YOUR_FALLBACK_SECRET_KEY_CHANGE_ME_IN_PROD")
	if secretKey == "YOUR_FALLBACK_SECRET_KEY_CHANGE_ME_IN_PROD" {
		log.Println("WARNING: Using a default JWT_SECRET_KEY. This is INSECURE. Set a strong secret in your environment for production.")
	}
	return JWTConfig{SecretKey: secretKey}
}

// ServiceURLs holds URLs for downstream services.
type ServiceURLs struct {
	EmailServiceURL        string
	PaymentServiceURL      string
	TripServiceURL         string
	TicketServiceURL       string
	UserServiceURL         string
	BankServiceURL         string
	NewsServiceURL         string
	ShipServiceURL         string
	NotificationServiceURL string
	QrServiceURL           string
	ChatServiceURL         string
	DashboardURL           string
}

// LoadServiceURLs loads service URLs from environment variables.
func LoadServiceURLs() ServiceURLs {
	return ServiceURLs{
		EmailServiceURL:        GetEnv("EMAIL_SERVICE_URL", "http://localhost:8085"),
		PaymentServiceURL:      GetEnv("PAYMENT_SERVICE_URL", "http://localhost:8083"),
		TripServiceURL:         GetEnv("TRIP_SERVICE_URL", "http://localhost:8082"),
		TicketServiceURL:       GetEnv("TICKET_SERVICE_URL", "http://localhost:8084"),
		UserServiceURL:         GetEnv("USER_SERVICE_URL", "http://localhost:8081"),
		BankServiceURL:         GetEnv("BANK_SERVICE_URL", "http://localhost:8086"),
		NewsServiceURL:         GetEnv("NEWS_SERVICE_URL", "http://localhost:8087"),
		ShipServiceURL:         GetEnv("SHIP_SERVICE_URL", "http://localhost:8088"),
		NotificationServiceURL: GetEnv("NOTIFICATION_SERVICE_URL", "http://localhost:8089"),
		QrServiceURL:           GetEnv("QR_SERVICE_URL", "http://localhost:8090"),
		ChatServiceURL:         GetEnv("CHAT_BOT_SERVICE_URL", "http://localhost:8091"),
		DashboardURL:           GetEnv("DASHBOARD_SERVICE_URL", "http://localhost:8091"),
	}
}

// ServerConfig holds server related configuration.
type ServerConfig struct {
	Port string
}

// LoadServerConfig loads server configuration from environment variables.
func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Port: GetEnv("PORT", "8000"),
	}
}

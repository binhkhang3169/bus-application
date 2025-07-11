package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application.
type Config struct {
	DBDriver                string
	DBSource                string
	ServerAddress           string
	BaseRatePerKg           float64
	DimensionalWeightFactor float64
	ItemTypeMultipliers     map[string]float64
}

// LoadConfig reads configuration from environment variables or a .env file.
func LoadConfig(path string) (*Config, error) {
	// Attempt to load .env file. If it doesn't exist, continue with environment variables.
	if err := godotenv.Load(path); err != nil {
		log.Printf("Warning: Could not load .env file from %s. Using environment variables: %v", path, err)
	}

	cfg := &Config{
		ItemTypeMultipliers: make(map[string]float64),
	}

	cfg.DBDriver = getEnv("DB_DRIVER", "postgres")
	// Construct DB_DSN from individual parts if DB_DSN is not set directly
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "postgres")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "service")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	defaultDsn := fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBDriver, dbUser, dbPassword, dbHost, dbPort, dbName, dbSSLMode)
	cfg.DBSource = defaultDsn

	cfg.ServerAddress = getEnv("SERVER_ADDRESS", "0.0.0.0:8080")

	baseRateStr := getEnv("BASE_RATE_PER_KG", "5.0")
	baseRate, err := strconv.ParseFloat(baseRateStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid BASE_RATE_PER_KG '%s': %w", baseRateStr, err)
	}
	cfg.BaseRatePerKg = baseRate

	dimFactorStr := getEnv("DIMENSIONAL_WEIGHT_FACTOR", "5000.0")
	dimFactor, err := strconv.ParseFloat(dimFactorStr, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid DIMENSIONAL_WEIGHT_FACTOR '%s': %w", dimFactorStr, err)
	}
	cfg.DimensionalWeightFactor = dimFactor

	// Load item type multipliers from environment variables prefixed with ITEM_TYPE_MULTIPLIER_
	// Example: ITEM_TYPE_MULTIPLIER_DOCUMENT=1.0
	for _, envVar := range os.Environ() {
		parts := strings.SplitN(envVar, "=", 2)
		if len(parts) == 2 {
			key, valueStr := parts[0], parts[1]
			if strings.HasPrefix(key, "ITEM_TYPE_MULTIPLIER_") {
				itemType := strings.ToLower(strings.TrimPrefix(key, "ITEM_TYPE_MULTIPLIER_"))
				multiplier, err := strconv.ParseFloat(valueStr, 64)
				if err != nil {
					log.Printf("Warning: Invalid multiplier for %s ('%s'): %v. Skipping.", itemType, valueStr, err)
					continue
				}
				cfg.ItemTypeMultipliers[itemType] = multiplier
			}
		}
	}
	// Ensure default types from model are present if not set by env
	defaultTypes := []string{"document", "electronics", "furniture"}
	for _, itemType := range defaultTypes {
		if _, exists := cfg.ItemTypeMultipliers[itemType]; !exists {
			log.Printf("Warning: Multiplier for default item type '%s' not found in environment. Consider setting ITEM_TYPE_MULTIPLIER_%s.", itemType, strings.ToUpper(itemType))
			// You might want to set a default value here, e.g., 1.0, or make it an error
			// For now, it will be missing, and price calculation will fail for this type.
		}
	}

	log.Printf("Configuration loaded: BaseRate=%.2f, DimFactor=%.2f, Multipliers=%v", cfg.BaseRatePerKg, cfg.DimensionalWeightFactor, cfg.ItemTypeMultipliers)
	return cfg, nil
}

// getEnv retrieves an environment variable or returns a fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Environment variable %s not set, using fallback: %s", key, fallback)
	return fallback
}

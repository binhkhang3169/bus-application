package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

type Utils struct {
}

func NewUtils() *Utils {
	return &Utils{}
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func (u *Utils) init() {
	rand.Seed(time.Now().UnixNano())
}

func (u *Utils) GenerateRandomID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
func (u *Utils) GetCurrentTimeString() time.Time {
	return time.Now()
}
func (u *Utils) ValidateJWTWithSpring(token string) (*int, error) {
	if token == "" {
		return nil, nil
	}

	// Gửi sang service Spring Boot
	req, err := http.NewRequest("GET", "http://spring-auth-service/api/validate", nil)
	if err != nil {
		// Log this error, as it's unexpected for a static request setup
		// For the caller, this still means validation failed.
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Authorization", token)

	client := &http.Client{
		Timeout: 5 * time.Second, // Example timeout
	}
	resp, err := client.Do(req)
	if err != nil {
		// Network error or other error during client.Do
		return nil, fmt.Errorf("error calling auth service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Auth service responded, but token is invalid or other issue.
		// Consider logging resp.StatusCode and potentially body for debugging.
		return nil, fmt.Errorf("auth service returned status %d", resp.StatusCode)
	}

	var result struct {
		CustomerID int `json:"customer_id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode auth response: %w", err)
	}

	return &result.CustomerID, nil // Success
}

func ToNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}

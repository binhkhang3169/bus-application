package qrclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"ticket-service/pkg/utils" // Assuming utils.Logger is here
	"time"
)

// QRServiceResponse matches the expected JSON structure from the QR service
type QRServiceResponse struct {
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	CloudinaryURL string `json:"cloudinary_url"`
	PublicID      string `json:"public_id"`
	Content       string `json:"content"`
}

// GenerateQRCode calls the external QR code generation service.
func GenerateQRCode(ctx context.Context, content string, logger utils.Logger, qrServiceURL string) (string, error) {
	if qrServiceURL == "" {
		return "", fmt.Errorf("QR service URL is not configured")
	}

	payload := map[string]string{"content": content}
	jsonData, err := json.Marshal(payload)
	if err != nil {
		logger.Error("Failed to marshal QR request for content '%s': %v", content, err)
		return "", fmt.Errorf("failed to marshal QR request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", qrServiceURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("Failed to create QR request for content '%s': %v", content, err)
		return "", fmt.Errorf("failed to create QR request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second} // Consider making timeout configurable
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("QR service request failed for content '%s': %v", content, err)
		return "", fmt.Errorf("QR service request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		logger.Error("QR service returned non-OK status (%d) for content '%s'. Body: %s", resp.StatusCode, content, string(bodyBytes))
		return "", fmt.Errorf("QR service returned non-OK status: %d. Body: %s", resp.StatusCode, string(bodyBytes))
	}

	var qrResp QRServiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&qrResp); err != nil {
		logger.Error("Failed to decode QR service response for content '%s': %v", content, err)
		return "", fmt.Errorf("failed to decode QR service response: %w", err)
	}

	if !qrResp.Success || qrResp.CloudinaryURL == "" {
		logger.Error("QR service indicated failure or empty URL for content '%s': %s", content, qrResp.Message)
		return "", fmt.Errorf("QR service indicated failure or empty URL: %s", qrResp.Message)
	}

	logger.Info("QR code generated for content '%s': %s", content, qrResp.CloudinaryURL)
	return qrResp.CloudinaryURL, nil
}

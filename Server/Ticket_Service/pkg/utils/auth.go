package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims repres ents the structure of our JWT claims
type CustomClaims struct {
	jwt.RegisteredClaims
	ID   int    `json:"id"`
	Role string `json:"role"`
}

// TokenInfo holds the extracted token information
type TokenInfo struct {
	UserID    int       `json:"user_id"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

// Auth struct for authentication methods
type Auth struct {
	SecretKey string
}

// NewAuth creates a new Auth instance
func NewAuth(secretKey string) *Auth {
	return &Auth{
		SecretKey: secretKey,
	}
}

// GetCustomerIDFromJWT validates the token and returns the customer ID
func (a *Auth) GetCustomerIDFromJWT(token string) (*int, error) {
	// Validate the token
	tokenInfo, err := a.validateToken(token)
	if err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// Check if the user has customer role
	// if tokenInfo.Role != "ROLE_CUSTOMER" {
	// 	return nil, errors.New("insufficient permissions")
	// }

	// Return the user ID
	return &tokenInfo.UserID, nil
}

// validateToken checks the validity of a JWT token and returns extracted information
func (a *Auth) validateToken(tokenString string) (*TokenInfo, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key for verification
		return []byte(a.SecretKey), nil
	})

	// Check for parsing errors
	if err != nil {
		return nil, fmt.Errorf("token parsing error: %w", err)
	}

	// Extract claims
	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Check token validity
	if claims.ExpiresAt == nil {
		return nil, errors.New("token has no expiration")
	}

	// Check if token is expired
	now := time.Now()
	if claims.ExpiresAt.Time.Before(now) {
		return nil, errors.New("token has expired")
	}

	// Validate token was issued in the past
	if claims.IssuedAt != nil && claims.IssuedAt.Time.After(now) {
		return nil, errors.New("token issued in the future")
	}

	// Create and return token info
	return &TokenInfo{
		UserID:    claims.ID,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

// IsTokenExpired checks if a token is expired
func (a *Auth) IsTokenExpired(tokenString string) bool {
	_, err := a.validateToken(tokenString)
	return err != nil
}

// HasRole checks if the token has a specific role
func (a *Auth) HasRole(tokenString string, expectedRole string) bool {
	tokenInfo, err := a.validateToken(tokenString)
	if err != nil {
		return false
	}
	return tokenInfo.Role == expectedRole
}

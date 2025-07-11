package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims represents the structure of our JWT claims
// It embeds jwt.RegisteredClaims for standard claims like 'exp', 'iat', 'sub'.
// It also includes custom claims 'id' and 'role'.
type CustomClaims struct {
	jwt.RegisteredClaims
	ID   int    `json:"id"`   // Maps to the "id" field in the JWT payload
	Role string `json:"role"` // Maps to the "role" field in the JWT payload
}

// TokenInfo holds the extracted and validated token information
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

// ValidateToken checks the validity of a JWT token and returns extracted information
func (a *Auth) ValidateToken(tokenString string) (*TokenInfo, error) {
	// Parse the token with our custom claims
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (as expected for HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Return the secret key for verification
		return []byte(a.SecretKey), nil
	})

	// Handle parsing errors (which include validation errors like expiration, signature)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, errors.New("token is malformed")
		} else if errors.Is(err, jwt.ErrTokenSignatureInvalid) {
			return nil, errors.New("token signature is invalid")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token has expired")
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			return nil, errors.New("token not yet valid (check 'nbf' or 'iat' claims against current time)")
		}
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok {
		return nil, errors.New("invalid token claims: could not cast to CustomClaims")
	}

	if !token.Valid {
		// This case implies that although parsing was successful, one of the standard validations (exp, nbf, iat) failed.
		// The specific errors are usually caught above. This is a fallback.
		return nil, errors.New("invalid token (token.Valid is false for an undetermined reason post-parsing)")
	}

	// All checks passed
	return &TokenInfo{
		UserID:    claims.ID,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
	}, nil
}

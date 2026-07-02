package auth

import (
	"errors"
	"fmt"
	"os"
)

// TokenVerifier defines the interface for JWT verification.
// TODO v0.3.0: Implement JWKS RS256 verification with key rotation.
type TokenVerifier interface {
	Verify(tokenString string) (*Claims, error)
}

// HS256Verifier implements TokenVerifier using the existing ValidateToken logic.
// Uses MITRAN_JWT_SECRET from environment.
type HS256Verifier struct{}

// NewHS256Verifier creates a verifier, ensuring MITRAN_JWT_SECRET is set.
func NewHS256Verifier() (*HS256Verifier, error) {
	secret := os.Getenv("MITRAN_JWT_SECRET")
	if secret == "" {
		return nil, errors.New("MITRAN_JWT_SECRET not set")
	}
	return &HS256Verifier{}, nil
}

// Verify validates an HS256 JWT and returns decoded claims.
func (v *HS256Verifier) Verify(tokenString string) (*Claims, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("verify: %w", err)
	}
	return claims, nil
}

// TODO v0.3.0: JWKSVerifier struct implementing TokenVerifier
// - Fetches JWKS from MITRAN_JWKS_URL
// - Supports RS256 with key rotation
// - Caches keys with configurable TTL
// type JWKSVerifier struct { ... }
// func NewJWKSVerifier(jwksURL string) (*JWKSVerifier, error) { ... }
// func (v *JWKSVerifier) Verify(tokenString string) (*Claims, error) { ... }

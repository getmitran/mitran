package auth

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	jwksCacheMu    sync.Mutex
	jwksCacheData  = make(map[string]*jwksCacheEntry)
)

type jwksCacheEntry struct {
	jwks      *JWKS
	fetchedAt time.Time
}

const jwksCacheTTL = 1 * time.Hour

type JWK struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

func FetchJWKS(issuerURL string) (*JWKS, error) {
	jwksCacheMu.Lock()
	if entry, ok := jwksCacheData[issuerURL]; ok && time.Since(entry.fetchedAt) < jwksCacheTTL {
		jwksCacheMu.Unlock()
		return entry.jwks, nil
	}
	jwksCacheMu.Unlock()

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(issuerURL + "/.well-known/jwks.json")
	if err != nil {
		return nil, fmt.Errorf("jwks fetch failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("jwks fetch returned status %d", resp.StatusCode)
	}
	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("jwks decode failed: %w", err)
	}

	jwksCacheMu.Lock()
	jwksCacheData[issuerURL] = &jwksCacheEntry{jwks: &jwks, fetchedAt: time.Now()}
	jwksCacheMu.Unlock()

	return &jwks, nil
}

func FindKey(jwks *JWKS, kid string) *JWK {
	for _, k := range jwks.Keys {
		if k.Kid == kid {
			return &k
		}
	}
	return nil
}

// VerifyJWT verifies an RS256-signed JWT using the issuer's JWKS endpoint.
// Returns the decoded payload claims on success.
func VerifyJWT(tokenString, issuerURL string) (map[string]interface{}, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	// Parse header to extract kid
	headerBytes, err := base64URLDecode(parts[0])
	if err != nil {
		return nil, fmt.Errorf("decode header: %w", err)
	}
	var header struct {
		Alg string `json:"alg"`
		Kid string `json:"kid"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, fmt.Errorf("parse header: %w", err)
	}
	if header.Alg != "RS256" {
		return nil, fmt.Errorf("unsupported algorithm: %s", header.Alg)
	}

	// Fetch JWKS and find matching key
	jwks, err := FetchJWKS(issuerURL)
	if err != nil {
		return nil, err
	}
	jwk := FindKey(jwks, header.Kid)
	if jwk == nil {
		return nil, fmt.Errorf("key not found for kid: %s", header.Kid)
	}
	if jwk.Kty != "RSA" {
		return nil, fmt.Errorf("unsupported key type: %s", jwk.Kty)
	}

	// Reconstruct RSA public key from JWK
	pubKey, err := jwkToRSAPublicKey(jwk)
	if err != nil {
		return nil, fmt.Errorf("reconstruct public key: %w", err)
	}

	// Verify RS256 signature
	signingInput := parts[0] + "." + parts[1]
	signature, err := base64URLDecode(parts[2])
	if err != nil {
		return nil, fmt.Errorf("decode signature: %w", err)
	}
	hash := sha256.Sum256([]byte(signingInput))
	if err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hash[:], signature); err != nil {
		return nil, errors.New("signature verification failed")
	}

	// Decode and return payload
	payloadBytes, err := base64URLDecode(parts[1])
	if err != nil {
		return nil, fmt.Errorf("decode payload: %w", err)
	}
	var claims map[string]interface{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, fmt.Errorf("parse payload: %w", err)
	}
	return claims, nil
}

func jwkToRSAPublicKey(jwk *JWK) (*rsa.PublicKey, error) {
	nBytes, err := base64URLDecode(jwk.N)
	if err != nil {
		return nil, fmt.Errorf("decode modulus: %w", err)
	}
	eBytes, err := base64URLDecode(jwk.E)
	if err != nil {
		return nil, fmt.Errorf("decode exponent: %w", err)
	}
	n := new(big.Int).SetBytes(nBytes)
	e := new(big.Int).SetBytes(eBytes)
	if !e.IsInt64() {
		return nil, errors.New("exponent too large")
	}
	return &rsa.PublicKey{N: n, E: int(e.Int64())}, nil
}

func base64URLDecode(s string) ([]byte, error) {
	// Add padding if needed
	switch len(s) % 4 {
	case 2:
		s += "=="
	case 3:
		s += "="
	}
	return base64.URLEncoding.DecodeString(s)
}

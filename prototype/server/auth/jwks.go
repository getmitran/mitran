package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// TODO v0.2.0: Full RSA signature verification using crypto/rsa + big.Int

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
	resp, err := http.Get(issuerURL + "/.well-known/jwks.json")
	if err != nil {
		return nil, fmt.Errorf("jwks fetch failed: %w", err)
	}
	defer resp.Body.Close()
	var jwks JWKS
	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return nil, fmt.Errorf("jwks decode failed: %w", err)
	}
	return &jwks, nil
}

func FindKey(jwks *JWKS, kid string) *JWK {
	for _, k := range jwks.Keys {
		if k.Kid == kid {
			log.Printf("[jwks] best-effort: found key kid=%s kty=%s alg=%s", k.Kid, k.Kty, k.Alg)
			return &k
		}
	}
	return nil
}

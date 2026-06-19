package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"
	"github.com/getmitran/mitran/server/apierr"
)

var jwtSecret = []byte(os.Getenv("MITRAN_JWT_SECRET"))

type Role string

const (
	RoleAdmin    Role = "admin"
	RoleOperator Role = "operator"
	RoleViewer   Role = "viewer"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   Role   `json:"role"`
	Exp    int64  `json:"exp"`
}

type APIKey struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	Role   Role   `json:"role"`
	Active bool   `json:"active"`
}

func GenerateToken(userID string, role Role, ttl time.Duration) (string, error) {
	claims := Claims{UserID: userID, Role: role, Exp: time.Now().Add(ttl).Unix()}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload, _ := json.Marshal(claims)
	payloadB64 := base64.RawURLEncoding.EncodeToString(payload)
	sig := sign(header + "." + payloadB64)
	return header + "." + payloadB64 + "." + sig, nil
}

func ValidateToken(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}
	if sign(parts[0]+"."+parts[1]) != parts[2] {
		return nil, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	if time.Now().Unix() > claims.Exp {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

func sign(data string) string {
	h := hmac.New(sha256.New, jwtSecret)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// AuthMiddleware validates JWT tokens or passes through in dev mode.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/health" || r.URL.Path == "/" {
			next.ServeHTTP(w, r)
			return
		}
		if len(jwtSecret) == 0 {
			next.ServeHTTP(w, r)
			return
		}

		auth := r.Header.Get("Authorization")
		if auth == "" {
			apierr.WriteError(w, apierr.Unauthorized("unauthorized"))
			return
		}

		token := strings.TrimPrefix(auth, "Bearer ")
		claims, err := ValidateToken(token)
		if err != nil {
			apierr.WriteError(w, apierr.Unauthorized(err.Error()))
			return
		}

		r.Header.Set("X-User-ID", claims.UserID)
		r.Header.Set("X-User-Role", string(claims.Role))
		next.ServeHTTP(w, r)
	})
}

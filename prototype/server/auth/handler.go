package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	PassHash string `json:"-"`
	Role     Role   `json:"role"`
}

type userStore struct {
	mu    sync.RWMutex
	users map[string]*User // keyed by username
	keys  map[string]*APIKey
}

var store = &userStore{
	users: make(map[string]*User),
	keys:  make(map[string]*APIKey),
}

func init() {
	// Seed default admin user (password: admin)
	h := sha256.Sum256([]byte("admin"))
	store.users["admin"] = &User{
		ID:       "usr_admin",
		Username: "admin",
		PassHash: hex.EncodeToString(h[:]),
		Role:     RoleAdmin,
	}
}

func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return hex.EncodeToString(h[:])
}

// RegisterRoutes mounts auth endpoints on the given mux.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", handleLogin)
	mux.HandleFunc("POST /api/v1/auth/token", handleGenerateAPIKey)
	mux.HandleFunc("GET /api/v1/auth/me", handleMe)
	mux.HandleFunc("POST /api/v1/auth/register", handleRegister)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	store.mu.RLock()
	user, ok := store.users[req.Username]
	store.mu.RUnlock()

	if !ok || user.PassHash != hashPassword(req.Password) {
		http.Error(w, "invalid credentials", 401)
		return
	}

	token, err := GenerateToken(user.ID, user.Role, 24*time.Hour)
	if err != nil {
		http.Error(w, "token generation failed", 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token":   token,
		"user_id": user.ID,
		"role":    string(user.Role),
	})
}

func handleGenerateAPIKey(w http.ResponseWriter, r *http.Request) {
	role := Role(r.Header.Get("X-User-Role"))
	if role != RoleAdmin {
		http.Error(w, "admin only", 403)
		return
	}

	var req struct {
		Name string `json:"name"`
		Role Role   `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}

	b := make([]byte, 32)
	rand.Read(b)
	key := fmt.Sprintf("mk_%s", hex.EncodeToString(b))

	apiKey := &APIKey{Key: key, Name: req.Name, Role: req.Role, Active: true}
	store.mu.Lock()
	store.keys[key] = apiKey
	store.mu.Unlock()

	json.NewEncoder(w).Encode(apiKey)
}

func handleMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	role := r.Header.Get("X-User-Role")
	if userID == "" {
		http.Error(w, "not authenticated", 401)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{
		"user_id": userID,
		"role":    role,
	})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	callerRole := Role(r.Header.Get("X-User-Role"))
	if callerRole != RoleAdmin {
		http.Error(w, "admin only", 403)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Role     Role   `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid body", 400)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "username and password required", 400)
		return
	}

	store.mu.Lock()
	defer store.mu.Unlock()
	if _, exists := store.users[req.Username]; exists {
		http.Error(w, "user already exists", 409)
		return
	}

	b := make([]byte, 8)
	rand.Read(b)
	user := &User{
		ID:       fmt.Sprintf("usr_%s", hex.EncodeToString(b)),
		Username: req.Username,
		PassHash: hashPassword(req.Password),
		Role:     req.Role,
	}
	store.users[req.Username] = user

	json.NewEncoder(w).Encode(map[string]string{
		"id":       user.ID,
		"username": user.Username,
		"role":     string(user.Role),
	})
}

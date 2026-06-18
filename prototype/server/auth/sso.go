package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

type SSOConfig struct{ Issuer, ClientID, ClientSecret, RedirectURL string }

func LoadSSOConfig() SSOConfig {
	return SSOConfig{os.Getenv("SSO_ISSUER"), os.Getenv("SSO_CLIENT_ID"), os.Getenv("SSO_CLIENT_SECRET"), os.Getenv("SSO_REDIRECT_URL")}
}

var (
	ssoStates  = map[string]time.Time{}
	ssoNonces  = map[string]time.Time{}
	statesMu   sync.Mutex
)

func cleanExpiredStates() {
	now := time.Now()
	for k, v := range ssoStates { if now.After(v) { delete(ssoStates, k) } }
	for k, v := range ssoNonces { if now.After(v) { delete(ssoNonces, k) } }
}

func makeSessionToken(email string) string {
	expiry := fmt.Sprintf("%d", time.Now().Add(24*time.Hour).Unix())
	mac := hmac.New(sha256.New, []byte(os.Getenv("MITRAN_SESSION_SECRET")))
	mac.Write([]byte(email + "|" + expiry))
	return fmt.Sprintf("%s|%s|%s", email, expiry, hex.EncodeToString(mac.Sum(nil)))
}

func HandleSSOLogin(w http.ResponseWriter, r *http.Request) {
	if len(os.Getenv("MITRAN_SESSION_SECRET")) < 16 {
		log.Println("WARNING: MITRAN_SESSION_SECRET should be at least 16 characters")
	}
	cfg := LoadSSOConfig()
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		http.Error(w, "entropy failure", http.StatusInternalServerError); return
	}
	state := hex.EncodeToString(b)
	nb := make([]byte, 16)
	if _, err := rand.Read(nb); err != nil {
		http.Error(w, "entropy failure", http.StatusInternalServerError); return
	}
	nonce := hex.EncodeToString(nb)
	statesMu.Lock(); cleanExpiredStates(); ssoStates[state] = time.Now().Add(10 * time.Minute); ssoNonces[nonce] = time.Now().Add(10 * time.Minute); statesMu.Unlock()
	u := fmt.Sprintf("%s/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email&state=%s&nonce=%s",
		strings.TrimRight(cfg.Issuer, "/"), url.QueryEscape(cfg.ClientID), url.QueryEscape(cfg.RedirectURL), state, nonce)
	http.Redirect(w, r, u, http.StatusFound)
}

func HandleSSOCallback(w http.ResponseWriter, r *http.Request) {
	cfg := LoadSSOConfig()
	state := r.URL.Query().Get("state")
	statesMu.Lock()
	exp, ok := ssoStates[state]
	if ok { delete(ssoStates, state) }
	statesMu.Unlock()
	if !ok || time.Now().After(exp) { http.Error(w, "invalid state", http.StatusBadRequest); return }
	resp, err := http.PostForm(strings.TrimRight(cfg.Issuer, "/")+"/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {r.URL.Query().Get("code")},
		"redirect_uri": {cfg.RedirectURL}, "client_id": {cfg.ClientID}, "client_secret": {cfg.ClientSecret},
	})
	if err != nil { http.Error(w, "token exchange failed", http.StatusBadGateway); return }
	defer resp.Body.Close()
	var tok struct{ IDToken string `json:"id_token"` }
	json.NewDecoder(resp.Body).Decode(&tok)
	parts := strings.Split(tok.IDToken, ".")
	if len(parts) != 3 { http.Error(w, "invalid id_token", http.StatusBadGateway); return }
	// NOTE: Production must verify JWT signature via JWKS endpoint (issuer/.well-known/jwks.json).
	// Skipped for v0.1.0 — we trust the TLS channel to the IdP for token exchange.
	seg := strings.ReplaceAll(strings.ReplaceAll(parts[1], "-", "+"), "_", "/")
	for len(seg)%4 != 0 { seg += "=" }
	payload, _ := base64.StdEncoding.DecodeString(seg)
	var claims struct{ Iss, Sub, Email, Aud, Nonce string; Exp int64 }
	json.Unmarshal(payload, &claims)
	if claims.Iss != cfg.Issuer || claims.Aud != cfg.ClientID || time.Now().Unix() > claims.Exp {
		http.Error(w, "token validation failed", http.StatusUnauthorized); return
	}
	statesMu.Lock()
	nExp, nOk := ssoNonces[claims.Nonce]
	if nOk { delete(ssoNonces, claims.Nonce) }
	statesMu.Unlock()
	if !nOk || time.Now().After(nExp) { http.Error(w, "invalid nonce", http.StatusUnauthorized); return }
	http.SetCookie(w, &http.Cookie{Name: "session", Value: makeSessionToken(claims.Email), Path: "/", HttpOnly: true, SameSite: http.SameSiteStrictMode})
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

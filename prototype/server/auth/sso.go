package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type SSOConfig struct {
	Issuer, ClientID, ClientSecret, RedirectURL string
}

func LoadSSOConfig() SSOConfig {
	return SSOConfig{os.Getenv("SSO_ISSUER"), os.Getenv("SSO_CLIENT_ID"), os.Getenv("SSO_CLIENT_SECRET"), os.Getenv("SSO_REDIRECT_URL")}
}

var ssoStates = map[string]time.Time{}

func HandleSSOLogin(w http.ResponseWriter, r *http.Request) {
	cfg := LoadSSOConfig()
	b := make([]byte, 16)
	rand.Read(b)
	state := hex.EncodeToString(b)
	ssoStates[state] = time.Now().Add(10 * time.Minute)
	u := fmt.Sprintf("%s/authorize?client_id=%s&redirect_uri=%s&response_type=code&scope=openid+email&state=%s",
		strings.TrimRight(cfg.Issuer, "/"), url.QueryEscape(cfg.ClientID), url.QueryEscape(cfg.RedirectURL), state)
	http.Redirect(w, r, u, http.StatusFound)
}

func HandleSSOCallback(w http.ResponseWriter, r *http.Request) {
	cfg := LoadSSOConfig()
	state := r.URL.Query().Get("state")
	exp, ok := ssoStates[state]
	if !ok || time.Now().After(exp) {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}
	delete(ssoStates, state)
	resp, err := http.PostForm(strings.TrimRight(cfg.Issuer, "/")+"/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {r.URL.Query().Get("code")},
		"redirect_uri": {cfg.RedirectURL}, "client_id": {cfg.ClientID}, "client_secret": {cfg.ClientSecret},
	})
	if err != nil {
		http.Error(w, "token exchange failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()
	var tok struct{ IDToken string `json:"id_token"` }
	json.NewDecoder(resp.Body).Decode(&tok)
	parts := strings.Split(tok.IDToken, ".")
	if len(parts) != 3 {
		http.Error(w, "invalid id_token", http.StatusBadGateway)
		return
	}
	seg := strings.ReplaceAll(strings.ReplaceAll(parts[1], "-", "+"), "_", "/")
	for len(seg)%4 != 0 {
		seg += "="
	}
	payload, _ := base64.StdEncoding.DecodeString(seg)
	var claims struct {
		Iss, Sub, Email string
		Aud             string `json:"aud"`
		Exp             int64  `json:"exp"`
	}
	json.Unmarshal(payload, &claims)
	if claims.Iss != cfg.Issuer || claims.Aud != cfg.ClientID || time.Now().Unix() > claims.Exp {
		http.Error(w, "token validation failed", http.StatusUnauthorized)
		return
	}
	session := CreateSession(claims.Email, "sso")
	http.SetCookie(w, &http.Cookie{Name: "session", Value: session, Path: "/", HttpOnly: true})
	http.Redirect(w, r, "/dashboard", http.StatusFound)
}

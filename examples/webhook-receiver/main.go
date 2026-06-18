package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	secret := os.Getenv("WEBHOOK_SECRET")
	if secret == "" {
		fmt.Fprintln(os.Stderr, "WEBHOOK_SECRET env var required")
		os.Exit(1)
	}
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		body, _ := io.ReadAll(r.Body)
		sig := r.Header.Get("X-Signature")
		mac := hmac.New(sha256.New, []byte(secret))
		mac.Write(body)
		expected := hex.EncodeToString(mac.Sum(nil))
		if !hmac.Equal([]byte(sig), []byte(expected)) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}
		fmt.Printf("Event: %s\n", body)
		w.WriteHeader(http.StatusOK)
	})
	fmt.Println("Listening on :9090")
	http.ListenAndServe(":9090", nil)
}

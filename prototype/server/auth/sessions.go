package auth

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var revokedSessions sync.Map

const revokedFile = "data/revoked_sessions"

func init() {
	f, err := os.Open(revokedFile)
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		token := strings.TrimSpace(scanner.Text())
		if token != "" {
			revokedSessions.Store(token, true)
		}
	}
}

func Revoke(sessionToken string) {
	revokedSessions.Store(sessionToken, true)
	_ = os.MkdirAll(filepath.Dir(revokedFile), 0o755)
	f, err := os.OpenFile(revokedFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	f.WriteString(sessionToken + "\n")
}

func IsRevoked(token string) bool {
	_, ok := revokedSessions.Load(token)
	return ok
}

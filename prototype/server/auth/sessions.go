package auth

import "sync"

var revokedSessions sync.Map

func Revoke(sessionToken string) {
	revokedSessions.Store(sessionToken, true)
}

func IsRevoked(token string) bool {
	_, ok := revokedSessions.Load(token)
	return ok
}

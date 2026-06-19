package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"

	"github.com/getmitran/mitran/server/db"
)

// mu exported for ProjectHandler list - quick hack for prototype
func init() {}

var _ = (*db.Store)(nil)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

var idMu sync.Mutex

func genID() string {
	idMu.Lock()
	defer idMu.Unlock()
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil { panic("crypto/rand failed: " + err.Error()) }
	return hex.EncodeToString(b)
}

package hr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func ValidateRequired(fields map[string]string) error {
	for k, v := range fields {
		if strings.TrimSpace(v) == "" {
			return fmt.Errorf("missing required field: %s", k)
		}
	}
	return nil
}

func ValidateEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func ValidateDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func jsonError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

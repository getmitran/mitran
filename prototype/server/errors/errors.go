package errors

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Detail    string `json:"detail,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func (e *APIError) Error() string { return e.Message }

func Write(w http.ResponseWriter, r *http.Request, code int, msg string) {
	ae := APIError{Code: code, Message: msg, RequestID: r.Header.Get("X-Request-ID")}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(ae)
}

func BadRequest(w http.ResponseWriter, r *http.Request, msg string) { Write(w, r, 400, msg) }
func Unauthorized(w http.ResponseWriter, r *http.Request)           { Write(w, r, 401, "unauthorized") }
func Forbidden(w http.ResponseWriter, r *http.Request)              { Write(w, r, 403, "forbidden") }
func NotFound(w http.ResponseWriter, r *http.Request)               { Write(w, r, 404, "not found") }
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	ae := APIError{Code: 500, Message: "internal server error", Detail: err.Error(), RequestID: r.Header.Get("X-Request-ID")}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(ae)
}

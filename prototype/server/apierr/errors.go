package apierr

import (
	"encoding/json"
	"errors"
	"net/http"
)

type APIError struct {
	Code    int    `json:"-"`
	Type    string `json:"error"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

func (e *APIError) Error() string { return e.Message }

func NotFound(msg string) *APIError {
	return &APIError{Code: 404, Type: "not_found", Message: msg}
}

func BadRequest(msg string) *APIError {
	return &APIError{Code: 400, Type: "bad_request", Message: msg}
}

func Unauthorized(msg string) *APIError {
	return &APIError{Code: 401, Type: "unauthorized", Message: msg}
}

func Forbidden(msg string) *APIError {
	return &APIError{Code: 403, Type: "forbidden", Message: msg}
}

func Internal(msg string) *APIError {
	return &APIError{Code: 500, Type: "internal_error", Message: msg}
}

func MethodNotAllowed(msg string) *APIError {
	return &APIError{Code: 405, Type: "method_not_allowed", Message: msg}
}

func Conflict(msg string) *APIError {
	return &APIError{Code: 409, Type: "conflict", Message: msg}
}

func WriteError(w http.ResponseWriter, err error) {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(apiErr.Code)
		json.NewEncoder(w).Encode(apiErr)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)
	json.NewEncoder(w).Encode(&APIError{Code: 500, Type: "internal_error", Message: err.Error()})
}

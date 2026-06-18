package apierr

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotFound(t *testing.T) {
	err := NotFound("item not found")
	if err.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", err.Code)
	}
	if err.Message != "item not found" {
		t.Fatalf("expected 'item not found', got %q", err.Message)
	}
	if err.Type != "not_found" {
		t.Fatalf("expected type 'not_found', got %q", err.Type)
	}
}

func TestBadRequest(t *testing.T) {
	err := BadRequest("invalid input")
	if err.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", err.Code)
	}
	if err.Message != "invalid input" {
		t.Fatalf("expected 'invalid input', got %q", err.Message)
	}
}

func TestWriteErrorWithAPIError(t *testing.T) {
	w := httptest.NewRecorder()
	apiErr := NotFound("resource missing")
	WriteError(w, apiErr)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["error"] != "not_found" {
		t.Fatalf("expected error 'not_found', got %v", body["error"])
	}
	if body["message"] != "resource missing" {
		t.Fatalf("expected message 'resource missing', got %v", body["message"])
	}
}

func TestWriteErrorWithPlainError(t *testing.T) {
	w := httptest.NewRecorder()
	WriteError(w, errors.New("something broke"))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["error"] != "internal_error" {
		t.Fatalf("expected error 'internal_error', got %v", body["error"])
	}
	if body["message"] != "something broke" {
		t.Fatalf("expected message 'something broke', got %v", body["message"])
	}
}

func TestErrorMethod(t *testing.T) {
	err := NotFound("gone")
	if err.Error() != "gone" {
		t.Fatalf("Error() should return message, got %q", err.Error())
	}
}

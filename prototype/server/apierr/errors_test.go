package apierr

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewNotFound(t *testing.T) {
	err := NewNotFound("item not found")
	if err.Status != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", err.Status)
	}
	if err.Message != "item not found" {
		t.Fatalf("expected 'item not found', got %q", err.Message)
	}
}

func TestNewBadRequest(t *testing.T) {
	err := NewBadRequest("invalid input")
	if err.Status != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", err.Status)
	}
	if err.Message != "invalid input" {
		t.Fatalf("expected 'invalid input', got %q", err.Message)
	}
}

func TestWriteErrorWithAPIError(t *testing.T) {
	w := httptest.NewRecorder()
	apiErr := NewNotFound("resource missing")
	WriteError(w, apiErr)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
	var body map[string]interface{}
	if err := json.NewDecoder(w.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode body: %v", err)
	}
	if body["error"] != "resource missing" {
		t.Fatalf("expected error 'resource missing', got %v", body["error"])
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
	if body["error"] != "something broke" {
		t.Fatalf("expected error 'something broke', got %v", body["error"])
	}
}

func TestErrorMethod(t *testing.T) {
	err := NewNotFound("gone")
	msg := err.Error()
	if !strings.Contains(msg, "not_found") && !strings.Contains(msg, "404") {
		t.Fatalf("Error() should include type info, got %q", msg)
	}
	if !strings.Contains(msg, "gone") {
		t.Fatalf("Error() should include message, got %q", msg)
	}
}

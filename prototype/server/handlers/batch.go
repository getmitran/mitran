package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

type BatchRequest struct {
	Method string          `json:"method"`
	Path   string          `json:"path"`
	Body   json.RawMessage `json:"body,omitempty"`
}

type BatchResponse struct {
	Status int             `json:"status"`
	Body   json.RawMessage `json:"body"`
}

func BatchHandler(mux http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var requests []BatchRequest
		if err := json.NewDecoder(r.Body).Decode(&requests); err != nil {
			http.Error(w, "invalid batch request", 400)
			return
		}
		if len(requests) > 20 {
			http.Error(w, "max 20 requests per batch", 400)
			return
		}
		responses := make([]BatchResponse, len(requests))
		for i, req := range requests {
			var body *bytes.Reader
			if req.Body != nil {
				body = bytes.NewReader(req.Body)
			} else {
				body = bytes.NewReader(nil)
			}
			subReq := httptest.NewRequest(req.Method, req.Path, body)
			subReq.Header.Set("Content-Type", "application/json")
			for k, v := range r.Header {
				subReq.Header[k] = v
			}
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, subReq)
			responses[i] = BatchResponse{Status: rec.Code, Body: rec.Body.Bytes()}
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(responses)
	}
}

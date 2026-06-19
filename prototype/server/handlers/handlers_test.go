package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/getmitran/mitran/server/db"
)

func testStore(t *testing.T) *db.Store {
	t.Helper()
	s, err := db.NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestHealthz(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(Healthz))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}

func TestAgentsList(t *testing.T) {
	store := testStore(t)
	handler := &AgentHandler{Store: store}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/api/v1/agents"
		handler.ServeHTTP(w, r)
	}))
	defer srv.Close()

	resp, err := http.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	var agents []interface{}
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		t.Fatal(err)
	}
	if agents == nil {
		t.Fatal("expected JSON array, got nil")
	}
}

func TestCreateTask(t *testing.T) {
	store := testStore(t)
	handler := &TaskHandler{Store: store}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/api/v1/tasks"
		r.Method = "POST"
		handler.ServeHTTP(w, r)
	}))
	defer srv.Close()

	body := `{"project_id":"p1","title":"Test task","description":"desc","agent_type":"dev","priority":1}`
	resp, err := http.Post(srv.URL, "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	var task map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		t.Fatal(err)
	}
	if task["title"] != "Test task" {
		t.Fatalf("expected title 'Test task', got %v", task["title"])
	}
	if task["status"] != "queued" {
		t.Fatalf("expected status 'queued', got %v", task["status"])
	}
}

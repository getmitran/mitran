package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/getmitran/mitran/server/db"
	"github.com/getmitran/mitran/server/handlers"
	"github.com/getmitran/mitran/server/middleware"
	"github.com/getmitran/mitran/server/scheduler"
)

func main() {
	// Data directory
	dataDir := filepath.Join(".", ".mitran")
	if d := os.Getenv("MITRAN_DATA_DIR"); d != "" {
		dataDir = d
	}

	store, err := db.NewStore(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	// Handlers
	mux := http.NewServeMux()

	mux.Handle("/api/v1/projects", &handlers.ProjectHandler{Store: store})
	mux.Handle("/api/v1/projects/", &handlers.ProjectHandler{Store: store})
	mux.Handle("/api/v1/tasks", &handlers.TaskHandler{Store: store})
	mux.Handle("/api/v1/tasks/", &handlers.TaskHandler{Store: store})
	mux.Handle("/api/v1/checkpoints", &handlers.CheckpointHandler{Store: store})
	mux.Handle("/api/v1/checkpoints/", &handlers.CheckpointHandler{Store: store})
	mux.Handle("/api/v1/agents", &handlers.AgentHandler{Store: store})
	mux.Handle("/api/v1/agents/", &handlers.AgentHandler{Store: store})
	mux.Handle("/api/v1/init", &handlers.InitHandler{Store: store})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Start scheduler
	sched := &scheduler.Scheduler{Store: store}
	sched.Start()

	port := "7777"
	if p := os.Getenv("MITRAN_PORT"); p != "" {
		port = p
	}

	handler := middleware.CORS(mux)
	fmt.Printf("Mitran Core Engine running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

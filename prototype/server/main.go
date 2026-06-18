package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getmitran/mitran/server/db"
	"github.com/getmitran/mitran/server/handlers"
	"github.com/getmitran/mitran/server/middleware"
	"github.com/getmitran/mitran/server/scheduler"
	"github.com/getmitran/mitran/server/workspace"
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

	// Project workspaces
	wsDir := filepath.Join(dataDir, "projects")
	wsMgr, err := workspace.New(wsDir)
	if err != nil {
		log.Fatalf("Failed to initialize workspace manager: %v", err)
	}

	// Handlers
	mux := http.NewServeMux()

	mux.Handle("/api/v1/tasks", &handlers.TaskHandler{Store: store})
	mux.Handle("/api/v1/tasks/", &handlers.TaskHandler{Store: store})
	mux.Handle("/api/v1/checkpoints", &handlers.CheckpointHandler{Store: store})
	mux.Handle("/api/v1/checkpoints/", &handlers.CheckpointHandler{Store: store})
	mux.Handle("/api/v1/agents", &handlers.AgentHandler{Store: store})
	mux.Handle("/api/v1/agents/", &handlers.AgentHandler{Store: store})
	mux.Handle("/api/v1/init", &handlers.InitHandler{Store: store, Workspace: wsMgr})

	// Environment & Pipeline handlers (project sub-resources)
	// These need to be registered BEFORE the catch-all /api/v1/projects/ 
	// Since Go 1.22 ServeMux doesn't help here with catch-all prefix, 
	// we route via a wrapper that checks for sub-resource paths
	envHandler := &handlers.EnvironmentHandler{Store: store}
	pipeHandler := &handlers.PipelineHandler{Store: store}
	projectHandler := &handlers.ProjectHandler{Store: store}

	mux.Handle("/api/v1/projects", projectHandler)
	mux.HandleFunc("/api/v1/projects/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if strings.Contains(path, "/environments") || strings.Contains(path, "/deployments") {
			envHandler.ServeHTTP(w, r)
		} else if strings.Contains(path, "/pipelines") {
			pipeHandler.ServeHTTP(w, r)
		} else {
			projectHandler.ServeHTTP(w, r)
		}
	})

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Start scheduler
	sched := &scheduler.Scheduler{Store: store, WorkspaceDir: wsDir}
	sched.Start()

	port := "7777"
	if p := os.Getenv("MITRAN_PORT"); p != "" {
		port = p
	}

	handler := middleware.CORS(mux)
	fmt.Printf("Mitran Core Engine running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

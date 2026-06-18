package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/getmitran/mitran/server/db"
	"github.com/getmitran/mitran/server/logging"
	"github.com/getmitran/mitran/server/handlers"
	"github.com/getmitran/mitran/server/memory"
	"github.com/getmitran/mitran/server/middleware"
	"github.com/getmitran/mitran/server/scheduler"
	"github.com/getmitran/mitran/server/settings"
	"github.com/getmitran/mitran/server/shutdown"
	"github.com/getmitran/mitran/server/queue"
	"github.com/getmitran/mitran/server/websocket"
	"github.com/getmitran/mitran/server/workspace"
)

func main() {
	logger := logging.New()

	// Data directory
	dataDir := filepath.Join(".", ".mitran")
	if d := os.Getenv("MITRAN_DATA_DIR"); d != "" {
		dataDir = d
	}

	store, err := db.NewStore(dataDir)
	if err != nil {
		log.Fatalf("Failed to initialize store: %v", err)
	}

	// Load settings (user-configurable paths)
	cfg := settings.Load(dataDir)

	// Project workspaces (uses settings.ProjectsDir)
	wsMgr, err := workspace.New(cfg.ProjectsDir)
	if err != nil {
		log.Fatalf("Failed to initialize workspace manager: %v", err)
	}

	// Handlers
	mux := http.NewServeMux()

	// WebSocket hub for real-time events
	hub := websocket.NewHub()
	mux.HandleFunc("/api/v1/ws", websocket.HandleWS(hub))

	// Memory system
	memStore, err := memory.NewStore(filepath.Join(dataDir, "memory"))
	if err != nil {
		log.Fatalf("Failed to initialize memory: %v", err)
	}
	memHandler := &memory.Handler{Store: memStore}
	mux.Handle("/api/v1/memory/facts", memHandler)
	mux.Handle("/api/v1/memory/search", memHandler)
	mux.Handle("/api/v1/memory/episodes", memHandler)
	mux.Handle("/api/v1/memory/corrections", memHandler)
	mux.Handle("/api/v1/memory/corrections/", memHandler)

	taskQueue := queue.New()
	mux.Handle("/api/v1/tasks", &handlers.TaskHandler{Store: store, Queue: taskQueue})
	mux.Handle("/api/v1/tasks/", &handlers.TaskHandler{Store: store, Queue: taskQueue})
	mux.Handle("/api/v1/checkpoints", &handlers.CheckpointHandler{Store: store})
	mux.Handle("/api/v1/checkpoints/", &handlers.CheckpointHandler{Store: store})
	mux.Handle("/api/v1/agents", &handlers.AgentHandler{Store: store})
	mux.Handle("/api/v1/agents/", &handlers.AgentHandler{Store: store})
	mux.Handle("/api/v1/init", &handlers.InitHandler{Store: store, Workspace: wsMgr})

	// Settings
	mux.Handle("/api/v1/settings", &handlers.SettingsHandler{Config: &cfg})
	mux.Handle("/api/v1/settings/", &handlers.SettingsHandler{Config: &cfg})

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

	// K8s probes
	mux.HandleFunc("/healthz", handlers.Healthz)
	mux.HandleFunc("/readyz", handlers.Readyz)

	// Root — show available routes
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name":"mitran","version":"0.1.0","endpoints":["/health","/api/v1/init","/api/v1/tasks","/api/v1/agents","/api/v1/checkpoints","/api/v1/projects","/api/v1/settings"],"dashboard":"http://localhost:5173"}`))
	})

	// Start scheduler
	sched := &scheduler.Scheduler{Store: store, WorkspaceDir: cfg.ProjectsDir}
	sched.Start()

	port := "7780"
	if p := os.Getenv("MITRAN_PORT"); p != "" {
		port = p
	}

	handler := middleware.CORS(middleware.RequestLogger(mux))
	logger.Info("server starting", "port", port, "data_dir", dataDir)
	fmt.Printf("\n  ╔══════════════════════════════════════════╗\n")
	fmt.Printf("  ║   Mitran Core Engine v0.1.0              ║\n")
	fmt.Printf("  ╚══════════════════════════════════════════╝\n\n")
	fmt.Printf("  Engine:    http://localhost:%s\n", port)
	fmt.Printf("  Dashboard: http://localhost:5173 (start separately)\n\n")
	fmt.Printf("  API Routes:\n")
	fmt.Printf("    POST /api/v1/init          Create project + task plan\n")
	fmt.Printf("    GET  /api/v1/tasks         List tasks\n")
	fmt.Printf("    GET  /api/v1/agents        List agents\n")
	fmt.Printf("    GET  /api/v1/checkpoints   Pending checkpoints\n")
	fmt.Printf("    GET  /api/v1/projects      List projects\n")
	fmt.Printf("    GET  /api/v1/settings      Configuration\n")
	fmt.Printf("    GET  /api/v1/memory/facts   Memory facts\n")
	fmt.Printf("    WS   /api/v1/ws            WebSocket (real-time)\n")
	fmt.Printf("    GET  /health               Health check\n\n")
	fmt.Printf("  Data: %s\n\n", dataDir)
	if err := shutdown.ListenAndServeGraceful(":"+port, handler, 15*time.Second); err != nil {
		log.Fatal(err)
	}
}

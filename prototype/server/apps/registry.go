package apps

import (
	"encoding/json"
	"net/http"
	"strings"
)

func RegisterRoutes(mux *http.ServeMux, installer *Installer, appsDir string) {
	mux.HandleFunc("GET /api/v1/apps", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(installer.List())
	})

	mux.HandleFunc("GET /api/v1/apps/{name}", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		app, ok := installer.Get(name)
		if !ok {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(app)
	})

	mux.HandleFunc("POST /api/v1/apps/install", func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Path string `json:"path"`
		}
		json.NewDecoder(r.Body).Decode(&body)
		if body.Path == "" {
			http.Error(w, "path required", 400)
			return
		}
		app, err := installer.Install(body.Path)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(app)
	})

	mux.HandleFunc("POST /api/v1/apps/{name}/enable", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		workspace := r.URL.Query().Get("workspace")
		if workspace == "" {
			workspace = "default"
		}
		if err := installer.Enable(name, workspace); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(200)
	})

	mux.HandleFunc("POST /api/v1/apps/{name}/disable", func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		workspace := r.URL.Query().Get("workspace")
		if workspace == "" {
			workspace = "default"
		}
		if err := installer.Disable(name, workspace); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(200)
	})

	mux.HandleFunc("DELETE /api/v1/apps/{name}", func(w http.ResponseWriter, r *http.Request) {
		if err := installer.Uninstall(r.PathValue("name")); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		w.WriteHeader(204)
	})

	mux.HandleFunc("POST /api/v1/apps/init", HandleScaffold(appsDir))
}

// SearchMarketplace is a placeholder for future remote app registry
func SearchMarketplace(query string) []map[string]string {
	_ = strings.ToLower(query)
	// TODO: connect to remote marketplace API
	return []map[string]string{}
}

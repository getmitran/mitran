package artifacts

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/getmitran/mitran/server/apierr"
)

// RegisterRoutes registers artifact REST handlers on the given mux.
func RegisterRoutes(mux *http.ServeMux, store *ArtifactStore) {
	mux.HandleFunc("/api/v1/projects/", func(w http.ResponseWriter, r *http.Request) {
		// Parse: /api/v1/projects/{id}/artifacts[/{slug}[/versions]]
		path := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")
		parts := strings.SplitN(path, "/", 4)

		if len(parts) < 2 || parts[1] != "artifacts" {
			http.NotFound(w, r)
			return
		}

		projectID := parts[0]

		switch {
		case len(parts) == 2:
			handleArtifacts(w, r, store, projectID)
		case len(parts) == 3:
			handleArtifact(w, r, store, projectID, parts[2])
		case len(parts) == 4 && parts[3] == "versions":
			handleVersions(w, r, store, projectID, parts[2])
		default:
			http.NotFound(w, r)
		}
	})
}

func handleArtifacts(w http.ResponseWriter, r *http.Request, store *ArtifactStore, projectID string) {
	switch r.Method {
	case http.MethodGet:
		kind := r.URL.Query().Get("kind")
		tag := r.URL.Query().Get("tag")
		query := r.URL.Query().Get("q")

		artifacts, err := store.List(projectID, kind, tag, query)
		if err != nil {
			apierr.WriteError(w, apierr.Internal(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(artifacts)

	case http.MethodPost:
		var req struct {
			Name        string   `json:"name"`
			Kind        string   `json:"kind"`
			Content     string   `json:"content"`
			Tags        []string `json:"tags"`
			Description string   `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apierr.WriteError(w, apierr.BadRequest("invalid JSON"))
			return
		}
		if req.Name == "" || req.Content == "" {
			apierr.WriteError(w, apierr.BadRequest("name and content required"))
			return
		}
		if req.Kind == "" {
			req.Kind = "widget"
		}

		tagsStr := strings.Join(req.Tags, ",")
		slug, err := store.Save(projectID, req.Name, req.Kind, req.Content, tagsStr, req.Description)
		if err != nil {
			apierr.WriteError(w, apierr.Conflict(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(map[string]string{"slug": slug})

	default:
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
	}
}

func handleArtifact(w http.ResponseWriter, r *http.Request, store *ArtifactStore, projectID, slug string) {
	switch r.Method {
	case http.MethodGet:
		version := 0
		if v := r.URL.Query().Get("version"); v != "" {
			if parsed, err := strconv.Atoi(v); err == nil {
				version = parsed
			}
		}

		artifact, err := store.Get(projectID, slug, version)
		if err != nil {
			apierr.WriteError(w, apierr.NotFound("artifact not found"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(artifact)

	case http.MethodPut:
		var req struct {
			Content     string   `json:"content"`
			Name        string   `json:"name"`
			Tags        []string `json:"tags"`
			Description string   `json:"description"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			apierr.WriteError(w, apierr.BadRequest("invalid JSON"))
			return
		}
		if req.Content == "" {
			apierr.WriteError(w, apierr.BadRequest("content required"))
			return
		}

		tagsStr := strings.Join(req.Tags, ",")
		if err := store.Update(projectID, slug, req.Content, req.Name, tagsStr, req.Description); err != nil {
			apierr.WriteError(w, apierr.NotFound(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "updated"})

	case http.MethodDelete:
		if err := store.Delete(projectID, slug); err != nil {
			apierr.WriteError(w, apierr.NotFound(err.Error()))
			return
		}
		w.WriteHeader(204)

	default:
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
	}
}

func handleVersions(w http.ResponseWriter, r *http.Request, store *ArtifactStore, projectID, slug string) {
	if r.Method != http.MethodGet {
		apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
		return
	}

	versions, err := store.Versions(projectID, slug)
	if err != nil {
		apierr.WriteError(w, apierr.NotFound("artifact not found"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"slug": slug, "versions": versions})
}

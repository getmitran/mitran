package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/getmitran/mitran/server/settings"
)

type SettingsHandler struct {
	Config *settings.Config
}

func (h *SettingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		h.get(w, r)
	case "PUT", "PATCH":
		h.update(w, r)
	default:
		http.NotFound(w, r)
	}
}

// GET /api/v1/settings — return full config
// GET /api/v1/settings/projects/:id/path — return project workspace path
func (h *SettingsHandler) get(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/settings")

	if path == "" || path == "/" {
		writeJSON(w, http.StatusOK, h.Config)
		return
	}

	// GET /api/v1/settings/projects/<id>/path
	if strings.HasPrefix(path, "/projects/") {
		parts := strings.Split(strings.TrimPrefix(path, "/projects/"), "/")
		if len(parts) >= 1 {
			projectID := parts[0]
			writeJSON(w, http.StatusOK, map[string]string{
				"project_id":     projectID,
				"workspace_path": h.Config.ProjectWorkspace(projectID),
			})
			return
		}
	}

	http.NotFound(w, r)
}

// PUT /api/v1/settings — update global settings
// PUT /api/v1/settings/projects/:id/path — set custom project path
func (h *SettingsHandler) update(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/settings")

	// Set custom project workspace path
	if strings.HasPrefix(path, "/projects/") {
		parts := strings.Split(strings.TrimPrefix(path, "/projects/"), "/")
		if len(parts) >= 1 {
			projectID := parts[0]
			var body struct {
				Path string `json:"path"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				writeErr(w, http.StatusBadRequest, err.Error())
				return
			}
			h.Config.SetProjectPath(projectID, body.Path)
			if err := h.Config.Save(); err != nil {
				writeErr(w, http.StatusInternalServerError, err.Error())
				return
			}
			writeJSON(w, http.StatusOK, map[string]string{
				"project_id":     projectID,
				"workspace_path": body.Path,
				"status":         "updated",
			})
			return
		}
	}

	// Update global settings
	var updates settings.Config
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		writeErr(w, http.StatusBadRequest, err.Error())
		return
	}

	if updates.RootDir != "" {
		h.Config.RootDir = updates.RootDir
	}
	if updates.ProjectsDir != "" {
		h.Config.ProjectsDir = updates.ProjectsDir
	}
	if updates.Modules.Source != "" {
		h.Config.Modules.Source = updates.Modules.Source
	}
	if updates.Modules.Docs != "" {
		h.Config.Modules.Docs = updates.Modules.Docs
	}
	if updates.Modules.Wiki != "" {
		h.Config.Modules.Wiki = updates.Modules.Wiki
	}
	if updates.Modules.SOPs != "" {
		h.Config.Modules.SOPs = updates.Modules.SOPs
	}
	if updates.Modules.Monitoring != "" {
		h.Config.Modules.Monitoring = updates.Modules.Monitoring
	}
	if updates.Modules.CICD != "" {
		h.Config.Modules.CICD = updates.Modules.CICD
	}
	if updates.Modules.Tickets != "" {
		h.Config.Modules.Tickets = updates.Modules.Tickets
	}
	if updates.Modules.HR != "" {
		h.Config.Modules.HR = updates.Modules.HR
	}

	if err := h.Config.Save(); err != nil {
		writeErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, h.Config)
}

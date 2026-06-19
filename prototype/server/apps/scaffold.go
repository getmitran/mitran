package apps

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/getmitran/mitran/server/apierr"
)

func HandleScaffold(appsDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			apierr.WriteError(w, apierr.MethodNotAllowed("method not allowed"))
			return
		}

		name := r.URL.Query().Get("name")
		if name == "" {
			var body struct {
				Name string `json:"name"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				name = body.Name
			}
		}
		name = strings.TrimSpace(name)
		if name == "" {
			apierr.WriteError(w, apierr.BadRequest("name is required"))
			return
		}

		dir := filepath.Join(appsDir, name)
		if _, err := os.Stat(dir); err == nil {
			apierr.WriteError(w, apierr.Conflict("app already exists"))
			return
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			apierr.WriteError(w, apierr.Internal(err.Error()))
			return
		}

		files := map[string]string{
			"app.json":  appJSON(name),
			"entry.py":  entryPy(name),
			"SKILL.md":  skillMD(name),
			"cron.json": cronJSON(name),
			"README.md": readmeMD(name),
		}

		for fname, content := range files {
			if err := os.WriteFile(filepath.Join(dir, fname), []byte(content), 0644); err != nil {
				apierr.WriteError(w, apierr.Internal("writing file: "+err.Error()))
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"name":  name,
			"path":  dir,
			"files": []string{"app.json", "entry.py", "SKILL.md", "cron.json", "README.md"},
		})
	}
}

func appJSON(name string) string {
	m := Manifest{
		Name:       name,
		Version:    "0.1.0",
		EntryPoint: "entry.py",
		Crons: []CronSpec{{
			Name:     name + "-heartbeat",
			Schedule: "0 */6 * * *",
			Handler:  "entry:heartbeat",
		}},
		Skills: []SkillSpec{{
			Name: name,
			File: "SKILL.md",
		}},
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b) + "\n"
}

func entryPy(name string) string {
	return fmt.Sprintf(`"""Entry point for %s app."""


def run(context):
    """Main handler called by Mitran runtime."""
    return {"status": "ok", "app": "%s"}


def heartbeat(context):
    """Cron handler for periodic tasks."""
    pass
`, name, name)
}

func skillMD(name string) string {
	return fmt.Sprintf(`# %s

## Description

TODO: Describe what this app does.

## Usage

TODO: Add usage instructions.
`, name)
}

func cronJSON(name string) string {
	c := []map[string]string{{
		"name":     name + "-heartbeat",
		"schedule": "0 */6 * * *",
		"handler":  "entry:heartbeat",
	}}
	b, _ := json.MarshalIndent(c, "", "  ")
	return string(b) + "\n"
}

func readmeMD(name string) string {
	return fmt.Sprintf(`# %s

A Mitran app.

## Setup

1. Install: `+"`mitran app install .`"+`
2. Enable: `+"`mitran app enable %s`"+`

## Development

- `+"`entry.py`"+` — Main entry point
- `+"`SKILL.md`"+` — Agent skill definition
- `+"`cron.json`"+` — Scheduled tasks
- `+"`app.json`"+` — App manifest
`, name, name)
}

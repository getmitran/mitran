package settings

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Config holds all user-configurable paths and preferences.
// Loaded from .mitran/settings.json or overridden via env vars.
type Config struct {
	// Root directory for all Mitran data. Default: ./.mitran
	RootDir string `json:"root_dir"`

	// Base directory for project workspaces. Default: <root_dir>/projects
	ProjectsDir string `json:"projects_dir"`

	// Per-project workspace overrides. Key: project ID, Value: custom path
	ProjectPaths map[string]string `json:"project_paths"`

	// Module-specific output directories (override defaults)
	Modules ModulePaths `json:"modules"`

	// Engine settings
	EnginePort string `json:"engine_port"`
	WorkerPort string `json:"worker_port"`
}

// ModulePaths allows overriding where each module stores its files.
// If empty, defaults to <project_workspace>/<module_name>/
type ModulePaths struct {
	Source     string `json:"source,omitempty"`     // Default: <project>/src
	Docs       string `json:"docs,omitempty"`       // Default: <project>/docs
	Wiki       string `json:"wiki,omitempty"`       // Default: <project>/wiki
	SOPs       string `json:"sops,omitempty"`       // Default: <project>/sops
	Monitoring string `json:"monitoring,omitempty"` // Default: <project>/monitoring
	CICD       string `json:"cicd,omitempty"`       // Default: <project>/cicd
	Tickets    string `json:"tickets,omitempty"`    // Default: <project>/tickets
	HR         string `json:"hr,omitempty"`         // Default: <project>/hr
}

// Default returns a Config with all defaults applied.
func Default() Config {
	return Config{
		RootDir:      ".mitran",
		ProjectsDir:  "",
		ProjectPaths: make(map[string]string),
		Modules:      ModulePaths{},
		EnginePort:   "7777",
		WorkerPort:   "8888",
	}
}

// Load reads settings from disk, applying env var overrides.
func Load(rootDir string) Config {
	cfg := Default()
	cfg.RootDir = rootDir

	// Try loading from file
	path := filepath.Join(rootDir, "settings.json")
	data, err := os.ReadFile(path)
	if err == nil {
		json.Unmarshal(data, &cfg)
	}

	// Env overrides
	if v := os.Getenv("MITRAN_ROOT_DIR"); v != "" {
		cfg.RootDir = v
	}
	if v := os.Getenv("MITRAN_PROJECTS_DIR"); v != "" {
		cfg.ProjectsDir = v
	}
	if v := os.Getenv("MITRAN_PORT"); v != "" {
		cfg.EnginePort = v
	}

	// Resolve ProjectsDir default
	if cfg.ProjectsDir == "" {
		cfg.ProjectsDir = filepath.Join(cfg.RootDir, "projects")
	}

	return cfg
}

// Save persists the config to .mitran/settings.json.
func (c *Config) Save() error {
	os.MkdirAll(c.RootDir, 0755)
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(c.RootDir, "settings.json"), data, 0644)
}

// ProjectWorkspace returns the workspace path for a project.
// Priority: per-project override > global projects_dir > default
func (c *Config) ProjectWorkspace(projectID string) string {
	if custom, ok := c.ProjectPaths[projectID]; ok {
		return custom
	}
	return filepath.Join(c.ProjectsDir, projectID)
}

// ModulePath returns where a module should write files for a project.
// Priority: module override > <project_workspace>/<module>/
func (c *Config) ModulePath(projectID, module string) string {
	ws := c.ProjectWorkspace(projectID)

	switch module {
	case "source", "src", "dev":
		if c.Modules.Source != "" {
			return filepath.Join(c.Modules.Source, projectID)
		}
		return filepath.Join(ws, "src")
	case "docs":
		if c.Modules.Docs != "" {
			return filepath.Join(c.Modules.Docs, projectID)
		}
		return filepath.Join(ws, "docs")
	case "wiki":
		if c.Modules.Wiki != "" {
			return filepath.Join(c.Modules.Wiki, projectID)
		}
		return filepath.Join(ws, "wiki")
	case "sops":
		if c.Modules.SOPs != "" {
			return filepath.Join(c.Modules.SOPs, projectID)
		}
		return filepath.Join(ws, "sops")
	case "monitoring", "ops":
		if c.Modules.Monitoring != "" {
			return filepath.Join(c.Modules.Monitoring, projectID)
		}
		return filepath.Join(ws, "monitoring")
	case "cicd":
		if c.Modules.CICD != "" {
			return filepath.Join(c.Modules.CICD, projectID)
		}
		return filepath.Join(ws, "cicd")
	case "tickets":
		if c.Modules.Tickets != "" {
			return filepath.Join(c.Modules.Tickets, projectID)
		}
		return filepath.Join(ws, "tickets")
	case "hr":
		if c.Modules.HR != "" {
			return filepath.Join(c.Modules.HR, projectID)
		}
		return filepath.Join(ws, "hr")
	default:
		return filepath.Join(ws, module)
	}
}

// SetProjectPath sets a custom workspace path for a specific project.
func (c *Config) SetProjectPath(projectID, path string) {
	if c.ProjectPaths == nil {
		c.ProjectPaths = make(map[string]string)
	}
	c.ProjectPaths[projectID] = path
}

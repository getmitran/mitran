package apps

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type AppState struct {
	Name      string            `json:"name"`
	Version   string            `json:"version"`
	Enabled   map[string]bool   `json:"enabled"` // workspace -> enabled
	Manifest  *Manifest         `json:"manifest"`
}

type Installer struct {
	mu       sync.RWMutex
	storeDir string
	apps     map[string]*AppState
}

func NewInstaller(storeDir string) (*Installer, error) {
	if err := os.MkdirAll(storeDir, 0755); err != nil {
		return nil, err
	}
	inst := &Installer{storeDir: storeDir, apps: make(map[string]*AppState)}
	inst.loadRegistry()
	return inst, nil
}

func (i *Installer) Install(srcDir string) (*AppState, error) {
	manifest, err := ParseManifest(filepath.Join(srcDir, "app.json"))
	if err != nil {
		return nil, fmt.Errorf("install: %w", err)
	}

	i.mu.Lock()
	defer i.mu.Unlock()

	destDir := filepath.Join(i.storeDir, manifest.Name)
	if err := copyDir(srcDir, destDir); err != nil {
		return nil, fmt.Errorf("install copy: %w", err)
	}

	state := &AppState{
		Name:     manifest.Name,
		Version:  manifest.Version,
		Enabled:  make(map[string]bool),
		Manifest: manifest,
	}
	i.apps[manifest.Name] = state
	return state, i.saveRegistry()
}

func (i *Installer) Uninstall(name string) error {
	i.mu.Lock()
	defer i.mu.Unlock()

	if _, ok := i.apps[name]; !ok {
		return fmt.Errorf("app %q not installed", name)
	}
	if err := os.RemoveAll(filepath.Join(i.storeDir, name)); err != nil {
		return fmt.Errorf("uninstall cleanup: %w", err)
	}
	delete(i.apps, name)
	return i.saveRegistry()
}

func (i *Installer) Enable(name, workspace string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	app, ok := i.apps[name]
	if !ok {
		return fmt.Errorf("app %q not installed", name)
	}
	app.Enabled[workspace] = true
	return i.saveRegistry()
}

func (i *Installer) Disable(name, workspace string) error {
	i.mu.Lock()
	defer i.mu.Unlock()
	app, ok := i.apps[name]
	if !ok {
		return fmt.Errorf("app %q not installed", name)
	}
	delete(app.Enabled, workspace)
	return i.saveRegistry()
}

func (i *Installer) IsEnabled(name, workspace string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()
	app, ok := i.apps[name]
	if !ok {
		return false
	}
	return app.Enabled[workspace]
}

func (i *Installer) List() []*AppState {
	i.mu.RLock()
	defer i.mu.RUnlock()
	result := make([]*AppState, 0, len(i.apps))
	for _, a := range i.apps {
		result = append(result, a)
	}
	return result
}

func (i *Installer) Get(name string) (*AppState, bool) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	a, ok := i.apps[name]
	return a, ok
}

func (i *Installer) GetRoutes(workspace string) []UIComponent {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var routes []UIComponent
	for _, app := range i.apps {
		if app.Enabled[workspace] && app.Manifest != nil {
			routes = append(routes, app.Manifest.UIComponents...)
		}
	}
	return routes
}

func (i *Installer) GetCrons(workspace string) []CronSpec {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var crons []CronSpec
	for _, app := range i.apps {
		if app.Enabled[workspace] && app.Manifest != nil {
			crons = append(crons, app.Manifest.Crons...)
		}
	}
	return crons
}

func (i *Installer) GetSkills(workspace string) []SkillSpec {
	i.mu.RLock()
	defer i.mu.RUnlock()
	var skills []SkillSpec
	for _, app := range i.apps {
		if app.Enabled[workspace] && app.Manifest != nil {
			skills = append(skills, app.Manifest.Skills...)
		}
	}
	return skills
}

// --- persistence ---

func (i *Installer) registryPath() string {
	return filepath.Join(i.storeDir, "registry.json")
}

func (i *Installer) saveRegistry() error {
	data, err := json.MarshalIndent(i.apps, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(i.registryPath(), data, 0644)
}

func (i *Installer) loadRegistry() {
	data, err := os.ReadFile(i.registryPath())
	if err != nil {
		return
	}
	json.Unmarshal(data, &i.apps)
}

// --- file copy helpers ---

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if info.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

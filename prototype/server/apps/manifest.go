package apps

import (
	"encoding/json"
	"fmt"
	"os"
)

type UIComponent struct {
	Name  string `json:"name"`
	Route string `json:"route"`
	Icon  string `json:"icon,omitempty"`
}

type CronSpec struct {
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Handler  string `json:"handler"`
}

type SkillSpec struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	File        string `json:"file"`
}

type Manifest struct {
	Name         string        `json:"name"`
	Version      string        `json:"version"`
	Description  string        `json:"description,omitempty"`
	Author       string        `json:"author,omitempty"`
	Permissions  []string      `json:"permissions,omitempty"`
	UIComponents []UIComponent `json:"ui_components,omitempty"`
	Crons        []CronSpec    `json:"crons,omitempty"`
	Skills       []SkillSpec   `json:"skills,omitempty"`
	EntryPoint   string        `json:"entry_point"`
}

func ParseManifest(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading manifest: %w", err)
	}
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("manifest: name is required")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest: version is required")
	}
	if m.EntryPoint == "" {
		return fmt.Errorf("manifest: entry_point is required")
	}
	for i, c := range m.Crons {
		if c.Name == "" || c.Schedule == "" || c.Handler == "" {
			return fmt.Errorf("manifest: cron[%d] requires name, schedule, and handler", i)
		}
	}
	for i, s := range m.Skills {
		if s.Name == "" || s.File == "" {
			return fmt.Errorf("manifest: skill[%d] requires name and file", i)
		}
	}
	for i, u := range m.UIComponents {
		if u.Name == "" || u.Route == "" {
			return fmt.Errorf("manifest: ui_component[%d] requires name and route", i)
		}
	}
	return nil
}

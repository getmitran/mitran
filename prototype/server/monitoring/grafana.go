package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type GrafanaConfig struct {
	ProvisioningDir string
	PrometheusURL   string
	DashboardTitle  string
}

func DefaultGrafanaConfig() GrafanaConfig {
	return GrafanaConfig{
		ProvisioningDir: "/etc/grafana/provisioning",
		PrometheusURL:   "http://localhost:9090",
		DashboardTitle:  "Mitran Agent Metrics",
	}
}

func GenerateDatasource(cfg GrafanaConfig) error {
	dir := filepath.Join(cfg.ProvisioningDir, "datasources")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create datasource dir: %w", err)
	}

	yaml := fmt.Sprintf(`apiVersion: 1
datasources:
  - name: Prometheus
    type: prometheus
    access: proxy
    url: %s
    isDefault: true
    editable: false
`, cfg.PrometheusURL)

	path := filepath.Join(dir, "mitran-prometheus.yaml")
	if err := os.WriteFile(path, []byte(yaml), 0644); err != nil {
		return fmt.Errorf("write datasource: %w", err)
	}
	return nil
}

func GenerateDashboard(cfg GrafanaConfig) error {
	dir := filepath.Join(cfg.ProvisioningDir, "dashboards")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create dashboard dir: %w", err)
	}

	// Write dashboard provider YAML
	providerYAML := fmt.Sprintf(`apiVersion: 1
providers:
  - name: Mitran
    folder: Mitran
    type: file
    options:
      path: %s/json
`, dir)

	if err := os.WriteFile(filepath.Join(dir, "mitran.yaml"), []byte(providerYAML), 0644); err != nil {
		return fmt.Errorf("write dashboard provider: %w", err)
	}

	// Write dashboard JSON
	jsonDir := filepath.Join(dir, "json")
	if err := os.MkdirAll(jsonDir, 0755); err != nil {
		return fmt.Errorf("create json dir: %w", err)
	}

	dashboard := map[string]interface{}{
		"title":         cfg.DashboardTitle,
		"uid":           "mitran-agents",
		"schemaVersion": 39,
		"timezone":      "browser",
		"panels": []map[string]interface{}{
			{
				"title": "Agent Task Duration",
				"type":  "timeseries",
				"gridPos": map[string]int{"h": 8, "w": 12, "x": 0, "y": 0},
				"targets": []map[string]string{
					{"expr": "histogram_quantile(0.95, rate(mitran_agent_task_duration_seconds_bucket[5m]))", "legendFormat": "p95 {{agent}}"},
				},
			},
			{
				"title": "Agent Tasks Total",
				"type":  "stat",
				"gridPos": map[string]int{"h": 8, "w": 12, "x": 12, "y": 0},
				"targets": []map[string]string{
					{"expr": "sum(mitran_agent_tasks_total) by (agent, status)", "legendFormat": "{{agent}} {{status}}"},
				},
			},
			{
				"title": "Active Agents",
				"type":  "gauge",
				"gridPos": map[string]int{"h": 8, "w": 8, "x": 0, "y": 8},
				"targets": []map[string]string{
					{"expr": "mitran_agents_active", "legendFormat": "active"},
				},
			},
			{
				"title": "LLM Token Usage",
				"type":  "timeseries",
				"gridPos": map[string]int{"h": 8, "w": 16, "x": 8, "y": 8},
				"targets": []map[string]string{
					{"expr": "rate(mitran_llm_tokens_total[5m])", "legendFormat": "{{agent}} {{direction}}"},
				},
			},
			{
				"title": "Error Rate",
				"type":  "timeseries",
				"gridPos": map[string]int{"h": 8, "w": 24, "x": 0, "y": 16},
				"targets": []map[string]string{
					{"expr": "rate(mitran_agent_errors_total[5m])", "legendFormat": "{{agent}} {{error_type}}"},
				},
			},
		},
	}

	data, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal dashboard: %w", err)
	}

	if err := os.WriteFile(filepath.Join(jsonDir, "mitran-agents.json"), data, 0644); err != nil {
		return fmt.Errorf("write dashboard json: %w", err)
	}
	return nil
}

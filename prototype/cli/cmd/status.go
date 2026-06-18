package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Uptime    string `json:"uptime"`
	UptimeSec int64  `json:"uptime_seconds"`
}

type PoolStats struct {
	RunningAgents int `json:"running_agents"`
	QueuedTasks   int `json:"queued_tasks"`
	MemoryEntries int `json:"memory_entries"`
	WorkerHealth  string `json:"worker_health"`
}

type CronJob struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Schedule string `json:"schedule"`
	Status   string `json:"status"`
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show live system stats",
	Run:   runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	base := fmt.Sprintf("http://localhost:%d", enginePort)

	health, err := fetchJSON[HealthResponse](base + "/health")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Engine unreachable at %s: %v\n", base, err)
		os.Exit(1)
	}

	pool, _ := fetchJSON[PoolStats](base + "/api/v1/pool/stats")
	crons, _ := fetchJSON[[]CronJob](base + "/api/v1/cron")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "\n  Mitran Engine Status\n")
	fmt.Fprintf(w, "  %s\n\n", "────────────────────────────")
	fmt.Fprintf(w, "  Status:\t%s\n", health.Status)
	fmt.Fprintf(w, "  Uptime:\t%s\n", formatUptime(health.UptimeSec))
	fmt.Fprintf(w, "  Running Agents:\t%d\n", pool.RunningAgents)
	fmt.Fprintf(w, "  Queued Tasks:\t%d\n", pool.QueuedTasks)
	fmt.Fprintf(w, "  Memory Entries:\t%d\n", pool.MemoryEntries)
	fmt.Fprintf(w, "  Worker Health:\t%s\n", pool.WorkerHealth)
	fmt.Fprintf(w, "  Cron Jobs:\t%d\n", len(*crons))
	fmt.Fprintln(w)
	w.Flush()

	if crons != nil && len(*crons) > 0 {
		fmt.Println("  Cron Jobs:")
		cw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(cw, "    ID\tNAME\tSCHEDULE\tSTATUS\n")
		for _, c := range *crons {
			fmt.Fprintf(cw, "    %s\t%s\t%s\t%s\n", c.ID, c.Name, c.Schedule, c.Status)
		}
		cw.Flush()
		fmt.Println()
	}
}

func fetchJSON[T any](url string) (*T, error) {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result T
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return &result, nil
}

func formatUptime(seconds int64) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}

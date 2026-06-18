package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Manage scheduled cron jobs",
}

var cronAddCmd = &cobra.Command{
	Use:   "add [name] [message]",
	Short: "Add a new cron job",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		every, _ := cmd.Flags().GetInt("every")
		cronExpr, _ := cmd.Flags().GetString("cron")
		agent, _ := cmd.Flags().GetString("agent")

		body := map[string]interface{}{
			"name":    args[0],
			"message": args[1],
		}
		if every > 0 {
			body["every"] = every
		}
		if cronExpr != "" {
			body["cron_expr"] = cronExpr
		}
		if agent != "" {
			body["agent"] = agent
		}

		resp, err := postJSON("/api/v1/cron", body)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created job: %s\n", resp["id"])
	},
}

var cronListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all cron jobs",
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := getJSON("/api/v1/cron")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		jobs, ok := resp["jobs"].([]interface{})
		if !ok || len(jobs) == 0 {
			fmt.Println("No cron jobs.")
			return
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSCHEDULE\tAGENT\tSTATUS")
		for _, j := range jobs {
			job := j.(map[string]interface{})
			schedule := ""
			if v, ok := job["every"]; ok && v != nil {
				schedule = fmt.Sprintf("every %vs", v)
			} else if v, ok := job["cron_expr"]; ok && v != nil {
				schedule = fmt.Sprintf("%v", v)
			}
			agent := ""
			if v, ok := job["agent"]; ok && v != nil {
				agent = fmt.Sprintf("%v", v)
			}
			status := "active"
			if v, ok := job["status"]; ok && v != nil {
				status = fmt.Sprintf("%v", v)
			}
			fmt.Fprintf(w, "%v\t%v\t%s\t%s\t%s\n", job["id"], job["name"], schedule, agent, status)
		}
		w.Flush()
	},
}

var cronPauseCmd = &cobra.Command{
	Use:   "pause [job-id]",
	Short: "Pause a cron job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := postJSON(fmt.Sprintf("/api/v1/cron/%s/pause", args[0]), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Paused job %s\n", args[0])
	},
}

var cronResumeCmd = &cobra.Command{
	Use:   "resume [job-id]",
	Short: "Resume a paused cron job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := postJSON(fmt.Sprintf("/api/v1/cron/%s/resume", args[0]), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Resumed job %s\n", args[0])
	},
}

var cronTriggerCmd = &cobra.Command{
	Use:   "trigger [job-id]",
	Short: "Trigger immediate execution of a cron job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		_, err := postJSON(fmt.Sprintf("/api/v1/cron/%s/trigger", args[0]), nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Triggered job %s\n", args[0])
	},
}

var cronRemoveCmd = &cobra.Command{
	Use:   "remove [job-id]",
	Short: "Remove a cron job",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		req, _ := http.NewRequest("DELETE", baseURL()+fmt.Sprintf("/api/v1/cron/%s", args[0]), nil)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			body, _ := io.ReadAll(resp.Body)
			fmt.Fprintf(os.Stderr, "Error: %s\n", body)
			os.Exit(1)
		}
		fmt.Printf("Removed job %s\n", args[0])
	},
}

func postJSON(path string, payload interface{}) (map[string]interface{}, error) {
	var body io.Reader
	if payload != nil {
		data, _ := json.Marshal(payload)
		body = bytes.NewReader(data)
	}
	resp, err := http.Post(baseURL()+path, "application/json", body)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error: %s", raw)
	}
	var result map[string]interface{}
	json.Unmarshal(raw, &result)
	return result, nil
}

func getJSON(path string) (map[string]interface{}, error) {
	resp, err := http.Get(baseURL() + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("server error: %s", raw)
	}
	var result map[string]interface{}
	json.Unmarshal(raw, &result)
	return result, nil
}

func baseURL() string {
	if url := os.Getenv("MITRAN_ENGINE_URL"); url != "" {
		return url
	}
	return "http://localhost:7780"
}

func init() {
	cronAddCmd.Flags().Int("every", 0, "Interval in seconds")
	cronAddCmd.Flags().String("cron", "", "Cron expression (5-field)")
	cronAddCmd.Flags().String("agent", "", "Agent name for this job")

	cronCmd.AddCommand(cronAddCmd)
	cronCmd.AddCommand(cronListCmd)
	cronCmd.AddCommand(cronPauseCmd)
	cronCmd.AddCommand(cronResumeCmd)
	cronCmd.AddCommand(cronTriggerCmd)
	cronCmd.AddCommand(cronRemoveCmd)

	rootCmd.AddCommand(cronCmd)
}

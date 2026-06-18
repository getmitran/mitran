package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const (
	colorReset  = "\033[0m"
	colorGreen  = "\033[32m"
	colorCyan   = "\033[36m"
	colorYellow = "\033[33m"
	colorBold   = "\033[1m"
	colorRed    = "\033[31m"
)

var engineURL = "http://localhost:7780"

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Mitran for your company",
	RunE:  runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

type initRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Languages   []string `json:"languages"`
	TeamSize    int      `json:"team_size"`
}

type taskResponse struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	AgentType string `json:"agent_type"`
	Status    string `json:"status"`
	Priority  int    `json:"priority"`
}

func runInit(cmd *cobra.Command, args []string) error {
	fmt.Printf("\n%s%s🚀 Mitran Init — Let's set up your AI engineering team%s\n\n", colorBold, colorCyan, colorReset)

	name := prompt("Company/project name?")
	description := prompt("What does your company do?")
	engineers := prompt("How many engineers?")
	stack := prompt("What languages/frameworks do you use?")

	fmt.Printf("\n%s%s📡 Sending to Mitran Engine (%s)...%s\n", colorBold, colorCyan, engineURL, colorReset)

	// Call the engine API
	payload := initRequest{
		Name:        name,
		Description: description,
		Languages:   splitLanguages(stack),
		TeamSize:    parseIntOr(engineers, 5),
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(engineURL+"/api/v1/init", "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Printf("\n%s%s✗ Cannot reach engine at %s%s\n", colorBold, colorRed, engineURL, colorReset)
		fmt.Printf("  Start the engine first: cd prototype/server && go run .\n\n")
		return nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		fmt.Printf("\n%s%s✗ Engine error (%d): %s%s\n\n", colorBold, colorRed, resp.StatusCode, string(respBody), colorReset)
		return nil
	}

	// Parse task plan from response
	var result struct {
		ProjectID string         `json:"project_id"`
		Tasks     []taskResponse `json:"tasks"`
	}
	json.Unmarshal(respBody, &result)

	fmt.Printf("\n%s%s📋 Task Plan (project: %s):%s\n\n", colorBold, colorYellow, result.ProjectID, colorReset)
	for i, t := range result.Tasks {
		fmt.Printf("  %s%d. [%-12s]%s %s\n", colorCyan, i+1, t.AgentType, colorReset, t.Title)
	}
	fmt.Println()

	approved := promptYN("Approve this plan? The scheduler will begin executing tasks")
	if !approved {
		fmt.Printf("\n%sPlan created but not started. Tasks remain queued.%s\n", colorYellow, colorReset)
		fmt.Printf("  View in dashboard: http://localhost:5173\n\n")
		return nil
	}

	fmt.Printf("\n%s%s⚡ Plan approved. Scheduler is now executing...%s\n\n", colorBold, colorGreen, colorReset)
	fmt.Printf("  Dashboard:  http://localhost:5173\n")
	fmt.Printf("  Tasks API:  %s/api/v1/tasks\n", engineURL)
	fmt.Printf("  Checkpoints: %s/api/v1/checkpoints\n\n", engineURL)

	// Poll for first checkpoint
	fmt.Printf("%sWaiting for first agent to complete...%s\n", colorCyan, colorReset)
	for i := 0; i < 60; i++ {
		time.Sleep(2 * time.Second)
		cps := fetchCheckpoints()
		if len(cps) > 0 {
			fmt.Printf("\n%s%s✅ First checkpoint ready!%s\n", colorBold, colorGreen, colorReset)
			fmt.Printf("  Agent: %s\n", cps[0].AgentType)
			fmt.Printf("  Review at: http://localhost:5173 (Checkpoints tab)\n\n")
			return nil
		}
		fmt.Printf(".")
	}
	fmt.Printf("\n\n  Still running. Check dashboard for progress.\n\n")
	return nil
}

func fetchCheckpoints() []taskResponse {
	resp, err := http.Get(engineURL + "/api/v1/checkpoints")
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	var cps []taskResponse
	body, _ := io.ReadAll(resp.Body)
	json.Unmarshal(body, &cps)
	return cps
}

func splitLanguages(s string) []string {
	parts := strings.Split(s, ",")
	var out []string
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

func parseIntOr(s string, fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return fallback
	}
	return n
}

func prompt(label string) string {
	p := promptui.Prompt{Label: label}
	result, err := p.Run()
	if err != nil {
		return ""
	}
	return result
}

func promptYN(label string) bool {
	p := promptui.Prompt{Label: label, IsConfirm: true}
	_, err := p.Run()
	return err == nil
}

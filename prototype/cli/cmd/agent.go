package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Manage Mitran agents",
}

var agentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configured agents",
	RunE:  runAgentList,
}

var agentCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new agent",
	RunE:  runAgentCreate,
}

var agentDeleteCmd = &cobra.Command{
	Use:   "delete [name]",
	Short: "Delete an agent",
	Args:  cobra.ExactArgs(1),
	RunE:  runAgentDelete,
}

var (
	agentName  string
	agentType  string
	agentModel string
)

func init() {
	agentCreateCmd.Flags().StringVar(&agentName, "name", "", "Agent name (required)")
	agentCreateCmd.Flags().StringVar(&agentType, "type", "dev", "Agent type")
	agentCreateCmd.Flags().StringVar(&agentModel, "model", "", "LLM model override")
	agentCreateCmd.MarkFlagRequired("name")

	agentCmd.AddCommand(agentListCmd)
	agentCmd.AddCommand(agentCreateCmd)
	agentCmd.AddCommand(agentDeleteCmd)
	rootCmd.AddCommand(agentCmd)
}

func engineBase() string {
	if u := os.Getenv("MITRAN_ENGINE_URL"); u != "" {
		return u
	}
	return "http://localhost:7780"
}

func runAgentList(cmd *cobra.Command, args []string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(engineBase() + "/api/v1/agents/config")
	if err != nil {
		return fmt.Errorf("engine unreachable: %w", err)
	}
	defer resp.Body.Close()

	var agents []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if len(agents) == 0 {
		fmt.Println("No agents configured.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tTYPE\tMODEL\tSTATUS\n")
	for _, a := range agents {
		name, _ := a["name"].(string)
		typ, _ := a["type"].(string)
		model, _ := a["model"].(string)
		status, _ := a["status"].(string)
		if status == "" {
			status = "active"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", name, typ, model, status)
	}
	w.Flush()
	return nil
}

func runAgentCreate(cmd *cobra.Command, args []string) error {
	payload, _ := json.Marshal(map[string]string{
		"name":  agentName,
		"type":  agentType,
		"model": agentModel,
	})

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(
		engineBase()+"/api/v1/agents/config",
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return fmt.Errorf("engine unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("✅ Agent '%s' created (type: %s)\n", agentName, agentType)
	return nil
}

func runAgentDelete(cmd *cobra.Command, args []string) error {
	name := args[0]
	req, _ := http.NewRequest(http.MethodDelete, engineBase()+"/api/v1/agents/config/"+name, nil)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("engine unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed (%d): %s", resp.StatusCode, string(body))
	}

	fmt.Printf("🗑️  Agent '%s' deleted.\n", name)
	return nil
}

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

var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage Mitran workspaces",
}

var wsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List workspaces",
	RunE:  runWsList,
}

var wsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new workspace",
	RunE:  runWsCreate,
}

var wsSwitchCmd = &cobra.Command{
	Use:   "switch [name]",
	Short: "Switch active workspace",
	Args:  cobra.ExactArgs(1),
	RunE:  runWsSwitch,
}

var (
	wsName string
	wsPath string
)

func init() {
	wsCreateCmd.Flags().StringVar(&wsName, "name", "", "Workspace name (required)")
	wsCreateCmd.Flags().StringVar(&wsPath, "path", "", "Workspace directory path")
	wsCreateCmd.MarkFlagRequired("name")

	workspaceCmd.AddCommand(wsListCmd)
	workspaceCmd.AddCommand(wsCreateCmd)
	workspaceCmd.AddCommand(wsSwitchCmd)
	rootCmd.AddCommand(workspaceCmd)
}

func runWsList(cmd *cobra.Command, args []string) error {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(engineBase() + "/api/v1/workspaces")
	if err != nil {
		return fmt.Errorf("engine unreachable: %w", err)
	}
	defer resp.Body.Close()

	var workspaces []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&workspaces); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if len(workspaces) == 0 {
		fmt.Println("No workspaces found.")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "NAME\tPATH\tACTIVE\n")
	for _, ws := range workspaces {
		name, _ := ws["name"].(string)
		path, _ := ws["path"].(string)
		active, _ := ws["active"].(bool)
		marker := ""
		if active {
			marker = "✓"
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", name, path, marker)
	}
	w.Flush()
	return nil
}

func runWsCreate(cmd *cobra.Command, args []string) error {
	payload, _ := json.Marshal(map[string]string{
		"name": wsName,
		"path": wsPath,
	})

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Post(
		engineBase()+"/api/v1/workspaces",
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

	fmt.Printf("✅ Workspace '%s' created.\n", wsName)
	return nil
}

func runWsSwitch(cmd *cobra.Command, args []string) error {
	name := args[0]
	payload, _ := json.Marshal(map[string]string{"name": name})

	req, _ := http.NewRequest(http.MethodPut, engineBase()+"/api/v1/workspaces/active", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

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

	fmt.Printf("🔄 Switched to workspace '%s'.\n", name)
	return nil
}

package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var agentChatCmd = &cobra.Command{
	Use:   "chat [agent-name]",
	Short: "Interactive chat with a specific agent",
	Long:  "Start a conversational REPL with a named Mitran agent. Reads stdin, POSTs to the engine, and prints responses.",
	Args:  cobra.ExactArgs(1),
	RunE:  runAgentChat,
}

func init() {
	agentCmd.AddCommand(agentChatCmd)
}

func runAgentChat(cmd *cobra.Command, args []string) error {
	agent := args[0]
	base := engineBase()

	fmt.Printf("🤖 Mitran Agent Chat (agent: %s)\n", agent)
	fmt.Println("Type 'exit' or Ctrl+D to quit.\n")

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("you> ")
		if !scanner.Scan() {
			fmt.Println()
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}
		if input == "exit" || input == "quit" {
			break
		}

		resp, err := postTask(base, agent, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			continue
		}
		fmt.Printf("agent> %s\n\n", resp)
	}
	return nil
}

func postTask(base, agent, message string) (string, error) {
	payload, _ := json.Marshal(map[string]string{
		"agent":   agent,
		"message": message,
	})

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Post(
		base+"/api/v1/tasks",
		"application/json",
		bytes.NewReader(payload),
	)
	if err != nil {
		return "", fmt.Errorf("engine unreachable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("engine returned %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		TaskID   string `json:"task_id"`
		Status   string `json:"status"`
		Response string `json:"response"`
	}

	body, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("decode error: %w", err)
	}

	// If response is immediate, return it
	if result.Response != "" {
		return result.Response, nil
	}

	// Otherwise poll for result
	if result.TaskID != "" {
		return pollTaskResult(base, result.TaskID)
	}

	return string(body), nil
}

func pollTaskResult(base, taskID string) (string, error) {
	client := &http.Client{Timeout: 5 * time.Second}
	for i := 0; i < 60; i++ {
		time.Sleep(time.Second)

		resp, err := client.Get(fmt.Sprintf("%s/api/v1/tasks/%s", base, taskID))
		if err != nil {
			continue
		}

		var result struct {
			Status   string `json:"status"`
			Response string `json:"response"`
		}
		json.NewDecoder(resp.Body).Decode(&result)
		resp.Body.Close()

		if result.Status == "completed" || result.Status == "done" {
			return result.Response, nil
		}
		if result.Status == "failed" || result.Status == "error" {
			return "", fmt.Errorf("task failed: %s", result.Response)
		}
	}
	return "", fmt.Errorf("timeout waiting for task %s", taskID)
}

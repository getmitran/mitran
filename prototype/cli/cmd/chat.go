package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var chatAgent string

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Interactive chat with a Mitran agent",
	Long:  "Start a REPL-mode conversational interface with streaming responses from a Mitran agent.",
	RunE:  runChat,
}

func init() {
	chatCmd.Flags().StringVar(&chatAgent, "agent", "dev", "Agent to chat with")
	rootCmd.AddCommand(chatCmd)
}

func runChat(cmd *cobra.Command, args []string) error {
	engineURL := os.Getenv("MITRAN_ENGINE_URL")
	if engineURL == "" {
		engineURL = "http://localhost:7780"
	}

	fmt.Printf("🤖 Mitran Chat (agent: %s)\n", chatAgent)
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

		if err := streamChat(engineURL, chatAgent, input); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
		fmt.Println()
	}
	return nil
}

func streamChat(engineURL, agent, message string) error {
	payload, _ := json.Marshal(map[string]string{
		"agent":   agent,
		"message": message,
	})

	resp, err := http.Post(
		engineURL+"/api/v1/chat",
		"application/json",
		strings.NewReader(string(payload)),
	)
	if err != nil {
		return fmt.Errorf("failed to connect to engine: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("engine returned %d: %s", resp.StatusCode, string(body))
	}

	fmt.Print("agent> ")
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if len(line) > 0 {
			s := strings.TrimSpace(string(line))
			if strings.HasPrefix(s, "data: ") {
				token := strings.TrimPrefix(s, "data: ")
				if token == "[DONE]" {
					break
				}
				fmt.Print(token)
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("stream read error: %w", err)
		}
	}
	fmt.Println()
	return nil
}

package cmd

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var logsFollow bool

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Stream Mitran engine events in real-time",
	Long:  "Connects to the engine's SSE endpoint (/api/v1/events) and streams events to stdout.",
	RunE:  runLogs,
}

func init() {
	logsCmd.Flags().BoolVarP(&logsFollow, "follow", "f", false, "Follow mode (stream continuously)")
	rootCmd.AddCommand(logsCmd)
}

func runLogs(cmd *cobra.Command, args []string) error {
	base := engineBase()
	url := base + "/api/v1/events"

	fmt.Fprintf(os.Stderr, "📡 Connecting to %s ...\n", url)

	client := &http.Client{Timeout: 0} // no timeout for SSE
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned %d", resp.StatusCode)
	}

	fmt.Fprintf(os.Stderr, "✅ Connected. Streaming events...\n\n")

	scanner := bufio.NewScanner(resp.Body)
	var eventType string

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "event: ") {
			eventType = strings.TrimPrefix(line, "event: ")
			continue
		}

		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			ts := time.Now().Format("15:04:05")
			if eventType != "" {
				fmt.Printf("[%s] %s: %s\n", ts, eventType, data)
				eventType = ""
			} else {
				fmt.Printf("[%s] %s\n", ts, data)
			}
			continue
		}

		// Empty line = end of event
		if line == "" {
			eventType = ""
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	fmt.Fprintf(os.Stderr, "\n⚡ Connection closed.\n")
	return nil
}

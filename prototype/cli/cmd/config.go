package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Mitran configuration",
}

var configGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get configuration values",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := http.Get(getAPIURL() + "/api/v1/settings")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(os.Stderr, "Error: %s\n", body)
			os.Exit(1)
		}

		var settings map[string]interface{}
		json.Unmarshal(body, &settings)

		if len(args) == 0 {
			out, _ := json.MarshalIndent(settings, "", "  ")
			fmt.Println(string(out))
			return
		}

		val := getNestedValue(settings, args[0])
		if val == nil {
			fmt.Fprintf(os.Stderr, "Key not found: %s\n", args[0])
			os.Exit(1)
		}
		switch v := val.(type) {
		case string:
			fmt.Println(v)
		default:
			out, _ := json.MarshalIndent(v, "", "  ")
			fmt.Println(string(out))
		}
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		key, value := args[0], args[1]

		// Get current settings
		resp, err := http.Get(getAPIURL() + "/api/v1/settings")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		var settings map[string]interface{}
		json.Unmarshal(body, &settings)

		// Set nested value
		setNestedValue(settings, key, value)

		// PUT updated settings
		payload, _ := json.Marshal(settings)
		req, _ := http.NewRequest(http.MethodPut, getAPIURL()+"/api/v1/settings", strings.NewReader(string(payload)))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp2, err := client.Do(req)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		defer resp2.Body.Close()

		if resp2.StatusCode != http.StatusOK {
			b, _ := io.ReadAll(resp2.Body)
			fmt.Fprintf(os.Stderr, "Error: %s\n", b)
			os.Exit(1)
		}
		fmt.Printf("Set %s = %s\n", key, value)
	},
}

func getAPIURL() string {
	if url := os.Getenv("MITRAN_API_URL"); url != "" {
		return url
	}
	return "http://localhost:7780"
}

func getNestedValue(m map[string]interface{}, key string) interface{} {
	parts := strings.Split(key, ".")
	var current interface{} = m
	for _, p := range parts {
		cm, ok := current.(map[string]interface{})
		if !ok {
			return nil
		}
		current, ok = cm[p]
		if !ok {
			return nil
		}
	}
	return current
}

func setNestedValue(m map[string]interface{}, key, value string) {
	parts := strings.Split(key, ".")
	current := m
	for i, p := range parts {
		if i == len(parts)-1 {
			current[p] = value
			return
		}
		next, ok := current[p].(map[string]interface{})
		if !ok {
			next = make(map[string]interface{})
			current[p] = next
		}
		current = next
	}
}

func init() {
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configSetCmd)
	rootCmd.AddCommand(configCmd)
}

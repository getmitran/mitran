package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var Version = "0.1.0"

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Check for newer Mitran CLI version",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := &http.Client{CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}}
		resp, err := client.Get("https://api.github.com/repos/getmitran/mitran/releases/latest")
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
		}

		var release struct {
			TagName string `json:"tag_name"`
			HTMLURL string `json:"html_url"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		latest := strings.TrimPrefix(release.TagName, "v")
		if latest != Version && latest > Version {
			fmt.Printf("Update available: %s -> %s\nDownload: %s\n", Version, latest, release.HTMLURL)
		} else {
			fmt.Println("Already up to date.")
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
}

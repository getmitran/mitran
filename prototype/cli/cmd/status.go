package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show agent statuses",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\033[1m\033[36m📊 Mitran Agent Status\033[0m\n")
		statuses := []struct{ name, state, task string }{
			{"Dev Agent", "🟢 idle", "—"},
			{"Docs Agent", "🟢 idle", "—"},
			{"Ops Agent", "🟢 idle", "—"},
			{"Review Agent", "🟢 idle", "—"},
			{"HR Agent", "🟢 idle", "—"},
			{"CI/CD Agent", "🟢 idle", "—"},
			{"Tickets Agent", "🟢 idle", "—"},
			{"Wiki Agent", "🟢 idle", "—"},
		}
		fmt.Printf("  %-14s %-10s %s\n", "AGENT", "STATUS", "CURRENT TASK")
		fmt.Printf("  %-14s %-10s %s\n", "─────", "──────", "────────────")
		for _, s := range statuses {
			fmt.Printf("  %-14s %-10s %s\n", s.name, s.state, s.task)
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the Mitran engine",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\033[1m\033[32m🟢 Starting Mitran engine on :7777\033[0m")
		fmt.Println("   Dashboard: http://localhost:7777")
		fmt.Println("   API:       http://localhost:7777/api")
		fmt.Println("\n   Press Ctrl+C to stop")
		select {}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

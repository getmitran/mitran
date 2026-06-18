package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Start Mitran in demo mode (no cloud credentials required)",
	Run: func(cmd *cobra.Command, args []string) {
		os.Setenv("MITRAN_LLM_PROVIDER", "mock")
		os.Setenv("MITRAN_WORKER_URL", "")
		fmt.Println("Starting Mitran demo (self-contained, no external deps)...")
		fmt.Println("Dashboard: http://localhost:7780")
		fmt.Println("Try: curl http://localhost:7780/api/v1/health")
		serveCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

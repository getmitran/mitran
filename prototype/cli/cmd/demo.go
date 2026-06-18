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
		fmt.Println("Starting Mitran in demo mode (no cloud credentials required)...")
		os.Setenv("MITRAN_LLM_PROVIDER", "mock")
		serveCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
}

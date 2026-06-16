package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "mitran",
	Short: "Mitran — your AI engineering team",
	Long:  "Mitran is an agent-native platform that runs your internal tooling with AI agents.\nHumans set direction. Agents execute.",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

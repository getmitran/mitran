package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generate shell completions",
}

var bashCompletionCmd = &cobra.Command{
	Use:   "bash",
	Short: "Generate bash completions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenBashCompletion(os.Stdout)
	},
}

var zshCompletionCmd = &cobra.Command{
	Use:   "zsh",
	Short: "Generate zsh completions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenZshCompletion(os.Stdout)
	},
}

var fishCompletionCmd = &cobra.Command{
	Use:   "fish",
	Short: "Generate fish completions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return rootCmd.GenFishCompletion(os.Stdout, true)
	},
}

func init() {
	completionCmd.AddCommand(bashCompletionCmd, zshCompletionCmd, fishCompletionCmd)
	rootCmd.AddCommand(completionCmd)
}

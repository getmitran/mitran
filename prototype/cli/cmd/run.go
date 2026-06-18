package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run [task-file]",
	Short: "Run a task from a spec file",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			fmt.Printf("Running task: %s\n", args[0])
		}
		fmt.Println("Task execution coming in v0.2.0")
	},
}

func init() { rootCmd.AddCommand(runCmd) }

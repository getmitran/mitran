package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "Stream Mitran service logs",
	Run: func(cmd *cobra.Command, args []string) {
		c := exec.Command("docker", "compose", "logs", "--tail=50", "-f")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}

func init() { rootCmd.AddCommand(logsCmd) }

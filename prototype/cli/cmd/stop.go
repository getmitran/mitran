package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop Mitran services",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Stopping Mitran services...")
		c := exec.Command("docker", "compose", "down")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}

func init() { rootCmd.AddCommand(stopCmd) }

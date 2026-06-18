package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start Mitran services via Docker Compose",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting Mitran services...")
		c := exec.Command("docker", "compose", "up", "-d")
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		c.Run()
	},
}

func init() { rootCmd.AddCommand(upCmd) }

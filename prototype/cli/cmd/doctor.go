package cmd

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check system dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		checks := []struct{ name, bin, arg string }{
			{"Go", "go", "version"},
			{"Python", "python3", "--version"},
			{"Node", "node", "--version"},
			{"Docker", "docker", "--version"},
		}
		found := 0
		for _, c := range checks {
			out, err := exec.Command(c.bin, c.arg).Output()
			if err == nil {
				fmt.Printf("✓ %s: %s\n", c.name, strings.TrimSpace(string(out)))
				found++
			} else {
				fmt.Printf("✗ %s: not found\n", c.name)
			}
		}
		fmt.Printf("\n%d/4 dependencies found\n", found)
	},
}

func init() { rootCmd.AddCommand(doctorCmd) }

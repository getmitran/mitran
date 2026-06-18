package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var companyName string

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a Mitran workspace",
	RunE: func(cmd *cobra.Command, args []string) error {
		config := fmt.Sprintf(`company: %s
agents: [dev, ops, tickets, wiki, hr]
llm_provider: bedrock
engine_port: 7780
`, companyName)
		if err := os.WriteFile("mitran.yaml", []byte(config), 0644); err != nil {
			return fmt.Errorf("failed to write mitran.yaml: %w", err)
		}
		if err := os.MkdirAll("data", 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
		fmt.Printf("Initialized Mitran workspace for %s. Run `mitran serve` to start.\n", companyName)
		return nil
	},
}

func init() {
	initCmd.Flags().StringVar(&companyName, "name", "my-company", "Company name")
	rootCmd.AddCommand(initCmd)
}

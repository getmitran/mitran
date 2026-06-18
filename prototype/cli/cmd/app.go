package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage apps/plugins (install, enable, disable, list)",
}

var appInstallCmd = &cobra.Command{
	Use:   "install [path-or-url]",
	Short: "Install an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := postJSON("/api/v1/plugins/install", map[string]interface{}{
			"source": args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to install app: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

var appEnableCmd = &cobra.Command{
	Use:   "enable [name]",
	Short: "Enable an installed app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := postJSON("/api/v1/plugins/enable", map[string]interface{}{
			"name": args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to enable app: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

var appDisableCmd = &cobra.Command{
	Use:   "disable [name]",
	Short: "Disable an app",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := postJSON("/api/v1/plugins/disable", map[string]interface{}{
			"name": args[0],
		})
		if err != nil {
			return fmt.Errorf("failed to disable app: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed apps",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := getJSON("/api/v1/plugins")
		if err != nil {
			return fmt.Errorf("failed to list apps: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

func init() {
	appCmd.AddCommand(appInstallCmd)
	appCmd.AddCommand(appEnableCmd)
	appCmd.AddCommand(appDisableCmd)
	appCmd.AddCommand(appListCmd)
	rootCmd.AddCommand(appCmd)
}

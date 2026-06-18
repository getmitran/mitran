package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Manage snapshots (create, restore, list)",
}

var snapshotCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new snapshot",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := postJSON("/api/v1/snapshot", map[string]interface{}{})
		if err != nil {
			return fmt.Errorf("failed to create snapshot: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

var snapshotRestoreCmd = &cobra.Command{
	Use:   "restore [snapshot-id]",
	Short: "Restore from a snapshot",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		components, _ := cmd.Flags().GetString("components")
		body := map[string]interface{}{
			"snapshot_id": args[0],
		}
		if components != "" {
			body["components"] = strings.Split(components, ",")
		}
		resp, err := postJSON("/api/v1/restore", body)
		if err != nil {
			return fmt.Errorf("failed to restore snapshot: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

var snapshotListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available snapshots",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := getJSON("/api/v1/snapshots")
		if err != nil {
			return fmt.Errorf("failed to list snapshots: %w", err)
		}
		fmt.Println(resp)
		return nil
	},
}

func init() {
	snapshotRestoreCmd.Flags().String("components", "", "Comma-separated components to restore")
	snapshotCmd.AddCommand(snapshotCreateCmd)
	snapshotCmd.AddCommand(snapshotRestoreCmd)
	snapshotCmd.AddCommand(snapshotListCmd)
	rootCmd.AddCommand(snapshotCmd)
}

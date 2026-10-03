package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var authorizationsCmd = &cobra.Command{
	Use:   "authorizations",
	Short: "Manage device authorizations",
}
var authId uint32

var authorizationsListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all authorizations stored on the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			entries, err := flow.GetAuthorizationEntries(ctx, 0, 100)
			if err != nil {
				return fmt.Errorf("failed to list authorizations: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(entries)
			}
			t := table.New().Headers("Auth ID", "Type", "Name", "Enabled", "Remote", "Lock Count")
			for _, e := range entries {
				t = t.Row(
					fmt.Sprintf("%d", e.AuthorizationID),
					fmt.Sprintf("%s", e.IDType),
					e.Name,
					fmt.Sprintf("%t", e.Enabled),
					fmt.Sprintf("%t", e.RemoteAllowed),
					fmt.Sprintf("%d", e.LockCount),
				)
			}
			fmt.Println(t)
			return nil
		})
	},
}

var authorizationsRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an authorization from the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.RemoveAuthorizationEntry(ctx, authId); err != nil {
				return fmt.Errorf("failed to remove authorization: %w", err)
			}
			fmt.Printf("Authorization %d removed.\n", authId)
			return nil
		})
	},
}

var (
	authUpdateName          string
	authUpdateEnabled       bool
	authUpdateRemoteAllowed bool
)

var authorizationsUpdateCmd = &cobra.Command{
	Use:     "update",
	Short:   "Update an authorization entry on the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			req := &blecommands.UpdateAuthorizationEntry{
				AuthorizationID: authId,
				Name:            authUpdateName,
				Enabled:         authUpdateEnabled,
				RemoteAllowed:   authUpdateRemoteAllowed,
			}
			if err := flow.UpdateAuthorizationEntry(ctx, req); err != nil {
				return fmt.Errorf("failed to update authorization: %w", err)
			}
			fmt.Printf("Authorization %d updated.\n", authId)
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(authorizationsCmd)
	authorizationsCmd.AddCommand(authorizationsListCmd)
	authorizationsCmd.AddCommand(authorizationsRemoveCmd)
	authorizationsCmd.AddCommand(authorizationsUpdateCmd)

	authorizationsRemoveCmd.Flags().Uint32Var(&authId, "auth-id", 0, "Authorization ID to remove")
	_ = authorizationsRemoveCmd.MarkFlagRequired("auth-id")

	authorizationsUpdateCmd.Flags().Uint32Var(&authId, "auth-id", 0, "Authorization ID to update")
	authorizationsUpdateCmd.Flags().StringVar(&authUpdateName, "name", "", "New name")
	authorizationsUpdateCmd.Flags().BoolVar(&authUpdateEnabled, "enabled", true, "Enable or disable this authorization")
	authorizationsUpdateCmd.Flags().BoolVar(&authUpdateRemoteAllowed, "remote-allowed", true, "Allow remote access")
	_ = authorizationsUpdateCmd.MarkFlagRequired("auth-id")
	_ = authorizationsUpdateCmd.MarkFlagRequired("name")
}

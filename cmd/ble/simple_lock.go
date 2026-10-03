package cmd

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var simpleLockCmd = &cobra.Command{
	Use:     "simple-lock",
	Short:   "Lock the device using SimpleLockAction (no PIN required)",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.PerformSimpleLockAction(ctx, blecommands.Lock); err != nil {
				return fmt.Errorf("simple lock failed: %w", err)
			}
			fmt.Println("Lock action completed.")
			return nil
		})
	},
}

var simpleUnlockCmd = &cobra.Command{
	Use:     "simple-unlock",
	Short:   "Unlock the device using SimpleLockAction (no PIN required)",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.PerformSimpleLockAction(ctx, blecommands.Unlock); err != nil {
				return fmt.Errorf("simple unlock failed: %w", err)
			}
			fmt.Println("Unlock action completed.")
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(simpleLockCmd)
	bleCmd.AddCommand(simpleUnlockCmd)
}

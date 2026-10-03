package cmd

import (
	"context"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var unlatchCmd = &cobra.Command{
	Use:     "unlatch",
	Short:   "Unlatch the door (open latch without unlocking deadbolt)",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			return flow.PerformLockOperation(ctx, blecommands.Unlatch)
		})
	},
}

var lockNGoCmd = &cobra.Command{
	Use:     "lock-n-go",
	Short:   "Unlock, wait for door to open, then automatically lock",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			return flow.PerformLockOperation(ctx, blecommands.LockAndGo)
		})
	},
}

func init() {
	bleCmd.AddCommand(unlatchCmd)
	bleCmd.AddCommand(lockNGoCmd)
}

package cmd

import (
	"context"
	"fmt"
	"time"

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
		// The device reports completion only after relocking, so the exchange lasts at
		// least the configured lock 'n' go timeout (up to 255s).
		return withAuthenticatedFlowTimeout(2*bleTimeout+255*time.Second, func(ctx context.Context, flow *bleflows.Flow) error {
			cfgCtx, cancel := context.WithTimeout(ctx, bleTimeout)
			adv, err := flow.GetAdvancedConfig(cfgCtx)
			cancel()
			if err != nil {
				return fmt.Errorf("failed to read lock 'n' go timeout: %w", err)
			}
			opCtx, cancel := context.WithTimeout(ctx, time.Duration(adv.LockNGoTimeout)*time.Second+bleTimeout)
			defer cancel()
			return flow.PerformLockOperation(opCtx, blecommands.LockAndGo)
		})
	},
}

func init() {
	bleCmd.AddCommand(unlatchCmd)
	bleCmd.AddCommand(lockNGoCmd)
}

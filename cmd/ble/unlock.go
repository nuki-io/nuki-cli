package cmd

import (
	"context"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:     "unlock",
	Short:   "Unlock a device via Bluetooth",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			return flow.PerformLockOperation(ctx, blecommands.Unlock)
		})
	},
}

func init() {
	bleCmd.AddCommand(unlockCmd)
}

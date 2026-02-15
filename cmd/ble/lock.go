package cmd

import (
	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:     "lock",
	Short:   "Lock a device via Bluetooth",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		withAuthenticatedFlow(func(flow *bleflows.Flow) error {
			return flow.PerformLockOperation(blecommands.Lock)
		})
	},
}

func init() {
	bleCmd.AddCommand(lockCmd)
}

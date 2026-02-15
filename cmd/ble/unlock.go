package cmd

import (
	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

// unlockCmd represents the unlock command
var unlockCmd = &cobra.Command{
	Use:     "unlock",
	Short:   "Unlock a device via Bluetooth",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		withAuthenticatedFlow(func(flow *bleflows.Flow) error {
			return flow.PerformLockOperation(blecommands.Unlock)
		})
	},
}

func init() {
	bleCmd.AddCommand(unlockCmd)
}

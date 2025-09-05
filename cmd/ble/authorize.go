package cmd

import (
	"time"

	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/cobra"
)

var pin string

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:     "authorize",
	Short:   "Authorizes and pairs this machine with the given Nuki device",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		err = ble.ScanForDevice(deviceId, 10*time.Second)
		if err != nil {
			c.Logger.Error("Failed to scan", "error", err.Error())
			return
		}
		flow, err := bleflows.NewUnauthenticatedFlow(ble, deviceId)
		if flow == nil {
			c.Logger.Error("Failed to create BLE flow", "error", err.Error())
			return
		}
		defer flow.DisconnectDevice()

		err = flow.Authorize(pin)
		if err != nil {
			c.Logger.Error("Failed to authorize", "error", err.Error())
			return
		}
	},
}

func init() {
	bleCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringVarP(&pin, "pin", "p", "", "The PIN code to use for authorization.")
}

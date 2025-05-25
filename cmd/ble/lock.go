package cmd

import (
	"runtime"
	"time"

	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/cobra"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:     "lock",
	Short:   "Lock a device via Bluetooth",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		if runtime.GOOS == "linux" {
			err = ble.ScanForDevice(deviceId, 10*time.Second)
			if err != nil {
				c.Logger.Error("Failed to scan", "error", err.Error())
				return
			}
		}
		flow := bleflows.NewFlow(ble)

		err = flow.PerformLockOperation(deviceId, blecommands.Lock)
		if err != nil {
			c.Logger.Error("Failed to perform lock operation", "error", err.Error())
			return
		}
	},
}

func init() {
	bleCmd.AddCommand(lockCmd)
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"
	"time"

	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"

	"github.com/spf13/cobra"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:     "state",
	Short:   "Gets the current lock state of the device",
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
		flow := bleflows.NewAuthenticatedFlow(ble, deviceId)
		defer flow.DisconnectDevice()

		status, err := flow.GetStatus()
		if err != nil {
			c.Logger.Error("Failed to get status", "error", err.Error())
			return
		}
		c.Logger.Info("Current lock state", "state", status.LockState.String(), "battery", fmt.Sprintf("%d%%", status.BatteryPercentage))
	},
}

func init() {
	bleCmd.AddCommand(stateCmd)
}

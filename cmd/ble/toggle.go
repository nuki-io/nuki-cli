/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
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

var repeats int

// toggleCmd represents the toggle command
var toggleCmd = &cobra.Command{
	Use:   "toggle",
	Short: "Toggles the current lock state",
	Long:  `Depending on the lock's current state, this command either locks or unlocks.`,
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
		var startAt int
		if status.LockState == blecommands.LockStateLocked {
			startAt = 0
		} else {
			startAt = 1
		}
		for i := range repeats {
			if i%2 == startAt {
				err = flow.PerformLockOperation(blecommands.Unlock)
			} else {
				err = flow.PerformLockOperation(blecommands.Lock)
			}
			// TODO: although we received the StatusComplete at this point, we apparently need to wait a bit longer
			time.Sleep(500 * time.Millisecond)
		}

		if err != nil {
			c.Logger.Error("Failed to perform lock operation", "error", err.Error())
			return
		}

	},
}

func init() {
	bleCmd.AddCommand(toggleCmd)
	toggleCmd.Flags().IntVarP(&repeats, "repeats", "n", 1, "The number of times to repeat the lock operation. Default is 1, which means it will toggle once. If set to 2, it will toggle twice, etc. This is useful for testing purposes.")
}

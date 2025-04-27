/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"runtime"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	c "go.nuki.io/nuki/nukictl/cmd"
	"go.nuki.io/nuki/nukictl/pkg/blecommands"
	"go.nuki.io/nuki/nukictl/pkg/bleflows"
	"go.nuki.io/nuki/nukictl/pkg/nukible"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the given (already paired) device",

	Run: func(cmd *cobra.Command, args []string) {
		// TODO: handle with cobra builtin functionality
		if deviceId == "" {
			c.Logger.Error("Please specify a device-id")
			return
		}
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

		action, _ := strconv.Atoi(args[0])
		err = flow.PerformLockOperation(deviceId, blecommands.Action(action))
		if err != nil {
			c.Logger.Error("Failed to perform lock operation", "error", err.Error(), "action", action)
			return
		}
	},
}

func init() {
	devicesCmd.AddCommand(lockCmd)
}

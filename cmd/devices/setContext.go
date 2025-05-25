/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// setContextCmd represents the setContext command
var setContextCmd = &cobra.Command{
	Use:   "set-context",
	Short: "Sets the active device used for device specific commands.",
	Long: `Instead of always specifying the device-id, you can set the active device with this command.
This is useful for commands that require a device-id, but you don't want to specify it every time.
The device-id is stored in the config file and used for all commands that require a device-id.`,
	Example: `nukictl devices set-context 1234567890abcdef`,
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("activeContext", deviceId)
		err := viper.WriteConfig()
		if err != nil {
			c.Logger.Error("Failed to write config file", "error", err)
			return
		}
		c.Logger.Info("Set active device to", "deviceId", deviceId)
	},
}

func init() {
	devicesCmd.AddCommand(setContextCmd)
}

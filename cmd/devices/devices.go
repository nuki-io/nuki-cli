/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.nuki.io/nuki/nukictl/cmd"
)

var (
	deviceId string
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:              "devices",
	Short:            "Command that interacts with devices through BLE.",
	Long:             ``,
	PersistentPreRun: preRun,
}

func init() {
	cmd.RootCmd.AddCommand(devicesCmd)
	devicesCmd.PersistentFlags().StringVarP(&deviceId, "device-id", "d", "", "The device to use. If not set, the device set by set-context command is used. This is ignored for some commands.")
	// viper.BindPFlag("activeContext", devicesCmd.PersistentFlags().Lookup("device-id"))
}

func preRun(cmd *cobra.Command, args []string) {
	if viper.IsSet("activecontext") {
		deviceId = viper.GetString("activecontext")
	}
	// TODO: The following "should" work. Check why it doesn't.
	// viper.BindPFlag("activeContext", cmd.PersistentFlags().Lookup("device-id"))
}

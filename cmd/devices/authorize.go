/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"time"

	"github.com/spf13/cobra"
	c "go.nuki.io/nuki/nukictl/cmd"
	"go.nuki.io/nuki/nukictl/pkg/bleflows"
	"go.nuki.io/nuki/nukictl/pkg/nukible"
)

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorizes and pairs this machine with the given Nuki device",
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		err = ble.ScanForDevice(args[0], 10*time.Second)
		if err != nil {
			c.Logger.Error("Failed to scan", "error", err.Error())
			return
		}
		flow := bleflows.NewFlow(ble)

		err = flow.Authorize(args[0])
		if err != nil {
			c.Logger.Error("Failed to authorize", "error", err.Error())
			return
		}
	},
}

func init() {
	devicesCmd.AddCommand(authorizeCmd)
}

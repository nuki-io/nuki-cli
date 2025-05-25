/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/nuki-io/nuki-cli/cmd"
	"github.com/spf13/cobra"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Command to interact with devices through BLE or the Nuki Web API",
}

func init() {
	cmd.RootCmd.AddCommand(devicesCmd)
}

/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// devicesCmd represents the devices command
var devicesCmd = &cobra.Command{
	Use:   "devices",
	Short: "Command that interacts with devices through BLE.",
	Long:  ``,
}

func init() {
	rootCmd.AddCommand(devicesCmd)
}

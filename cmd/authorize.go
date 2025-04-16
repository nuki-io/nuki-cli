/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
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
			log.Printf("Failed to start scan. %s\n", err.Error())
			return
		}
		ble.ScanForDevice(args[0], 10*time.Second)
		flow := bleflows.NewFlow(ble)

		err = flow.Authorize(args[0])
		if err != nil {
			log.Printf("Failed to authorize. %s\n", err.Error())
			return
		}
	},
}

func init() {
	devicesCmd.AddCommand(authorizeCmd)
}

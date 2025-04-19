/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"go.nuki.io/nuki/nukictl/pkg/blecommands"
	"go.nuki.io/nuki/nukictl/pkg/bleflows"
	"go.nuki.io/nuki/nukictl/pkg/nukible"
)

// lockCmd represents the lock command
var lockCmd = &cobra.Command{
	Use:   "lock",
	Short: "Lock the given (already paired) device",

	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			log.Printf("Failed to start scan. %s\n", err.Error())
			return
		}
		ble.ScanForDevice(args[0], 10*time.Second)
		flow := bleflows.NewFlow(ble)

		action, _ := strconv.Atoi(args[1])
		err = flow.PerformLockOperation(args[0], blecommands.Action(action))
		if err != nil {
			log.Printf("Failed to lock. %s\n", err.Error())
			return
		}
	},
}

func init() {
	devicesCmd.AddCommand(lockCmd)

}

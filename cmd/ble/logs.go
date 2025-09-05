/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "Get the activity log for a device",
	PreRunE: mustDeviceId,

	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		flow, err := bleflows.NewAuthenticatedFlow(ble, deviceId)
		if flow == nil {
			c.Logger.Error("Failed to create BLE flow", "error", err.Error())
			return
		}
		defer flow.DisconnectDevice()

		res, err := flow.GetLogs(0, 10)
		if err != nil {
			c.Logger.Error("Failed to read log entries", "error", err.Error())
			return
		}

		t := table.New().Headers("Timestamp", "Action")
		for _, e := range res {
			t = t.Row(e.Time.Local().String(), e.Type.String())
		}
		fmt.Println(t)

	},
}

func init() {
	bleCmd.AddCommand(logsCmd)
}

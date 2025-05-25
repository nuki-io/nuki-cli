/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/internal"
	"github.com/spf13/cobra"
)

// listDevicesCmd represents the listDevices command
var listDevicesCmd = &cobra.Command{
	Use:     "listDevices",
	Short:   "List all devices registered in Nuki Web",
	Aliases: []string{"ls"},
	Run: func(cmd *cobra.Command, args []string) {
		cl := internal.NewWebApiClient(apiKey)
		res, err := cl.GetDevices()
		if err != nil {
			c.Logger.Error("Failed to get account details", "error", err)
			return
		}
		devices := make([][]string, 0, len(res))
		for _, v := range res {
			devices = append(devices, []string{v.Name, fmt.Sprintf("%X", int32(v.SmartlockId)), fmt.Sprintf("%X", v.AuthId)})
		}

		t := table.New().Rows(devices...).Headers("Name", "Device ID", "Auth ID")
		fmt.Println(t)

	},
}

func init() {
	webCmd.AddCommand(listDevicesCmd)
}

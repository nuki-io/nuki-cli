/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss/table"
	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listDevicesCmd represents the listDevices command
var listDevicesCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Lists all known devices",
	Long:    `Lists all devices, either paired locally or registered in Nuki Web.`,
	Run: func(cmd *cobra.Command, args []string) {
		auths := viper.Get("authorizations").(map[string]any)
		devices := make([][]string, 0, len(auths))
		for k, v := range auths {
			values := v.(map[string]any)

			devices = append(devices, []string{
				values["name"].(string),
				k,
				strings.ToUpper(values["appid"].(string)),
				strings.ToUpper(values["authid"].(string)),
				"BLE",
			})
		}
		apiKey := viper.GetString("web.apiKey")
		if apiKey != "" {
			cl := internal.NewWebApiClient(apiKey)
			res, err := cl.GetDevices()
			if err != nil {
				c.Logger.Error("Failed to get account details", "error", err)
				return
			}
			for _, v := range res {
				devices = append(devices, []string{
					v.Name,
					fmt.Sprintf("%X", int32(v.SmartlockId)),
					"",
					fmt.Sprintf("%X", v.AuthId),
					"Web",
				})
			}
		}
		t := table.New().Rows(devices...).Headers("Name", "Device ID", "App ID", "Auth ID", "Connection")
		fmt.Println(t)
	},
}

func init() {
	devicesCmd.AddCommand(listDevicesCmd)
}

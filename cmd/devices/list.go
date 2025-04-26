/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all authorized devices",
	Run: func(cmd *cobra.Command, args []string) {
		auths := viper.Get("authorizations").(map[string]interface{})
		devices := make([][]string, 0, len(auths))
		for k, v := range auths {
			values := v.(map[string]interface{})

			devices = append(devices, []string{values["name"].(string), k, values["appid"].(string), values["authid"].(string)})
		}
		t := table.New().Rows(devices...).Headers("Name", "Device ID", "App ID", "Auth ID")
		fmt.Println(t)
	},
}

func init() {
	devicesCmd.AddCommand(listCmd)
}

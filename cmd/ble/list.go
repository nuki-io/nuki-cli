package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all authorized devices",
	RunE: func(cmd *cobra.Command, args []string) error {
		auths := viper.GetStringMap("authorizations")

		type entry struct {
			Name     string `json:"name"`
			DeviceID string `json:"deviceId"`
			AppID    string `json:"appId"`
			AuthID   string `json:"authId"`
		}

		entries := make([]entry, 0, len(auths))
		for k, v := range auths {
			values, ok := v.(map[string]any)
			if !ok {
				return fmt.Errorf("malformed authorization entry for device %q", k)
			}
			name, _ := values["name"].(string)
			appid, _ := values["appid"].(string)
			authid, _ := values["authid"].(string)
			entries = append(entries, entry{Name: name, DeviceID: k, AppID: appid, AuthID: authid})
		}

		if outputFormat == "json" {
			return printJSON(entries)
		}

		rows := make([][]string, len(entries))
		for i, e := range entries {
			rows[i] = []string{e.Name, e.DeviceID, e.AppID, e.AuthID}
		}
		t := table.New().Rows(rows...).Headers("Name", "Device ID", "App ID", "Auth ID")
		fmt.Println(t)
		return nil
	},
}

func init() {
	bleCmd.AddCommand(listCmd)
}

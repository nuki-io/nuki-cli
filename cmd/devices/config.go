/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/spf13/cobra"
	c "go.nuki.io/nuki/nukictl/cmd"
	"go.nuki.io/nuki/nukictl/pkg/bleflows"
	"go.nuki.io/nuki/nukictl/pkg/nukible"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:     "config",
	Short:   "Retrieves and display the configuration of the device",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		flow := bleflows.NewFlow(ble)

		cfg, err := flow.GetConfig(deviceId)
		if err != nil {
			c.Logger.Error("Failed to read config", "error", err.Error())
			return
		}
		t := table.New().Rows(
			[]string{"Nuki ID", fmt.Sprintf("%X", cfg.NukiID)},
			[]string{"Name", cfg.Name},
			[]string{"Latitude", fmt.Sprintf("%f", cfg.Latitude)},
			[]string{"Longitude", fmt.Sprintf("%f", cfg.Longitude)},
			[]string{"Auto Unlatch", fmt.Sprintf("%t", cfg.AutoUnlatch)},
			[]string{"Pairing enabled", fmt.Sprintf("%t", cfg.PairingEnabled)},
			[]string{"Button enabled", fmt.Sprintf("%t", cfg.ButtonEnabled)},
			[]string{"Led enabled", fmt.Sprintf("%t", cfg.LedEnabled)},
			[]string{"Led Brightness", fmt.Sprintf("%d", cfg.LedBrightness)},
			[]string{"Current Time", cfg.CurrentTime.String()},
			[]string{"Timezone Offset", fmt.Sprintf("%d", cfg.TimezoneOffset)},
			[]string{"DST Mode", fmt.Sprintf("%d", cfg.DstMode)},
			[]string{"Timezone", cfg.GetTimezoneLocation().String()},
			[]string{"Has Fob", fmt.Sprintf("%t", cfg.HasFob)},
			[]string{"Fob Action 1", fmt.Sprintf("%d", cfg.FobAction1)},
			[]string{"Fob Action 2", fmt.Sprintf("%d", cfg.FobAction2)},
			[]string{"Fob Action 3", fmt.Sprintf("%d", cfg.FobAction3)},
			[]string{"Has Keypad", fmt.Sprintf("%t", cfg.HasKeypad)},
			[]string{"Has Keypad2", fmt.Sprintf("%t", cfg.HasKeypad2)},
			[]string{"Single Lock", fmt.Sprintf("%t", cfg.SingleLock)},
			[]string{"Advertising Mode", fmt.Sprintf("%d", cfg.AdvertisingMode)},
			[]string{"Firmware Version", cfg.FirmwareVersion},
			[]string{"Hardware Revision", cfg.HardwareRevision},
			[]string{"HomeKit Status", fmt.Sprintf("%d", cfg.HomeKitStatus)},
			[]string{"Device Type", fmt.Sprintf("%d", cfg.DeviceType)},
			[]string{"Capabilities", fmt.Sprintf("%d", cfg.Capabilities)},
			[]string{"Matter Status", fmt.Sprintf("%d", cfg.MatterStatus)},
		)
		fmt.Println(t)
	},
}

func init() {
	devicesCmd.AddCommand(configCmd)
}

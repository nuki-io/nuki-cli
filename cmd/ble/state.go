/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"runtime"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"

	"github.com/spf13/cobra"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:     "state",
	Short:   "Gets the current lock state of the device",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		if runtime.GOOS == "linux" {
			err = ble.ScanForDevice(deviceId, 10*time.Second)
			if err != nil {
				c.Logger.Error("Failed to scan", "error", err.Error())
				return
			}
		}
		flow, err := bleflows.NewAuthenticatedFlow(ble, deviceId)
		if flow == nil {
			c.Logger.Error("Failed to create BLE flow", "error", err.Error())
			return
		}
		defer flow.DisconnectDevice()

		status, err := flow.GetStatus()
		if err != nil {
			c.Logger.Error("Failed to get status", "error", err.Error())
			return
		}
		style := lipgloss.NewStyle().PaddingLeft(1).PaddingRight(1)
		table := table.New().Headers("Property", "Value").StyleFunc(func(row, col int) lipgloss.Style { return style })
		table.
			Row("Nuki State", status.NukiState.String()).
			Row("LockState", status.LockState.String()).
			Row("Trigger", status.Trigger.String()).
			Row("Current Time", status.CurrentTime.String()).
			Row("Timezone Offset", fmt.Sprintf("%v", status.TimezoneOffset)).
			Row("Battery critical", fmt.Sprintf("%v", status.BatteryStateCritical)).
			Row("Charging", fmt.Sprintf("%v", status.Charging)).
			Row("Battery %", fmt.Sprintf("%d%%", status.BatteryPercentage)).
			Row("Config Update Count", fmt.Sprintf("%v", status.ConfigUpdateCount)).
			Row("Lock'n'Go Timer", fmt.Sprintf("%v", status.LockNGoTimer)).
			Row("Last Lock Action", fmt.Sprintf("%v", status.LastLockAction)).
			Row("Last Lock Action Trigger", fmt.Sprintf("%v", status.LastLockActionTrigger)).
			Row("Last Lock Action Completion Status", fmt.Sprintf("%v", status.LastLockActionCompletionStatus)).
			Row("Door Sensor State", fmt.Sprintf("%v", status.DoorSensorState)).
			Row("Nightmode active", fmt.Sprintf("%v", status.NightmodeActive)).
			Row("Accessory Battery State", fmt.Sprintf("%v", status.AccessoryBatteryState)).
			Row("Remote Access Status", fmt.Sprintf("%v", status.RemoteAccessStatus)).
			Row("BLE Connection Strength", fmt.Sprintf("%v", status.BleConnectionStrength)).
			Row("Wifi Connection Strength", fmt.Sprintf("%v", status.WifiConnectionStrength)).
			Row("Wifi Connection Status", fmt.Sprintf("%v", status.WifiConnectionStatus)).
			Row("Mqtt Connection Status", fmt.Sprintf("%v", status.MqttConnectionStatus)).
			Row("Thread Connection Status", fmt.Sprintf("%v", status.ThreadConnectionStatus))
		fmt.Println(table.Render())
	},
}

func init() {
	bleCmd.AddCommand(stateCmd)
}

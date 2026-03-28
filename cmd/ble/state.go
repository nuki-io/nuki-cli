package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

// stateCmd represents the state command
var stateCmd = &cobra.Command{
	Use:     "state",
	Short:   "Gets the current lock state of the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			status, err := flow.GetStatus(ctx)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(status)
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
			return nil
		})
	},

}

func init() {
	bleCmd.AddCommand(stateCmd)
}

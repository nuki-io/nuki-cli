package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var batteryCmd = &cobra.Command{
	Use:     "battery",
	Short:   "Get the battery report of the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			report, err := flow.GetBatteryReport(ctx)
			if err != nil {
				return fmt.Errorf("failed to get battery report: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(report)
			}
			t := table.New().Rows(
				[]string{"Battery Drain", fmt.Sprintf("%d mWs", report.BatteryDrain)},
				[]string{"Battery Voltage", fmt.Sprintf("%d mV", report.BatteryVoltage)},
				[]string{"Critical State", fmt.Sprintf("%t", report.CriticalState)},
				[]string{"Last Lock Action", report.LockAction.String()},
				[]string{"Start Voltage", fmt.Sprintf("%d mV", report.StartVoltage)},
				[]string{"Lowest Voltage", fmt.Sprintf("%d mV", report.LowestVoltage)},
				[]string{"Lock Distance", fmt.Sprintf("%d°", report.LockDistance)},
				[]string{"Start Temperature", fmt.Sprintf("%d°C", report.StartTemperature)},
				[]string{"Max Turn Current", fmt.Sprintf("%d", report.MaxTurnCurrent)},
				[]string{"Battery Resistance", fmt.Sprintf("%d", report.BatteryResistance)},
			)
			fmt.Println(t)
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(batteryCmd)
}

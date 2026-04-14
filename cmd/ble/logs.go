package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var (
	logsStart int
	logsCount int
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "Get the activity log for a device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			// always get LogEntryCount
			res, count, err := flow.GetLogs(ctx, logsStart, logsCount, true)
			if err != nil {
				return fmt.Errorf("failed to read log entries: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(res)
			}
			t := table.New().StyleFunc(styleLogEntryCount)
			if count != nil {
				t.
					Row("Logging Enabled", boolToIcon(count.LoggingEnabled)).
					Row("Doorsensor Enabled", boolToIcon(count.DoorSensorEnabled)).
					Row("Doorsensor Logging", boolToIcon(count.DoorSensorLoggingEnabled)).
					Row("Total Count", fmt.Sprintf("%d", count.Count))
			}
			t.Row("Start", fmt.Sprintf("%d", logsStart))
			t.Row("Count", fmt.Sprintf("%d", logsCount))
			fmt.Println(t)
			t = table.New().Headers("Index", "Timestamp", "Log")
			for _, e := range res {
				t = t.Row(fmt.Sprintf("%d", e.Index), e.Time.Local().String(), e.String())
			}
			fmt.Println(t)
			return nil
		})
	},
}

func boolToIcon(v bool) string {
	if v {
		return colorGreen("✓")
	}
	return colorRed("✗")
}

func styleLogEntryCount(row, col int) lipgloss.Style {
	if col == 1 {
		return styleCenter
	}
	return emptyStyle
}

func init() {
	bleCmd.AddCommand(logsCmd)
	logsCmd.Flags().IntVarP(&logsStart, "start", "s", 0, "Index where to start reading log entries from")
	logsCmd.Flags().IntVarP(&logsCount, "count", "n", 10, "Number of log entries to read")
}

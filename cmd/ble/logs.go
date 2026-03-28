package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

// logsCmd represents the logs command
var logsCmd = &cobra.Command{
	Use:     "logs",
	Short:   "Get the activity log for a device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			res, err := flow.GetLogs(ctx, 0, 10)
			if err != nil {
				return fmt.Errorf("failed to read log entries: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(res)
			}
			t := table.New().Headers("Timestamp", "Log")
			for _, e := range res {
				t = t.Row(e.Time.Local().String(), e.String())
			}
			fmt.Println(t)
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(logsCmd)
}

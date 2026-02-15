package cmd

import (
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
	Run: func(cmd *cobra.Command, args []string) {
		withAuthenticatedFlow(func(flow *bleflows.Flow) error {
			res, err := flow.GetLogs(0, 10)
			if err != nil {
				return fmt.Errorf("failed to read log entries: %w", err)
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

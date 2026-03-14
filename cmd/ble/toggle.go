package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/blecommands"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var repeats int

// toggleCmd represents the toggle command
var toggleCmd = &cobra.Command{
	Use:     "toggle",
	Short:   "Toggles the current lock state",
	Long:    `Depending on the lock's current state, this command either locks or unlocks.`,
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			status, err := flow.GetStatus(ctx)
			if err != nil {
				return fmt.Errorf("failed to get status: %w", err)
			}
			var startAt int
			if status.LockState == blecommands.LockStateLocked {
				startAt = 0
			} else {
				startAt = 1
			}
			for i := range repeats {
				if i%2 == startAt {
					err = flow.PerformLockOperation(ctx, blecommands.Unlock)
				} else {
					err = flow.PerformLockOperation(ctx, blecommands.Lock)
				}
				// TODO: although we received the StatusComplete at this point, we apparently need to wait a bit longer
				time.Sleep(500 * time.Millisecond)
			}
			return err
		})
	},
}

func init() {
	bleCmd.AddCommand(toggleCmd)
	toggleCmd.Flags().IntVarP(&repeats, "repeats", "n", 1, "The number of times to repeat the lock operation. Default is 1, which means it will toggle once. If set to 2, it will toggle twice, etc. This is useful for testing purposes.")
}

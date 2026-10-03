package cmd

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var verifyPinCmd = &cobra.Command{
	Use:          "verify-pin",
	Short:        "Verify the stored security PIN against the device",
	PreRunE:      mustDeviceId,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.VerifyPIN(ctx); err != nil {
				return fmt.Errorf("PIN verification failed: %w", err)
			}
			fmt.Println("Security PIN verified successfully.")
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(verifyPinCmd)
}

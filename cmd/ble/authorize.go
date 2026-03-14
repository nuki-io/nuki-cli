package cmd

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var pin string

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:   "authorize",
	Short: "Authorizes and pairs this machine with the given Nuki device",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := mustDeviceId(cmd, args); err != nil {
			return err
		}
		if len(pin) != 4 && len(pin) != 6 {
			return fmt.Errorf("--pin must be exactly 4 or 6 digits, got %q", pin)
		}
		for _, c := range pin {
			if c < '0' || c > '9' {
				return fmt.Errorf("--pin must contain only digits, got %q", pin)
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return withUnauthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			return flow.Authorize(ctx, pin)
		})
	},
}

func init() {
	bleCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().StringVarP(&pin, "pin", "p", "", "The PIN code to use for authorization (4 or 6 digits).")
	authorizeCmd.MarkFlagRequired("pin")
}

package cmd

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var pinPattern = regexp.MustCompile(`^\d{4}(\d{2})?$`)

var calibrateCmd = &cobra.Command{
	Use:     "calibrate",
	Short:   "Request calibration of the lock motor",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.Calibrate(ctx); err != nil {
				return fmt.Errorf("calibration failed: %w", err)
			}
			fmt.Println("Calibration completed successfully.")
			return nil
		})
	},
}

var rebootCmd = &cobra.Command{
	Use:     "reboot",
	Short:   "Request a reboot of the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.Reboot(ctx); err != nil {
				return fmt.Errorf("reboot failed: %w", err)
			}
			fmt.Println("Reboot requested successfully.")
			return nil
		})
	},
}

var updateTimeCmd = &cobra.Command{
	Use:   "update-time",
	Short: "Sync the device clock to the current system time",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return mustDeviceId(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			t := time.Now().UTC()
			if err := flow.UpdateTime(ctx, t); err != nil {
				return fmt.Errorf("update time failed: %w", err)
			}
			fmt.Printf("Device time updated to %s.\n", t.Format(time.RFC3339))
			return nil
		})
	},
}

var newPin string

var setPinCmd = &cobra.Command{
	Use:   "set-pin",
	Short: "Change the security PIN of the device",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if err := mustDeviceId(cmd, args); err != nil {
			return err
		}
		if !pinPattern.MatchString(newPin) {
			return fmt.Errorf("--new-pin must be 4 or 6 digits")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.SetSecurityPIN(ctx, newPin); err != nil {
				return fmt.Errorf("set PIN failed: %w", err)
			}
			fmt.Println("Security PIN updated successfully.")
			return nil
		})
	},
}

var enableLoggingCmd = &cobra.Command{
	Use:     "enable-logging",
	Short:   "Enable activity logging on the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.EnableLogging(ctx, true); err != nil {
				return fmt.Errorf("enable logging failed: %w", err)
			}
			fmt.Println("Logging enabled.")
			return nil
		})
	},
}

var disableLoggingCmd = &cobra.Command{
	Use:     "disable-logging",
	Short:   "Disable activity logging on the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			if err := flow.EnableLogging(ctx, false); err != nil {
				return fmt.Errorf("disable logging failed: %w", err)
			}
			fmt.Println("Logging disabled.")
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(calibrateCmd)
	bleCmd.AddCommand(rebootCmd)
	bleCmd.AddCommand(updateTimeCmd)
	bleCmd.AddCommand(setPinCmd)
	bleCmd.AddCommand(enableLoggingCmd)
	bleCmd.AddCommand(disableLoggingCmd)

	setPinCmd.Flags().StringVar(&newPin, "new-pin", "", "New security PIN (4 or 6 digits)")
	_ = setPinCmd.MarkFlagRequired("new-pin")
}

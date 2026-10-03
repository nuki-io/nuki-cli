package cmd

import (
	"context"
	"fmt"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var advancedConfigCmd = &cobra.Command{
	Use:     "advanced-config",
	Short:   "Retrieve and display the advanced configuration of the device",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			cfg, err := flow.GetAdvancedConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to read advanced config: %w", err)
			}
			if outputFormat == "json" {
				return printJSON(cfg)
			}
			t := table.New().Rows(
				[]string{"Total Degrees", fmt.Sprintf("%d", cfg.TotalDegrees)},
				[]string{"Unlocked Position Offset", fmt.Sprintf("%d°", cfg.UnlockedPositionOffsetDegrees)},
				[]string{"Locked Position Offset", fmt.Sprintf("%d°", cfg.LockedPositionOffsetDegrees)},
				[]string{"Single Locked Position Offset", fmt.Sprintf("%d°", cfg.SingleLockedPositionOffsetDegrees)},
				[]string{"Unlocked→Locked Transition Offset", fmt.Sprintf("%d°", cfg.UnlockedToLockedTransitionOffset)},
				[]string{"Lock'n'Go Timeout", fmt.Sprintf("%ds", cfg.LockNGoTimeout)},
				[]string{"Single Button Press Action", fmt.Sprintf("%d", cfg.SingleButtonPressAction)},
				[]string{"Double Button Press Action", fmt.Sprintf("%d", cfg.DoubleButtonPressAction)},
				[]string{"Detached Cylinder", fmt.Sprintf("%t", cfg.DetachedCylinder)},
				[]string{"Battery Type", fmt.Sprintf("%d", cfg.BatteryType)},
				[]string{"Auto Battery Type Detection", fmt.Sprintf("%t", cfg.AutomaticBatteryTypeDetection)},
				[]string{"Unlatch Duration", fmt.Sprintf("%ds", cfg.UnlatchDuration)},
				[]string{"Auto Lock Timeout", fmt.Sprintf("%ds", cfg.AutoLockTimeout)},
				[]string{"Auto Unlock Disabled", fmt.Sprintf("%t", cfg.AutoUnlockDisabled)},
				[]string{"Nightmode Enabled", fmt.Sprintf("%t", cfg.NightmodeEnabled)},
				[]string{"Nightmode Start", fmt.Sprintf("%02d:%02d", cfg.NightmodeStartTime[0], cfg.NightmodeStartTime[1])},
				[]string{"Nightmode End", fmt.Sprintf("%02d:%02d", cfg.NightmodeEndTime[0], cfg.NightmodeEndTime[1])},
				[]string{"Nightmode Auto Lock", fmt.Sprintf("%t", cfg.NightmodeAutoLockEnabled)},
				[]string{"Nightmode Auto Unlock Disabled", fmt.Sprintf("%t", cfg.NightmodeAutoUnlockDisabled)},
				[]string{"Nightmode Immediate Lock On Start", fmt.Sprintf("%t", cfg.NightmodeImmediateLockOnStart)},
				[]string{"Auto Lock Enabled", fmt.Sprintf("%t", cfg.AutoLockEnabled)},
				[]string{"Immediate Auto Lock Enabled", fmt.Sprintf("%t", cfg.ImmediateAutoLockEnabled)},
				[]string{"Auto Update Enabled", fmt.Sprintf("%t", cfg.AutoUpdateEnabled)},
				[]string{"Motor Speed", fmt.Sprintf("%d", cfg.MotorSpeed)},
				[]string{"Slow Speed During Nightmode", fmt.Sprintf("%t", cfg.EnableSlowSpeedDuringNightMode)},
			)
			fmt.Println(t)
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(advancedConfigCmd)
}

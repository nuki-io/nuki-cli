package cmd

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var (
	sacUnlockedOffset   int16
	sacLockedOffset     int16
	sacSingleOffset     int16
	sacTransitionOffset int16
	sacLockNGoTimeout   uint8
	sacSingleAction     uint8
	sacDoubleAction     uint8
	sacDetached         bool
	sacBattType         uint8
	sacAutoBatt         bool
	sacUnlatchDuration  uint8
	sacAutoLockTimeout  uint16
	sacAutoUnlockDis    bool
	sacNightmode        bool
	sacNightStart       string
	sacNightEnd         string
	sacNightAutoLock    bool
	sacNightAutoUnlDis  bool
	sacNightImmLock     bool
	sacAutoLock         bool
	sacImmAutoLock      bool
	sacAutoUpdate       bool
)

var setAdvancedConfigCmd = &cobra.Command{
	Use:     "set-advanced-config",
	Short:   "Update advanced device configuration",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			cur, err := flow.GetAdvancedConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to read current advanced config: %w", err)
			}
			req := cur.ToSetAdvancedConfig()
			if cmd.Flags().Changed("unlocked-offset") {
				req.UnlockedPositionOffsetDegrees = sacUnlockedOffset
			}
			if cmd.Flags().Changed("locked-offset") {
				req.LockedPositionOffsetDegrees = sacLockedOffset
			}
			if cmd.Flags().Changed("single-locked-offset") {
				req.SingleLockedPositionOffsetDegrees = sacSingleOffset
			}
			if cmd.Flags().Changed("transition-offset") {
				req.UnlockedToLockedTransitionOffset = sacTransitionOffset
			}
			if cmd.Flags().Changed("lock-n-go-timeout") {
				req.LockNGoTimeout = sacLockNGoTimeout
			}
			if cmd.Flags().Changed("single-button-action") {
				req.SingleButtonPressAction = sacSingleAction
			}
			if cmd.Flags().Changed("double-button-action") {
				req.DoubleButtonPressAction = sacDoubleAction
			}
			if cmd.Flags().Changed("detached-cylinder") {
				req.DetachedCylinder = sacDetached
			}
			if cmd.Flags().Changed("battery-type") {
				req.BatteryType = sacBattType
			}
			if cmd.Flags().Changed("auto-battery-detection") {
				req.AutomaticBatteryTypeDetection = sacAutoBatt
			}
			if cmd.Flags().Changed("unlatch-duration") {
				req.UnlatchDuration = sacUnlatchDuration
			}
			if cmd.Flags().Changed("auto-lock-timeout") {
				req.AutoLockTimeout = sacAutoLockTimeout
			}
			if cmd.Flags().Changed("auto-unlock-disabled") {
				req.AutoUnlockDisabled = sacAutoUnlockDis
			}
			if cmd.Flags().Changed("nightmode") {
				req.NightmodeEnabled = sacNightmode
			}
			if cmd.Flags().Changed("nightmode-start") {
				h, m, err := parseTime(sacNightStart)
				if err != nil {
					return err
				}
				req.NightmodeStartTime = [2]byte{h, m}
			}
			if cmd.Flags().Changed("nightmode-end") {
				h, m, err := parseTime(sacNightEnd)
				if err != nil {
					return err
				}
				req.NightmodeEndTime = [2]byte{h, m}
			}
			if cmd.Flags().Changed("nightmode-auto-lock") {
				req.NightmodeAutoLockEnabled = sacNightAutoLock
			}
			if cmd.Flags().Changed("nightmode-auto-unlock-disabled") {
				req.NightmodeAutoUnlockDisabled = sacNightAutoUnlDis
			}
			if cmd.Flags().Changed("nightmode-immediate-lock") {
				req.NightmodeImmediateLockOnStart = sacNightImmLock
			}
			if cmd.Flags().Changed("auto-lock") {
				req.AutoLockEnabled = sacAutoLock
			}
			if cmd.Flags().Changed("immediate-auto-lock") {
				req.ImmediateAutoLockEnabled = sacImmAutoLock
			}
			if cmd.Flags().Changed("auto-update") {
				req.AutoUpdateEnabled = sacAutoUpdate
			}
			if err := flow.SetAdvancedConfig(ctx, req); err != nil {
				return fmt.Errorf("failed to set advanced config: %w", err)
			}
			fmt.Println("Advanced configuration updated successfully.")
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(setAdvancedConfigCmd)
	setAdvancedConfigCmd.Flags().Int16Var(&sacUnlockedOffset, "unlocked-offset", 0, "Unlocked position offset in degrees")
	setAdvancedConfigCmd.Flags().Int16Var(&sacLockedOffset, "locked-offset", 0, "Locked position offset in degrees")
	setAdvancedConfigCmd.Flags().Int16Var(&sacSingleOffset, "single-locked-offset", 0, "Single locked position offset in degrees")
	setAdvancedConfigCmd.Flags().Int16Var(&sacTransitionOffset, "transition-offset", 0, "Unlocked-to-locked transition offset in degrees")
	setAdvancedConfigCmd.Flags().Uint8Var(&sacLockNGoTimeout, "lock-n-go-timeout", 20, "Lock'n'Go timeout in seconds")
	setAdvancedConfigCmd.Flags().Uint8Var(&sacSingleAction, "single-button-action", 0, "Single button press action")
	setAdvancedConfigCmd.Flags().Uint8Var(&sacDoubleAction, "double-button-action", 0, "Double button press action")
	setAdvancedConfigCmd.Flags().BoolVar(&sacDetached, "detached-cylinder", false, "Detached cylinder mode")
	setAdvancedConfigCmd.Flags().Uint8Var(&sacBattType, "battery-type", 0, "Battery type (0=alkaline, 1=accumulator, 2=lithium)")
	setAdvancedConfigCmd.Flags().BoolVar(&sacAutoBatt, "auto-battery-detection", true, "Automatic battery type detection")
	setAdvancedConfigCmd.Flags().Uint8Var(&sacUnlatchDuration, "unlatch-duration", 3, "Unlatch duration in seconds")
	setAdvancedConfigCmd.Flags().Uint16Var(&sacAutoLockTimeout, "auto-lock-timeout", 0, "Auto lock timeout in seconds (0=disabled)")
	setAdvancedConfigCmd.Flags().BoolVar(&sacAutoUnlockDis, "auto-unlock-disabled", false, "Disable auto unlock")
	setAdvancedConfigCmd.Flags().BoolVar(&sacNightmode, "nightmode", false, "Enable nightmode")
	setAdvancedConfigCmd.Flags().StringVar(&sacNightStart, "nightmode-start", "22:00", "Nightmode start time (HH:MM)")
	setAdvancedConfigCmd.Flags().StringVar(&sacNightEnd, "nightmode-end", "07:00", "Nightmode end time (HH:MM)")
	setAdvancedConfigCmd.Flags().BoolVar(&sacNightAutoLock, "nightmode-auto-lock", false, "Auto lock at nightmode start")
	setAdvancedConfigCmd.Flags().BoolVar(&sacNightAutoUnlDis, "nightmode-auto-unlock-disabled", false, "Disable auto unlock during nightmode")
	setAdvancedConfigCmd.Flags().BoolVar(&sacNightImmLock, "nightmode-immediate-lock", false, "Immediate lock when nightmode starts")
	setAdvancedConfigCmd.Flags().BoolVar(&sacAutoLock, "auto-lock", false, "Enable auto lock")
	setAdvancedConfigCmd.Flags().BoolVar(&sacImmAutoLock, "immediate-auto-lock", false, "Immediate auto lock after unlock")
	setAdvancedConfigCmd.Flags().BoolVar(&sacAutoUpdate, "auto-update", true, "Enable automatic firmware updates")
}

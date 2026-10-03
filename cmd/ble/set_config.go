package cmd

import (
	"context"
	"fmt"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var (
	scName           string
	scLat            float32
	scLon            float32
	scAutoUnlatch    bool
	scPairingEnabled bool
	scButtonEnabled  bool
	scLedEnabled     bool
	scLedBrightness  uint8
	scTzOffset       int16
	scDstMode        uint8
	scFobAction1     uint8
	scFobAction2     uint8
	scFobAction3     uint8
	scSingleLock     bool
	scAdvMode        uint8
	scTzID           uint16
)

var setConfigCmd = &cobra.Command{
	Use:     "set-config",
	Short:   "Update device configuration",
	PreRunE: mustDeviceId,
	RunE: func(cmd *cobra.Command, args []string) error {
		return withAuthenticatedFlow(func(ctx context.Context, flow *bleflows.Flow) error {
			cur, err := flow.GetConfig(ctx)
			if err != nil {
				return fmt.Errorf("failed to read current config: %w", err)
			}
			req := cur.ToSetConfig()
			if cmd.Flags().Changed("name") {
				req.Name = scName
			}
			if cmd.Flags().Changed("lat") {
				req.Latitude = scLat
			}
			if cmd.Flags().Changed("lon") {
				req.Longitude = scLon
			}
			if cmd.Flags().Changed("auto-unlatch") {
				req.AutoUnlatch = scAutoUnlatch
			}
			if cmd.Flags().Changed("pairing-enabled") {
				req.PairingEnabled = scPairingEnabled
			}
			if cmd.Flags().Changed("button-enabled") {
				req.ButtonEnabled = scButtonEnabled
			}
			if cmd.Flags().Changed("led-enabled") {
				req.LedEnabled = scLedEnabled
			}
			if cmd.Flags().Changed("led-brightness") {
				req.LedBrightness = scLedBrightness
			}
			if cmd.Flags().Changed("timezone-offset") {
				req.TimezoneOffset = scTzOffset
			}
			if cmd.Flags().Changed("dst-mode") {
				req.DstMode = scDstMode
			}
			if cmd.Flags().Changed("fob-action-1") {
				req.FobAction1 = scFobAction1
			}
			if cmd.Flags().Changed("fob-action-2") {
				req.FobAction2 = scFobAction2
			}
			if cmd.Flags().Changed("fob-action-3") {
				req.FobAction3 = scFobAction3
			}
			if cmd.Flags().Changed("single-lock") {
				req.SingleLock = scSingleLock
			}
			if cmd.Flags().Changed("advertising-mode") {
				req.AdvertisingMode = scAdvMode
			}
			if cmd.Flags().Changed("timezone-id") {
				req.TimezoneID = scTzID
			}
			if err := flow.SetConfig(ctx, req); err != nil {
				return fmt.Errorf("failed to set config: %w", err)
			}
			fmt.Println("Configuration updated successfully.")
			return nil
		})
	},
}

func init() {
	bleCmd.AddCommand(setConfigCmd)
	setConfigCmd.Flags().StringVar(&scName, "name", "", "Device name (max 32 chars)")
	setConfigCmd.Flags().Float32Var(&scLat, "lat", 0, "Latitude")
	setConfigCmd.Flags().Float32Var(&scLon, "lon", 0, "Longitude")
	setConfigCmd.Flags().BoolVar(&scAutoUnlatch, "auto-unlatch", false, "Auto unlatch on unlock")
	setConfigCmd.Flags().BoolVar(&scPairingEnabled, "pairing-enabled", true, "Allow pairing new devices")
	setConfigCmd.Flags().BoolVar(&scButtonEnabled, "button-enabled", true, "Enable the button on the device")
	setConfigCmd.Flags().BoolVar(&scLedEnabled, "led-enabled", true, "Enable the LED")
	setConfigCmd.Flags().Uint8Var(&scLedBrightness, "led-brightness", 3, "LED brightness (0-5)")
	setConfigCmd.Flags().Int16Var(&scTzOffset, "timezone-offset", 0, "Timezone offset in minutes")
	setConfigCmd.Flags().Uint8Var(&scDstMode, "dst-mode", 0, "Daylight saving time mode (0=off, 1=on)")
	setConfigCmd.Flags().Uint8Var(&scFobAction1, "fob-action-1", 0, "Fob button action 1")
	setConfigCmd.Flags().Uint8Var(&scFobAction2, "fob-action-2", 0, "Fob button action 2")
	setConfigCmd.Flags().Uint8Var(&scFobAction3, "fob-action-3", 0, "Fob button action 3")
	setConfigCmd.Flags().BoolVar(&scSingleLock, "single-lock", false, "Enable single lock mode")
	setConfigCmd.Flags().Uint8Var(&scAdvMode, "advertising-mode", 0, "Advertising mode")
	setConfigCmd.Flags().Uint16Var(&scTzID, "timezone-id", 0, "Timezone ID (see timezone list)")
}

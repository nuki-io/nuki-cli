package cmd

import (
	"context"
	"fmt"
	"runtime"
	"time"

	"github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bleTimeout is the maximum time allowed for a BLE command exchange after
// the device connection has been established.
const bleTimeout = 30 * time.Second

var (
	deviceId string
)

// bleCmd represents the bleCmd command
var bleCmd = &cobra.Command{
	Use:              "ble",
	Short:            "Command to interact with devices through BLE",
	Long:             ``,
	PersistentPreRun: preRun,
}

func init() {
	cmd.RootCmd.AddCommand(bleCmd)
	bleCmd.PersistentFlags().StringVarP(&deviceId, "device-id", "d", "", "The device to use. If not set, the device set by set-context command is used. This is ignored for some commands.")
	// viper.BindPFlag("activeContext", bleCmd.PersistentFlags().Lookup("device-id"))
}

func preRun(cmd *cobra.Command, args []string) {
	if deviceId == "" && viper.IsSet("activecontext") {
		deviceId = viper.GetString("activecontext")
	}
	// TODO: The following "should" work. Check why it doesn't.
	// viper.BindPFlag("activeContext", cmd.PersistentFlags().Lookup("device-id"))
}

func mustDeviceId(cmd *cobra.Command, args []string) error {
	if deviceId == "" {
		return fmt.Errorf("either --device-id flag must be set or a device ID must set with set-context")
	}
	return nil
}

// withAuthenticatedFlow creates a BLE adapter, establishes an authenticated flow,
// and calls fn with a timeout-bounded context. The device is disconnected after fn returns.
// Errors from setup or fn are logged.
func withAuthenticatedFlow(fn func(ctx context.Context, flow *bleflows.Flow) error) {
	ble, err := nukible.NewNukiBle()
	if err != nil {
		cmd.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
		return
	}
	if runtime.GOOS == "linux" {
		err = ble.ScanForDevice(deviceId, 10*time.Second)
		if err != nil {
			cmd.Logger.Error("Failed to scan for device", "error", err.Error())
			return
		}
	}
	flow, err := bleflows.NewAuthenticatedFlow(ble, deviceId)
	if err != nil {
		cmd.Logger.Error("Failed to create BLE flow", "error", err.Error())
		return
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), bleTimeout)
	defer cancel()
	if err := fn(ctx, flow); err != nil {
		cmd.Logger.Error("Command failed", "error", err.Error())
	}
}

// withUnauthenticatedFlow creates a BLE adapter, scans for the device (since it is not yet known),
// establishes an unauthenticated flow, and calls fn with a timeout-bounded context.
// The flow is disconnected after fn returns.
// Errors from setup or fn are logged.
func withUnauthenticatedFlow(fn func(ctx context.Context, flow *bleflows.Flow) error) {
	ble, err := nukible.NewNukiBle()
	if err != nil {
		cmd.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
		return
	}
	err = ble.ScanForDevice(deviceId, 10*time.Second)
	if err != nil {
		cmd.Logger.Error("Failed to scan for device", "error", err.Error())
		return
	}
	flow, err := bleflows.NewUnauthenticatedFlow(ble, deviceId)
	if err != nil {
		cmd.Logger.Error("Failed to create BLE flow", "error", err.Error())
		return
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), bleTimeout)
	defer cancel()
	if err := fn(ctx, flow); err != nil {
		cmd.Logger.Error("Command failed", "error", err.Error())
	}
}

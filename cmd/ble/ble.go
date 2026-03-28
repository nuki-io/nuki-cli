package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"time"

	parentcmd "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bleTimeout is the maximum time allowed for a BLE command exchange after
// the device connection has been established.
const bleTimeout = 30 * time.Second

var (
	deviceId     string
	outputFormat string
)

// bleCmd represents the bleCmd command
var bleCmd = &cobra.Command{
	Use:              "ble",
	Short:            "Command to interact with devices through BLE",
	Long:             ``,
	PersistentPreRun: preRun,
	SilenceUsage:     true,
}

func init() {
	parentcmd.RootCmd.AddCommand(bleCmd)
	bleCmd.PersistentFlags().StringVarP(&deviceId, "device-id", "d", "", "The device to use. If not set, the device set by set-context command is used. This is ignored for some commands.")
	bleCmd.PersistentFlags().StringVar(&outputFormat, "format", "table", "Output format: table or json")
	// viper.BindPFlag("activeContext", bleCmd.PersistentFlags().Lookup("device-id"))
}

// printJSON writes v as indented JSON to stdout.
func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
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
func withAuthenticatedFlow(fn func(ctx context.Context, flow *bleflows.Flow) error) error {
	ble, err := nukible.NewNukiBle()
	if err != nil {
		return fmt.Errorf("failed to enable bluetooth: %w", err)
	}
	if runtime.GOOS == "linux" {
		if err = ble.ScanForDevice(deviceId, 10*time.Second); err != nil {
			return fmt.Errorf("failed to scan for device: %w", err)
		}
	}
	flow, err := bleflows.NewAuthenticatedFlow(ble, deviceId, viperAuthStore{})
	if err != nil {
		return fmt.Errorf("failed to create BLE flow: %w", err)
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), bleTimeout)
	defer cancel()
	return fn(ctx, flow)
}

// withUnauthenticatedFlow creates a BLE adapter, scans for the device (since it is not yet known),
// establishes an unauthenticated flow, and calls fn with a timeout-bounded context.
// The flow is disconnected after fn returns.
func withUnauthenticatedFlow(fn func(ctx context.Context, flow *bleflows.Flow) error) error {
	ble, err := nukible.NewNukiBle()
	if err != nil {
		return fmt.Errorf("failed to enable bluetooth: %w", err)
	}
	if err = ble.ScanForDevice(deviceId, 10*time.Second); err != nil {
		return fmt.Errorf("failed to scan for device: %w", err)
	}
	flow, err := bleflows.NewUnauthenticatedFlow(ble, deviceId, viperAuthStore{})
	if err != nil {
		return fmt.Errorf("failed to create BLE flow: %w", err)
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), bleTimeout)
	defer cancel()
	return fn(ctx, flow)
}

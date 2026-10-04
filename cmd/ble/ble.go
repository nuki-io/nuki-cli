package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/lipgloss"
	parentcmd "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/internal/authstore"
	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// bleTimeout is the maximum time allowed for a BLE command exchange after
// the device connection has been established.
const bleTimeout = 30 * time.Second

var (
	deviceId     string
	outputFormat string

	emptyStyle  = lipgloss.NewStyle()
	styleCenter = lipgloss.NewStyle().AlignHorizontal(lipgloss.Center)
	colorRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render
	colorGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render
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

// withAuthenticatedFlow establishes an authenticated flow and calls fn with a timeout-bounded context. The device is disconnected after fn returns.
func withAuthenticatedFlow(fn func(ctx context.Context, flow *bleflows.Flow) error) error {
	return withAuthenticatedFlowTimeout(bleTimeout, fn)
}

// withAuthenticatedFlowTimeout is withAuthenticatedFlow for commands that may need longer than bleTimeout.
func withAuthenticatedFlowTimeout(timeout time.Duration, fn func(ctx context.Context, flow *bleflows.Flow) error) error {
	flow, err := bleflows.ConnectAuthenticated(deviceId, authstore.New(viper.GetViper()))
	if err != nil {
		return err
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return fn(ctx, flow)
}

// withUnauthenticatedFlow establishes an unauthenticated flow for pairing and calls fn
// with a timeout-bounded context. The flow is disconnected after fn returns.
func withUnauthenticatedFlow(fn func(ctx context.Context, flow *bleflows.Flow) error) error {
	flow, err := bleflows.ConnectUnauthenticated(deviceId, authstore.New(viper.GetViper()))
	if err != nil {
		return err
	}
	defer flow.DisconnectDevice()
	ctx, cancel := context.WithTimeout(context.Background(), bleTimeout)
	defer cancel()
	return fn(ctx, flow)
}

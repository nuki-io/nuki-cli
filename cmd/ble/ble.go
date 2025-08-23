package cmd

import (
	"fmt"

	"github.com/nuki-io/nuki-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

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

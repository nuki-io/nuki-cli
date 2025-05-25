package cmd

import (
	"time"

	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/pkg/nukible"
	"github.com/spf13/cobra"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for Nuki devices using Bluetooth",
	Long:  `Scans for Nuki devices using your local default Bluetooth adapter.`,
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			c.Logger.Error("Failed to enable bluetooth device", "error", err.Error())
			return
		}
		err = ble.Scan(10 * time.Second)
		if err != nil {
			c.Logger.Error("Failed to scan", "error", err.Error())
			return
		}
	},
}

func init() {
	devicesCmd.AddCommand(scanCmd)
}

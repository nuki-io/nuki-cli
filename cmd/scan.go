package cmd

import (
	"log"
	"time"

	"github.com/spf13/cobra"
	"go.nuki.io/nuki/nukictl/pkg/nukible"
)

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for Nuki devices using Bluetooth",
	Long:  `Scans for Nuki devices using your local default Bluetooth adapter.`,
	Run: func(cmd *cobra.Command, args []string) {
		ble, err := nukible.NewNukiBle()
		if err != nil {
			log.Printf("Failed to start scan. %s\n", err.Error())
			return
		}
		ble.Scan(10 * time.Second)
	},
}

func init() {
	devicesCmd.AddCommand(scanCmd)
}

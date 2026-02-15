package cmd

import (
	"strconv"

	"github.com/nuki-io/nuki-cli/pkg/bleflows"
	"github.com/spf13/cobra"
)

var pin int

// authorizeCmd represents the authorize command
var authorizeCmd = &cobra.Command{
	Use:     "authorize",
	Short:   "Authorizes and pairs this machine with the given Nuki device",
	PreRunE: mustDeviceId,
	Run: func(cmd *cobra.Command, args []string) {
		withUnauthenticatedFlow(func(flow *bleflows.Flow) error {
			return flow.Authorize(strconv.Itoa(pin))
		})
	},
}

func init() {
	bleCmd.AddCommand(authorizeCmd)
	authorizeCmd.Flags().IntVarP(&pin, "pin", "p", 0, "The PIN code to use for authorization.")
}

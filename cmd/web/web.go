package cmd

import (
	"fmt"

	"github.com/nuki-io/nuki-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	apiKey string
)

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:               "web",
	Short:             "Command to interact with devices and resources of the Nuki Web API",
	PersistentPreRunE: mustApiKey,
}

func init() {
	cmd.RootCmd.AddCommand(webCmd)
	webCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "The API key to use. If not set, the one configured through web login command is used.")
}

func mustApiKey(cmd *cobra.Command, args []string) error {
	if apiKey == "" && viper.IsSet("web.apiKey") {
		apiKey = viper.GetString("web.apiKey")
	}
	if apiKey == "" {
		return fmt.Errorf("either --api-key flag must be set or an API key must set with login")
	}
	return nil

}

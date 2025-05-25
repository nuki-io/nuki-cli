package cmd

import (
	c "github.com/nuki-io/nuki-cli/cmd"
	"github.com/nuki-io/nuki-cli/internal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to the Nuki Web API",
	Long: `Login with the given API key and persist it in the configuration.
The key is then used to authenticate all subsequent calls to the Nuki Web API.`,
	Run: func(cmd *cobra.Command, args []string) {
		viper.Set("web.apiKey", apiKey)

		cl := internal.NewWebApiClient(apiKey)
		res, err := cl.GetMyAccount()
		if err != nil {
			c.Logger.Error("Failed to get account details", "error", err)
			return
		}

		err = viper.WriteConfig()
		if err != nil {
			c.Logger.Error("Failed to write config file", "error", err)
			return
		}
		c.Logger.Info("API key stored in config file", "email", res.Email)
	},
}

func init() {
	webCmd.AddCommand(loginCmd)
}

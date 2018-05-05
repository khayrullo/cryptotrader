package cmd

import (
	"github.com/spf13/cobra"
	"github.com/crankykernel/cryptotrader/binance"
	"github.com/spf13/viper"
	"github.com/crankykernel/cryptotrader/cmd/common"
)

var binancePostCmd = &cobra.Command{
	Use: "post",
	Run: func(cmd *cobra.Command, args []string) {
		var clientConfig *binance.RestClientAuth = nil
		auth, _ := cmd.Flags().GetBool("auth")
		if auth {
			clientConfig = &binance.RestClientAuth{
				ApiKey:    viper.GetString("binance.api.key"),
				ApiSecret: viper.GetString("binance.api.secret"),
			}
		}
		common.Post(binance.NewClient(clientConfig), args)
	},
}

func init() {
	flags := binancePostCmd.Flags()

	flags.Bool("auth", false, "Send authenticated request")

	binanceCmd.AddCommand(binancePostCmd)
}

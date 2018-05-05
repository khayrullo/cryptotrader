package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var binanceCmd = &cobra.Command{
	Use:   "binance",
	Short: "Binance tools",
}

func init() {
	flags := binanceCmd.PersistentFlags()

	flags.String("api-key", "", "Binance API key")
	viper.BindPFlag("binance.api.key", flags.Lookup("api-key"))
	viper.BindEnv("binance.api.key", "BINANCE_API_KEY")

	flags.String("api-secret", "", "Binance API secret")
	viper.BindPFlag("binance.api.secret", flags.Lookup("api-secret"))
	viper.BindEnv("binance.api.secret", "BINANCE_API_SECRET")

	rootCmd.AddCommand(binanceCmd)
}

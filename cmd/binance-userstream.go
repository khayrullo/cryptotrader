package cmd

import (
	"github.com/spf13/cobra"
	"cryptotrader/binance"
)

var binanceUserStreamCmd = &cobra.Command{
	Use: "user-stream",
	Run: func(cmd *cobra.Command, args []string) {
		binance.BinanceUserStreamCommand()
	},
}

func init() {
	binanceCmd.AddCommand(binanceUserStreamCmd)
}

package cmd

import (
	"github.com/khayrullo/cryptotrader/cmd/binance"
	"github.com/spf13/cobra"
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

package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/crankykernel/cryptotrader/cmd/binance"
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

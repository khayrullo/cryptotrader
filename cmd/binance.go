package cmd

import (
	"github.com/spf13/cobra"
)

var binanceCmd = &cobra.Command{
	Use:   "binance",
	Short: "Binance tools",
}

func init() {
	rootCmd.AddCommand(binanceCmd)
}

// The MIT License (MIT)
//
// Copyright (c) 2018 Cranky Kernel
//
// Permission is hereby granted, free of charge, to any person
// obtaining a copy of this software and associated documentation
// files (the "Software"), to deal in the Software without
// restriction, including without limitation the rights to use, copy,
// modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS
// BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN
// ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cmd

import (
	"github.com/spf13/cobra"
	"gitlab.com/crankykernel/cryptotrader/cmd/kucoin"
)

var kucoinTradesCmd = &cobra.Command{
	Use: "trades",
	Short: "Get trades",
	Long: `Print trades in descending order.

This command will print your KuCoin trades with the most recent being
being displayed first.

By default the 20 most recent trades will be printed. To print all trades
use the --all flag (or set --limit to 0).

Available output formats:
  - tab
  - csv
  - raw
  - default
`,
	Run: func(cmd *cobra.Command, args []string) {
		kucoin.GetTrades()
	},
}

func init() {
	kucoinCmd.AddCommand(kucoinTradesCmd)

	flags := kucoinTradesCmd.Flags()
	flags.StringVar(&kucoin.GetTradesFlags.Format, "format", "",
		"Display format (raw, tab, csv, ...)")
	flags.UintVar(&kucoin.GetTradesFlags.Limit, "limit", 20,
		"Limit the result to a number of trades")
	flags.BoolVar(&kucoin.GetTradesFlags.All, "all", false, "Print all trades")
}

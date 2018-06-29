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
	"gitlab.com/crankykernel/cryptotrader/cmd/tickerlogger"
)

var tickerLoggerCmd = &cobra.Command{
	Use: "ticker-logger",
	Long: `Log tickers across exchanges

By default the tickers will be updated every minute on the minute. This can be
disabled by setting a non-0 interval value.

Example:

  cryptotrader ticker-logging kraken:xbtusd,xbtcad binance:bnbusdt binance:btcusdt
`,
	Run: func(cmd *cobra.Command, args []string) {
		tickerlogger.TickerLoggerCommand(args)
	},
}

func init() {
	rootCmd.AddCommand(tickerLoggerCmd)

	flags := tickerLoggerCmd.Flags()
	flags.Int64Var(&tickerlogger.Flags.Interval, "interval", 0,
		"Seconds between updates")
	flags.BoolVar(&tickerlogger.Flags.OnTheMinute, "on-minute", true,
		"Update on the minute")
}

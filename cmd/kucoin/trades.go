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

package kucoin

import (
	"log"
	"fmt"
	"gitlab.com/crankykernel/cryptotrader/kucoin"
	"encoding/json"
	"strings"
)

var GetTradesFlags struct {
	Format string
	Limit  uint
	All    bool
}

func GetTrades() {
	page := 0
	count := 0

	if GetTradesFlags.All {
		// For limit to 0.
		GetTradesFlags.Limit = 0
	}

Loop:
	for {
		// Grab in batches of 100, even though the limit may be less.
		response, err := getClient().GetDealtOrders(100, page)
		if err != nil {
			log.Fatal(err)
		}
		if !response.Success {
			log.Fatal("error: ", response.Raw)
		}

		for _, trade := range response.Data.Trades {
			switch GetTradesFlags.Format {
			case "raw":
				renderRaw(trade)
			case "csv":
				renderDelim(count, trade, ",")
			case "tab":
				renderDelim(count, trade, "\t")
			case "default":
				fallthrough
			default:
				renderDefault(trade)
			}

			count += 1

			if GetTradesFlags.Limit > 0 && count == int(GetTradesFlags.Limit) {
				break Loop
			}
		}

		if int64(count) == response.Data.Total {
			break
		}

		page += 1
	}
}

func renderDelim(i int, trade *kucoin.Trade, delim string) {
	if i == 0 {
		header := []string{
			"timestamp",
			"type",
			"pair",
			"cost",
			"fee",
			"volume",
		}
		fmt.Printf("%s\n", strings.Join(header, delim))
	}
	parts := []string{
		trade.Timestamp.Format("2006-01-02 15:04:05"),
		trade.Direction,
		fmt.Sprintf("%s/%s", trade.CoinType, trade.CoinTypePair),
		fmt.Sprintf("%.8f", trade.DealValue),
		fmt.Sprintf("%.8f", trade.Fee),
		fmt.Sprintf("%.8f", trade.Amount),
	}
	fmt.Printf("%s\n", strings.Join(parts, delim))
}

func renderDefault(trade *kucoin.Trade) {
	var feeCoin string
	if trade.Direction == "BUY" {
		feeCoin = trade.CoinType
	} else {
		feeCoin = trade.CoinTypePair
	}

	fmt.Printf(
		"Timestamp: %s; "+
			"Action: %-4s; "+
			"Pair: %s/%s; "+
			"Amount: %.8f %s; "+
			"Cost: %.8f %s; "+
			"Fee: %.8f %s; "+
			"\n",
		trade.Timestamp.Format("2006-01-02 15:04:05"),
		strings.Title(strings.ToLower(trade.Direction)),
		trade.CoinType,
		trade.CoinTypePair,
		trade.Amount, trade.CoinType,
		trade.DealValue, trade.CoinTypePair,
		trade.Fee, feeCoin)
}

func renderRaw(trade *kucoin.Trade) {
	buf, err := json.Marshal(trade)
	if err != nil {
		log.Fatal("error: failed to render trade: %v", err)
	}
	fmt.Printf("%s\n", buf)
}

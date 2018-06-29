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
	"gitlab.com/crankykernel/cryptotrader/util"
	"fmt"
)

func Transfers(args []string) {
	coins := args

	client := getClient()

	if len(coins) == 0 {
		ticks, err := client.GetTick()
		if err != nil {
			log.Fatal("error: failed to get coins: ", err)
		}
		coinTypes := map[string]bool{}
		for _, tick := range ticks.Entries {
			coinTypes[tick.CoinType] = true
		}
		for coinType, _ := range coinTypes {
			coins = append(coins, coinType)
		}
	}

	for _, coin := range coins {
		page := 1
		for {
			response, err := client.WalletRecords(coin, page)
			if err != nil {
				log.Fatal("error: ", err)
			}
			if len(response.Data.Entries) == 0 {
				break;
			}

			for _, entry := range response.Data.Entries {
				timestamp := util.MillisToTime(entry.CreatedAtMillis)
				fmt.Printf("Timestamp: %s, "+
					"Coin: %s, Type: %s, "+
					"Status: %s, "+
					"Amount: %f, Fee: %f\n",
					timestamp.Format("2006-01-02 15:04:05"),
					entry.CoinType,
					entry.Status,
					entry.Type,
					entry.Amount,
					entry.Fee,
				)
			}

			page += 1
		}
	}
}

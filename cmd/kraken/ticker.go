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

package kraken

import (
	"log"
	"gitlab.com/crankykernel/cryptotrader/kraken"
	"github.com/spf13/viper"
	"fmt"
	"encoding/json"
	"time"
)

const INTERVAL = 3

var TickerFlags struct {
	Once     bool
	Interval int64
}

func Ticker(args []string) {

	if len(args) == 0 {
		log.Fatal("error: no pairs provided")
	}

	client := kraken.NewClient(viper.GetString("kraken.api.key"),
		viper.GetString("kraken.api.secret"))

	for {
		ticker, err := client.Ticker(args...)
		if err != nil {
			log.Fatal("error: ", err)
		}
		renderTicker(ticker)

		if TickerFlags.Once {
			break
		}

		time.Sleep(time.Duration(TickerFlags.Interval) * time.Second)
	}
}

func renderTicker(ticker map[string]kraken.Ticker) {
	txt, _ := json.Marshal(ticker)
	fmt.Printf("%s\n", txt)
}

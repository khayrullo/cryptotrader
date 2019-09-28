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

package ticker

import (
	"encoding/json"
	"fmt"
	"github.com/khayrullo/cryptotrader/gdax"
	"log"
	"time"
)

func Main(args []string) {
	products := args
	if len(products) == 0 {
		log.Println("No products provided, will default to BTC-USD")
		products = append(products, "BTC-USD")
	}
	client := gdax.NewFeedClient()
	for {
		if err := client.Connect(); err != nil {
			log.Printf("error: failed to connect: ", err)
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}
	if err := client.Subscribe(gdax.TickerChannel(products)); err != nil {
		log.Fatal("error: ", err)
	}
	for {
		message, err := client.Next()
		if err != nil {
			log.Fatal("error: ", err)
		}
		printable, _ := json.Marshal(message)
		fmt.Println(string(printable))
	}
}

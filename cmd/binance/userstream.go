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

package binance

import (
	"github.com/spf13/viper"
	"log"
	"gitlab.com/crankykernel/cryptotrader/binance"
	"encoding/json"
)

func BinanceUserStreamCommand() {
	apiKey := viper.GetString("binance.api.key")

	if apiKey == "" {
		log.Fatal("error: this command requires an api key")
	}

	restClient := binance.NewAuthenticatedClient(apiKey, "")

	streamClient, err := binance.OpenUserStream(restClient)
	if err != nil {
		log.Fatal("error: failed to open user stream: %v", err)
	}

	for {
		_, body, err := streamClient.Next()
		if err != nil {
			log.Fatal("error: failed to read next message: %v", err)
		}

		var rawOrderUpdate binance.StreamExecutionReport
		err = json.Unmarshal(body, &rawOrderUpdate)
		if err == nil {
			x := binance.StreamExecutionReport(rawOrderUpdate)
			log.Printf("%+v\n", rawOrderUpdate)
			log.Printf("%+v\n", x)
		}
	}
}

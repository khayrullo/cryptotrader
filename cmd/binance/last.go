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
	"gitlab.com/crankykernel/cryptotrader/binance"
	"log"
	"fmt"
	"strings"
	"encoding/json"
)

func LastCommand(args []string) {
	client := binance.NewAnonymousClient()
	response, err := client.GetAllPriceTicker()
	if err != nil {
		log.Fatal("error: ", err)
	}
	for _, last := range response {
		if matchArgs(last.Symbol, args) {
			buf, _ := json.Marshal(last)
			fmt.Println(string(buf))
		}
	}
}

func matchArgs(symbol string, args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, arg := range args {
		if strings.Index(symbol, arg) > -1 {
			return true
		}
	}
	return false
}

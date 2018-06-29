// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tickerlogger

import (
	"gitlab.com/crankykernel/cryptotrader/binance"
	"gitlab.com/crankykernel/cryptotrader/kraken"
	"strings"
	"log"
	"time"
	"fmt"
	"sync"
	"encoding/json"
)

var Flags struct {
	Interval    int64
	OnTheMinute bool
}

var Binance struct {
	Client  *binance.RestClient
	Symbols []string
}

var Kraken struct {
	Client  *kraken.Client
	Symbols []string
}

type NormalizedTicker struct {
	Timestamp time.Time
	Exchange  string
	Symbol    string
	Price     float64
}

func TickerLoggerCommand(args []string) {

	for _, arg := range args {
		parts := strings.SplitN(arg, ":", 2)
		if len(parts) < 2 {
			log.Fatalf("error: invalid symbols %s (format: exchange:symbols)",
				arg)
		}
		exchange := parts[0]
		symbols := parts[1]
		switch strings.ToUpper(exchange) {
		case "KRAKEN":
			if Kraken.Client == nil {
				Kraken.Client = kraken.NewClient("", "")
			}
			for _, symbol := range strings.Split(symbols, ",") {
				Kraken.Symbols = append(Kraken.Symbols, strings.ToUpper(symbol))
			}
		case "BINANCE":
			if Binance.Client == nil {
				Binance.Client = binance.NewAnonymousClient()
			}
			for _, symbol := range strings.Split(symbols, ",") {
				Binance.Symbols = append(Binance.Symbols, strings.ToUpper(symbol))
			}
		default:
			log.Fatalf("error: exchange not supported: %s", exchange)
		}
	}

	wg := sync.WaitGroup{}
	logChannel := make(chan NormalizedTicker)

	// Log listener.
	go func() {
		for {
			msg := <-logChannel
			buf, _ := json.Marshal(msg)
			fmt.Println(string(buf))
		}
	}()

	// Kraken loop.
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			sleep()
			now := time.Now()
			tickers, err := Kraken.Client.Ticker(Kraken.Symbols...)
			if err != nil {
				log.Println("kraken error: ", err)
			} else {
				for _, tick := range tickers {
					normalizedTicker := NormalizedTicker{
						Timestamp: now,
						Exchange:  "Kraken",
						Symbol:    tick.Pair,
						Price:     tick.Last,
					}
					logChannel <- normalizedTicker
				}
			}
		}
	}()

	// Binance loop.
	wg.Add(1)
	go func() {
		defer wg.Done()
		if len(Binance.Symbols) < 1 {
			return
		}
		for {
			sleep()
			now := time.Now()
			tickers, err := Binance.Client.GetAllPriceTicker()
			if err != nil {
				log.Println("binance error: ", err)
			} else {
				for _, tick := range tickers {
					if !symbolMatch(tick.Symbol, Binance.Symbols) {
						continue
					}
					normalizedTicker := NormalizedTicker{
						Timestamp: now,
						Exchange:  "Binance",
						Symbol:    tick.Symbol,
						Price:     tick.Price,
					}
					logChannel <- normalizedTicker
				}
			}
		}
	}()

	wg.Wait()
}

func symbolMatch(symbol string, symbols []string) bool {
	for _, s := range symbols {
		if symbol == s {
			return true
		}
	}
	return false
}

func sleep() {
	if Flags.Interval > 0 {
		time.Sleep(time.Duration(Flags.Interval) * time.Second)
	} else {
		sleepTime := time.Now().Truncate(time.Minute).Add(time.Minute).Sub(time.Now())
		time.Sleep(sleepTime)
	}
}

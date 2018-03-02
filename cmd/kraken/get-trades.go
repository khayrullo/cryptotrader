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
	"gitlab.com/crankykernel/ctrader/kraken"
	"github.com/spf13/viper"
	"log"
	"encoding/json"
	"math"
	"time"
	"sort"
	"github.com/spf13/pflag"
	"fmt"
	"strings"
)

var normalizedPairs map[string]string

func init() {
	normalizedPairs = map[string]string{
		"XXBTZUSD": "BTC/USD",
		"XETHZUSD": "ETH/USD",
		"XLTCZUSD": "LTC/USD",
		"XXMRZUSD": "XMR/USD",
		"XETCZUSD": "ETC/USD",
	}
}

func normalizePair(pair string) string {
	normalized, ok := normalizedPairs[pair]
	if ok {
		return normalized
	}
	if strings.HasSuffix(pair, "USD") {
		pair = strings.Replace(pair, "USD", "/USD", -1)
	}
	return pair
}

func parseTimestamp(value json.Number) (time.Time, error) {
	v64, err := value.Float64()
	if err != nil {
		return time.Time{}, err
	}
	sec, subsec := math.Modf(v64)
	return time.Unix(int64(sec), int64(subsec*(1e9))), nil
}

type Trade struct {
	Timestamp time.Time
	Type      string
	Pair      string
	Cost      string
	Fee       string
	Volume    string
}

func KrakenGetTrades(opts *pflag.FlagSet, args []string) {

	format, _ := opts.GetString("format")

	client := kraken.NewClientWithAuth(viper.GetString("kraken.api.key"),
		viper.GetString("kraken.api.secret"))

	params := map[string]interface{}{}

	//params["start"] = 1515980360
	//params["end"] = 1517276414
	//params["ofs"] = 2

	response, err := client.Post("/0/private/TradesHistory", params)
	if err != nil {
		log.Fatal(err)
	}

	body := map[string]interface{}{}

	decoder := json.NewDecoder(response.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&body); err != nil {
		log.Fatal("error: failed to decode response: ", err)
	}

	// Check for errors.
	errors, ok := body["error"].([]interface{})
	if ok {
		if len(errors) > 0 {
			log.Fatal("error: ", errors)
		}
	}

	rawTrades, ok := body["result"].(map[string]interface{})["trades"]
	if !ok {
		log.Fatal("error: failed to find trades in response")
	}

	trades := []Trade{}

	for key, trade := range rawTrades.(map[string]interface{}) {
		trade := trade.(map[string]interface{})

		timestamp, err := parseTimestamp(trade["time"].(json.Number))
		if err != nil {
			log.Fatal("error: failed to parse timestamp: %v", trade["time"])
		}

		trade["id"] = key
		trade["timestamp"] = timestamp

		xtrade := Trade{
			Timestamp: timestamp,
			Pair:      normalizePair(trade["pair"].(string)),
			Type:      strings.Title(trade["type"].(string)),
			Cost:      trade["cost"].(string),
			Fee:       trade["fee"].(string),
			Volume:    trade["vol"].(string),
		}

		trades = append(trades, xtrade)
	}

	sort.Slice(trades, func(i, j int) bool {
		if reverse, _ := opts.GetBool("reverse"); reverse {
			i, j = j, i
		}
		return trades[i].Timestamp.Before(trades[j].Timestamp)
	})

	for i, trade := range trades {
		switch format {
		case "":
			fallthrough
		case "json":
			buf, err := json.Marshal(trade)
			if err != nil {
				log.Println("error: ", err)
			} else {
				fmt.Println(string(buf))
			}
		case "csv":
			if i == 0 || i%10 == 0 {
				fmt.Printf("timestamp,type,pair,cost,fee,volume\n")
			}
			fmt.Printf("%s,%s,%s,%s,%s,%s\n",
				trade.Timestamp.Format("2006-01-02 15:04:05"),
				trade.Type,
				trade.Pair,
				trade.Cost,
				trade.Fee,
				trade.Volume)
		case "tab":
			if i == 0 || i%10 == 0 {
				fmt.Printf("timestamp\ttype\tpair\tcost\tfee\tvolume\n")
			}
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\n",
				trade.Timestamp.Format("2006-01-02 15:04:05"),
				trade.Type,
				trade.Pair,
				trade.Cost,
				trade.Fee,
				trade.Volume)
		}
	}

}

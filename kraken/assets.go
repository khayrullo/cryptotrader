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

import "strings"

var pairSubsTable [][]string

var assetSubsTable [][]string

func init() {
	pairSubsTable = [][]string{
		{"XXBTZ", "BTC/"},
		{"XETHZ", "ETH/"},
		{"XLTCZ", "LTC/"},
		{"XXMRZ", "XMR/"},
		{"XETCZ", "ETC/"},
		{"BCH", "BCH/"},
	}

	assetSubsTable = [][]string{
		{"XXBT", "BTC"},
		{"XLTC", "LTC"},
		{"XXMR", "XMR"},
		{"XETH", "ETH"},
		{"XZEC", "ZEC"},
		{"ZUSD", "USD"},
	}
}

func GetNormalizePairName(pair string) string {
	for _, sub := range pairSubsTable {
		pair = strings.Replace(pair, sub[0], sub[1], -1)
	}

	return pair
}

func NormalizeAssetName(name string) string {
	for _, sub := range assetSubsTable {
		name = strings.Replace(name, sub[0], sub[1], -1)
	}

	return name
}
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
	"gitlab.com/crankykernel/cryptotrader/kraken"
	"github.com/spf13/viper"
	"log"
	"fmt"
	"sort"
	"strings"
	"math"
	"encoding/json"
)

var KrakenLedgerFlags struct {
	Format string
	Merged bool
	Count  int
	Type   string
}

func KrakenLedgerCmd() {
	client := kraken.NewClient(
		viper.GetString("kraken.api.key"),
		viper.GetString("kraken.api.secret"))

	count := KrakenLedgerFlags.Count
	if count > 0 && KrakenLedgerFlags.Merged {
		count = count * 2
	}

	options := kraken.GetLedgerOptions{
		Count: count,
		Type:  KrakenLedgerFlags.Type,
	}

	ledger, err := kraken.NewLedgerService(client).Ledger(options)
	if err != nil {
		log.Fatal("error: ", err)
	}

	if KrakenLedgerFlags.Merged {
		merged := mergeEntries(ledger)
		switch KrakenLedgerFlags.Format {
		case "tab":
			printMergedDelim(merged, "\t")
		case "csv":
			printMergedDelim(merged, ",")
		case "":
			printMergedPretty(merged)
		default:
			log.Fatal(fmt.Sprintf(
				"error: unsupported format for --merged: %s",
				KrakenLedgerFlags.Format))
		}
		return
	}

	// TODO Support formats non-merged printing.
	for _, entry := range ledger {
		buf, _ := json.Marshal(entry)
		fmt.Println(string(buf))
	}
}

func printMergedPretty(merged []MergedEntry) {
	for _, entry := range merged {
		var fee float64
		var feeAsset string

		out := entry.Out
		in_ := entry.In

		if out.Fee > 0 {
			fee = out.Fee
			feeAsset = fmt.Sprintf(" %s", out.Asset)
		} else if in_.Fee > 0 {
			fee = in_.Fee
			feeAsset = fmt.Sprintf(" %s", in_.Asset)
		}

		fmt.Printf("Timestamp: %s; "+
			"In: %.8f %s; "+
			"Out: %.8f %s; "+
			"Fee: %.8f%s; "+
			"In Balance: %.8f %s; "+
			"Out Balance: %.8f %s; "+
			"\n",
			out.Timestamp.Format("2006-01-02 15:04:05"),
			in_.Amount, in_.Asset,
			out.Amount, out.Asset,
			fee, feeAsset,
			in_.Balance, in_.Asset,
			out.Balance, out.Asset)
	}
}

func printMergedDelim(entries []MergedEntry, delim string) {
	header := []string{
		"timestamp",
		"type",
		"in",
		"out",
		"in_amount",
		"in_fee",
		"out_amount",
		"out_fee",
		"in_balance",
		"out_balance",
	}
	fmt.Printf("%s\n", strings.Join(header, delim))
	for _, e := range entries {
		parts := []string{}
		if e.Type == "trade" {
			parts = []string{
				e.Out.Timestamp.Format("2006-01-02 15:04:05"),
				"trade",
				kraken.NormalizeAssetName(e.In.Asset),
				kraken.NormalizeAssetName(e.Out.Asset),
				fmt.Sprintf("%.8f", e.In.Amount),
				fmt.Sprintf("%.8f", e.In.Fee),
				fmt.Sprintf("%.8f", math.Abs(e.Out.Amount)),
				fmt.Sprintf("%.8f", e.Out.Fee),
				fmt.Sprintf("%.8f", e.In.Balance),
				fmt.Sprintf("%.8f", e.Out.Balance),
			}
		} else if e.Type == "deposit" {
			parts = []string{
				e.In.Timestamp.Format("2006-01-02 15:04:05"),
				"deposit",
				kraken.NormalizeAssetName(e.In.Asset),
				"",
				fmt.Sprintf("%.8f", e.In.Amount),
				fmt.Sprintf("%.8f", e.In.Fee),
				"",
				"",
				fmt.Sprintf("%.8f", e.In.Balance),
				"",
			}
		} else if e.Type == "withdrawal" {
			parts = []string{
				e.Out.Timestamp.Format("2006-01-02 15:04:05"),
				"withdrawal",
				"",
				e.Out.Asset,
				"",
				"",
				fmt.Sprintf("%.8f", e.Out.Amount),
				fmt.Sprintf("%.8f", e.Out.Fee),
				"",
				fmt.Sprintf("%.8f", e.Out.Balance),
			}
		} else {
			log.Fatal("error: unsupport ledger entry: ", e)
		}
		fmt.Printf("%s\n", strings.Join(parts, delim))
	}
}

// MergedEntry contains the pair of ledge entries where one is the asset in,
// and the other is the asset out.
type MergedEntry struct {
	Type string

	In  kraken.LedgerEntry
	Out kraken.LedgerEntry

	HasIn  bool
	HasOut bool
}

func mergeEntries(entries []kraken.LedgerEntry) []MergedEntry {

	mergeSet := map[string]MergedEntry{}

	for _, entry := range entries {
		referenceId := entry.ReferenceID

		var mergedEntry MergedEntry
		mergedEntry, ok := mergeSet[referenceId]
		if !ok {
			mergedEntry = MergedEntry{}
		} else {
		}

		if entry.Type == "deposit" {
			mergedEntry.In = entry
			mergedEntry.Type = "deposit"
		} else if entry.Type == "withdrawal" {
			mergedEntry.Out = entry
			mergedEntry.Type = "withdrawal"
		} else if entry.Type == "trade" {
			mergedEntry.Type = "trade"
			if entry.Amount < 0 {
				mergedEntry.Out = entry
				mergedEntry.HasOut = true
			} else {
				mergedEntry.In = entry
				mergedEntry.HasIn = true
			}
		} else {
			log.Fatal("error: unknown ledger type: ", entry.Type)
		}

		mergeSet[referenceId] = mergedEntry
	}

	// Return back to a list.
	merged := []MergedEntry{}
	for _, entry := range mergeSet {
		merged = append(merged, entry)
	}

	sort.Slice(merged, func(i, j int) bool {
		return merged[i].Out.Timestamp.Before(merged[j].Out.Timestamp)
	})

	return merged
}

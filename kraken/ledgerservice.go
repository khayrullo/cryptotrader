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
	"time"
	"cryptotrader/util"
	"strconv"
	"sort"
	"fmt"
	"log"
)

type LedgerEntry struct {
	LedgerID    string
	ReferenceID string
	Timestamp   time.Time
	Type        string
	AssetClass  string
	Asset       string
	Amount      float64
	Fee         float64
	Balance     float64
}

func NewLedgerEntryFromRaw(id string, raw RawLedgerEntry) LedgerEntry {
	entry := LedgerEntry{}
	entry.LedgerID = id
	entry.ReferenceID = raw.RefID
	entry.Timestamp = util.Float64ToTime(raw.Time)
	entry.Type = raw.Type
	entry.AssetClass = raw.AssetClass
	entry.Asset = NormalizeAssetName(raw.Asset)
	entry.Amount, _ = strconv.ParseFloat(raw.Amount, 64)
	entry.Fee, _ = strconv.ParseFloat(raw.Fee, 64)
	entry.Balance, _ = strconv.ParseFloat(raw.Balance, 64)
	return entry
}

type RawLedgerEntry struct {
	RefID      string  `json:"refid"`
	Time       float64 `json:"time"`
	Type       string  `json:"type"`
	AssetClass string  `json:"aclass"`
	Asset      string  `json:"asset"`
	Amount     string  `json:"amount"`
	Fee        string  `json:"fee"`
	Balance    string  `json:"balance"`
}

type RawLedgerResponse struct {
	Error []string `json:"error"`
	Result struct {
		Count  int64                     `json:"count"`
		Ledger map[string]RawLedgerEntry `json:"ledger""`
	} `json:"result"`
	Raw string `json:"-"`
}

func (r *RawLedgerResponse) SetRaw(raw string) {
	r.Raw = raw
}

type LedgerService struct {
	client *Client
}

func NewLedgerService(client *Client) *LedgerService {
	return &LedgerService{client}
}

func (s *LedgerService) RawLedger(params map[string]interface{}) (*RawLedgerResponse, error) {
	endpoint := "/0/private/Ledgers"
	httpResponse, err := s.client.Post(endpoint, params)
	if err != nil {
		return nil, err
	}
	response := RawLedgerResponse{}
	if err := decodeBody(httpResponse, &response); err != nil {
		return nil, nil
	}
	return &response, nil
}

type GetLedgerOptions struct {
	Count int
	Type  string
}

func (s *LedgerService) Ledger(options GetLedgerOptions) ([]LedgerEntry, error) {
	entries := []LedgerEntry{}

	var end *time.Time
	lastCount := 0

	for {
		params := map[string]interface{}{}

		if end != nil {
			params["end"] = end.Unix() + 1
		}

		if options.Type != "" {
			params["type"] = options.Type
		}

		response, err := s.RawLedger(params)
		if err != nil {
			return nil, err
		}

		if len(response.Error) > 0 {
			if response.Error[0] == "EAPI:Rate limit exceeded" {
				log.Println("warning: rate limit exceeded, sleeping for 5s")
				time.Sleep(5 * time.Second)
				continue
			}
			return entries, fmt.Errorf("%s", response.Error[0])
		}

		if len(response.Result.Ledger) == 0 {
			break
		}

		for ledgerId, v := range response.Result.Ledger {
			entry := NewLedgerEntryFromRaw(ledgerId, v)
			entries = append(entries, entry)
		}

		sort.Slice(entries, func(i, j int) bool {
			return entries[i].Timestamp.Before(entries[j].Timestamp)
		})

		// Dedupe.
		entries = dedupe(entries)

		if len(entries) == lastCount {
			break
		}

		lastCount = len(entries)

		if options.Count > 0 && len(entries) >= options.Count {
			entries = entries[len(entries)-options.Count:]
			break
		}

		end = &entries[0].Timestamp
	}

	return entries, nil
}

func dedupe(entries []LedgerEntry) []LedgerEntry {
	deduped := []LedgerEntry{}
	for _, entry := range entries {
		if len(deduped) == 0 {
			deduped = append(deduped, entry)
		} else if (deduped[len(deduped)-1].LedgerID != entry.LedgerID) {
			deduped = append(deduped, entry)
		} else {
			// Dupe.
		}
	}
	return deduped
}

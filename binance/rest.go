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
	"fmt"
	"net/http"
	"sort"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"strconv"
)

const API_ROOT = "https://api.binance.com"

type RestClient struct {
}

func NewClient() *RestClient {
	client := RestClient{}
	return &client
}

func (c *RestClient) Get(endpoint string, params map[string]interface{}) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params != nil {
		queryString = c.BuildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}

func (c *RestClient) BuildQueryString(params map[string]interface{}) string {
	queryString := ""

	keys := func() []string {
		keys := []string{}
		for key, _ := range params {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys
	}()

	for _, key := range keys {
		if queryString != "" {
			queryString = fmt.Sprintf("%s&", queryString)
		}
		queryString = fmt.Sprintf("%s%s=%v", queryString, key, params[key])
	}

	return queryString
}

type LastResponseRaw struct {
	Symbol string `json:"symbol"`
	Price  string `json:"price"`
}

type LastResponse struct {
	Symbol string
	Price  float64
}

// Return the latest prices for all symbols.
func (c *RestClient) Last() ([]LastResponse, error) {
	endpoint := "/api/v3/ticker/price"
	httpResponse, err := c.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}
	responseRaw := []LastResponseRaw{}
	_, err = c.decodeBody(httpResponse, &responseRaw)
	if err != nil {
		return nil, err
	}

	response := []LastResponse{}
	for _, last := range responseRaw {
		price, err := strconv.ParseFloat(last.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse as float64: %s: %v",
				last.Price, err)
		}
		response = append(response, LastResponse{
			Symbol: last.Symbol,
			Price:  price,
		})
	}

	return response, nil
}

func (c *RestClient) decodeBody(r *http.Response, v interface{}) ([]byte, error) {
	raw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(v); err != nil {
		return nil, err
	}
	return raw, nil
}

func (c *RestClient) GetAllSymbols() ([]string, error) {
	lastTrades, err := c.Last()
	if err != nil {
		return nil, err
	}
	symbols := []string{}
	for _, trade := range lastTrades {
		symbols = append(symbols, trade.Symbol)
	}
	return symbols, nil
}

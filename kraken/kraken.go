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
	"net/http"
	"fmt"
	"sort"
	"time"
	"encoding/base64"
	"log"
	"crypto/sha256"
	"crypto/hmac"
	"crypto/sha512"
	"strings"
	"encoding/json"
	"strconv"
)

const API_ROOT = "https://api.kraken.com"

type Client struct {
	apiKey    string
	apiSecret []byte
}

func NewClient(apiKey string, apiSecret string) *Client {
	var decodedApiSecret []byte
	var err error
	if apiSecret != "" {
		decodedApiSecret, err = base64.StdEncoding.DecodeString(apiSecret)
		if err != nil {
			log.Fatal("error: failed to base64 decode kraken api secret")
		}
	} else {
		decodedApiSecret = nil
	}

	return &Client{
		apiKey:    apiKey,
		apiSecret: decodedApiSecret,
	}
}

// HasAuth returns true if client has authentication information.
func (c *Client) HasAuth() bool {
	if c.apiKey != "" && c.apiSecret != nil {
		return true
	}
	return false
}

func (c *Client) Get(endpoint string, params map[string]interface{}) (*http.Response, error) {

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

func (c *Client) Post(endpoint string, params map[string]interface{}) (*http.Response, error) {

	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	nonce := c.getNonce()
	if params == nil {
		params = map[string]interface{}{
			"nonce": nonce,
		}
	} else {
		params["nonce"] = nonce
	}

	if params != nil {
		queryString = c.BuildQueryString(params)
	}

	request, err := http.NewRequest("POST", url, strings.NewReader(queryString))
	if err != nil {
		return nil, err
	}

	if c.apiKey != "" && c.apiSecret != nil {
		c.authenticateRequest(request, endpoint, nonce, queryString)
	}

	return http.DefaultClient.Do(request)
}

func (c *Client) authenticateRequest(request *http.Request, endpoint string, nonce int64, postData string) {
	s256 := sha256.New()
	s256.Write([]byte(fmt.Sprintf("%d%s", nonce, postData)))

	mac := hmac.New(sha512.New, c.apiSecret)
	mac.Write([]byte(endpoint))
	mac.Write(s256.Sum(nil))

	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	request.Header.Add("API-Key", c.apiKey)
	request.Header.Add("API-Sign", signature)
}

func (c *Client) BuildQueryString(params map[string]interface{}) string {
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

func (c *Client) getNonce() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type Ticker struct {
	Pair      string
	Timestamp time.Time
	Ask       float64
	Bid       float64
	Last      float64
}

type RawTickerPair struct {
	A []string `json:"a"` // Ask array.
	B []string `json:"b"` // Bid array.
	C []string `json:"c"` // Last array.
	V []string `json:"v"`
	P []string `json:"p"`
	T []int64  `json:"t"`
	H []string `json:"h"`
	O string   `json:"o"`
}

type RawTickerResponse struct {
	Error  []string                 `json:"error"`
	Result map[string]RawTickerPair `json:"result"`
}

func (c *Client) Ticker(pairs ... string) (tickers map[string]Ticker, err error) {
	endpoint := "/0/public/Ticker"

	if len(pairs) == 0 {
		return nil, fmt.Errorf("no pairs")
	}

	params := map[string]interface{}{
		"pair": strings.Join(pairs, ","),
	}

	var r *http.Response

	if (c.HasAuth()) {
		r, err = c.Post(endpoint, params)
	} else {
		r, err = c.Get(endpoint, params)
	}

	now := time.Now()

	decoder := json.NewDecoder(r.Body)
	decoder.UseNumber()
	rawResponse := RawTickerResponse{}
	if err = decoder.Decode(&rawResponse); err != nil {
		return nil, err
	}

	if len(rawResponse.Error) > 0 {
		return nil, fmt.Errorf("%s", rawResponse.Error[0])
	}

	tickers = map[string]Ticker{}

	for pair, val := range rawResponse.Result {
		normalizedPair := GetNormalizePairName(pair)
		ticker := Ticker{}
		ticker.Timestamp = now
		ticker.Pair = normalizedPair
		ticker.Ask, _ = strconv.ParseFloat(val.A[0], 64)
		ticker.Bid, _ = strconv.ParseFloat(val.B[0], 64)
		ticker.Last, _ = strconv.ParseFloat(val.C[0], 64)
		tickers[normalizedPair] = ticker
	}

	return tickers, nil
}

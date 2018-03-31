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

package quadriga

import (
	"fmt"
	"net/http"
	"bytes"
	"encoding/json"
	"sort"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
)

const API_ROOT = "https://api.quadrigacx.com"

type Client struct {
	clientId  string
	apiKey    string
	apiSecret string
}

func NewClient(clientId interface{}, apiKey string, apiSecret string) *Client {
	quadriga := &Client{
		clientId:  fmt.Sprintf("%s", clientId),
		apiKey:    apiKey,
		apiSecret: apiSecret,
	}
	return quadriga
}

func NewAnonymousClient() *Client {
	return &Client{}
}

func (c *Client) Post(endpoint string, params map[string]interface{}) (*http.Response, error) {
	url := c.buildURL(endpoint)
	body, err := json.Marshal(c.authenticateParams(params))
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(request)
}

func (c *Client) Get(endpoint string, params map[string]interface{}) (*http.Response, error) {
	url := c.buildURL(endpoint)
	queryString := c.buildQueryString(params)
	if queryString != "" {
		url = fmt.Sprintf("%s?%s", url, queryString)
	}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(request)
}

func (c *Client) authenticateParams(params map[string]interface{}) map[string]interface{} {
	// First make a new set of params and copy in the provided ones. We don't
	// worry about making a deep copy as we'll only add fields to the top
	// level.
	outParams := map[string]interface{}{}
	for k, v := range params {
		outParams[k] = v
	}

	nonce := c.getNonce()
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(fmt.Sprintf("%d%s%s", nonce, c.clientId, c.apiKey)))
	sig := mac.Sum(nil)

	outParams["nonce"] = nonce
	outParams["key"] = c.apiKey
	outParams["signature"] = hex.EncodeToString(sig)

	return outParams
}

func (c *Client) buildURL(endpoint string) string {
	return fmt.Sprintf("%s/%s", API_ROOT, endpoint)
}

// buildQueryString converts the params map into a query string sorted by
// field names as required by some (but not all APIs).
func (c *Client) buildQueryString(params map[string]interface{}) string {
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

// Return the available order books.
func (c *Client) Books() ([]string, error) {
	response, err := c.Get("/v2/ticker", map[string]interface{}{
		"book": "all",
	})
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("%s: %s", response.Status, string(body))
	}

	ticker := map[string]interface{}{}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	if err := decoder.Decode(&ticker); err != nil {
		return nil, err
	}

	books := []string{}

	for key, _ := range ticker {
		books = append(books, key)
	}

	return books, nil
}

// RequestEngineOrders calls an undocumented URL that the Quadriga frontend
// userse to get the order book.
func (c *Client) RequestEngineOrders(book string) (*http.Response, error) {
	return http.DefaultClient.Get(fmt.Sprintf(
		"https://www.quadrigacx.com/engine/orders/%s", book))
}

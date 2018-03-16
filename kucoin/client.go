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

package kucoin

import (
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
	"sort"
)

const API_ROOT = "https://api.kucoin.com"

type Client struct {
	apiKey    string
	apiSecret string
}

func NewClient(key string, secret string) *Client {
	client := Client{}
	client.apiKey = key
	client.apiSecret = secret
	return &client
}

func (c *Client) Get(endpoint string, params map[string]interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	queryString := ""

	if params != nil {
		queryString = buildQueryString(params)
		if queryString != "" {
			url = fmt.Sprintf("%s?%s", url, queryString)
		}
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" && c.apiSecret != "" {
		c.authenticateRequest(request, endpoint, queryString)
	}
	return http.DefaultClient.Do(request)
}

func (c *Client) authenticateRequest(request *http.Request, endpoint string,
	queryString string) error {
	nonce := c.getNonce()
	signature := c.getSignature(endpoint, nonce, queryString)
	request.Header.Add("KC-API-SIGNATURE", signature)
	request.Header.Add("KC-API-NONCE", fmt.Sprintf("%d", nonce))
	request.Header.Add("KC-API-KEY", c.apiKey)
	return nil
}

func (c *Client) getNonce() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func (c *Client) getSignature(endpoint string, nonce int64, queryString string) string {
	signature := base64.StdEncoding.EncodeToString(
		[]byte(fmt.Sprintf("%s/%d/%s", endpoint, nonce, queryString)))
	mac := hmac.New(sha256.New, []byte(c.apiSecret))
	mac.Write([]byte(signature))
	return hex.EncodeToString(mac.Sum(nil))
}

// Build the query string as KuCoin expects it, in alphabetical order. This
// is important as its part of the signature.
func buildQueryString(params map[string]interface{}) string {
	queryString := ""

	// Create a sorted list of the query parameter field names...
	keys := func() []string {
		keys := []string{}
		for key, _ := range params {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		return keys
	}()

	// Build the query string.
	for _, key := range keys {
		if queryString != "" {
			queryString = fmt.Sprintf("%s&", queryString)
		}
		queryString = fmt.Sprintf("%s%s=%v", queryString, key, params[key])
	}

	return queryString
}

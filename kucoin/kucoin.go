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
	"net/http"
	"fmt"
	"time"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strings"
	"sort"
)

const API_ROOT = "https://api.kucoin.com"

type Client struct {
	apiKey    string
	apiSecret string
}

func NewClient() *Client {
	client := &Client{}
	return client
}

func NewClientWithAuth(key string, secret string) *Client {
	client := NewClient()
	client.apiKey = key
	client.apiSecret = secret
	return client
}

type UserInfo struct {
	Raw string

	Success   bool                   `json:"success"`
	Code      string                 `json:"code"`
	Msg       string                 `json:"msg"`
	Timestamp int64                  `json:"timestamp"`
	Data      map[string]interface{} `json:"data"`
}

func (t *UserInfo) PrintPretty() string {
	decoder := json.NewDecoder(strings.NewReader(t.Raw))
	decoder.UseNumber()
	generic := map[string]interface{}{}
	decoder.Decode(&generic)
	output, _ := json.MarshalIndent(generic, "", "  ")
	return string(output)
}

func (c *Client) GetUserInfo() (*UserInfo, error) {
	endPoint := "/v1/user/info"
	response, err := c.Get(endPoint, nil)

	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(rawBody))
	decoder.UseNumber()
	userInfo := UserInfo{Raw: string(rawBody),}
	err = decoder.Decode(&userInfo)
	if err != nil {
		return nil, err
	}
	fmt.Println(userInfo.PrintPretty())
	return &userInfo, nil
}

type Trade struct {
	CoinType        string  `json:"coinType"`
	CreatedAtMillis int64   `json:"createdAt"`
	Amount          float64 `json:"amount"`
	DealValue       float64 `json:"dealValue"`
	Fee             float64 `json:"fee"`
	DealDirection   string  `json:"dealDirection"`
	CoinTypePair    string  `json:"coinTypePair"`
	OID             string  `json:"oid"`
	DealPrice       float64 `json:"dealPrice"`
	OrderID         string  `json:"orderOid"`
	FeeRate         float64 `json:"feeRate"`
	Direction       string  `json:"direction"`

	Timestamp time.Time `json:"-"`
}

type DealtOrdersResponse struct {
	Success         bool   `json:"success"`
	Code            string `json:"code"`
	Message         string `json:"msg"`
	TimestampMillis int64  `json:"timestamp"`
	Data struct {
		Total  int64    `json:"total"`
		Limit  int64    `json:"limit"`
		Page   int64    `json:"page"`
		Trades []*Trade `json:"datas"`
	} `json:"data"`

	Raw string
}

func (c *Client) GetDealtOrders(limit int, page int) (*DealtOrdersResponse, error) {
	endPoint := "/v1/order/dealt"

	params := map[string]interface{}{
		"page": page + 1,
	}
	if limit == 0 {
		params["limit"] = 20
	} else {
		params["limit"] = limit
	}

	response, err := c.Get(endPoint, params)
	if err != nil {
		return nil, err
	}
	rawBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	orders := DealtOrdersResponse{}
	if err := decode(rawBody, &orders); err != nil {
		return nil, err
	}
	orders.Raw = string(rawBody)

	for _, trade := range orders.Data.Trades {
		trade.Timestamp = time.Unix(trade.CreatedAtMillis/1000, 0)
	}

	return &orders, nil
}

func decode(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
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

func buildQueryString(params map[string]interface{}) string {
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

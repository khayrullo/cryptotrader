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
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strings"
	"gitlab.com/crankykernel/cryptotrader/util"
)

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

type WalletRecordEntry struct {
	Address         string      `json:"address"`
	Amount          float64     `json:"amount"`
	CoinType        string      `json:"coinType"`
	Confirmation    int64       `json:"confirmation"`
	CreatedAtMillis int64       `json:"createdAt"`
	Fee             float64     `json:"fee"`
	OID             string      `json:"oid"`
	OuterWalletTxID string      `json:"outerWalletTxid"`
	Remark          interface{} `json:"remark"`
	Status          string      `json:"status"`
	Type            string      `json:"type"`
	UpdatedAtMillis int64       `json:"updateAt"`
}

type WalletRecordsResponse struct {
	Success         bool   `json:"success"`
	Code            string `json:"code"`
	Message         string `json:"msg"`
	TimestampMillis int64  `json:"timestamp"`
	Data struct {
		CurrPageNo int64               `json:"currPageNo"`
		FirstPage  bool                `json:"firstPage"`
		LastPage   bool                `json:"lastPage"`
		Total      int64               `json:"total"`
		Limit      int64               `json:"limit"`
		Entries    []WalletRecordEntry `json:"datas"`
	} `json:"data"`
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
	defer response.Body.Close()
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

func (c *Client) WalletRecords(coin string, page int) (*WalletRecordsResponse, error) {
	endpoint := fmt.Sprintf("/v1/account/%s/wallet/records",
		strings.ToUpper(coin))

	params := map[string]interface{}{
		"page": page,
	}

	httpResponse, err := c.Get(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	var response WalletRecordsResponse
	if err := decode(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

type TickResponse struct {
	Success         bool        `json:"success"`
	Code            string      `json:"code"`
	Message         string      `json:"msg"`
	TimestampMillis int64       `json:"timestamp"`
	Entries         []TickEntry `json:"data"`
	Raw             string
}

func (t *TickResponse) GetTimestamp() time.Time {
	return util.MillisToTime(t.TimestampMillis)
}

type TickEntry struct {
	CoinType       string  `json:"coinType"`
	Trading        bool    `json:"trading"`
	Symbol         string  `json:"symbol"`
	LastDealPrice  float64 `json:"lastDealPrice"`
	Buy            float64 `json:"buy"`
	Sell           float64 `json:"sell"`
	Change         float64 `json:"change"`
	CoinTypePair   string  `json:"coinTypePair"`
	Sort           int64   `json:"sort"`
	FeeRate        float64 `json:"feeRate"`
	VolValue       float64 `json:"volValue"`
	High           float64 `json:"high"`
	DateTimeMillis int64   `json:"datetime"`
	Vol            float64 `json:"vol"`
	Low            float64 `json:"low"`
	ChangeRate     float64 `json:"changeRate"`
}

func (c *Client) GetTick() (*TickResponse, error) {
	endpoint := "/v1/open/tick"

	httpResponse, err := c.Get(endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	body, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode != 200 {
		return nil, fmt.Errorf("%v", httpResponse.Status)
	}

	var response TickResponse
	if err := decode(body, &response); err != nil {
		return nil, err
	}

	response.Raw = string(body)

	return &response, nil
}

func decode(data []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	return decoder.Decode(v)
}

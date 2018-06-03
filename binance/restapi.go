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
	"io/ioutil"
	"net/http"
	"fmt"
	"log"
)

type RestApiError struct {
	StatusCode int
	Body       []byte
}

func NewRestApiErrorFromResponse(r *http.Response) *RestApiError {
	body, _ := ioutil.ReadAll(r.Body)
	return &RestApiError{
		StatusCode: r.StatusCode,
		Body:       body,
	}
}

func (e *RestApiError) Error() string {
	return string(e.Body)
}

type UserDataStreamResponse struct {
	ListenKey string `json:"listenKey"`
}

// GetUserDataStream makes the get request for a user data stream listen key.
func (c *RestClient) GetUserDataStream() (string, error) {
	httpResponse, err := c.PostWithApiKey("/api/v1/userDataStream", nil)
	if err != nil {
		return "", err
	}

	if httpResponse.StatusCode >= 400 {
		return "", NewRestApiErrorFromResponse(httpResponse)
	}

	var response UserDataStreamResponse
	if _, err = c.decodeBody(httpResponse, &response); err != nil {
		return "", err
	}

	return response.ListenKey, nil
}

func (c *RestClient) PutUserStreamKeepAlive(listenKey string) error {
	queryString := c.BuildQueryString(map[string]interface{}{
		"listenKey": listenKey,
	})
	path := fmt.Sprintf("/api/v1/userDataStream?%s", queryString)
	log.Println(path)
	httpResponse, err := c.DoPut(path)
	if err != nil {
		return err
	}
	if httpResponse.StatusCode != http.StatusOK {
		return NewRestApiErrorFromResponse(httpResponse)
	}
	return nil
}

type OrderSide string

const (
	BUY  OrderSide = "BUY"
	SELL OrderSide = "SELL"
)

type OrderType string

const (
	LIMIT  OrderType = "LIMIT"
	MARKET OrderType = "MARKET"
)

type TimeInForce string

const (
	GTC TimeInForce = "GTC"
	IOC TimeInForce = "IOC"
	FOK TimeInForce = "FOK"
)

// Order status / execution type.
type OrderStatus string

const (
	OrderStatusNew      OrderStatus = "NEW"
	OrderStatusCanceled OrderStatus = "CANCELED"
)

type OrderParameters struct {
	Symbol           string
	Side             OrderSide
	Type             OrderType
	TimeInForce      TimeInForce
	Quantity         float64
	Price            float64
	NewClientOrderId string
}

func (c *RestClient) PostOrder(order OrderParameters) (*http.Response, error) {
	params := map[string]interface{}{}
	params["symbol"] = order.Symbol
	params["side"] = order.Side
	params["type"] = order.Type
	params["quantity"] = order.Quantity
	params["price"] = fmt.Sprintf("%.8f", order.Price)
	params["newClientOrderId"] = order.NewClientOrderId
	params["timeInForce"] = order.TimeInForce

	response, err := c.Post("/api/v3/order", params)
	if err != nil {
		return nil, err
	}
	if response.StatusCode >= 400 {
		return response, NewRestApiErrorFromResponse(response)
	}
	return response, nil
}

func (c *RestClient) CancelOrder(symbol string, orderId int64) (*http.Response, error) {
	params := map[string]interface{}{}
	params["symbol"] = symbol
	params["orderId"] = orderId

	return c.Delete("/api/v3/order", params)
}

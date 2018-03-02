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

package gdax

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"bytes"
)

const API_ROOT = "https://api.gdax.com"

type ApiClient struct {
}

func NewApiClient() *ApiClient {
	return &ApiClient{}
}

func (c *ApiClient) Get(endpoint string) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", API_ROOT, endpoint)
	return http.Get(url)
}

type Product struct {
	Id                string `json:"id"`
	DisplayName       string `json:"display_name"`
	QuoteIncrementRaw string `json:"quote_increment"`
}

func (c *ApiClient) Products() ([]Product, error) {
	endpoint := "/products"
	response, err := c.Get(endpoint)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(response.Status)
	}

	var products []Product

	_, err = parseResponse(response, &products)
	if err != nil {
		return nil, err
	}

	return products, nil
}

func parseResponse(response *http.Response, v interface{}) (string, error) {
	raw, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	err = decoder.Decode(v)
	return string(raw), err
}

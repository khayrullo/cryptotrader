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

type SymbolFilterResponse struct {
	FilterType  string  `json:"filterType"`
	MinPrice    float64 `json:"minPrice,string"`
	MaxPrice    float64 `json:"maxPrice,string"`
	TickSize    float64 `json:"tickSize,string"`
	MinQty      float64 `json:"minQty,string"`
	MaxQty      float64 `json:"maxQty,string"`
	StepSize    float64 `json:"stepSize,string"`
	MinNotional float64 `json:"minNotional,string"`
}

type SymbolInfoResponse struct {
	Symbol              string                 `json:"symbol"`
	Status              string                 `json:"status"`
	BaseAsset           string                 `json:"baseAsset"`
	BaseAssetPrecision  int64                  `json:"baseAssetPrecision`
	QuoteAsset          string                 `json:"quoteAsset"`
	QuoteAssetPrecision int64                  `json:"quoteAssetPrecision"`
	OrderTypes          []string               `json:"orderTypes"`
	IcebergAllowed      bool                   `json:"icebergAllowed"`
	Filters             []SymbolFilterResponse `json:"filters"`
}

type ExchangeInfoResponse struct {
	Timezone         string `json:"timezone"`
	ServerTimeMillis int64  `json:"serverTime"`
	RateLimits []struct {
		RateLimitType     string `json:"rateLimitType"`
		RateLimitInterval string `json"rateLimitInterval"`
		Limit             int64  `json"limit"`
	}
	Symbols []SymbolInfoResponse `json:"symbols"`

	RawResponse []byte `json:"-"`
}

type CancelOrderResponse struct {
	Symbol            string `json:"symbol"`
	OrigClientOrderID string `json:"origClientOrderId"`
	OrderID           int64  `json:"orderId"`
	ClientOrderID     string `json:"clientOrderId"`
}

package binance

import (
	"time"
	"github.com/crankykernel/cryptotrader/util"
	"strconv"
)

type RawAggTrade struct {
	EventType       string  `json:"e"`
	EventTimeMillis float64 `json:"E"`
	Symbol          string  `json:"s"`
	TradeID         int64   `json:"a"`
	Price           string  `json:"p"`
	Quantity        string  `json:"q"`
	FirstTradeID    int64   `json:"f"`
	LastTradeID     int64   `json:"l"`
	TradeTimeMillis float64 `json:"T"`
	BuyerMaker      bool    `json:"m"`
	Ignored         bool    `json:"M"`
}

type AggTrade struct {
	Timestamp  time.Time
	Symbol     string
	Price      float64

	// Quantity of the symbol.
	Quantity   float64

	QuoteQuantity float64

	// If the buyer made the order, then this was a sell.
	BuyerMaker bool
}

func NewAggTradeFromRaw(raw RawAggTrade) AggTrade {
	trade := AggTrade{}

	trade.Timestamp = util.Float64ToTime(raw.TradeTimeMillis / 1000)
	trade.Symbol = raw.Symbol
	trade.Price, _ = strconv.ParseFloat(raw.Price, 64)
	trade.Quantity, _ = strconv.ParseFloat(raw.Quantity, 64)
	trade.BuyerMaker = raw.BuyerMaker
	trade.QuoteQuantity = trade.Quantity * trade.Price

	return trade
}

func (t *AggTrade) IsBuy() bool {
	if t.BuyerMaker {
		return false
	}
	return true
}
package binance

import (
	"time"
	"strings"
	"encoding/json"
	"fmt"
	"gitlab.com/crankykernel/cryptotrader/util"
	"strconv"
)

// RawTicker24 is the deconstruction of the raw 24hr JSON ticker as sent from
// Binance.
type RawTicker24 struct {
	EventType            string `json:"e"`
	EventTime            int64  `json:"E"`
	Symbol               string `json:"s"`
	PriceChange          string `json:"p"`
	PriceChangePercent   string `json:"P"`
	WeightedAveragePrice string `json:"w"`
	PreviousDayClose     string `json:"x"`
	CurrentDayClose      string `json:"c"`
	CloseTradeQuantity   string `json:"Q"`
	Bid                  string `json:"b"`
	BidQuantity          string `json:"B"`
	Ask                  string `json:"a"`
	AskQuantity          string `json:"A"`
	OpenPrice            string `json:"o"`
	HighPrice            string `json:"h"`
	LowPrice             string `json:"l"`
	TotalBaseVolume      string `json:"v"`
	TotalQuoteVolume     string `json:"q"`
	StatsOpenTime        int64  `json:"O"`
	StatsCloseTime       int64  `json:"C"`
	FirstTradeID         int64  `json:"F"`
	LastTradeID          int64  `json:"L"`
	TotalNumberTrades    int64  `json:"n"`
}

// Ticker24 is a "nicer" to use variation of RawTicker24.
type Ticker24 struct {
	Timestamp      time.Time `json:"timestamp"`
	Symbol         string    `json:"symbol"`
	Close          float64   `json:"close"`
	PriceChange    float64   `json:"price_change"`
	PriceChangePct float64   `json:"price_change_pct"`
	LowPrice       float64   `json:"low"`
	HighPrice      float64   `json:"high"`
	BaseVolume     float64   `json:"base_volume"`
	QuoteVolume    float64   `json:"quote_volume"`
	Bid            float64   `json:"bid"`
	Ask            float64   `json:"ask"`
}

// Converts a raw ticker to a more useable ticker.
func NewTicker24FromRawTicker24(raw RawTicker24) Ticker24 {
	ticker := Ticker24{}

	ticker.Timestamp = util.MillisToTime(raw.EventTime)
	ticker.Symbol = raw.Symbol
	ticker.Close, _ = strconv.ParseFloat(raw.CurrentDayClose, 64)
	ticker.PriceChange, _ = strconv.ParseFloat(raw.PriceChange, 64)
	ticker.PriceChangePct, _ = strconv.ParseFloat(raw.PriceChangePercent, 64)
	ticker.LowPrice, _ = strconv.ParseFloat(raw.LowPrice, 64)
	ticker.HighPrice, _ = strconv.ParseFloat(raw.HighPrice, 64)
	ticker.BaseVolume, _ = strconv.ParseFloat(raw.TotalBaseVolume, 64)
	ticker.QuoteVolume, _ = strconv.ParseFloat(raw.TotalQuoteVolume, 64)
	ticker.Ask, _ = strconv.ParseFloat(raw.Ask, 64)
	ticker.Bid, _ = strconv.ParseFloat(raw.Bid, 64)

	return ticker
}

type RawStreamAllMarketTicker24 struct {
	Stream  string        `json:"stream"`
	Tickers []RawTicker24 `json:"data"`
}

// Useful for decoding AggTrade's when they are coming from a combined stream,
// and you know all message types will be an "aggTrade".
type RawStreamAggTrade struct {
	Stream   string      `json:"stream"`
	AggTrade RawAggTrade `json:"data"`
}

type RawStreamData struct {
	Stream string      `json:"stream"`
	Data   interface{} `json:"data"`
}

type RawStreamMessage struct {
	Stream string `json:"stream"`

	// Data for !ticker@arr messages.
	Tickers []RawTicker24

	// Data for <symbol>@aggTrade messages.
	AggTrade RawAggTrade

	// For a stream that is unknown, decode the data into an interface{}.
	UnknownData interface{}

	RawData []byte
}

func (r *RawStreamMessage) UnmarshalJSON(b []byte) error {
	r.RawData = b
	prefix := string(b[0:40])
	if strings.HasPrefix(prefix, `{"stream":"!ticker@arr"`) {
		var message RawStreamAllMarketTicker24
		if err := json.Unmarshal(b, &message); err != nil {
			return err
		}
		r.Stream = message.Stream
		r.Tickers = message.Tickers
		return nil
	} else if strings.Index(prefix, "@aggTrade") > -1 {
		var message RawStreamAggTrade
		if err := json.Unmarshal(b, &message); err != nil {
			return err
		}
		r.Stream = message.Stream
		r.AggTrade = message.AggTrade
		return nil
	} else if strings.HasPrefix(prefix, `{"stream":"`) {
		var message RawStreamData
		if err := json.Unmarshal(b, &message); err != nil {
			return err
		}
		r.Stream = message.Stream
		r.UnknownData = message.Data
		return nil
	}
	return fmt.Errorf("not part of a multi-stream")
}

func DecodeRawStreamMessage(b []byte) (RawStreamMessage, error) {
	var message RawStreamMessage
	err := json.Unmarshal(b, &message)
	return message, err
}
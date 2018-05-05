package binance

type AccountInfoBalance struct {
	Asset  string  `json:"a"`
	Free   float64 `json:"f,string"`
	Locked float64 `json:"l,string"`
}

type OutboundAccountInfo struct {
	EventType             string               `json:"e"`
	EventTimeMillis       int64                `json:"E"`
	MakerCommissionRate   int64                `json:"m"`
	TakerCommissionRate   int64                `json:"t"`
	BuyerCommissionRate   int64                `json:"b"`
	SellerCommissionRate  int64                `json:"s"`
	CanTrade              bool                 `json:"T"`
	CanWithdraw           bool                 `json:"W"`
	CanDeposit            bool                 `json:"D"`
	LastAccountUpdateTime int64                `json:"u"`
	Balances              []AccountInfoBalance `json:"B"`
}

type OrderUpdate struct {
	EventType                string  `json:"e"`
	EventTimeMillis          int64   `json:"E"`
	Symbol                   string  `json:"s"`
	ClientOrderID            string  `json:"c"`
	Side                     string  `json:"S"`
	OrderType                string  `json:"o"`
	TimeInForce              string  `json:"f"`
	Quantity                 float64 `json:"q,string"`
	Price                    float64 `json:"p,string"`
	StopPrice                float64 `json:"P,string"`
	IcebergQuantity          float64 `json:"F,string"`
	OriginalClientOrderID    string  `json:"C"`
	CurrentExecutionType     string  `json:"x"`
	CurrentOrderStatus       string  `json:"X"`
	OrderRejectReason        string  `json:"r"`
	OrderID                  int64   `json:"i"`
	LastExecutedQuantity     float64 `json:"l,string"`
	CumulativeFilledQuantity float64 `json:"z,string"`
	LastExecutedPrice        float64 `json:"L,string"`
	CommissionAmount         string  `json:"n"`
	CommissionAsset          string  `json:"N"`
	TransactionTimeMillis    int64   `json:"T"`
	TradeID                  int64   `json:"t"`
	IsWorking                bool    `json:"w"`
	IsMaker                  bool    `json:"m"`

	// We have to include this otherwise it will try to decode the "O" field
	// as the order type, which is "o".
	Ignore0 int64 `json:"O,omit"`
}

func OpenUserStream(restClient *RestClient) (*StreamClient, error) {
	listenKey, err := restClient.GetUserDataStream()
	if err != nil {
		return nil, err
	}

	streamClient := NewStreamClient()
	if err := streamClient.ConnectSingle(listenKey); err != nil {
		return nil, err
	}

	return streamClient, nil
}

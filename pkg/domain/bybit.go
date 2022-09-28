package domain

import "time"

type BybitPerpResponseOrder struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		OrderID      string `json:"orderId"`
		OrderLinkID  string `json:"orderLinkId"`
		Mmp          bool   `json:"mmp"`
		Symbol       string `json:"symbol"`
		OrderType    string `json:"orderType"`
		Side         string `json:"side"`
		OrderQty     string `json:"orderQty"`
		OrderPrice   string `json:"orderPrice"`
		Iv           string `json:"iv"`
		TimeInForce  string `json:"timeInForce"`
		OrderStatus  string `json:"orderStatus"`
		CreatedAt    string `json:"createdAt"`
		BasePrice    string `json:"basePrice"`
		TriggerPrice string `json:"triggerPrice"`
		TakeProfit   string `json:"takeProfit"`
		StopLoss     string `json:"stopLoss"`
		SlTriggerBy  string `json:"slTriggerBy"`
		TpTriggerBy  string `json:"tpTriggerBy"`
	} `json:"result"`
}

type BybitDerivTicker struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Category string `json:"category"`
		List     []struct {
			Symbol                 string `json:"symbol"`
			BidPrice               string `json:"bidPrice"`
			AskPrice               string `json:"askPrice"`
			LastPrice              string `json:"lastPrice"`
			LastTickDirection      string `json:"lastTickDirection"`
			PrevPrice24H           string `json:"prevPrice24h"`
			Price24HPcnt           string `json:"price24hPcnt"`
			HighPrice24H           string `json:"highPrice24h"`
			LowPrice24H            string `json:"lowPrice24h"`
			PrevPrice1H            string `json:"prevPrice1h"`
			MarkPrice              string `json:"markPrice"`
			IndexPrice             string `json:"indexPrice"`
			OpenInterest           string `json:"openInterest"`
			Turnover24H            string `json:"turnover24h"`
			Volume24H              string `json:"volume24h"`
			FundingRate            string `json:"fundingRate"`
			NextFundingTime        string `json:"nextFundingTime"`
			PredictedDeliveryPrice string `json:"predictedDeliveryPrice"`
			BasisRate              string `json:"basisRate"`
			DeliveryFeeRate        string `json:"deliveryFeeRate"`
			DeliveryTime           string `json:"deliveryTime"`
		} `json:"list"`
	} `json:"result"`
	RetExtInfo interface{} `json:"retExtInfo"`
	Time       int64       `json:"time"`
}

type BybitPerpTicker struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
	Result  struct {
		Symbol                 string    `json:"symbol"`
		Bid                    string    `json:"bid"`
		BidIv                  string    `json:"bidIv"`
		BidSize                string    `json:"bidSize"`
		Ask                    string    `json:"ask"`
		AskIv                  string    `json:"askIv"`
		AskSize                string    `json:"askSize"`
		LastPrice              string    `json:"lastPrice"`
		OpenInterest           string    `json:"openInterest"`
		IndexPrice             string    `json:"indexPrice"`
		MarkPrice              string    `json:"markPrice"`
		MarkPriceIv            string    `json:"markPriceIv"`
		Change24H              string    `json:"change24h"`
		High24H                string    `json:"high24h"`
		Low24H                 string    `json:"low24h"`
		Volume24H              string    `json:"volume24h"`
		Turnover24H            string    `json:"turnover24h"`
		TotalVolume            string    `json:"totalVolume"`
		TotalTurnover          string    `json:"totalTurnover"`
		FundingRate            string    `json:"fundingRate"`
		PredictedFundingRate   string    `json:"predictedFundingRate"`
		NextFundingTime        time.Time `json:"nextFundingTime"`
		CountdownHour          string    `json:"countdownHour"`
		PredictedDeliveryPrice string    `json:"predictedDeliveryPrice"`
		UnderlyingPrice        string    `json:"underlyingPrice"`
		Delta                  string    `json:"delta"`
		Gamma                  string    `json:"gamma"`
		Vega                   string    `json:"vega"`
		Theta                  string    `json:"theta"`
	} `json:"result"`
}

type BybitUSDCPositions struct {
	Result struct {
		Cursor          string `json:"cursor"`
		ResultTotalSize int    `json:"resultTotalSize"`
		DataList        []struct {
			Symbol              string `json:"symbol"`
			Leverage            string `json:"leverage"`
			OccClosingFee       string `json:"occClosingFee"`
			LiqPrice            string `json:"liqPrice"`
			PositionValue       string `json:"positionValue"`
			TakeProfit          string `json:"takeProfit"`
			RiskID              string `json:"riskId"`
			TrailingStop        string `json:"trailingStop"`
			UnrealisedPnl       string `json:"unrealisedPnl"`
			CreatedAt           string `json:"createdAt"`
			MarkPrice           string `json:"markPrice"`
			CumRealisedPnl      string `json:"cumRealisedPnl"`
			PositionMM          string `json:"positionMM"`
			PositionIM          string `json:"positionIM"`
			UpdatedAt           string `json:"updatedAt"`
			TpSLMode            string `json:"tpSLMode"`
			Side                string `json:"side"`
			BustPrice           string `json:"bustPrice"`
			DeleverageIndicator int    `json:"deleverageIndicator"`
			EntryPrice          string `json:"entryPrice"`
			Size                string `json:"size"`
			SessionRPL          string `json:"sessionRPL"`
			PositionStatus      string `json:"positionStatus"`
			SessionUPL          string `json:"sessionUPL"`
			StopLoss            string `json:"stopLoss"`
			OrderMargin         string `json:"orderMargin"`
			SessionAvgPrice     string `json:"sessionAvgPrice"`
		} `json:"dataList"`
	} `json:"result"`
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

type BybitLinearClosedPnLResponse struct {
	RetCode int    `json:"ret_code"`
	RetMsg  string `json:"ret_msg"`
	Result  struct {
		Data        []BybitLinearClosedPnL `json:"data"`
		CurrentPage int                    `json:"current_page"`
	} `json:"result"`
	ExtCode          string `json:"ext_code"`
	ExtInfo          string `json:"ext_info"`
	TimeNow          string `json:"time_now"`
	RateLimitStatus  int    `json:"rate_limit_status"`
	RateLimitResetMs int64  `json:"rate_limit_reset_ms"`
	RateLimit        int    `json:"rate_limit"`
}

type BybitLinearClosedPnL struct {
	ID            int     `json:"id"`
	UserID        int     `json:"user_id"`
	Symbol        string  `json:"symbol"`
	OrderID       string  `json:"order_id"`
	Side          string  `json:"side"`
	Qty           float64 `json:"qty"`
	OrderPrice    float64 `json:"order_price"`
	OrderType     string  `json:"order_type"`
	ExecType      string  `json:"exec_type"`
	ClosedSize    float64 `json:"closed_size"`
	CumEntryValue float64 `json:"cum_entry_value"`
	AvgEntryPrice float64 `json:"avg_entry_price"`
	CumExitValue  float64 `json:"cum_exit_value"`
	AvgExitPrice  float64 `json:"avg_exit_price"`
	ClosedPnl     float64 `json:"closed_pnl"`
	FillCount     int     `json:"fill_count"`
	Leverage      int     `json:"leverage"`
	CreatedAt     int64   `json:"created_at"`
}

type BybitWallet struct {
	Result struct {
		WalletBalance    string `json:"walletBalance"`
		AccountMM        string `json:"accountMM"`
		Bonus            string `json:"bonus"`
		AccountIM        string `json:"accountIM"`
		TotalSessionRPL  string `json:"totalSessionRPL"`
		Equity           string `json:"equity"`
		TotalRPL         string `json:"totalRPL"`
		MarginBalance    string `json:"marginBalance"`
		AvailableBalance string `json:"availableBalance"`
		TotalSessionUPL  string `json:"totalSessionUPL"`
	} `json:"result"`
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

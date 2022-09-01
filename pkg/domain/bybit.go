package domain

import "time"

type PerpResponseOrder struct {
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

type DerivTicker struct {
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

type PerpTicker struct {
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

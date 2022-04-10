package domain

type TV struct {
	IsTestNet    bool    `json:"is_test_net"`
	APIKey       string  `json:"api_key" binding:"required"`
	APISecretKey string  `json:"api_secret_key" binding:"required"`
	Order        TVOrder `json:"order"`
}

type TVOrder struct {
	Type   string  `json:"type" binding:"required"`   // "Market" or "Limit"
	Symbol string  `json:"symbol" binding:"required"` // eg: BTCUSDT
	Side   string  `json:"side" binding:"required"`   // "Buy" or "Sell"
	Price  float64 `json:"price"`                     // Set 0 if order_type is Market
	QTY    float64 `json:"qty" binding:"required"`
	TP     float64 `json:"tp"`
	SL     float64 `json:"sl"`
}

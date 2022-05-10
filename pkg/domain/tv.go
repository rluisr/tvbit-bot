/*
tvbit-bot
Copyright (C) 2022  rluisr(Takuya Hasegawa)

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package domain

import "github.com/frankrap/bybit-api/rest"

type TV struct {
	IsTestNet    bool    `json:"is_test_net"`
	APIKey       string  `json:"api_key" binding:"required"`
	APISecretKey string  `json:"api_secret_key" binding:"required"`
	Order        TVOrder `json:"order"`
}

type TVOrder struct {
	Type   string      `json:"type" binding:"required"`   // "Market" or "Limit"
	Symbol string      `json:"symbol" binding:"required"` // eg: BTCUSDT
	Side   string      `json:"side" binding:"required"`   // "Buy" or "Sell"
	Price  float64     `json:"price"`                     // Set 0 if order_type is Market
	QTY    float64     `json:"qty" binding:"required"`
	TP     interface{} `json:"tp"`
	SL     interface{} `json:"sl"`
}

type TVOrderResponse struct {
	Success bool        `json:"successful" binding:"required"`
	Reason  string      `json:"reason,omitempty" binding:"required"`
	Order   *rest.Order `json:"order"`
}

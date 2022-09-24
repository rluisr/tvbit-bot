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

import (
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type TV struct {
	IsTestNet    bool    `json:"is_test_net"`
	APIKey       string  `json:"api_key" binding:"required"`
	APISecretKey string  `json:"api_secret_key" binding:"required"`
	Order        TVOrder `json:"order"`
}

type TVOrder struct {
	gorm.Model
	CEX        string          `gorm:"type:varchar(255);not null" json:"-"`
	SettingID  uint64          `gorm:"type:uint;not null" json:"-"`
	OrderID    string          `gorm:"type:varchar(255);uniqueIndex:order_id;not null"`
	Type       string          `gorm:"type:varchar(255)" json:"type" binding:"required"`   // "Market" or "Limit"
	Symbol     string          `gorm:"type:varchar(255)" json:"symbol" binding:"required"` // eg: BTCUSDT
	Side       string          `gorm:"type:varchar(255)" json:"side" binding:"required"`   // "Buy" or "Sell"
	Price      float64         `gorm:"-" json:"price"`                                     // Set 0 if order_type is Market
	EntryPrice decimal.Decimal `gorm:"type:decimal(10,4)" json:"-"`
	QTY        float64         `gorm:"type:float" json:"qty" binding:"required"`
	TP         interface{}     `gorm:"type:float" json:"tp"`
	SL         interface{}     `gorm:"type:float" json:"sl"`
}

type Tabler interface {
	TableName() string
}

func (TVOrder) TableName() string {
	return "orders"
}

type TVOrderResponse struct {
	Success bool     `json:"successful" binding:"required"`
	Reason  string   `json:"reason,omitempty" binding:"required"`
	Order   *TVOrder `json:"order"`
}

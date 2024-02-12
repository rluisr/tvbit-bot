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

type Order struct {
	gorm.Model
	Name       string          `gorm:"type:varchar(255);default:null" json:"name"` // alert name, description or something
	CEX        string          `gorm:"type:varchar(255);not null" json:"-"`
	OrderID    string          `gorm:"type:varchar(255);uniqueIndex:order_id;not null"`
	Type       string          `gorm:"type:varchar(255)" json:"type" binding:"required"`   // "Market" or "Limit"
	Symbol     string          `gorm:"type:varchar(255)" json:"symbol" binding:"required"` // eg: BTCUSDT
	Side       string          `gorm:"type:varchar(255)" json:"side" binding:"required"`   // "Buy" or "Sell"
	Price      string          `gorm:"-" json:"price"`                                     // Set 0 if order_type is Market
	EntryPrice decimal.Decimal `gorm:"type:decimal(10,4)" json:"-"`
	QTY        string          `gorm:"type:float" json:"qty" binding:"required"`
	TP         string          `gorm:"type:float" json:"tp"`
	SL         string          `gorm:"type:float" json:"sl"`
	PL         decimal.Decimal `gorm:"type:decimal(10,4)" json:"-"`
}

type TVOrderResponse struct {
	Success bool   `json:"successful" binding:"required"`
	Reason  string `json:"reason,omitempty" binding:"required"`
	Order   *Order `json:"order"`
}

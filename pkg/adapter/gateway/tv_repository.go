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

package gateway

import (
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"gorm.io/gorm"
)

type (
	TVRepository struct {
		RWDB *gorm.DB
		RODB *gorm.DB
	}
)

func (r *TVRepository) UpdateOrder(order *domain.Order) error {
	return r.RWDB.Save(&order).Error
}

func (r *TVRepository) SaveOrder(order *domain.Order) error {
	order.CEX = "bybit"
	return r.RWDB.Save(order).Error
}

func (r *TVRepository) SaveWalletHistories(histories []domain.WalletHistory) error {
	return r.RWDB.Create(&histories).Error
}

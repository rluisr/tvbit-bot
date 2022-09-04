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

package interfaces

import (
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
)

type TVRepository interface {
	IsOK(domain.TV) (bool, error)
	SaveOrder(domain.TV, *domain.TVOrder) error
	GetSettings() ([]domain.Setting, error)
	SaveWalletHistories([]domain.WalletHistory) error
}

type BybitRepository interface {
	Set(domain.TV)
	CreateOrder(domain.TV) (string, error)
	CalculateTPSL(domain.TV, interface{}, string) (float64, error)
	FetchOrder(req *domain.TV, orderID string) error
	GetWalletInfoUSDC() (*domain.BybitWallet, error)
	GetWalletInfoDeriv() (*rest.Balance, error)
}

type SettingRepository interface {
	Get(domain.Setting) (domain.Setting, error)
	GetByID(uint64) (domain.Setting, error)
	Set(domain.Setting) (domain.Setting, error)
}

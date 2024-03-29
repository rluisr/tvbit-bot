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
	"github.com/hirokisan/bybit/v2"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/external/logging"
)

type TVRepository interface {
	Logging() *logging.Logging
	SaveOrder(*domain.Order) error
	SaveClosedPnL(closedPnL []*domain.ClosedPnL) error
	UpdateOrder(*domain.Order) error
	GetUniqueSymbol() ([]string, error)
	TruncateClosedPnL() error
}

type BybitRepository interface {
	CreateOrder(*domain.Order) error
	CancelOrder(req *domain.Order) error
	FetchOpenOrder(*domain.Order) error
	GetOpenOrders() (*bybit.V5GetOrdersResponse, error)
	GetClosedPNL(param bybit.V5GetClosedPnLParam) (*bybit.V5GetClosedPnLResponse, error)
	CalculateTPSL(*domain.Order) error
	GetWalletBalance() (float64, error)
}

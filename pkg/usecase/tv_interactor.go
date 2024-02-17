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

package usecase

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hirokisan/bybit/v2"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
	"github.com/rluisr/tvbit-bot/utils"
	"github.com/shopspring/decimal"
)

type TVInteractor struct {
	TVRepository    interfaces.TVRepository
	BybitRepository interfaces.BybitRepository
}

func (i *TVInteractor) CreateOrder(c *gin.Context) (domain.TVOrderResponse, error) {
	var req domain.Order
	err := c.ShouldBindJSON(&req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	err = i.BybitRepository.CalculateTPSL(&req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	err = i.BybitRepository.CreateOrder(&req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	err = i.BybitRepository.FetchOpenOrder(&req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   &req,
		}, err
	}

	err = i.TVRepository.SaveOrder(&req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   &req,
		}, err
	}

	i.TVRepository.Logging().Info(fmt.Sprintf("created order: %+v", req))

	return domain.TVOrderResponse{Success: true, Order: &req}, nil
}

func (i *TVInteractor) FetchPL() error {
	// TODO TRUNCATE せずに増分更新する
	err := i.TVRepository.TruncateClosedPnL()
	if err != nil {
		return err
	}

	symbols, err := i.TVRepository.GetUniqueSymbol()
	if err != nil {
		return err
	}

	var closedOrders []bybit.V5GetClosedPnLItem

	for _, symbol := range symbols {
		bSymbol := bybit.SymbolV5(symbol)
		limit := 100

		var temp []bybit.V5GetClosedPnLItem

		resp, gcPNLErr := i.BybitRepository.GetClosedPNL(bybit.V5GetClosedPnLParam{
			Category: bybit.CategoryV5Linear,
			Symbol:   &bSymbol,
			Limit:    &limit,
		})
		if gcPNLErr != nil {
			return gcPNLErr
		}
		temp = append(temp, resp.Result.List...)

		if resp.Result.NextPageCursor != "" {
			for {
				resp, err = i.BybitRepository.GetClosedPNL(bybit.V5GetClosedPnLParam{
					Category: bybit.CategoryV5Linear,
					Symbol:   &bSymbol,
					Limit:    &limit,
					Cursor:   &resp.Result.NextPageCursor,
				})
				if err != nil {
					return err
				}
				temp = append(temp, resp.Result.List...)

				if resp.Result.NextPageCursor == "" {
					break
				}
			}
		}

		closedOrders = append(closedOrders, temp...)
	}

	var closedPnLs []*domain.ClosedPnL
	for _, order := range closedOrders {
		orderPrice, fErr := decimal.NewFromString(order.OrderPrice)
		if fErr != nil {
			return fErr
		}

		closedSize, fErr := decimal.NewFromString(order.ClosedSize)
		if fErr != nil {
			return fErr
		}

		cumEntryValue, fErr := decimal.NewFromString(order.CumEntryValue)
		if fErr != nil {
			return fErr
		}

		avgEntryPrice, fErr := decimal.NewFromString(order.AvgEntryPrice)
		if fErr != nil {
			return fErr
		}

		cumExitValue, fErr := decimal.NewFromString(order.CumExitValue)
		if fErr != nil {
			return fErr
		}

		avgExitPrice, fErr := decimal.NewFromString(order.AvgExitPrice)
		if fErr != nil {
			return fErr
		}

		closedPnL, fErr := decimal.NewFromString(order.ClosedPnl)
		if fErr != nil {
			return fErr
		}

		createTime, fErr := utils.TimestampMSToTime(order.CreatedTime)
		if fErr != nil {
			return fErr
		}

		updateTime, fErr := utils.TimestampMSToTime(order.UpdatedTime)
		if fErr != nil {
			return fErr
		}

		a := &domain.ClosedPnL{
			OrderID:       order.OrderID,
			Symbol:        string(order.Symbol),
			Side:          string(order.Side),
			Qty:           order.Qty,
			OrderPrice:    orderPrice,
			ClosedSize:    closedSize,
			CumEntryValue: cumEntryValue,
			AvgEntryPrice: avgEntryPrice,
			CumExitValue:  cumExitValue,
			AvgExitPrice:  avgExitPrice,
			ClosedPnL:     closedPnL,
			CreatedAt:     createTime,
			UpdatedAt:     updateTime,
		}
		closedPnLs = append(closedPnLs, a)
	}

	return i.TVRepository.SaveClosedPnL(closedPnLs)
}

// InventoryCheck は約定していない注文で、n分経過したものをキャンセルする
func (i *TVInteractor) InventoryCheck(cancelAfter time.Duration) error {
	orders, err := i.BybitRepository.GetOpenOrders()
	if err != nil {
		return err
	}

	for _, order := range orders.Result.List {
		if order.OrderStatus == bybit.OrderStatusUntriggered {
			createdAt, cErr := utils.TimestampMSToTime(order.CreatedTime)
			if cErr != nil {
				return cErr
			}

			if createdAt.Add(cancelAfter).Before(time.Now()) {
				err = i.BybitRepository.CancelOrder(&domain.Order{
					OrderID: order.OrderID,
					Symbol:  string(order.Symbol),
				})
				if err != nil {
					return err
				}

				i.TVRepository.Logging().Info(fmt.Sprintf("canceled order: %+v", order))
			}
		}
	}

	return nil
}

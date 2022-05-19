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
	"errors"
	"fmt"
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/utils"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

type (
	TVRepository struct {
		RWDB *gorm.DB
		RODB *gorm.DB
	}
)

func (r *TVRepository) CreateOrder(req domain.TV, bybitClient *rest.ByBit) (*rest.Order, error) {
	isOK, err := r.isOK(req)
	if err != nil {
		return nil, err
	}

	if isOK {
		_, _, order, err := bybitClient.LinearCreateOrder(req.Order.Side, req.Order.Type, req.Order.Price, req.Order.QTY, "GoodTillCancel", req.Order.TP.(float64), req.Order.SL.(float64), false, false, "", req.Order.Symbol)
		return &order, err
	}

	return nil, nil
}

func (r *TVRepository) SaveOrder(req domain.TV, order *rest.Order) error {
	f64Price, _ := order.Price.Float64()
	f64Qty, _ := order.Qty.Float64()

	orderHistory := domain.TVOrder{
		Type:   order.OrderType,
		Symbol: order.Symbol,
		Side:   order.Side,
		Price:  f64Price,
		QTY:    f64Qty,
		TP:     req.Order.TP,
		SL:     req.Order.SL,
	}

	return r.RWDB.Save(&orderHistory).Error
}

// isOK: current time is between "start_time" and "stop_time"
func (r *TVRepository) isOK(req domain.TV) (bool, error) {
	var setting domain.Setting
	err := r.RODB.Where("api_key = ? and api_secret_key = ?", req.APIKey, req.APISecretKey).Take(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}

	utc := utils.UTC()
	now := time.Now().In(utc)

	hour := now.Hour()
	minute := now.Minute()
	n := fmt.Sprintf("%d%d", hour, minute)
	currentTime, err := strconv.Atoi(n)
	if err != nil {
		return false, err
	}

	userStartTimeStr := strings.Replace(setting.StartTime, ":", "", 1)
	userStopTimeStr := strings.Replace(setting.StopTime, ":", "", 1)

	userStartTime, err := strconv.Atoi(userStartTimeStr)
	if err != nil {
		return false, err
	}
	userStopTime, err := strconv.Atoi(userStopTimeStr)
	if err != nil {
		return false, err
	}

	if userStartTime < currentTime && currentTime < userStopTime {
		return true, nil
	}

	return false, nil
}

func (r *TVRepository) CalculateTPSL(req domain.TV, bybitClient *rest.ByBit, value interface{}, isType string) (float64, error) {
	var err error

	str, isOK := value.(string)
	if !isOK {
		return 0, fmt.Errorf("failed to assertion to string: %v", value)
	}

	if str == "" || str == "0" {
		return 0, nil
	}

	if strings.Contains(str, "%") {
		var price float64

		if req.Order.Price == 0 {
			price, err = r.getCurrentPrice(req.Order.Symbol, bybitClient)
			if err != nil {
				return 0, err
			}
		} else {
			price = req.Order.Price
		}

		tpslStr := strings.Replace(value.(string), "%", "", 1)
		tpsl, err := strconv.ParseFloat(tpslStr, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse to float64 %s", tpslStr)
		}
		tpsl = tpsl * 0.1

		switch req.Order.Side {
		case "Buy":
			if isType == "TP" {
				return (price * tpsl * req.Order.QTY) + price, nil
			}
			return math.Abs((price * tpsl * req.Order.QTY) - price), nil
		case "Sell":
			if isType == "SL" {
				return (price * tpsl * req.Order.QTY) + price, nil
			}
			return math.Abs((price * tpsl * req.Order.QTY) - price), nil
		default:
			return 0, fmt.Errorf("unknown request: %v / isType: %s", req, isType)
		}
	}

	f64, isOK := value.(float64)
	if !isOK {
		return 0, fmt.Errorf("failed to assertion to float64: %v", value)
	}

	return f64, nil

}

// getCurrentPrice returns close price one minute ago
func (r *TVRepository) getCurrentPrice(symbol string, bybitClient *rest.ByBit) (float64, error) {
	_, _, resp, err := bybitClient.LinearGetKLine(symbol, "1", time.Now().Unix()-60, 1)
	if err != nil {
		return 0, fmt.Errorf("failed to get current price err: %w", err)
	}

	if len(resp) == 0 {
		return 0, fmt.Errorf("failed to get current price. invalid query")
	}

	return resp[0].Close, nil
}

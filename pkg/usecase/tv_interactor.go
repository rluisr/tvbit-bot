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

	"github.com/frankrap/bybit-api/rest"

	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
	"golang.org/x/sync/errgroup"
)

type TVInteractor struct {
	TVRepository    interfaces.TVRepository
	BybitRepository interfaces.BybitRepository
}

func (i *TVInteractor) CreateOrder(req domain.TV) (domain.TVOrderResponse, error) {
	var (
		err       error
		setting   *domain.Setting
		positions *[]rest.LinearPosition
	)

	eg := errgroup.Group{}
	eg.Go(func() error {
		setting, err = i.TVRepository.GetSetting(req.APIKey, req.APISecretKey)
		return err
	})
	eg.Go(func() error {
		positions, err = i.BybitRepository.GetPositions(req.Order.Symbol)
		return err
	})
	eg.Go(func() error {
		req.Order.TP, err = i.BybitRepository.CalculateTPSL(req, req.Order.TP, "TP")
		return err
	})
	eg.Go(func() error {
		req.Order.SL, err = i.BybitRepository.CalculateTPSL(req, req.Order.SL, "SL")
		return err
	})
	err = eg.Wait()
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	activeOrder := i.BybitRepository.GetActiveOrderCount(&req, positions)
	if setting.MaxPosition.Valid && setting.MaxPosition.Int32 <= int32(activeOrder) {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  fmt.Errorf("your setting max_position is %d and a count of current active position is %d. your order is cancelled", setting.MaxPosition.Int32, activeOrder).Error(),
			Order:   nil,
		}, nil
	}

	orderID, err := i.BybitRepository.CreateOrder(req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	// Wait for ByBit server side reflection
	time.Sleep(5 * time.Second)

	err = i.BybitRepository.FetchOrder(&req, orderID)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   &req.Order,
		}, err
	}

	err = i.TVRepository.SaveOrder(req, &req.Order)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   &req.Order,
		}, err
	}

	return domain.TVOrderResponse{Success: true, Order: &req.Order}, nil
}

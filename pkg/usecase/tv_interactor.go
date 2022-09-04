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
	"time"

	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
)

type TVInteractor struct {
	TVRepository    interfaces.TVRepository
	BybitRepository interfaces.BybitRepository
}

func (i *TVInteractor) CreateOrder(req domain.TV) (domain.TVOrderResponse, error) {
	var err error

	req.Order.TP, err = i.BybitRepository.CalculateTPSL(req, req.Order.TP, "TP")
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	req.Order.SL, err = i.BybitRepository.CalculateTPSL(req, req.Order.SL, "SL")
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	isOK, err := i.TVRepository.IsOK(req)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}
	if !isOK {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
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

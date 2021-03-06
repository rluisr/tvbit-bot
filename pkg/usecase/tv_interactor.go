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
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
	"log"
)

type TVInteractor struct {
	TVRepository interfaces.TVRepository
}

func (i *TVInteractor) CreateOrder(req domain.TV, bybitClient *rest.ByBit) (domain.TVOrderResponse, error) {
	var err error

	req.Order.TP, err = i.TVRepository.CalculateTPSL(req, bybitClient, req.Order.TP, "TP")
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	req.Order.SL, err = i.TVRepository.CalculateTPSL(req, bybitClient, req.Order.SL, "SL")
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	order, err := i.TVRepository.CreateOrder(req, bybitClient)
	if err != nil {
		return domain.TVOrderResponse{
			Success: false,
			Reason:  err.Error(),
			Order:   nil,
		}, err
	}

	if order == nil {
		return domain.TVOrderResponse{Success: true, Reason: "order is cancelled by \"start_time\", \"stop_time\" setting", Order: nil}, nil
	}

	err = i.TVRepository.SaveOrder(req, order)
	if err != nil {
		// order is created, do not return err here.
		log.Printf("a order is created but failed to save order history to MySQL err: %s\n", err.Error())
	}

	return domain.TVOrderResponse{Success: true, Order: order}, nil
}

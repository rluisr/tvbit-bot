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

	"github.com/gin-gonic/gin"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
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

	err = i.BybitRepository.FetchOrder(&req)
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

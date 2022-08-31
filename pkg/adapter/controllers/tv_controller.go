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

package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/rluisr/tvbit-bot/pkg/adapter/gateway"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/external/bybit"
	"github.com/rluisr/tvbit-bot/pkg/usecase"
	"gorm.io/gorm"
	"net/http"
)

type TVController struct {
	Interactor usecase.TVInteractor
}

func NewTVController(rwDB, roDB *gorm.DB) *TVController {
	return &TVController{
		Interactor: usecase.TVInteractor{
			TVRepository: &gateway.TVRepository{
				RWDB: rwDB,
				RODB: roDB,
			},
		},
	}
}

func (controller *TVController) Handle(c *gin.Context) {
	var req domain.TV
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(http.StatusBadRequest, err))
		return
	}

	bybitClient := bybit.Init(req)

	order, err := controller.Interactor.CreateOrder(req, bybitClient)
	if err != nil {
		c.JSON(http.StatusInternalServerError, NewError(http.StatusInternalServerError, err))
		return
	}

	c.JSON(200, order)
}

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
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rluisr/tvbit-bot/pkg/adapter/gateway"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase"
	"gorm.io/gorm"
	"net/http"
)

type SettingController struct {
	Interactor usecase.SettingInteractor
}

func NewSettingController(rwDB, roDB *gorm.DB) *SettingController {
	return &SettingController{
		Interactor: usecase.SettingInteractor{
			SettingRepository: &gateway.SettingRepository{
				RWDB: rwDB,
				RODB: roDB,
			},
		},
	}
}

func (controller *SettingController) Get(c *gin.Context) {
	var setting domain.Setting
	err := c.ShouldBindJSON(&setting)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("bind error. \"api_key\" and \"api_secret_key\" is required as a json. err: %s", err.Error()))
		return
	}

	resp, err := controller.Interactor.Get(setting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Get error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (controller *SettingController) Set(c *gin.Context) {
	var setting domain.Setting
	err := c.ShouldBindJSON(&setting)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("bind error. invalid requestest json body. err: %s", err.Error()))
		return
	}

	resp, err := controller.Interactor.Set(setting)
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("Set error: %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, resp)
}

package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rluisr/tvbit-bot/pkg/adapter/gateway"
	"github.com/rluisr/tvbit-bot/pkg/adapter/interfaces"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/external/bybit"
	"github.com/rluisr/tvbit-bot/pkg/usecase"
	"net/http"
)

type TVController struct {
	Interactor usecase.TVInteractor
}

func NewTVController(logger interfaces.Logger) *TVController {
	return &TVController{
		Interactor: usecase.TVInteractor{
			TVRepository: &gateway.TVRepository{},
			Logger:       logger,
		},
	}
}

func (controller *TVController) Handle(c *gin.Context) {
	var req domain.TV
	err := c.BindJSON(&req)
	if err != nil {
		controller.Interactor.Logger.Log(fmt.Errorf("tv_controller: cannot handle. err: %w", err))
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("bind error: %s", err.Error()))
		return
	}

	bybitClient := bybit.Init(req)

	order, err := controller.Interactor.CreateOrder(req, bybitClient)
	if err != nil {
		controller.Interactor.Logger.Log(fmt.Errorf("tv_controller: cannot handle. err: %w", err))
		c.JSON(http.StatusInternalServerError, NewError(http.StatusInternalServerError, err.Error()))
		return
	}

	c.JSON(200, order)
}

package gateway

import (
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
)

type (
	TVRepository struct {
	}
)

func (r *TVRepository) CreateOrder(req domain.TV, bybitClient *rest.ByBit) (rest.Order, error) {
	_, _, order, err := bybitClient.LinearCreateOrder(req.Order.Side, req.Order.Type, req.Order.Price, req.Order.QTY, "GoodTillCancel", req.Order.TP, req.Order.SL, false, false, "", req.Order.Symbol)
	return order, err
}

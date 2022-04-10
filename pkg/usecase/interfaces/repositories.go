package interfaces

import (
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
)

type TVRepository interface {
	CreateOrder(domain.TV, *rest.ByBit) (rest.Order, error)
}

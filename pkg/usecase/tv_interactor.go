package usecase

import (
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
)

type TVInteractor struct {
	TVRepository interfaces.TVRepository
	Logger       interfaces.Logger
}

func (i *TVInteractor) CreateOrder(req domain.TV, bybitClient *rest.ByBit) (rest.Order, error) {
	return i.TVRepository.CreateOrder(req, bybitClient)
}

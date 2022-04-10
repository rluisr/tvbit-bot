package bybit

import (
	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
)

func Init(req domain.TV) *rest.ByBit {
	baseURL := "https://api.bybit.com/"
	if req.IsTestNet {
		baseURL = "https://api-testnet.bybit.com/"
	}

	return rest.New(nil, baseURL, req.APIKey, req.APISecretKey, false)
}

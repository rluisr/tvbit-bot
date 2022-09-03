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

package bybit

import (
	"net"
	"net/http"
	"time"

	"github.com/frankrap/bybit-api/rest"
	"github.com/rluisr/tvbit-bot/pkg/domain"
)

var BaseURL = "https://api.bybit.com/"

func Init(req domain.TV) (*rest.ByBit, string) {
	if req.IsTestNet {
		BaseURL = "https://api-testnet.bybit.com/"
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
			MaxIdleConns:          128,
			MaxIdleConnsPerHost:   100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		Timeout: 60 * time.Second,
	}

	return rest.New(httpClient, BaseURL, req.APIKey, req.APISecretKey, false), BaseURL
}

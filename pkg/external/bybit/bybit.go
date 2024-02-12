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
	"fmt"
	"net/http"

	"github.com/hirokisan/bybit/v2"
)

func Init(httpClient *http.Client) *bybit.Client {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("bybit.NewConfig err: %w", err))
	}

	var client *bybit.Client

	if config.IsTestnet {
		client = bybit.NewTestClient().WithAuth(config.APIKey, config.APISecret).WithHTTPClient(httpClient)
	} else {
		client = bybit.NewClient().WithAuth(config.APIKey, config.APISecret).WithHTTPClient(httpClient)
	}

	return client
}

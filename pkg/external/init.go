/*
 *
 * tvbit-bot
 * Copyright (C) 2022  rluisr(Takuya Hasegawa)
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 * /
 */

package external

import (
	"github.com/rluisr/tvbit-bot/pkg/adapter/controllers"
	tvbitBybit "github.com/rluisr/tvbit-bot/pkg/external/bybit"
	"github.com/rluisr/tvbit-bot/pkg/external/mysql"
)

var (
	tvController *controllers.TVController
)

func Init() (err error) {
	rwDB, roDB, err := mysql.Connect()
	if err != nil {
		return err
	}

	httpClient := NewHTTPClient()
	bybitClient := tvbitBybit.Init(httpClient)

	tvController = controllers.NewTVController(rwDB, roDB, bybitClient)

	return nil
}

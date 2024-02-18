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
	"os"
	"os/signal"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func Cron() {
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	_, err = s.NewJob(
		gocron.CronJob("0 * * * *", false),
		gocron.NewTask(func() {
			foErr := tvController.FetchOrder()
			if foErr != nil {
				tvController.Interactor.TVRepository.Logging().Error("FetchOrder", foErr.Error(), foErr)
				return
			}
			tvController.Interactor.TVRepository.Logging().Info("FetchOrder is done")
		}),
	)
	if err != nil {
		tvController.Interactor.TVRepository.Logging().Error("NewJob: FetchOrder", err.Error(), err)
	}

	_, err = s.NewJob(
		// KeepAlive を続けるために短い間隔で行う
		gocron.DurationJob(
			3*time.Second,
		),
		gocron.NewTask(func() {
			icErr := tvController.InventoryCheck(5 * time.Minute)
			if icErr != nil {
				tvController.Interactor.TVRepository.Logging().Error("InventoryCheck", icErr.Error(), icErr)
				return
			}
		}),
	)
	if err != nil {
		tvController.Interactor.TVRepository.Logging().Error("NewJob", err.Error(), err)
	}

	s.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	_ = s.Shutdown()
}

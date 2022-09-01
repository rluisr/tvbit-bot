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

package gateway

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/utils"
	"gorm.io/gorm"
)

type (
	TVRepository struct {
		RWDB *gorm.DB
		RODB *gorm.DB
	}
)

func (r *TVRepository) SaveOrder(order *domain.TVOrder) error {
	return r.RWDB.Save(order).Error
}

// IsOK current time is between "start_time" and "stop_time"
func (r *TVRepository) IsOK(req domain.TV) (bool, error) {
	var setting domain.Setting
	err := r.RODB.Where("api_key = ? and api_secret_key = ?", req.APIKey, req.APISecretKey).Take(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}

	if !setting.StartTime.Valid || !setting.StopTime.Valid {
		return true, nil
	}

	utc := utils.UTC()
	now := time.Now().In(utc)

	hour := now.Hour()
	minute := now.Minute()
	n := fmt.Sprintf("%d%d", hour, minute)
	currentTime, err := strconv.Atoi(n)
	if err != nil {
		return false, err
	}

	userStartTimeStr := strings.Replace(setting.StartTime.String, ":", "", 1)
	userStopTimeStr := strings.Replace(setting.StopTime.String, ":", "", 1)

	userStartTime, err := strconv.Atoi(userStartTimeStr)
	if err != nil {
		return false, err
	}
	userStopTime, err := strconv.Atoi(userStopTimeStr)
	if err != nil {
		return false, err
	}

	if userStartTime < currentTime && currentTime < userStopTime {
		return true, nil
	}

	return false, nil
}

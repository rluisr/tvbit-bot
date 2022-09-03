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

	"github.com/rluisr/tvbit-bot/pkg/domain"
	"gorm.io/gorm"
)

type (
	SettingRepository struct {
		RWDB *gorm.DB
		RODB *gorm.DB
	}
)

func (r *SettingRepository) Get(setting domain.Setting) (domain.Setting, error) {
	err := r.RODB.Where("api_key = ? AND api_secret_key = ?", setting.APIKey, setting.APISecretKey).Take(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Setting{}, fmt.Errorf("api_key and api_secret_key pair is not found. try PUT /setting")
	}

	return setting, err
}

func (r *SettingRepository) GetByID(id uint64) (domain.Setting, error) {
	var setting domain.Setting
	err := r.RODB.Where("id = ?", id).Take(&setting).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.Setting{}, fmt.Errorf("api_key and api_secret_key pair is not found. try PUT /setting")
	}

	return setting, err
}

func (r *SettingRepository) Set(setting domain.Setting) (domain.Setting, error) {
	// Todo check start_time and stop_time format is "xx:yy"
	return setting, r.RWDB.Where("api_key = ? AND api_secret_key = ?", setting.APIKey, setting.APISecretKey).Assign(setting).FirstOrCreate(&setting).Error
}

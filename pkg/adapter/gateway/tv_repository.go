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
	TVRepository struct {
		RWDB *gorm.DB
		RODB *gorm.DB
	}
)

func (r *TVRepository) SaveOrder(req domain.TV, order *domain.TVOrder) error {
	var setting domain.Setting
	err := r.RODB.Where("api_key = ? AND api_secret_key = ?", req.APIKey, req.APISecretKey).First(&setting).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed SaveOrder user is not found in setting table")
		} else {
			return err
		}
	}
	order.DEX = "bybit" // TODO we should support any other DEX
	order.SettingID = setting.ID
	return r.RWDB.Save(order).Error
}

func (r *TVRepository) GetSetting(apiKey, apiSecretKey string) (*domain.Setting, error) {
	var setting domain.Setting
	err := r.RODB.Where("api_key = ? AND api_secret_key = ?", apiKey, apiSecretKey).Take(&setting).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &setting, nil
}

func (r *TVRepository) GetSettings() ([]domain.Setting, error) {
	var settings []domain.Setting
	err := r.RODB.Find(&settings).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return settings, nil
}

func (r *TVRepository) SaveWalletHistories(histories []domain.WalletHistory) error {
	return r.RWDB.Create(&histories).Error
}

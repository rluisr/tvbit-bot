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

package usecase

import (
	"github.com/rluisr/tvbit-bot/pkg/domain"
	"github.com/rluisr/tvbit-bot/pkg/usecase/interfaces"
)

type SettingInteractor struct {
	SettingRepository interfaces.SettingRepository
}

func (i *SettingInteractor) Get(setting domain.Setting) (domain.Setting, error) {
	return i.SettingRepository.Get(setting)
}

func (i *SettingInteractor) Set(setting domain.Setting) (domain.Setting, error) {
	return i.SettingRepository.Set(setting)
}

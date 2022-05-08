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

package domain

type Setting struct {
	ID           uint64 `gorm:"primaryKey,autoIncrement"`
	APIKey       string `gorm:"type:varchar(255);unique;index:idx_api" json:"api_key" binding:"required"`
	APISecretKey string `gorm:"type:varchar(255);unique;index:idx_api" json:"api_secret_key" binding:"required"`
	StartTime    string `gorm:"type:char(5)" json:"start_time"` // eg "09:00"
	StopTime     string `gorm:"type:char(5)" json:"stop_time"`  // eg "21:00"
}

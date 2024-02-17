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

package utils

import (
	"strconv"
	"time"
)

func StringToFloat64(str string) float64 {
	f, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func TimestampMSToTime(timestampStr string) (time.Time, error) {
	timestampInt64, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return time.Now(), err
	}

	// Convert milliseconds to seconds
	timestampSec := timestampInt64 / 1000
	timestampNanoSec := (timestampInt64 % 1000) * int64(time.Millisecond)

	t := time.Unix(timestampSec, timestampNanoSec)

	return t, nil
}

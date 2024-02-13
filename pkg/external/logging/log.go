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

package logging

import (
	"fmt"
	"log"
	"os"

	"github.com/ic2hrmk/promtail"
)

const (
	envLocal = "local"
)

type Logging struct {
	Log promtail.Client
}

func New(source string) (*Logging, error) {
	config, err := NewConfig()
	if err != nil {
		panic(fmt.Errorf("logging.NewConfig err: %w", err))
	}

	identifiers := map[string]string{
		"source": source,
	}

	promtailClient, err := promtail.NewJSONv1Client(config.LokiURL, identifiers, promtail.WithBasicAuth(config.LokiBasicUser, config.LokiBasicPass))
	if err != nil {
		return nil, err
	}

	return &Logging{Log: promtailClient}, nil
}

func (l *Logging) Info(msg string) {
	if os.Getenv("SERVER_ENV") != envLocal {
		l.Log.Infof(msg)
	}

	log.Println(msg)
}

func (l *Logging) Error(funcName, msg string, err error) {
	body := fmt.Sprintf("func_name: %s, msg: %s, error: %s", funcName, msg, err.Error())

	if os.Getenv("SERVER_ENV") != envLocal {
		l.Log.Errorf(body)
	}

	log.Println(body)
}

func (l *Logging) Fatal(funcName, msg string, err error) {
	body := fmt.Sprintf("func_name: %s, msg: %s, error: %s", funcName, msg, err.Error())

	if os.Getenv("SERVER_ENV") != envLocal {
		l.Log.Errorf(body)
	}

	log.Fatal(body)
}

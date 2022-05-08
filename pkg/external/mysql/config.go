package mysql

import "github.com/Netflix/go-env"

type Config struct {
	MySQLHostRW string `env:"MYSQL_HOST_RW,required=true"`
	MySQLHostRO string `env:"MYSQL_HOST_RO,required=true"`
	MySQLUser   string `env:"MYSQL_USER,required=true"`
	MySQLPass   string `env:"MYSQL_PASS,required=true"`
	MySQLDBName string `env:"MYSQL_DB_NAME,required=true"`
}

func NewConfig() (*Config, error) {
	var config Config

	_, err := env.UnmarshalFromEnviron(&config)
	return &config, err
}

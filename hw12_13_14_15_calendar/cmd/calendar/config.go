package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Server   ServerConf
	Storage  StorageConf
	Postgres PostgresConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type string
}

type PostgresConf struct {
	Dsn string
}

func NewConfig(configFile string) (Config, error) {
	v := viper.New()
	v.SetConfigFile(configFile)

	var c Config

	if err := v.ReadInConfig(); err != nil {
		return c, err
	}

	if err := v.Unmarshal(&c); err != nil {
		return c, err
	}

	return c, nil
}

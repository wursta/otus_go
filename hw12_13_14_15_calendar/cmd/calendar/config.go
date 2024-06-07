package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger   LoggerConf
	Server   ServerConf
	HTTP     HTTPConf
	GRPC     GrpcConf
	Storage  StorageConf
	Postgres PostgresConf
}

type LoggerConf struct {
	Level string
}

type ServerConf struct {
	Type string
}

type HTTPConf struct {
	Host    string
	Port    string
	Timeout time.Duration
}

type GrpcConf struct {
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

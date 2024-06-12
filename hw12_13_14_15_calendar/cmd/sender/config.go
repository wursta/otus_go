package main

import (
	"github.com/spf13/viper"
)

type Config struct {
	Logger LoggerConf
	Rabbit RabbitConf
}

type LoggerConf struct {
	Level string
}

type RabbitConf struct {
	URI      string
	Exchange string
	Queue    string
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

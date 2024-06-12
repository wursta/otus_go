package main

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger    LoggerConf
	Scheduler SchedulerConf
	Storage   StorageConf
	Postgres  PostgresConf
	Rabbit    RabbitConf
}

type LoggerConf struct {
	Level string
}

type SchedulerConf struct {
	EventsNotifyCheckFrequency time.Duration
	OldEventsCleanerFrequency  time.Duration
}

type StorageConf struct {
	Type string
}

type PostgresConf struct {
	Dsn string
}

type RabbitConf struct {
	URI      string
	Exchange string
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

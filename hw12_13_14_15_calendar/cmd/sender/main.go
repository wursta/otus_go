package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/pflag"
	"github.com/streadway/amqp"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/logger"
	rabbit "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/queue/rabbit"
)

var configFile string

func init() {
	pflag.StringVar(&configFile, "config", "/etc/calendar/scheduler_config.toml", "Path to configuration file")
	pflag.Parse()
}

func main() {
	config, err := NewConfig(configFile)
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		return
	}

	log, err := logger.New(config.Logger.Level, os.Stderr)
	if err != nil {
		fmt.Printf("error creating logger: %v\n", err)
		return
	}

	log.Debug("created logger", log)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	consumer := rabbit.NewConsumer(config.Rabbit.URI, config.Rabbit.Exchange, config.Rabbit.Queue)
	defer consumer.Disconnect()
	log.Debug("create consumer", consumer)

	err = consumer.Connect()
	if err != nil {
		log.Error(fmt.Sprint("connect error:", err))
		return
	}
	log.Debug("connect consumer")

	deliveries, err := consumer.ConsumeEvents()
	if err != nil {
		fmt.Printf("error consuming events: %v\n", err)
		return
	}

	go handleEvents(log, deliveries)

	<-ctx.Done()
}

func handleEvents(log *logger.Logger, deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		log.Info(
			fmt.Sprintf(
				"got %dB delivery: [%v] %q",
				len(d.Body),
				d.DeliveryTag,
				d.Body,
			),
		)
		d.Ack(false)
	}
	log.Debug("handle: deliveries channel closed")
}

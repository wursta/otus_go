package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/pflag"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/logger"
	rabbit "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/queue/rabbit"
	sqlstorage "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage/sql"
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

	done := make(chan bool)

	var storage app.Storage

	switch config.Storage.Type {
	case "inmemory":
		log.Error("error creating storage: inmemory storage not allowed here")
		return
	case "postgres":
		sqlStorage := sqlstorage.New(config.Postgres.Dsn)

		ctx := context.Background()
		err = sqlStorage.Connect(ctx)
		if err != nil {
			log.Error(fmt.Sprintf("error connecting to database: %v\n", err))
			return
		}
		defer sqlStorage.Close(ctx)

		storage = sqlStorage

	default:
		log.Error("error creating storage: unknown storage type")
		return
	}

	log.Debug("create storage", storage)

	producer := rabbit.NewProducer(config.Rabbit.URI, config.Rabbit.Exchange)

	err = producer.Connect()
	if err != nil {
		log.Error(fmt.Sprint("connect error:", err))
		return
	}
	defer producer.Disconnect()

	log.Debug("create producer and connected", producer)

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		runScheduler(
			context.Background(),
			log,
			storage,
			config.Scheduler.EventsNotifyCheckFrequency,
			done,
			producer,
		)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		runOldEventsCleaner(
			context.Background(),
			log,
			storage,
			config.Scheduler.OldEventsCleanerFrequency,
			done,
		)
	}()

	<-ctx.Done()

	close(done)

	wg.Wait()
}

func runScheduler(
	ctx context.Context,
	log *logger.Logger,
	storage app.Storage,
	frequency time.Duration,
	doneCh <-chan bool,
	producer *rabbit.Producer,
) {
	checkEventsForNotify(ctx, log, storage, producer)

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-doneCh:
			log.Info("Stopping events notify checker")
			return
		case <-ticker.C:
			checkEventsForNotify(ctx, log, storage, producer)
		}
	}
}

func runOldEventsCleaner(
	ctx context.Context,
	log *logger.Logger,
	storage app.Storage,
	frequency time.Duration,
	doneCh <-chan bool,
) {
	removeOldEvents(ctx, log, storage)

	ticker := time.NewTicker(frequency)
	defer ticker.Stop()

	for {
		select {
		case <-doneCh:
			log.Info("Stopping old events cleaner")
			return
		case <-ticker.C:
			removeOldEvents(ctx, log, storage)
		}
	}
}

func checkEventsForNotify(ctx context.Context, log *logger.Logger, storage app.Storage, producer *rabbit.Producer) {
	events := storage.GetEventsForNotify(ctx, time.Now().Format(time.DateOnly))

	log.Info(fmt.Sprintf("Fetched events for notify: %d", len(events)))

	for i := range events {
		event := events[i]

		err := producer.ProduceEvent(event)

		if err != nil {
			log.Error(fmt.Sprint("error consume event:", err))
		} else {
			event.Notified = true

			err = storage.UpdateEvent(ctx, event.ID, event)
			if err != nil {
				log.Error(fmt.Sprint("error updating event:", err))
			}

			log.Info("Consume event: " + events[i].ID)
		}
	}
}

func removeOldEvents(ctx context.Context, log *logger.Logger, storage app.Storage) {
	yearAgo := time.Now().AddDate(-1, 0, 0)
	events := storage.GetEventsListByDates(ctx, nil, &yearAgo)

	log.Info(fmt.Sprintf("Fetched old events for clean: %d", len(events)))

	for i := range events {
		err := storage.DeleteEvent(ctx, events[i].ID)
		if err != nil {
			log.Error(fmt.Sprint("error while deleting old event:", err))
		}
	}
}

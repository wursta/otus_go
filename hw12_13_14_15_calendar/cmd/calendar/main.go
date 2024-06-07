package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		fmt.Printf("error reading config: %v\n", err)
		return
	}

	logg, err := logger.New(config.Logger.Level, os.Stderr)
	if err != nil {
		fmt.Printf("error creating logger: %v\n", err)
		return
	}

	logg.Debug("created logger", logg)

	var storage app.Storage
	switch config.Storage.Type {
	case "inmemory":
		storage = memorystorage.New()
	case "postgres":
		sqlStorage := sqlstorage.New(config.Postgres.Dsn)
		ctx := context.Background()
		err = sqlStorage.Connect(ctx)
		if err != nil {
			fmt.Printf("error connecting to database: %v\n", err)
			return
		}
		defer sqlStorage.Close(ctx)
		storage = sqlStorage
	default:
		fmt.Print("error creating storage: unknown storage type")
		return
	}
	logg.Debug("create storage", storage)

	calendar := app.New(logg, storage)
	logg.Debug("create calendar app", calendar)

	var server app.Server
	switch config.Server.Type {
	case "http":
		server = internalhttp.NewServer(
			logg,
			calendar,
			config.HTTP.Host,
			config.HTTP.Port,
			config.HTTP.Timeout,
		)
	case "grpc":
		server = internalgrpc.NewServer(
			logg,
			calendar,
			config.GRPC.Port,
		)
	default:
		fmt.Print("error creating server: unknown server type")
		return
	}

	logg.Debug("create server", server)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}

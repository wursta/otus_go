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

	httpServer := internalhttp.NewServer(
		logg,
		calendar,
		config.HTTP.Host,
		config.HTTP.Port,
		config.HTTP.Timeout,
	)
	logg.Debug("create http server", httpServer)

	grpcServer := internalgrpc.NewServer(
		logg,
		calendar,
		config.GRPC.Host,
		config.GRPC.Port,
	)
	logg.Debug("create grpc server", grpcServer)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		if err := httpServer.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			return
		}

		logg.Info("http server started...")
	}()

	go func() {
		if err := grpcServer.Start(ctx); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			return
		}

		logg.Info("grpc server started...")
	}()

	<-ctx.Done()

	timeoutCtx, timeoutCancel := context.WithTimeout(context.Background(), time.Second*3)
	defer timeoutCancel()

	if err := httpServer.Stop(timeoutCtx); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	if err := grpcServer.Stop(timeoutCtx); err != nil {
		logg.Error("failed to stop grpc server: " + err.Error())
	}
}

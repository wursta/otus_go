package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/spf13/pflag"
)

var (
	timeout time.Duration
	host    string
	port    string
)

const defaultTimeout = 10

func main() {
	pflag.Parse()

	parseArgs()

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	client := NewTelnetClient(net.JoinHostPort(host, port), timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Printf("Failed to connect: %v", err)
		return
	}

	go func() {
		defer cancel()

		if err := client.Send(); err != nil {
			return
		}

		log.Print("EOF")
	}()

	go func() {
		defer cancel()
		if err := client.Receive(); err != nil {
			return
		}
		log.Print("Connection was closed by peer")
	}()

	<-ctx.Done()

	if err := client.Close(); err != nil {
		log.Printf("Failed to close client: %v", err)
		return
	}
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetOutput(os.Stderr)
	log.SetPrefix("...")

	pflag.DurationVar(
		&timeout,
		"timeout",
		defaultTimeout*time.Second,
		fmt.Sprintf("connection timeout [default: %ss]", strconv.Itoa(defaultTimeout)),
	)
}

func parseArgs() {
	args := pflag.Args()
	if len(args) < 2 {
		pflag.Usage()
		os.Exit(1)
	}

	host = args[0]
	port = args[1]
}

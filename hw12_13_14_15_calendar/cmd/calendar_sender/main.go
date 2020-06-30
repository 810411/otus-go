package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/rmq"
	flag "github.com/spf13/pflag"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "/path/to/local.json")
}

func main() {
	flag.Parse()
	if configPath == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	conf, err := config.ParseConfig(configPath)
	if err != nil {
		log.Fatalf("can't parse config: %v", err)
	}

	err = logger.Configure(logger.Settings(conf.Log))
	if err != nil {
		log.Fatalf("can't config logger: %v", err)
	}
	logg := logger.Logger
	logg.Info("sender start")
	defer logg.Info("\nsender end")

	c := rmq.NewConsumer(conf.AMQP.URI, conf.AMQP.Queue)
	err = c.Connect()
	if err != nil {
		logg.Fatal(fmt.Sprintf("can't connect RMQ: %v", err))
	}
	defer c.Close()

	msgs, err := c.Consume()
	if err != nil {
		logg.Fatal(fmt.Sprintf("can't consume: %v", err))
	}

	sigChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-msgs:
				var notice repository.Notice
				if err := json.Unmarshal(msg.Body, &notice); err != nil {
					logg.Error(fmt.Sprintf("can't unmarshal notice: %v", err))
					continue
				}
				logg.Info(fmt.Sprintf(
					"ID: %d, title: \"%s\", datetime: %s, owner id: %d",
					notice.ID, notice.Title, notice.Datetime, notice.OwnerID,
				))
			}
		}
	}(ctx)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

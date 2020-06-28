package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/psql"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/rmq"
	flag "github.com/spf13/pflag"
	"go.uber.org/zap"
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

	period, err := time.ParseDuration(conf.Schedule.Period)
	if err != nil {
		log.Fatalf("wrong scan period: %v", err)
	}

	remindFor, err := time.ParseDuration(conf.Schedule.RemindFor)
	if err != nil {
		log.Fatalf("wrong remind_for: %v", err)
	}

	err = logger.Configure(logger.Settings(conf.Log))
	if err != nil {
		log.Fatalf("can't config logger: %v", err)
	}
	logg := logger.Logger
	logg.Info("scheduler start")
	defer logg.Info("\nscheduler end")

	ctx := context.Background()
	r := psql.New()
	err = r.Connect(ctx, conf.Repository.Dsn)
	if err != nil {
		logg.Fatal(fmt.Sprintf("can't connect DB: %v", err))
	}
	defer r.Close()

	p := rmq.NewProducer(conf.AMQP.URI, conf.AMQP.Queue)
	err = p.Connect()
	if err != nil {
		logg.Fatal(fmt.Sprintf("can't connect RMQ: %v", err))
	}
	defer p.Close()

	sigChan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go handle(ctx, remindFor, period, r, p, logg)

	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}

func handle(ctx context.Context, remindFor, period time.Duration, r *psql.Repo, p *rmq.Producer, logg *zap.Logger) {
	for {
		start := time.Now()

		notices, err := r.ListForScheduler(ctx, remindFor, period)
		if err != nil {
			logg.Error(fmt.Sprintf("can't get events: %v", err))
		}

		for _, v := range notices {
			b, err := json.Marshal(v)
			if err != nil {
				logg.Error(fmt.Sprintf("can't marshal notice: %v", err))
			}

			err = p.Publish(b)
			if err != nil {
				logg.Error(fmt.Sprintf("can't publish: %v", err))
			}
		}

		err = r.ClearMoreYearBefore(ctx)
		if err != nil {
			logg.Error(fmt.Sprintf("can't clear old events: %v", err))
		}

		timer := time.NewTimer(period - time.Since(start))
		select {
		case <-timer.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

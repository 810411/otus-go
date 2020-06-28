package main

import (
	"fmt"
	"log"
	"os"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/app/calendar"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
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
	logg.Info("calendar start")
	defer logg.Info("\ncalendar end")

	a, err := calendar.New(conf)
	if err != nil {
		logg.Fatal(fmt.Sprintf("can't create app: %v", err))
	}

	err = a.Run()
	if err != nil {
		logg.Fatal(fmt.Sprintf("when app run: %v", err))
	}
}

package main

import (
	"fmt"
	"log"
	"os"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/app"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/inmemory"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/psql"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/server"
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

	var r repository.EventsRepo
	switch conf.Repository.Type {
	case "psql":
		r = psql.New()
	default:
		r = inmemory.New()
	}

	s := server.New(server.Settings(conf.HTTP))

	a, err := app.New(r, s)
	if err != nil {
		logg.Fatal(fmt.Sprintf("can't create app: %v", err))
	}

	err = a.Run()
	if err != nil {
		logg.Fatal(fmt.Sprintf("when app run: %v", err))
	}
}
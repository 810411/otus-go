package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/api/rest"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/inmemory"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/psql"
)

const shutdownTime = 5 * time.Second

type App struct {
	ctx  context.Context
	conf *config.Config
	r    repository.EventsRepo
	s    *http.Server
}

func New(conf *config.Config) (*App, error) {
	var r repository.EventsRepo
	switch conf.Repository.Type {
	case "psql":
		r = psql.New()
	default:
		r = inmemory.New()
	}

	ctx := context.Background()

	s := rest.New(ctx, rest.Settings(conf.HTTP), r)

	return &App{
		ctx:  ctx,
		conf: conf,
		r:    r,
		s:    s,
	}, nil
}

func (a *App) Run() error {
	if r, ok := a.r.(repository.BaseRepo); ok {
		err := r.Connect(a.ctx, a.conf.Repository.Dsn)
		if err != nil {
			return err
		}
		defer r.Close()
	}

	eChan := make(chan error)
	sigChan := make(chan os.Signal, 1)

	go func() {
		err := a.s.ListenAndServe()
		if err != nil {
			eChan <- err
		}
	}()

	signal.Notify(sigChan, os.Interrupt)

	select {
	case err := <-eChan:
		return err
	case <-sigChan:
	}

	ctx, cancelFn := context.WithTimeout(a.ctx, shutdownTime)
	defer cancelFn()

	if err := a.s.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

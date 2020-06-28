package calendar

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/api"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/api/grpc"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/api/rest"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/config"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/inmemory"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/psql"
)

const httpShutdownTime = 5 * time.Second

type App struct {
	ctx  context.Context
	conf *config.Config
	r    repository.EventsRepo
	rest *http.Server
	grpc *grpc.Service
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

	s := rest.New(ctx, api.Settings(conf.HTTP), r)
	g := grpc.New(ctx, api.Settings(conf.GRPC), r)

	return &App{
		ctx:  ctx,
		conf: conf,
		r:    r,
		rest: s,
		grpc: g,
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
		err := a.rest.ListenAndServe()
		if err != nil {
			eChan <- err
		}
	}()

	go func() {
		err := a.grpc.ListenAndServe()
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

	a.grpc.GracefulStop()

	ctx, cancelFn := context.WithTimeout(a.ctx, httpShutdownTime)
	defer cancelFn()

	if err := a.rest.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

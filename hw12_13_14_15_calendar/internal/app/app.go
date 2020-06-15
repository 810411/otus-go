package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
)

const shutdownTime = 5 * time.Second

type App struct {
	r repository.EventsRepo
	s *http.Server
}

func New(r repository.EventsRepo, s *http.Server) (*App, error) {
	return &App{
		r: r,
		s: s,
	}, nil
}

func (a *App) Run() error {
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

	ctx, cancelFn := context.WithTimeout(context.Background(), shutdownTime)
	defer cancelFn()

	if err := a.s.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

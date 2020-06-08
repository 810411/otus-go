package app

import (
	"context"

	"github.com/810411/otus-go/hw_calendar/internal/repository"
)

type App struct {
	r repository.EventsRepo
}

func New(r repository.EventsRepo) (*App, error) {
	return &App{r: r}, nil
}

func (a *App) Run(ctx context.Context) error {
	return nil
}

package repository

import (
	"context"
	"time"
)

type EventsRepo interface {
	Create(ctx context.Context, event Event) (Event, error)
	Update(ctx context.Context, event Event) (Event, error)
	Delete(ctx context.Context, id EventID) (EventID, error)
	ListOfDay(ctx context.Context, from time.Time) ([]Event, error)
	ListOfWeek(ctx context.Context, from time.Time) ([]Event, error)
	ListOfMonth(ctx context.Context, from time.Time) ([]Event, error)
}

type BaseRepo interface {
	Connect(ctx context.Context, dsn string) error
	Close() error
	ClearMoreYearBefore(ctx context.Context) error
	ListForScheduler(ctx context.Context, remindFor time.Duration, period time.Duration) ([]Notice, error)
	EventsRepo
}

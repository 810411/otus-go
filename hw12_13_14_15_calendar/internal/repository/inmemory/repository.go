package inmemory

import (
	"context"
	"sync"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
)

type Repo struct {
	mx    sync.RWMutex
	Store map[repository.EventID]repository.Event
	idInc uint64
}

func New() *Repo {
	return &Repo{
		Store: make(map[repository.EventID]repository.Event),
	}
}

func (r *Repo) forEach(ctx context.Context, fn func(event repository.Event) error) error {
	for _, v := range r.Store {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := fn(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) Create(ctx context.Context, event repository.Event) (repository.Event, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	event.Datetime = event.Datetime.Round(time.Minute)
	err := r.forEach(ctx, func(v repository.Event) error {
		if event.Datetime.Equal(v.Datetime) && event.OwnerID == v.OwnerID {
			return repository.ErrTimeBusy
		}
		return nil
	})
	if err != nil {
		return event, err
	}

	event.ID = repository.EventID(r.idInc)
	r.Store[event.ID] = event
	r.idInc++

	return event, nil
}

func (r *Repo) Update(ctx context.Context, event repository.Event) (repository.Event, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.Store[event.ID]; !ok {
		return event, repository.ErrNotFound
	}

	event.Datetime = event.Datetime.Round(time.Minute)
	err := r.forEach(ctx, func(v repository.Event) error {
		if event.ID != v.ID && event.Datetime.Equal(v.Datetime) && event.OwnerID == v.OwnerID {
			return repository.ErrTimeBusy
		}
		return nil
	})
	if err != nil {
		return r.Store[event.ID], err
	}

	r.Store[event.ID] = event

	return event, nil
}

func (r *Repo) Delete(ctx context.Context, id repository.EventID) (repository.EventID, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	if _, ok := r.Store[id]; !ok {
		return id, repository.ErrNotFound
	}

	delete(r.Store, id)
	return id, nil
}

func (r *Repo) listOf(ctx context.Context, from time.Time, p repository.Period) (events []repository.Event, err error) {
	r.mx.RLock()
	defer r.mx.RUnlock()

	from, to := repository.GetTimeRange(from, p)

	err = r.forEach(ctx, func(event repository.Event) error {
		if event.Datetime.Unix() >= from.Unix() && event.Datetime.Before(to) {
			events = append(events, event)
		}
		return nil
	})

	return
}

func (r *Repo) ListOfDay(ctx context.Context, from time.Time) (events []repository.Event, err error) {
	return r.listOf(ctx, from, repository.Day)
}

func (r *Repo) ListOfWeek(ctx context.Context, from time.Time) (events []repository.Event, err error) {
	return r.listOf(ctx, from, repository.Week)
}

func (r *Repo) ListOfMonth(ctx context.Context, from time.Time) (events []repository.Event, err error) {
	return r.listOf(ctx, from, repository.Month)
}

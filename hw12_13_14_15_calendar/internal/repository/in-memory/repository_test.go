package memory //nolint:golint,stylecheck

import (
	"context"
	"github.com/810411/otus-go/hw_calendar/internal/repository"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
	"time"
)

func beforeEach() (*Repo, context.Context) {
	r := New()
	e1 := repository.Event{ID: 0, Title: "test", Datetime: time.Unix(0, 0), OwnerID: 1}
	e2 := repository.Event{ID: 1, Title: "test", Datetime: time.Unix(60, 0), OwnerID: 1}
	r.Store[e1.ID] = e1
	r.Store[e2.ID] = e2
	r.idInc = 2
	return r, context.Background()
}

func Test_Repo(t *testing.T) {
	t.Run("foreach storage", func(t *testing.T) {
		r, ctx := beforeEach()
		want := [2]string{"test0", "test1"}
		got := [2]string{}
		fn := func(v repository.Event) error {
			got[int64(v.ID)] = v.Title + strconv.Itoa(int(v.ID))
			return nil
		}

		err := r.forEach(ctx, fn)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("create event", func(t *testing.T) {
		r, ctx := beforeEach()
		event := repository.Event{Title: "test", Datetime: time.Now(), OwnerID: 1}

		got, err := r.Create(ctx, event)
		require.NoError(t, err)
		require.NotEqual(t, 0, got.ID)
		require.NotEqual(t, 1, got.ID)
		require.Equal(t, "test", got.Title)
		require.Equal(t, 3, len(r.Store))
	})

	t.Run("create event errors", func(t *testing.T) {
		r, ctx := beforeEach()
		event := repository.Event{Title: "test", Datetime: time.Unix(0, 0), OwnerID: 1}
		_, err := r.Create(ctx, event)
		require.Equal(t, repository.ErrTimeBusy, err)

		event.OwnerID = 0
		_, err = r.Create(ctx, event)
		require.NoError(t, err)

		ctx, cancel := context.WithCancel(ctx)
		cancel()
		event.Datetime = time.Now()
		_, err = r.Create(ctx, event)
		require.Equal(t, context.Canceled, err)
	})

	t.Run("update event", func(t *testing.T) {
		r, ctx := beforeEach()
		id := repository.EventID(0)
		event := r.Store[id]

		event.Description = "test"
		got, err := r.Update(ctx, event)
		require.NoError(t, err)
		require.Equal(t, "test", got.Description)

		_, err = r.Update(ctx, repository.Event{ID: 100})
		require.Equal(t, repository.ErrNotFound, err)

		event.Datetime = time.Unix(60, 0)
		_, err = r.Update(ctx, event)
		require.Equal(t, repository.ErrTimeBusy, err)
	})

	t.Run("delete event", func(t *testing.T) {
		r, ctx := beforeEach()
		errChan := make(chan error)
		vChan := make(chan repository.EventID)
		var vCount, errCount int

		for i := 0; i < 3; i++ {
			go func(id int, errChan chan<- error, vChan chan<- repository.EventID) {
				v, err := r.Delete(ctx, repository.EventID(id))
				if err != nil {
					errChan <- err
					return
				}
				vChan <- v
			}(i, errChan, vChan)
		}

		for i := 0; i < 3; i++ {
			select {
			case err := <-errChan:
				require.Equal(t, repository.ErrNotFound, err)
				errCount++
			case <-vChan:
				vCount++
			}
		}
		require.Empty(t, r.Store)
		require.Equal(t, vCount, 2)
		require.Equal(t, errCount, 1)
	})

	t.Run("list events", func(t *testing.T) {
		const hour = 60 * 60
		r, ctx := beforeEach()
		events := [4]repository.Event{
			{Datetime: time.Unix(25*hour, 0)},
			{Datetime: time.Unix(24*hour, 0)},
			{Datetime: time.Unix(36*hour, 0)},
			{Datetime: time.Unix(49*hour, 0)},
		}
		want := ""

		for _, v := range events {
			event, err := r.Create(ctx, v)
			require.NoError(t, err)
			want += event.Datetime.String()
		}
		require.Equal(t, 6, len(r.Store))

		eventsOfDay, err := r.ListOfDay(ctx, events[0].Datetime)
		require.NoError(t, err)
		require.Equal(t, 3, len(eventsOfDay))

		for i, _ := range eventsOfDay {
			require.Contains(t, want, eventsOfDay[i].Datetime.String())
		}
	})
}

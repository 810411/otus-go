package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository/inmemory"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var (
	testEvent = repository.Event{0, "test", time.Unix(60, 0), 10 * time.Hour, "desc", 1}
	upEvent   = repository.Event{0, "test", time.Unix(60, 0), 3 * time.Hour, "updated", 1}
)

func BeforeEach() {
	repo = inmemory.New()
	ctx = context.Background()
	_, err := repo.Create(ctx, testEvent)
	if err != nil {
		log.Fatal(err)
	}
}

func serveTest(t *testing.T, p string, m string, u string, fn func(w http.ResponseWriter, r *http.Request), b io.Reader) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc(p, fn).Methods(m)
	req, err := http.NewRequest(m, u, b)
	require.NoError(t, err)
	router.ServeHTTP(rr, req)

	return rr
}

func Test_GetEvents(t *testing.T) {
	t.Run("get events wrong period param", func(t *testing.T) {
		var restErr ErrorResult

		rr := serveTest(t, "/events/{period}", "GET", "/events/weekend", getEvents, nil)
		require.Equal(t, http.StatusBadRequest, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&restErr)
		require.Equal(t, statusError, restErr.Status)
		require.NoError(t, err)
		require.Equal(t, ErrWrongPeriod.Error(), restErr.Err.Message)
	})

	t.Run("get events missing datetime query", func(t *testing.T) {
		var restErr ErrorResult

		rr := serveTest(t, "/events/{period}", "GET", "/events/day", getEvents, nil)
		require.Equal(t, http.StatusBadRequest, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&restErr)
		require.Equal(t, statusError, restErr.Status)
		require.NoError(t, err)
		require.Equal(t, ErrDatetimeMissing.Error(), restErr.Err.Message)
	})

	t.Run("get events wrong datetime query value", func(t *testing.T) {
		var restErr ErrorResult
		BeforeEach()

		rr := serveTest(t, "/events/{period}", "GET", "/events/week?datetime=1971-01-01", getEvents, nil)
		require.Equal(t, http.StatusBadRequest, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&restErr)
		require.NoError(t, err)
		require.NotEmpty(t, restErr.Err.Message)
	})

	t.Run("get events ok", func(t *testing.T) {
		var got EventsResult
		want := EventsResult{statusOk, []repository.Event{testEvent}}
		BeforeEach()

		rr := serveTest(t, "/events/{period}", "GET", "/events/month?datetime=1970-01-01T00:00:00%2b0000", getEvents, nil)
		require.Equal(t, http.StatusOK, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})
}

func Test_DeleteEvent(t *testing.T) {
	t.Run("delete ok", func(t *testing.T) {
		var got DeleteResult
		want := DeleteResult{statusOk, repository.EventID(0)}
		BeforeEach()

		rr := serveTest(t, "/events/{id}", "DELETE", fmt.Sprintf("/events/%d", testEvent.ID), deleteEvent, nil)
		require.Equal(t, http.StatusOK, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&got)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("delete not found", func(t *testing.T) {
		var restErr ErrorResult
		BeforeEach()

		rr := serveTest(t, "/events/{id}", "DELETE", "/events/1", deleteEvent, nil)
		require.Equal(t, http.StatusNotFound, rr.Code)

		err := json.NewDecoder(rr.Body).Decode(&restErr)
		require.NoError(t, err)
		require.Equal(t, repository.ErrNotFound.Error(), restErr.Err.Message)
	})
}

func Test_UpdateEvent(t *testing.T) {
	t.Run("update ok", func(t *testing.T) {
		var got EventResult
		want := EventResult{statusOk, upEvent}
		BeforeEach()

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(&upEvent)
		require.NoError(t, err)
		rr := serveTest(t, "/events/{id}", "PUT", fmt.Sprintf("/events/%d", upEvent.ID), updateEvent, b)
		require.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&got)
		require.NoError(t, err)
		require.Equal(t, 0, int(got.Event.ID))
		require.Equal(t, want.Event.Duration, got.Event.Duration)
		require.Equal(t, "updated", got.Event.Description)
	})

	t.Run("update not ok", func(t *testing.T) {
		var restErr ErrorResult
		upEvent.ID = 13
		BeforeEach()

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(&upEvent)
		require.NoError(t, err)
		rr := serveTest(t, "/events/{id}", "PUT", fmt.Sprintf("/events/%d", upEvent.ID), updateEvent, b)
		require.Equal(t, http.StatusBadRequest, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&restErr)
		require.NoError(t, err)
		require.Equal(t, repository.ErrNotFound.Error(), restErr.Err.Message)
	})
}

func Test_CreateEvent(t *testing.T) {
	t.Run("create ok", func(t *testing.T) {
		var got EventResult
		BeforeEach()
		testEvent.Datetime = time.Unix(3660, 0)
		want := EventResult{statusOk, testEvent}

		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(&testEvent)
		require.NoError(t, err)
		rr := serveTest(t, "/events", "POST", "/events", createEvent, b)
		require.Equal(t, http.StatusOK, rr.Code)

		err = json.NewDecoder(rr.Body).Decode(&got)
		require.NoError(t, err)
		require.Equal(t, want.Status, got.Status)
		require.Equal(t, 1, int(got.Event.ID))

		arr, err := repo.ListOfDay(ctx, time.Unix(0, 0))
		require.NoError(t, err)
		require.Equal(t, 2, len(arr))
	})
}

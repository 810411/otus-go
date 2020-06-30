package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/gorilla/mux"
)

func handle404(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	e := Error{
		HTTPCode: code,
		Message:  http.StatusText(code),
	}
	writeError(w, e)
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var event repository.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var result EventResult
	result.Event, err = repo.Create(ctx, event)
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}
	result.Status = statusOk

	writeResponse(w, result)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var event repository.Event
	err = json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}
	event.ID = repository.EventID(id)

	var result EventResult
	result.Event, err = repo.Update(ctx, event)
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}
	result.Status = statusOk

	writeResponse(w, result)
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var result DeleteResult
	result.ID, err = repo.Delete(ctx, repository.EventID(id))
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusNotFound, Message: err.Error()})
		return
	}
	result.Status = statusOk

	writeResponse(w, result)
}

func getEvents(w http.ResponseWriter, r *http.Request) {
	var p repository.Period
	params := mux.Vars(r)
	switch params["period"] {
	case "day":
		p = repository.Day
	case "week":
		p = repository.Week
	case "month":
		p = repository.Month
	default:
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: ErrWrongPeriod.Error()})
		return
	}

	query := r.URL.Query()
	dt := query.Get("datetime")
	if dt == "" {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: ErrDatetimeMissing.Error()})
		return
	}

	from, err := time.Parse("2006-01-02T15:04:05-0700", dt)
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusBadRequest, Message: err.Error()})
		return
	}

	var result EventsResult
	switch p {
	case repository.Day:
		result.Events, err = repo.ListOfDay(ctx, from)
	case repository.Week:
		result.Events, err = repo.ListOfWeek(ctx, from)
	case repository.Month:
		result.Events, err = repo.ListOfMonth(ctx, from)
	}
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusOK, Message: err.Error()})
		return
	}
	result.Status = statusOk

	writeResponse(w, result)
}

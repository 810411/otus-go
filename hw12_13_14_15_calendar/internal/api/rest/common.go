package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
)

const (
	statusOk    = "ok"
	statusError = "error"
)

var (
	ErrWrongPeriod     = errors.New("wrong period param")
	ErrDatetimeMissing = errors.New("missing datetime query param")
)

type EventsResult struct {
	Status string             `json:"status"`
	Events []repository.Event `json:"events"`
}

type EventResult struct {
	Status string           `json:"status"`
	Event  repository.Event `json:"event"`
}

type DeleteResult struct {
	Status string             `json:"status"`
	ID     repository.EventID `json:"id"`
}

type Error struct {
	HTTPCode int    `json:"-"`
	Message  string `json:"message"`
}

type ErrorResult struct {
	Status string `json:"status"`
	Err    Error  `json:"error"`
}

func writeError(w http.ResponseWriter, e Error) {
	data := ErrorResult{statusError, e}

	b, err := json.Marshal(data)
	if err != nil {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	w.WriteHeader(e.HTTPCode)
	_, _ = fmt.Fprint(w, string(b))
}

func writeResponse(w http.ResponseWriter, result interface{}) {
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		writeError(w, Error{HTTPCode: http.StatusInternalServerError, Message: err.Error()})
		return
	}
}

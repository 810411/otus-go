package rest

import (
	"context"
	"net/http"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
	"github.com/gorilla/mux"
)

type Settings struct {
	Host string
	Port string
}

var (
	ctx  context.Context
	repo repository.EventsRepo
)

func New(c context.Context, s Settings, r repository.EventsRepo) *http.Server {
	ctx = c
	repo = r

	router := mux.NewRouter()
	router.HandleFunc("/events/{period}", getEvents).Methods("GET")
	router.HandleFunc("/events", createEvent).Methods("POST")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PUT")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	router.HandleFunc("/", handle404)
	router.Use(LogMiddleware)

	return &http.Server{
		Addr:    s.Host + ":" + s.Port,
		Handler: router,
	}
}

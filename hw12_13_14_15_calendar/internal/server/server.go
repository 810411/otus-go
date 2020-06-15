package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
)

type Settings struct {
	Host string
	Port string
}

func New(s Settings) *http.Server {
	handler := http.HandlerFunc(handle)
	http.Handle("/", LogMiddleware(handler))
	return &http.Server{
		Addr: s.Host + ":" + s.Port,
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	time.Sleep(1 * time.Second)
	_, _ = fmt.Fprint(w, r, time.Now())
}

func LogMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		h.ServeHTTP(w, r)
		latency := time.Since(startTime)

		info := fmt.Sprintf("%s [%s] %s %s %s %d %s \"%s\"",
			r.RemoteAddr,
			time.Now().Format("2006-01-02 15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			http.StatusOK,
			latency,
			r.UserAgent(),
		)
		logger.Logger.Info(info)
	})
}

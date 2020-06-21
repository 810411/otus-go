package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/logger"
)

func HeadersMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		h.ServeHTTP(w, r)
	})
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

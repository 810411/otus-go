package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/810411/otus-go/hw_calendar/internal/logger"
)

type Settings struct {
	Host string
	Port string
}

const shutdownTime = 5 * time.Second

func Start(ctx context.Context, s Settings) (err error) {
	e := make(chan error)

	handler := http.HandlerFunc(handle)
	http.Handle("/", logMiddleware(handler))
	srv := &http.Server{
		Addr: s.Host + ":" + s.Port,
	}

	go func(e chan<- error) {
		e <- srv.ListenAndServe()
	}(e)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	select {
	case err = <-e:
		return
	case <-c:
	}

	ctx, cancel := context.WithTimeout(ctx, shutdownTime)
	defer cancel()
	err = srv.Shutdown(ctx)

	return
}

func handle(w http.ResponseWriter, r *http.Request) {
	time.Sleep(3 * time.Second)
	_, _ = fmt.Fprint(w, r, time.Now())
}

func logMiddleware(h http.Handler) http.Handler {
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

package api

import (
	"log"
	"net/http"
	"time"
)

// NewHandler returns the fully-wired HTTP handler for the API, including
// request logging and panic recovery so a single bad request can't take
// down the process.
func NewHandler() http.Handler {
	return loggingMiddleware(recoverMiddleware(NewMux()))
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic handling %s %s: %v", r.Method, r.URL.Path, err)
				writeError(w, http.StatusInternalServerError, "internal server error")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// loggingMiddleware logs one line per request: method, path, remote address,
// resulting status code, and how long the request took.
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		log.Printf("-> %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(sw, r)
		log.Printf("<- %s %s %d (%s)", r.Method, r.URL.Path, sw.status, time.Since(start))
	})
}

// statusWriter records the status code passed to WriteHeader so it can be
// logged after the handler finishes.
type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

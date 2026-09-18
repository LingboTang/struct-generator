package api

import (
	"log"
	"net/http"
)

// NewHandler returns the fully-wired HTTP handler for the API, including
// panic recovery so a single bad request can't take down the process.
func NewHandler() http.Handler {
	return recoverMiddleware(NewMux())
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

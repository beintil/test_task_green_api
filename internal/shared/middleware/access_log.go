package middleware

import (
	"net/http"
	"time"

	"github.com/darkness/green_api/pkg/logger"

	"github.com/gorilla/mux"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func AccessLog(log logger.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{w, http.StatusOK}
			next.ServeHTTP(rw, r)
			log.Infof("[%s] %s %d %s", r.Method, r.URL.Path, rw.statusCode, time.Since(start))
		})
	}
}

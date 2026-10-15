package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func Logger(log *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			log.Info("HTTP  request",
				"method", r.Method,
				"path", r.URL.Path,
				"duration", time.Since(start),
			)
		})
	}
}

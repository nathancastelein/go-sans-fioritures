package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			middleware := middlewares[i]
			next = middleware(next)
		}
		return next
	}
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := s.login.LogUser(r.Context(), username, password)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		slog.Info("user logged", "user", user)
		next.ServeHTTP(w, r)
	})
}

type contextKey int

const (
	contextLoggerKey contextKey = 1
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.With(
			slog.String("http_method", r.Method),
			slog.String("http_url_path", r.URL.Path),
			slog.String("request_id", r.Header.Get("X-Request-Id")),
		)

		r = r.WithContext(context.WithValue(r.Context(), contextLoggerKey, logger))
		startedAt := time.Now()

		next.ServeHTTP(w, r)

		logger.Info("request end", slog.Int("elapsed_time_ms", int(time.Since(startedAt).Milliseconds())))
	})
}

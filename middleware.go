package main

import (
	"context"
	"log/slog"
	"net/http"
	"slices"
	"time"
)

type Middleware func(http.Handler) http.Handler

func MiddlewareChain(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for _, middleware := range slices.Backward(middlewares) {
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
			http.Error(w, "Forbidden", http.StatusForbidden)
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

		logger.Info("request end", slog.Duration("elapsed_time", time.Since(startedAt)))
	})
}

func RecoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				slog.Error("recover from panic", slog.Any("error", err))
				w.Header().Set("Content-Type", "application/plain")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

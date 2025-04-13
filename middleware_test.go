package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type LoginStub struct {
	username string
	password string
}

func NewLoginStub(username, password string) Login {
	return &LoginStub{
		username: username,
		password: password,
	}
}

func (l *LoginStub) LogUser(ctx context.Context, username string, password string) (*User, error) {
	if l.username == username && l.password == password {
		return &User{
			ID:        1,
			FirstName: "Peter",
			LastName:  "Parker",
		}, nil
	}

	return nil, errors.New("login failed")
}

func TestServer_AuthMiddleware(t *testing.T) {
	t.Parallel()
	// Arrange
	server := NewServer(nil, nil, NewLoginStub("Peter", "Parker"))

	t.Run("authentication is OK", func(t *testing.T) {
		t.Parallel()
		// Arrange
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetBasicAuth("Peter", "Parker")

		// Act
		server.AuthMiddleware(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			),
		).ServeHTTP(responseRecorder, request)

		// Assert
		if responseRecorder.Result().StatusCode != http.StatusOK {
			t.Fatalf("expected http code %d, got %d", http.StatusOK, responseRecorder.Result().StatusCode)
		}
	})

	t.Run("authentication failed", func(t *testing.T) {
		t.Parallel()
		// Arrange
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		request.SetBasicAuth("Spider", "Man")

		// Act
		server.AuthMiddleware(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			),
		).ServeHTTP(responseRecorder, request)

		// Assert
		if responseRecorder.Result().StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected http code %d, got %d", http.StatusUnauthorized, responseRecorder.Result().StatusCode)
		}
	})

	t.Run("no basic auth provided", func(t *testing.T) {
		t.Parallel()
		// Arrange
		responseRecorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// Act
		server.AuthMiddleware(
			http.HandlerFunc(
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
				},
			),
		).ServeHTTP(responseRecorder, request)

		// Assert
		if responseRecorder.Result().StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected http code %d, got %d", http.StatusUnauthorized, responseRecorder.Result().StatusCode)
		}
	})
}

func TestLoggerMiddleware(t *testing.T) {
	// Arrange
	type log struct {
		Message     string `json:"msg"`
		HttpMethod  string `json:"http_method"`
		HttpURLPath string `json:"http_url_path"`
		RequestID   string `json:"request_id"`
		ElapsedTime int64  `json:"elapsed_time_ms"`
	}

	var buffer bytes.Buffer
	defaultLogger := slog.Default()
	logger := slog.New(slog.NewJSONHandler(&buffer, nil))
	slog.SetDefault(logger)
	t.Cleanup(func() {
		slog.SetDefault(defaultLogger)
	})

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/test", nil)
	request.Header.Set("X-Request-Id", "foobar")

	// Act
	LoggerMiddleware(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	).ServeHTTP(responseRecorder, request)

	// Assert
	str, err := buffer.ReadBytes('\n')
	if err != nil {
		t.Fatalf("got an error while reading log buffer: %s", err)
	}

	receivedLog := &log{}
	if err := json.Unmarshal(str, receivedLog); err != nil {
		t.Fatalf("expected log to be JSON formatted, got an error while unmarshaling: %s", err)
	}

	if receivedLog.Message != "request end" {
		t.Fatalf("expected log 'request end', got: '%s'", receivedLog.Message)
	}

	if receivedLog.HttpMethod != http.MethodPost {
		t.Fatalf("expected http method 'POST', got: '%s'", receivedLog.HttpMethod)
	}

	if receivedLog.HttpURLPath != "/test" {
		t.Fatalf("expected http url path '/test', got: '%s'", receivedLog.HttpURLPath)
	}

	if receivedLog.RequestID != "foobar" {
		t.Fatalf("expected request id 'foobar', got: '%s'", receivedLog.RequestID)
	}

	if receivedLog.ElapsedTime > 0 {
		t.Fatalf("expected elapsed time > 0, got: %d", receivedLog.ElapsedTime)
	}
}

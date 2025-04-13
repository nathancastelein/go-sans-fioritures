package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	server := NewServer(
		NewInMemoryStoneRepository(),
		NewInMemoryReportRepository(),
		NewInMemoryLogin(),
	)

	slog.Info("API server starting on :8080")
	err := http.ListenAndServe(":8080", server)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}

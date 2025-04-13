package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type Server struct {
	http.Handler
	stone  StoneRepository
	report ReportRepository
	login  Login
}

func NewServer(
	stone StoneRepository,
	report ReportRepository,
	login Login,
) *Server {
	server := &Server{
		stone:  stone,
		report: report,
		login:  login,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /stones", server.getStones)
	mux.HandleFunc("GET /stones/{name}", server.getStoneByName)
	mux.HandleFunc("POST /stones/report", server.reportSuspiciousActivity)

	chain := MiddlewareChain(LoggerMiddleware, server.AuthMiddleware)
	server.Handler = chain(mux)

	return server
}

// Handlers
func (s *Server) getStones(w http.ResponseWriter, r *http.Request) {
	writeJSONResponse(w, http.StatusOK, s.stone.ListStones(r.Context()))
}

func (s *Server) getStoneByName(w http.ResponseWriter, r *http.Request) {
	stoneName := r.PathValue("name")
	stone, err := s.stone.GetStone(r.Context(), stoneName)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	writeJSONResponse(w, http.StatusOK, stone)
}

func (s *Server) reportSuspiciousActivity(w http.ResponseWriter, r *http.Request) {
	var report struct {
		Stone  string `json:"stone"`
		Report string `json:"report"`
	}

	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		slog.Warn("invalid report payload", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, data any) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

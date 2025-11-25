package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ChernykhITMO/Avito/internal/service"
)

type StatsHandler struct {
	service service.StatsService
}

func NewStatsHandler(s service.StatsService) *StatsHandler {
	return &StatsHandler{
		service: s}
}

func (h *StatsHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/stats", h.handleStats)
}

func (h *StatsHandler) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.service.GetStats(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "failed to encode stats", http.StatusInternalServerError)
		return
	}
}

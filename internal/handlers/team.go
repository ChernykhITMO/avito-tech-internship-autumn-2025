package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ChernykhITMO/Avito/internal/domain"
	"github.com/ChernykhITMO/Avito/internal/handlers/dto"
	"github.com/ChernykhITMO/Avito/internal/service"
)

type TeamHandler struct {
	serv service.TeamService
}

func NewTeamHandler(serv service.TeamService) *TeamHandler {
	return &TeamHandler{
		serv: serv,
	}
}

func (h *TeamHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/team/add", h.handleAddTeam)
	mux.HandleFunc("/team/get", h.handleGetTeam)
}

func (h *TeamHandler) handleAddTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var teamDTO dto.Team
	if err := json.NewDecoder(r.Body).Decode(&teamDTO); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	team := dto.TeamDTOToDomain(teamDTO)

	created, err := h.serv.CreateTeam(r.Context(), team.Name, team.Members)
	if err != nil {
		var derr *domain.Error
		if errors.As(err, &derr) {
			writeDomainError(w, derr)
			return
		}

		writeInternal(w)
		return
	}

	resp := struct {
		Team dto.Team `json:"team"`
	}{
		Team: dto.TeamToDTO(*created),
	}

	writeJSON(w, http.StatusCreated, resp)
}

func (h *TeamHandler) handleGetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("team_name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	team, err := h.serv.GetTeam(r.Context(), name)
	if err != nil {
		var derr *domain.Error
		if errors.As(err, &derr) {
			writeDomainError(w, derr)
			return
		}

		writeInternal(w)
		return
	}

	writeJSON(w, http.StatusOK, dto.TeamToDTO(*team))
}

package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ChernykhITMO/Avito/internal/domain"
	"github.com/ChernykhITMO/Avito/internal/handlers/dto"
	"github.com/ChernykhITMO/Avito/internal/service"
)

type UserHandler struct {
	serv service.UserService
}

func NewUserHandler(serv service.UserService) *UserHandler {
	return &UserHandler{
		serv: serv,
	}
}

func (h *UserHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/users/setIsActive", h.handleSetIsActive)
	mux.HandleFunc("/users/getReview", h.handleGetReview)
}

type setIsActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (h *UserHandler) handleSetIsActive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req setIsActiveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := h.serv.SetIsActive(r.Context(), req.UserID, req.IsActive)
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
		User dto.User `json:"user"`
	}{
		User: dto.UserToDTO(*user),
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) handleGetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	prs, err := h.serv.GetUserReviewPRs(r.Context(), userID)
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
		UserID       string                 `json:"user_id"`
		PullRequests []dto.PullRequestShort `json:"pull_requests"`
	}{
		UserID:       userID,
		PullRequests: make([]dto.PullRequestShort, 0, len(prs)),
	}

	for _, pr := range prs {
		resp.PullRequests = append(resp.PullRequests, *dto.PullRequestToShortDTO(pr))
	}

	writeJSON(w, http.StatusOK, resp)
}

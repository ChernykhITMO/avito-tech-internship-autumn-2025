package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ChernykhITMO/Avito/internal/domain"
	"github.com/ChernykhITMO/Avito/internal/handlers/dto"
	"github.com/ChernykhITMO/Avito/internal/service"
)

type PullRequestHandler struct {
	serv service.PullRequestService
}

func NewPullRequestHandler(serv service.PullRequestService) *PullRequestHandler {
	return &PullRequestHandler{serv: serv}
}

func (h *PullRequestHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/pullRequest/create", h.handleCreate)
	mux.HandleFunc("/pullRequest/merge", h.handleMerge)
	mux.HandleFunc("/pullRequest/reassign", h.handleReassign)
}

type createPRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

func (h *PullRequestHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req createPRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.PullRequestName == "" || req.AuthorID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, err := h.serv.Create(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
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
		PR *dto.PullRequest `json:"pr"`
	}{
		PR: dto.PullRequestToDTO(*pr),
	}

	writeJSON(w, http.StatusCreated, resp)
}

type mergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (h *PullRequestHandler) handleMerge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req mergePRRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, err := h.serv.Merge(r.Context(), req.PullRequestID)
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
		PR *dto.PullRequest `json:"pr"`
	}{
		PR: dto.PullRequestToDTO(*pr),
	}

	writeJSON(w, http.StatusOK, resp)
}

type reassignRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

func (h *PullRequestHandler) handleReassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req reassignRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.PullRequestID == "" || req.OldUserID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pr, replacedBy, err := h.serv.ReassignReviewer(r.Context(), req.PullRequestID, req.OldUserID)
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
		PR         *dto.PullRequest `json:"pr"`
		ReplacedBy string           `json:"replaced_by"`
	}{
		PR:         dto.PullRequestToDTO(*pr),
		ReplacedBy: replacedBy,
	}

	writeJSON(w, http.StatusOK, resp)
}

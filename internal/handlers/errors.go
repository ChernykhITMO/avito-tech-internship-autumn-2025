package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type errorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeAPIError(w http.ResponseWriter, status int, code, msg string) {
	var resp errorResponse
	resp.Error.Code = code
	resp.Error.Message = msg
	writeJSON(w, status, resp)
}

func writeInternal(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func statusByDomainCode(code domain.Code) int {
	switch code {
	case domain.ErrorCodeNotFound:
		return http.StatusNotFound
	case domain.ErrorCodeTeamExists:
		return http.StatusBadRequest
	case domain.ErrorCodePRExists:
		return http.StatusConflict
	case domain.ErrorCodePRMerged:
		return http.StatusConflict
	case domain.ErrorCodeNotAssigned:
		return http.StatusConflict
	case domain.ErrorCodeNoCandidate:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func writeDomainError(w http.ResponseWriter, derr *domain.Error) {
	status := statusByDomainCode(derr.Code)
	writeAPIError(w, status, string(derr.Code), derr.Message)
}

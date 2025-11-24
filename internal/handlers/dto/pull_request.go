package dto

import (
	"time"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type PullRequest struct {
	PullRequestID     string    `json:"pull_request_id"`
	PullRequestName   string    `json:"pull_request_name"`
	AuthorID          string    `json:"author_id"`
	Status            string    `json:"status"`
	AssignedReviewers []string  `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"createdAt,omitempty"`
	MergedAt          time.Time `json:"mergedAt,omitempty"`
}

type PullRequestShort struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func PullRequestToDTO(pr domain.PullRequest) *PullRequest {
	return &PullRequest{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: append([]string(nil), pr.Reviewers...),
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func PullRequestToShortDTO(pr domain.PullRequest) *PullRequestShort {
	return &PullRequestShort{
		PullRequestID:   pr.ID,
		PullRequestName: pr.Name,
		AuthorID:        pr.AuthorID,
		Status:          string(pr.Status),
	}
}

func PullRequestDTOToDomain(dto PullRequest) *domain.PullRequest {
	return &domain.PullRequest{
		ID:        dto.PullRequestID,
		Name:      dto.PullRequestName,
		AuthorID:  dto.AuthorID,
		Status:    domain.PRStatus(dto.Status),
		Reviewers: append([]string(nil), dto.AssignedReviewers...),
		CreatedAt: dto.CreatedAt,
		MergedAt:  dto.MergedAt,
	}
}

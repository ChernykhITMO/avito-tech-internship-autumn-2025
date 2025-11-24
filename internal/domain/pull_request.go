package domain

import (
	"context"
	"time"
)

type PullRequest struct {
	ID        string
	Name      string
	AuthorID  string
	Status    PRStatus
	Reviewers []string
	MergedAt  time.Time
	CreatedAt time.Time
}

type PRStatus string

const (
	PRStatusOpen   PRStatus = "OPEN"
	PRStatusMerged PRStatus = "MERGED"
)

type PRRepository interface {
	Create(ctx context.Context, request PullRequest) (*PullRequest, error)
	Get(ctx context.Context, id string) (*PullRequest, error)
	Update(ctx context.Context, id string, status PRStatus) error
	SetReviewers(ctx context.Context, id string, reviewers []string) error
	ListByReviewer(ctx context.Context, reviewerID string) ([]PullRequest, error)
	ListReviewers(ctx context.Context, prID string) ([]string, error)
}

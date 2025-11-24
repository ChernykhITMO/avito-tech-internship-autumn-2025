package service

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type PullRequestService interface {
	Create(ctx context.Context, id, name, authorID string) (*domain.PullRequest, error)
	Merge(ctx context.Context, id string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error)
}

var _ PullRequestService = (*pullRequestService)(nil)

type pullRequestService struct {
	prs   domain.PRRepository
	users domain.UserRepository
	teams domain.TeamRepository
}

func NewPullRequestService(prs domain.PRRepository, users domain.UserRepository, teams domain.TeamRepository) PullRequestService {
	return &pullRequestService{
		prs:   prs,
		users: users,
		teams: teams,
	}
}

func (s *pullRequestService) Create(ctx context.Context, id, name, authorID string) (*domain.PullRequest, error) {
	const maxReviewers = 2

	if id == "" {
		return nil, fmt.Errorf("create pull request: empty id")
	}
	if name == "" {
		return nil, fmt.Errorf("create pull request: empty name")
	}
	if authorID == "" {
		return nil, fmt.Errorf("create pull request: empty author id")
	}

	author, err := s.users.GetUserByID(ctx, authorID)
	if err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	if _, err := s.teams.GetByName(ctx, author.TeamName); err != nil {
		return nil, fmt.Errorf("create pull request: get team: %w", err)
	}

	candidates, err := s.users.ListReviewCandidates(ctx, author.TeamName, author.ID)
	if err != nil {
		return nil, fmt.Errorf("create pull request: list review candidates: %w", err)
	}

	if len(candidates) == 0 {
		return nil, domain.NewError(domain.ErrorCodeNoCandidate, "no review candidates found")
	}

	maxCandidates := maxReviewers
	if len(candidates) < maxCandidates {
		maxCandidates = len(candidates)
	}

	reviewerIDs := make([]string, 0, maxCandidates)
	for i := 0; i < maxCandidates; i++ {
		reviewerIDs = append(reviewerIDs, candidates[i].ID)
	}

	pr := domain.PullRequest{
		ID:       id,
		Name:     name,
		AuthorID: authorID,
		Status:   domain.PRStatusOpen,
	}

	created, err := s.prs.Create(ctx, pr)
	if err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	if err := s.prs.SetReviewers(ctx, created.ID, reviewerIDs); err != nil {
		return nil, fmt.Errorf("create pull request: set reviewers: %w", err)
	}
	created.Reviewers = reviewerIDs

	return created, nil
}

func (s *pullRequestService) Merge(ctx context.Context, id string) (*domain.PullRequest, error) {
	if id == "" {
		return nil, fmt.Errorf("merge pull request: empty id")
	}

	pr, err := s.prs.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("merge pull request: %w", err)
	}

	if pr.Status == domain.PRStatusMerged {
		return pr, nil
	}

	if err := s.prs.Update(ctx, id, domain.PRStatusMerged); err != nil {
		return nil, fmt.Errorf("merge pull request: %w", err)
	}

	mergedPR, err := s.prs.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("merge pull request: %w", err)
	}

	return mergedPR, nil
}

func (s *pullRequestService) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error) {
	if prID == "" {
		return nil, "", fmt.Errorf("reassign reviewer: empty pr id")
	}
	if oldReviewerID == "" {
		return nil, "", fmt.Errorf("reassign reviewer: empty old reviewer id")
	}

	pr, err := s.prs.Get(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("reassign reviewer: %w", err)
	}

	if pr.Status == domain.PRStatusMerged {
		return nil, "", domain.NewError(domain.ErrorCodePRMerged, "pull request already merged")
	}

	reviewers, err := s.prs.ListReviewers(ctx, prID)
	if err != nil {
		return nil, "", fmt.Errorf("reassign reviewer: list reviewers: %w", err)
	}

	const notFoundIndex = -1
	index := notFoundIndex

	for i, rID := range reviewers {
		if rID == oldReviewerID {
			index = i
			break
		}
	}

	if index == notFoundIndex {
		return nil, "", domain.NewError(domain.ErrorCodeNotAssigned, "reviewer is not assigned to this pull request")
	}

	author, err := s.users.GetUserByID(ctx, pr.AuthorID)
	if err != nil {
		return nil, "", fmt.Errorf("reassign reviewer: %w", err)
	}

	candidates, err := s.users.ListReviewCandidates(ctx, author.TeamName, author.ID)
	if err != nil {
		return nil, "", fmt.Errorf("reassign reviewer: list review candidates: %w", err)
	}

	var newReviewerID string

candidateLoop:
	for _, c := range candidates {
		if c.ID == oldReviewerID {
			continue
		}
		for _, existing := range reviewers {
			if existing == c.ID {
				continue candidateLoop
			}
		}
		newReviewerID = c.ID
		break
	}

	if newReviewerID == "" {
		return nil, "", domain.NewError(domain.ErrorCodeNoCandidate, "no replacement reviewer found")
	}

	newReviewers := make([]string, len(reviewers))
	copy(newReviewers, reviewers)
	newReviewers[index] = newReviewerID

	if err := s.prs.SetReviewers(ctx, pr.ID, newReviewers); err != nil {
		return nil, "", fmt.Errorf("reassign reviewer: set reviewers: %w", err)
	}

	pr.Reviewers = newReviewers

	return pr, newReviewerID, nil
}

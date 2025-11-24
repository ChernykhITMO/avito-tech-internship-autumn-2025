package service

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type UserService interface {
	SetIsActive(ctx context.Context, userID string, active bool) (*domain.User, error)
	GetUserReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error)
}

var _ UserService = (*userService)(nil)

type userService struct {
	users   domain.UserRepository
	pullReq domain.PRRepository
}

func NewUserService(users domain.UserRepository, pullReq domain.PRRepository) UserService {
	return &userService{
		users:   users,
		pullReq: pullReq,
	}
}

func (s *userService) SetIsActive(ctx context.Context, userID string, active bool) (*domain.User, error) {
	if userID == "" {
		return nil, fmt.Errorf("set user active: empty user id")
	}

	user, err := s.users.SetIsActive(ctx, userID, active)
	if err != nil {
		return nil, fmt.Errorf("set user active: %w", err)
	}

	return user, nil
}
func (s *userService) GetUserReviewPRs(ctx context.Context, userID string) ([]domain.PullRequest, error) {
	if userID == "" {
		return nil, fmt.Errorf("get user review PRs: empty user id")
	}

	if _, err := s.users.GetUserByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("get user review PRs: %w", err)
	}

	prs, err := s.pullReq.ListByReviewer(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user review PRs: %w", err)
	}

	return prs, nil
}

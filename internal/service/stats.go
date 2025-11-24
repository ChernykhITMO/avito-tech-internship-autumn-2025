package service

import (
	"context"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type StatsService interface {
	GetStats(ctx context.Context) (domain.StatsResponse, error)
}
type statsService struct {
	repo domain.StatsRepository
}

func NewStatsService(repo domain.StatsRepository) StatsService {
	return &statsService{
		repo: repo,
	}
}

func (s *statsService) GetStats(ctx context.Context) (domain.StatsResponse, error) {
	pr, err := s.repo.GetPRStats(ctx)
	if err != nil {
		return domain.StatsResponse{}, err
	}

	assignments, err := s.repo.GetAssignmentsStats(ctx)
	if err != nil {
		return domain.StatsResponse{}, err
	}

	return domain.StatsResponse{
		PRStats:            pr,
		AssignmentsPerUser: assignments,
	}, nil
}

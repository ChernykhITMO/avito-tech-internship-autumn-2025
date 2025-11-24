package service

import (
	"context"
	"testing"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type statsRepoMock struct {
	prStatsFn func(ctx context.Context) (domain.PRStats, error)
	assignFn  func(ctx context.Context) ([]domain.UserAssignmentStat, error)
}

func (m *statsRepoMock) GetPRStats(ctx context.Context) (domain.PRStats, error) {
	return m.prStatsFn(ctx)
}

func (m *statsRepoMock) GetAssignmentsStats(ctx context.Context) ([]domain.UserAssignmentStat, error) {
	return m.assignFn(ctx)
}

func TestStatsService_GetStats_OK(t *testing.T) {
	repo := &statsRepoMock{
		prStatsFn: func(ctx context.Context) (domain.PRStats, error) {
			return domain.PRStats{Total: 3, Open: 2, Merged: 1}, nil
		},
		assignFn: func(ctx context.Context) ([]domain.UserAssignmentStat, error) {
			return []domain.UserAssignmentStat{
				{UserID: "u1", Count: 2},
			}, nil
		},
	}

	svc := NewStatsService(repo)

	got, err := svc.GetStats(context.Background())
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if got.PRStats.Total != 3 || got.PRStats.Open != 2 || got.PRStats.Merged != 1 {
		t.Fatalf("wrong pr stats: %+v", got.PRStats)
	}
	if len(got.AssignmentsPerUser) != 1 || got.AssignmentsPerUser[0].UserID != "u1" {
		t.Fatalf("wrong assignments: %+v", got.AssignmentsPerUser)
	}
}

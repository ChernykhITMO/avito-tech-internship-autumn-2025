package service

import (
	"context"
	"errors"
	"testing"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type teamRepoMock struct {
	createFn    func(ctx context.Context, team *domain.Team) error
	getByNameFn func(ctx context.Context, name string) (*domain.Team, error)
}

func (m *teamRepoMock) Create(ctx context.Context, team *domain.Team) error {
	return m.createFn(ctx, team)
}

func (m *teamRepoMock) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	return m.getByNameFn(ctx, name)
}

type userRepoMock struct {
	saveAllFn func(ctx context.Context, users []domain.User) error
}

func (m *userRepoMock) SaveAll(ctx context.Context, users []domain.User) error {
	return m.saveAllFn(ctx, users)
}

func (m *userRepoMock) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	panic("not used")
}

func (m *userRepoMock) SetIsActive(ctx context.Context, id string, active bool) (*domain.User, error) {
	panic("not used")
}

func (m *userRepoMock) ListReviewCandidates(ctx context.Context, teamName, excludeUserID string) ([]domain.User, error) {
	panic("not used")
}

func TestTeamService_CreateTeam_OK(t *testing.T) {
	ctx := context.Background()

	teams := &teamRepoMock{
		createFn: func(ctx context.Context, team *domain.Team) error {
			if team.Name != "team1" {
				t.Fatalf("expected team name team1, got %s", team.Name)
			}
			if len(team.Members) != 2 {
				t.Fatalf("expected 2 members, got %d", len(team.Members))
			}
			return nil
		},
	}

	users := &userRepoMock{
		saveAllFn: func(ctx context.Context, members []domain.User) error {
			for _, u := range members {
				if u.TeamName != "team1" {
					t.Fatalf("expected TeamName=team1, got %s", u.TeamName)
				}
			}
			return nil
		},
	}

	svc := NewTeamService(teams, users)

	members := []domain.User{
		{ID: "1", Name: "A", IsActive: true},
		{ID: "2", Name: "B", IsActive: true},
	}

	team, err := svc.CreateTeam(ctx, "team1", members)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if team.Name != "team1" {
		t.Fatalf("wrong team name: %s", team.Name)
	}
}

func TestTeamService_CreateTeam_EmptyName(t *testing.T) {
	svc := NewTeamService(&teamRepoMock{
		createFn:    func(ctx context.Context, team *domain.Team) error { return nil },
		getByNameFn: func(ctx context.Context, name string) (*domain.Team, error) { return nil, nil },
	}, &userRepoMock{
		saveAllFn: func(ctx context.Context, users []domain.User) error { return nil },
	})

	_, err := svc.CreateTeam(context.Background(), "", []domain.User{{ID: "1"}})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestTeamService_CreateTeam_TeamRepoError(t *testing.T) {
	wantErr := errors.New("db fail")

	teams := &teamRepoMock{
		createFn:    func(ctx context.Context, team *domain.Team) error { return wantErr },
		getByNameFn: func(ctx context.Context, name string) (*domain.Team, error) { return nil, nil },
	}
	users := &userRepoMock{
		saveAllFn: func(ctx context.Context, users []domain.User) error { return nil },
	}

	svc := NewTeamService(teams, users)

	_, err := svc.CreateTeam(context.Background(), "team1", []domain.User{{ID: "1"}})
	if err == nil || !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped error %v, got %v", wantErr, err)
	}
}

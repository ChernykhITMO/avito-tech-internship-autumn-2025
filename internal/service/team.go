package service

import (
	"context"
	"fmt"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

type TeamService interface {
	CreateTeam(ctx context.Context, name string, members []domain.User) (*domain.Team, error)
	GetTeam(ctx context.Context, name string) (*domain.Team, error)
}

var _ TeamService = (*teamService)(nil)

type teamService struct {
	teams domain.TeamRepository
	users domain.UserRepository
}

func NewTeamService(teams domain.TeamRepository, users domain.UserRepository) TeamService {
	return &teamService{
		teams: teams,
		users: users,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, name string, members []domain.User) (*domain.Team, error) {
	if name == "" {
		return nil, fmt.Errorf("create team: empty team name")
	}

	if len(members) == 0 {
		return nil, fmt.Errorf("create team: empty team members")
	}

	for i := range members {
		members[i].TeamName = name
	}

	team := &domain.Team{
		Name:    name,
		Members: members,
	}

	if err := s.teams.Create(ctx, team); err != nil {
		return nil, fmt.Errorf("create team: %w", err)
	}

	if err := s.users.SaveAll(ctx, members); err != nil {
		return nil, fmt.Errorf("create team: save members: %w", err)
	}

	return team, nil
}
func (s *teamService) GetTeam(ctx context.Context, name string) (*domain.Team, error) {
	if name == "" {
		return nil, fmt.Errorf("get team: empty team name")
	}

	team, err := s.teams.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get team: %w", err)
	}

	return team, nil
}

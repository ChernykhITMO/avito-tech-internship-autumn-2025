package dto

import "github.com/ChernykhITMO/Avito/internal/domain"

type TeamMember struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func TeamToDTO(team domain.Team) Team {
	members := make([]TeamMember, 0, len(team.Members))

	for _, m := range team.Members {
		members = append(members, TeamMember{
			UserID:   m.ID,
			UserName: m.Name,
			IsActive: m.IsActive,
		})
	}

	return Team{
		TeamName: team.Name,
		Members:  members,
	}
}

func TeamDTOToDomain(teamDTO Team) domain.Team {
	members := make([]domain.User, 0, len(teamDTO.Members))

	for _, m := range teamDTO.Members {
		members = append(members, domain.User{
			ID:       m.UserID,
			Name:     m.UserName,
			TeamName: teamDTO.TeamName,
			IsActive: m.IsActive,
		})
	}

	return domain.Team{
		Name:    teamDTO.TeamName,
		Members: members,
	}
}

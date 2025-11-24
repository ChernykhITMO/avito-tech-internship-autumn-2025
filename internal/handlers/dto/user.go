package dto

import "github.com/ChernykhITMO/Avito/internal/domain"

type User struct {
	UserID   string `json:"user_id"`
	UserName string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

func UserToDTO(user domain.User) User {
	return User{
		UserID:   user.ID,
		UserName: user.Name,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func UserDTOToDomain(userDTO User) domain.User {
	return domain.User{
		ID:       userDTO.UserID,
		Name:     userDTO.UserName,
		TeamName: userDTO.TeamName,
		IsActive: userDTO.IsActive,
	}
}

package domain

import "context"

type User struct {
	ID       string
	Name     string
	TeamName string
	IsActive bool
}

type UserRepository interface {
	SaveAll(ctx context.Context, users []User) error
	GetUserByID(ctx context.Context, id string) (*User, error)
	SetIsActive(ctx context.Context, id string, active bool) (*User, error)
	ListReviewCandidates(ctx context.Context, teamName, excludeUserID string) ([]User, error)
}

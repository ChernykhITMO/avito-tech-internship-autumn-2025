package domain

import "context"

type Team struct {
	Name    string
	Members []User
}

type TeamRepository interface {
	Create(ctx context.Context, team *Team) error
	GetByName(ctx context.Context, name string) (*Team, error)
}

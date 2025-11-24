package domain

import (
	"context"
)

type PRStats struct {
	Total  int `json:"total"`
	Open   int `json:"open"`
	Merged int `json:"merged"`
}

type UserAssignmentStat struct {
	UserID string `json:"user_id"`
	Count  int    `json:"count"`
}
type StatsResponse struct {
	PRStats            PRStats              `json:"pr_stats"`
	AssignmentsPerUser []UserAssignmentStat `json:"assignments_per_user"`
}

type StatsRepository interface {
	GetPRStats(ctx context.Context) (PRStats, error)
	GetAssignmentsStats(ctx context.Context) ([]UserAssignmentStat, error)
}

package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

var _ domain.StatsRepository = (*StatsRepository)(nil)

type StatsRepository struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) domain.StatsRepository {
	return &StatsRepository{
		db: db,
	}
}
func (r *StatsRepository) GetPRStats(ctx context.Context) (domain.PRStats, error) {
	var s domain.PRStats

	query := `
        SELECT 
            COUNT(*) AS total,
            COUNT(*) FILTER (WHERE status='OPEN') AS open,
            COUNT(*) FILTER (WHERE status='MERGED') AS merged
        FROM pull_requests;
    `

	row := r.db.QueryRowContext(ctx, query)
	err := row.Scan(&s.Total, &s.Open, &s.Merged)
	return s, err
}

func (r *StatsRepository) GetAssignmentsStats(ctx context.Context) ([]domain.UserAssignmentStat, error) {
	query := `
        SELECT reviewer_id, COUNT(*)
        FROM pull_request_reviewers
        GROUP BY reviewer_id;
    `

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close error: %v", err)
		}
	}()

	var stats []domain.UserAssignmentStat
	for rows.Next() {
		var s domain.UserAssignmentStat
		if err := rows.Scan(&s.UserID, &s.Count); err != nil {
			return nil, err
		}
		stats = append(stats, s)
	}
	return stats, nil
}

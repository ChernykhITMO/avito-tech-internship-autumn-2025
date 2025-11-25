package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ChernykhITMO/Avito/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

const pgUniqueViolation = "23505"

var _ domain.TeamRepository = (*TeamRepository)(nil)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) domain.TeamRepository {
	return &TeamRepository{
		db: db,
	}
}

func (r *TeamRepository) Create(ctx context.Context, team *domain.Team) error {
	const query = `INSERT INTO teams(name) VALUES ($1)`

	_, err := r.db.ExecContext(ctx, query, team.Name)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgUniqueViolation {
			return domain.NewError(domain.ErrorCodeTeamExists, "team already exists")
		}
		return fmt.Errorf("create team: %w", err)
	}
	return nil
}

func (r *TeamRepository) GetByName(ctx context.Context, name string) (*domain.Team, error) {
	const (
		queryTeam    = `SELECT name FROM teams WHERE name = $1`
		queryMembers = `SELECT id, name, is_active FROM users WHERE team_name = $1`
	)

	var team domain.Team
	if err := r.db.QueryRowContext(ctx, queryTeam, name).Scan(&team.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewError(domain.ErrorCodeNotFound, "team not found")
		}
		return nil, fmt.Errorf("get team: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, queryMembers, team.Name)
	if err != nil {
		return nil, fmt.Errorf("get team members: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close error: %v", err)
		}
	}()

	var members []domain.User
	for rows.Next() {
		var user domain.User

		if err := rows.Scan(&user.ID, &user.Name, &user.IsActive); err != nil {
			return nil, fmt.Errorf("scan team member: %w", err)
		}
		user.TeamName = team.Name
		members = append(members, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate team members: %w", err)
	}

	team.Members = members
	return &team, nil
}

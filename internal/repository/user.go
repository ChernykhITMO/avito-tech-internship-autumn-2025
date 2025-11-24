package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

var _ domain.UserRepository = (*UserRepository)(nil)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) SaveAll(ctx context.Context, users []domain.User) error {
	const query = `
	INSERT INTO users (id, name, team_name, is_active) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (id) DO UPDATE
	SET 
	    name = EXCLUDED.name, 
		team_name = EXCLUDED.team_name, 
		is_active = EXCLUDED.is_active
	`

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx for save users: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("prepare save user stmt: %w", err)
	}
	defer stmt.Close()

	for _, u := range users {
		if _, err := stmt.ExecContext(ctx, u.ID, u.Name, u.TeamName, u.IsActive); err != nil {
			return fmt.Errorf("save user %s: %w", u.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx for save users: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	const query = `
	SELECT id, name, team_name, is_active 
	FROM users WHERE id = $1`

	var user domain.User

	if err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewError(domain.ErrorCodeNotFound, "user not found")
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) SetIsActive(ctx context.Context, id string, active bool) (*domain.User, error) {
	const query = `
	UPDATE users SET is_active = $2 
	WHERE id = $1 RETURNING id, name, team_name, is_active`

	var user domain.User
	if err := r.db.QueryRowContext(ctx, query, id, active).
		Scan(&user.ID, &user.Name, &user.TeamName, &user.IsActive); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewError(domain.ErrorCodeNotFound, "user not found")
		}
		return nil, fmt.Errorf("set user active: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) ListReviewCandidates(ctx context.Context, teamName, excludeUserID string) ([]domain.User, error) {
	const query = `
	SELECT id, name, team_name, is_active 
	FROM users 
	WHERE team_name = $1 AND is_active = true AND id <> $2`

	rows, err := r.db.QueryContext(ctx, query, teamName, excludeUserID)
	if err != nil {
		return nil, fmt.Errorf("list review candidates: %w", err)
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.TeamName, &u.IsActive); err != nil {
			return nil, fmt.Errorf("scan review candidate: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate review candidates: %w", err)
	}

	return users, nil
}

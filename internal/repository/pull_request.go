package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/ChernykhITMO/Avito/internal/domain"
)

var _ domain.PRRepository = (*PRRepository)(nil)

type PRRepository struct {
	db *sql.DB
}

func NewPRRepository(db *sql.DB) domain.PRRepository {
	return &PRRepository{
		db: db,
	}
}

func (r *PRRepository) Create(ctx context.Context, req domain.PullRequest) (*domain.PullRequest, error) {
	const query = `
    INSERT INTO pull_requests (
        pull_request_id,
        pull_request_name,
        author_id,
        status
    )
    VALUES ($1, $2, $3, $4)
    RETURNING pull_request_id, pull_request_name, author_id, status, created_at, merged_at
    `

	var (
		pr       domain.PullRequest
		mergedAt sql.NullTime
	)

	err := r.db.QueryRowContext(ctx, query,
		req.ID, req.Name, req.AuthorID, req.Status,
	).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&mergedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("create pull request: %w", err)
	}

	if mergedAt.Valid {
		pr.MergedAt = mergedAt.Time
	}

	return &pr, nil
}

func (r *PRRepository) Get(ctx context.Context, id string) (*domain.PullRequest, error) {
	const query = `
		SELECT pull_request_id, pull_request_name, author_id, status, created_at, merged_at
		FROM pull_requests
		WHERE pull_request_id = $1
	`

	var (
		pr       domain.PullRequest
		mergedAt sql.NullTime
	)

	if err := r.db.QueryRowContext(ctx, query, id).
		Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewError(domain.ErrorCodeNotFound, "pull request not found")
		}
		return nil, fmt.Errorf("get pull request: %w", err)
	}

	if mergedAt.Valid {
		pr.MergedAt = mergedAt.Time
	}

	return &pr, nil
}

func (r *PRRepository) Update(ctx context.Context, id string, status domain.PRStatus) error {
	const query = `
		UPDATE pull_requests
		SET	
		    status = $2,
		    merged_at = CASE WHEN $2 = 'MERGED' THEN NOW() ELSE merged_at END
		WHERE pull_request_id = $1
	`

	res, err := r.db.ExecContext(ctx, query, id, status)
	if err != nil {
		return fmt.Errorf("update pull request: %w", err)
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update pull request: %w", err)
	}

	if affected == 0 {
		return domain.NewError(domain.ErrorCodeNotFound, "pull request not found")
	}

	return nil
}

func (r *PRRepository) SetReviewers(ctx context.Context, id string, reviewers []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx for set reviewers: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	const queryDelete = `DELETE FROM pull_request_reviewers WHERE pull_request_id = $1`

	if _, err := tx.ExecContext(ctx, queryDelete, id); err != nil {
		return fmt.Errorf("delete reviewers: %w", err)
	}

	const queryInsert = `
	INSERT INTO pull_request_reviewers (pull_request_id, reviewer_id)
	VALUES ($1, $2)`

	for _, revID := range reviewers {
		if _, err := tx.ExecContext(ctx, queryInsert, id, revID); err != nil {
			return fmt.Errorf("insert reviewer: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx for set reviewers: %w", err)
	}

	return nil
}

func (r *PRRepository) ListByReviewer(ctx context.Context, reviewerID string) ([]domain.PullRequest, error) {
	const query = `
	SELECT pr.pull_request_id, pr.pull_request_name, pr.author_id, 
	pr.status, pr.created_at, pr.merged_at
	FROM pull_requests as pr
	JOIN pull_request_reviewers as r 
	ON r.pull_request_id = pr.pull_request_id
	WHERE r.reviewer_id = $1
`

	rows, err := r.db.QueryContext(ctx, query, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("select reviewers: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close error: %v", err)
		}
	}()

	var prs []domain.PullRequest
	for rows.Next() {
		var (
			pr       domain.PullRequest
			mergedAt sql.NullTime
		)
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &mergedAt); err != nil {
			return nil, fmt.Errorf("scan pr by reviewer: %w", err)
		}
		if mergedAt.Valid {
			pr.MergedAt = mergedAt.Time
		}
		prs = append(prs, pr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate by reviewer: %w", err)
	}

	return prs, nil
}

func (r *PRRepository) ListReviewers(ctx context.Context, prID string) ([]string, error) {
	const query = `
		SELECT reviewer_id
		FROM pull_request_reviewers
		WHERE pull_request_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, fmt.Errorf("list reviewers: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.Close error: %v", err)
		}
	}()

	var reviewers []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("scan reviewer: %w", err)
		}
		reviewers = append(reviewers, id)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate reviewers: %w", err)
	}

	return reviewers, nil
}

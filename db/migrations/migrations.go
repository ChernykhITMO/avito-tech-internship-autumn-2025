package migrations

import (
	"database/sql"
	"fmt"
)

func CreateTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            team_name TEXT,
            is_active BOOLEAN DEFAULT true
        )`,
		`CREATE TABLE IF NOT EXISTS teams (
            name TEXT PRIMARY KEY
        )`,
		`CREATE TABLE IF NOT EXISTS pull_requests (
            pull_request_id TEXT PRIMARY KEY,
            pull_request_name TEXT NOT NULL,
            author_id TEXT NOT NULL REFERENCES users(id),
            status TEXT NOT NULL,
            created_at TIMESTAMP DEFAULT NOW(),
            merged_at TIMESTAMP
        )`,
		`CREATE TABLE IF NOT EXISTS pull_request_reviewers (
            pull_request_id TEXT REFERENCES pull_requests(pull_request_id),
            reviewer_id TEXT REFERENCES users(id),
            PRIMARY KEY (pull_request_id, reviewer_id)
        )`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}
	return nil
}

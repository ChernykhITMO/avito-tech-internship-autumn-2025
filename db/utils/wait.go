package dbutils

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func WaitForDB(ctx context.Context, dsn string, maxAttempts int) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxAttempts; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			db, err = sql.Open("pgx", dsn)
			if err != nil {
				log.Printf("Failed to open database: %v, retrying...", err)
				time.Sleep(1 * time.Second)
				continue
			}

			if err = db.PingContext(ctx); err != nil {
				log.Printf("Failed to ping database (attempt %d/%d): %v", i+1, maxAttempts, err)
				db.Close()
				time.Sleep(1 * time.Second)
				continue
			}

			log.Println("Successfully connected to database")
			return db, nil
		}
	}

	return nil, err
}

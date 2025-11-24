package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ChernykhITMO/Avito/db/migrations"
	dbutils "github.com/ChernykhITMO/Avito/db/utils"

	"github.com/ChernykhITMO/Avito/internal/httpserver"
	"github.com/ChernykhITMO/Avito/internal/repository"
	"github.com/ChernykhITMO/Avito/internal/service"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const maxAttempts = 30

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		log.Fatal("DB_DSN is not set")
	}

	db, err := dbutils.WaitForDB(ctx, dsn, maxAttempts)
	if err != nil {
		log.Fatal("failed to connect to database after retries: ", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Running database migrations...")
	if err := migrations.CreateTables(db); err != nil {
		log.Fatal("Failed to run migrations: ", err)
	}
	log.Println("Migrations completed successfully")

	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPRRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	teamSvc := service.NewTeamService(teamRepo, userRepo)
	userSvc := service.NewUserService(userRepo, prRepo)
	prSvc := service.NewPullRequestService(prRepo, userRepo, teamRepo)
	statsSvc := service.NewStatsService(statsRepo)

	srv := httpserver.New(":8080", httpserver.Deps{
		TeamService:        teamSvc,
		UserService:        userSvc,
		PullRequestService: prSvc,
		StatsService:       statsSvc,
	})

	log.Println("Starting server on :8080")
	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}

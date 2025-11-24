package httpserver

import (
	"context"
	"net/http"
	"time"

	"github.com/ChernykhITMO/Avito/internal/handlers"
	"github.com/ChernykhITMO/Avito/internal/service"
)

type Deps struct {
	TeamService        service.TeamService
	UserService        service.UserService
	PullRequestService service.PullRequestService
	StatsService       service.StatsService
}

type Server struct {
	http *http.Server
}

func New(addr string, deps Deps) *Server {
	mux := http.NewServeMux()

	router := handlers.NewRouter(
		deps.TeamService,
		deps.UserService,
		deps.PullRequestService,
		deps.StatsService,
	)
	router.Register(mux)

	return &Server{
		http: &http.Server{
			Addr:    addr,
			Handler: mux,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = s.http.Shutdown(shutdownCtx)
	}()

	if err := s.http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

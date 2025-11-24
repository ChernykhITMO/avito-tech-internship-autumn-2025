package handlers

import (
	"net/http"

	"github.com/ChernykhITMO/Avito/internal/service"
)

type Router struct {
	team  *TeamHandler
	user  *UserHandler
	pr    *PullRequestHandler
	stats *StatsHandler
}

func NewRouter(
	teamSvc service.TeamService,
	userSvc service.UserService,
	prSvc service.PullRequestService,
	handler service.StatsService,
) *Router {
	return &Router{
		team:  NewTeamHandler(teamSvc),
		user:  NewUserHandler(userSvc),
		pr:    NewPullRequestHandler(prSvc),
		stats: NewStatsHandler(handler),
	}
}

func (r *Router) Register(mux *http.ServeMux) {
	r.team.Register(mux)
	r.user.Register(mux)
	r.pr.Register(mux)
	r.stats.Register(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
}

package team

import (
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

type Implementation struct {
	teamService service.TeamService
}

func NewTeamImplementation(s service.TeamService) *Implementation {
	return &Implementation{
		teamService: s,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	
	mux.HandleFunc("/team/add", i.CreateTeam)

	mux.HandleFunc("/team/get", i.GetTeam)
}

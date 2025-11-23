package stats

import (
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

type Implementation struct {
	s service.StatsService
}

func NewStatsHandler(s service.StatsService) *Implementation {
	return &Implementation{s: s}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/stats", i.GetStats)
}

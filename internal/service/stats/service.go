package stats

import (
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type serv struct {
	tr   *manager.Manager
	repo repository.StatsRepository
}

func NewStatsService(repo repository.StatsRepository, tr *manager.Manager) service.StatsService {
	return &serv{
		tr:   tr,
		repo: repo,
	}
}

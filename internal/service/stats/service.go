package stats

import (
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type serv struct {
	tr   service.TrManager
	repo repository.StatsRepository
}

func NewService(repo repository.StatsRepository, tr *manager.Manager) service.StatsService {
	return NewServiceWithTrManager(repo, service.NewManagerAdapter(tr))
}

func NewServiceWithTrManager(repo repository.StatsRepository, tr service.TrManager) service.StatsService {
	return &serv{
		tr:   tr,
		repo: repo,
	}
}

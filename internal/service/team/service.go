package team

import (
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type serv struct {
	teamRepo  repository.TeamRepository
	userRepo  repository.UserRepository
	trManager service.TrManager
}

func NewService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, trManager *manager.Manager) service.TeamService {
	return NewServiceWithTrManager(teamRepo, userRepo, service.NewManagerAdapter(trManager))
}

func NewServiceWithTrManager(teamRepo repository.TeamRepository, userRepo repository.UserRepository, trManager service.TrManager) service.TeamService {
	return &serv{
		teamRepo:  teamRepo,
		userRepo:  userRepo,
		trManager: trManager,
	}
}

package pr

import (
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type serv struct {
	trManager *manager.Manager
	teamRepo  repository.TeamRepository
	userRepo  repository.UserRepository
	prRepo    repository.PRRepository
}

func NewService(teamRepository repository.TeamRepository, userRepository repository.UserRepository, pRRepository repository.PRRepository, trManager *manager.Manager) service.PRService {
	return &serv{
		trManager: trManager,
		teamRepo:  teamRepository,
		userRepo:  userRepository,
		prRepo:    pRRepository,
	}
}

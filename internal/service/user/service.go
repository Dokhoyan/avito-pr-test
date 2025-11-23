package user

import (
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type serv struct {
	userRepo  repository.UserRepository
	trManager *manager.Manager
}

func NewService(userRepo repository.UserRepository, trManager *manager.Manager) service.UserService {
	return &serv{
		userRepo:  userRepo,
		trManager: trManager,
	}
}

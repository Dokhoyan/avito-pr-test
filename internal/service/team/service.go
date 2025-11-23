package team

import (
	"github.com/Dokhoyan/avito-pr-test/internal/repository"
	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

type TeamService struct {
	teamRepo  repository.TeamRepository
	userRepo  repository.UserRepository
	trManager *manager.Manager
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository, trManager *manager.Manager) *TeamService {
	return &TeamService{
		teamRepo:  teamRepo,
		userRepo:  userRepo,
		trManager: trManager,
	}
}
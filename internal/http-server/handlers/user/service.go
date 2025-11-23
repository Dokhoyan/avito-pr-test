package user

import (
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

type Implementation struct {
	s service.UserService
}

func NewImplementation(s service.UserService) *Implementation {
	return &Implementation{
		s: s,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/users/setIsActive", i.SetUserStatus)
	mux.HandleFunc("/users/getReview", i.GetReview)
}



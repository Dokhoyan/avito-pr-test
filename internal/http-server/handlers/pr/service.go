package pr

import (
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/service"
)

type Implementation struct {
	s service.PRService
}

func NewImplementation(s service.PRService) *Implementation {
	return &Implementation{
		s: s,
	}
}

func (i *Implementation) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/pullRequest/create", i.Create)
	mux.HandleFunc("/pullRequest/merge", i.Merge)
	mux.HandleFunc("/pullRequest/reassign", i.Reassign)
}

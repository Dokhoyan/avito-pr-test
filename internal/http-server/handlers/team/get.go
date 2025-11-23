package team

import (
	"encoding/json" 
	"errors"
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"go.uber.org/zap"
)

type GetTeamRequest struct {
	TeamName string 
}

func (i *Implementation) GetTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		logger.Warn("invalid HTTP method for Get Team",
			zap.String("method", r.Method),
			zap.String("expected", http.MethodGet))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(model.ErrorResponse{

			Error: model.ErrorBody{
				Code:    model.ErrorBadRequest,
				Message: "Method not allowed",

			},
		})
		return
	}

	var req GetTeamRequest

	req.TeamName = r.URL.Query().Get("team_name")
	if req.TeamName == "" {
		logger.Warn("invalid request payload for Get Team",
			zap.String("method", "GetTeam"))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{

			Error: model.ErrorBody{
				Code:    model.ErrorBadRequest,
				Message: "Invalid request payload",

			},
		})

		return
	}

	logger.Info("getting team",
		zap.String("team_name", req.TeamName))


	team, err := i.teamService.GetTeam(r.Context(), req.TeamName)
	if err != nil {
		if errors.Is(err, service.ErrTeamNotFound) {

			logger.Warn("team not found",
				zap.String("team_name", req.TeamName),
				zap.Error(err))
			
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    model.ErrorNotFound,
					Message: "resource not found",
				}})

			return
		}
		logger.Error("internal server error while getting team",
			zap.String("team_name", req.TeamName),
			zap.Error(err))
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{

			Error: model.ErrorBody{
				Code:    model.ErrorInternal,
				Message: "Internal server error",
				
			}})

		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(team)
}

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

type CreateTeamRequest struct {
	TeamName string             `json:"team_name"`
	Members  []model.TeamMember `json:"members"`
}

func (i *Implementation) CreateTeam(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		logger.Warn("invalid HTTP method for Create Team",
			zap.String("method", r.Method),
			zap.String("expected", http.MethodPost))

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

	var req CreateTeamRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		logger.Warn("invalid request payload for Create Team",
			zap.Error(err),
			zap.String("method", "CreateTeam"))

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

	logger.Info("creating team",
		zap.String("team_name", req.TeamName),
		zap.Int("members_count", len(req.Members)))

	team, err := i.teamService.CreateTeamWithMembers(r.Context(), req.TeamName, req.Members)
	if err != nil {
		if errors.Is(err, service.ErrTeamExists) {

			logger.Warn("team already exists",
				zap.String("team_name", req.TeamName),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorTeamExists,
					Message: "team_name already exists",
				}})

			return
		}
		logger.Error("internal server error while creating team",
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"team": team,
	})
}

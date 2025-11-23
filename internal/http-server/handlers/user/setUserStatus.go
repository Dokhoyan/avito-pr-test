package user

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"go.uber.org/zap"
)

type SetUserStatusRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

func (i *Implementation) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		logger.Warn("invalid HTTP method for Set User Status",
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

	var req SetUserStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		logger.Warn("invalid request payload for Set User Status",
			zap.Error(err),
			zap.String("method", "SetUserStatus"))

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

	logger.Info("updating user status",
		zap.String("user_id", req.UserID),
		zap.Bool("is_active", req.IsActive))

	user, err := i.s.UpdateUserStatus(r.Context(), req.UserID, req.IsActive)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {

			logger.Warn("user not found",
				zap.String("user_id", req.UserID),
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
		logger.Error("internal server error while updating user status",
			zap.String("user_id", req.UserID),
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": user,
	})
}

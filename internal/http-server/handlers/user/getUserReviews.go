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

type GetReviewsRequest struct {
	UserID string
}

func (i *Implementation) GetReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		logger.Warn("invalid HTTP method for Get Review",
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

	var req GetReviewsRequest

	req.UserID = r.URL.Query().Get("user_id")
	if req.UserID == "" {
		logger.Warn("invalid request payload for Get Review",
			zap.String("method", "GetReview"))

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

	logger.Info("getting user reviews",
		zap.String("user_id", req.UserID))

	reviews, err := i.s.GetUserReviews(r.Context(), req.UserID)
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
		logger.Error("internal server error while getting user reviews",
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
		"user_id":       req.UserID,
		"pull_requests": reviews,
	})
}

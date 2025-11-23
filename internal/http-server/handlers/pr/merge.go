package pr

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"github.com/Dokhoyan/avito-pr-test/internal/service"
	"go.uber.org/zap"
)

type MergePRRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

func (i *Implementation) Merge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		logger.Warn("invalid HTTP method for Merge PR",
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

	var req MergePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		logger.Warn("invalid request payload for Merge PR",
			zap.Error(err),
			zap.String("method", "Merge"))

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

	logger.Info("merging PR",
		zap.String("pull_request_id", req.PullRequestID))

	pr, err := i.s.MergePR(r.Context(), req.PullRequestID)
	if err != nil {
		if errors.Is(err, service.ErrPRNotFound) {

			logger.Warn("PR not found",
				zap.String("pull_request_id", req.PullRequestID),
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
		logger.Error("internal server error while merging PR",
			zap.String("pull_request_id", req.PullRequestID),
			zap.Error(err))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{

			Error: model.ErrorBody{
				Code:    model.ErrorInternal,
				Message: "Internal server error",
			},
		})

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": pr,
	})
}

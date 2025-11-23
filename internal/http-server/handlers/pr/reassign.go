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

type ReassignPRRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldReviewerID string `json:"old_reviewer_id"`
}

func (i *Implementation) Reassign(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		logger.Warn("invalid HTTP method for Reassign PR",
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

	var req ReassignPRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		logger.Warn("invalid request payload for Reassign PR",
			zap.Error(err),
			zap.String("method", "Reassign"))

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

	logger.Info("reassigning reviewer",
		zap.String("pull_request_id", req.PullRequestID),
		zap.String("old_reviewer_id", req.OldReviewerID))

	pr, newReviewerID, err := i.s.ReassignReviewer(r.Context(), req.PullRequestID, req.OldReviewerID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPRNotFound):

			logger.Warn("PR not found",
				zap.String("pull_request_id", req.PullRequestID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorNotFound,
					Message: "PR not found",
				},
			})
		case errors.Is(err, service.ErrUserNotFound):

			logger.Warn("user not found",
				zap.String("old_reviewer_id", req.OldReviewerID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.ErrorResponse{
				Error: model.ErrorBody{
					Code:    model.ErrorNotFound,
					Message: "User not found",
				},
			})
		case errors.Is(err, service.ErrPRMerged):

			logger.Warn("cannot reassign on merged PR",
				zap.String("pull_request_id", req.PullRequestID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorPRMerged,
					Message: "cannot reassign on merged PR",
				},
			})
		case errors.Is(err, service.ErrNotAssigned):

			logger.Warn("reviewer is not assigned to this PR",
				zap.String("pull_request_id", req.PullRequestID),
				zap.String("old_reviewer_id", req.OldReviewerID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorNotAssigned,
					Message: "reviewer is not assigned to this PR",
				},
			})
		case errors.Is(err, service.ErrNoCandidate):

			logger.Warn("no active replacement candidate in team",
				zap.String("pull_request_id", req.PullRequestID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorNoCandidate,
					Message: "no active replacement candidate in team",
				},
			})
		default:
			logger.Error("internal server error while reassigning reviewer",
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
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr":          pr,
		"replaced_by": newReviewerID,
	})
}

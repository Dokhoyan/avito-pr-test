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

type CreatePRRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

func (i *Implementation) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {

		logger.Warn("invalid HTTP method for Create PR",
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

	var req CreatePRRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {

		logger.Warn("invalid request payload for Create PR",
			zap.Error(err),
			zap.String("method", "Create"))

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

	logger.Info("creating PR",
		zap.String("pull_request_id", req.PullRequestID),
		zap.String("pull_request_name", req.PullRequestName),
		zap.String("author_id", req.AuthorID))

	pr, err := i.s.CreatePR(r.Context(), req.PullRequestID, req.PullRequestName, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrPRExists):

			logger.Warn("PR already exists",
				zap.String("pull_request_id", req.PullRequestID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorPRExists,
					Message: "PR already exists",
				}})
		case errors.Is(err, service.ErrUserNotFound):

			logger.Warn("author not found",
				zap.String("author_id", req.AuthorID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorNotFound,
					Message: "Author not found",
				}})
		case errors.Is(err, service.ErrTeamNotFound):

			logger.Warn("team not found",
				zap.String("author_id", req.AuthorID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorNotFound,
					Message: "Team not found",
				}})
		case errors.Is(err, service.ErrUserInactive):

			logger.Warn("author is not active",
				zap.String("author_id", req.AuthorID),
				zap.Error(err))

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(model.ErrorResponse{

				Error: model.ErrorBody{
					Code:    model.ErrorBadRequest,
					Message: "Author is not active",
				}})
		default:
			logger.Error("internal server error while creating PR",
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
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pr": pr,
	})
}

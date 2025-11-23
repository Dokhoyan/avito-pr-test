package stats

import (
	"encoding/json"
	"net/http"

	"github.com/Dokhoyan/avito-pr-test/internal/logger"
	"github.com/Dokhoyan/avito-pr-test/internal/model"
	"go.uber.org/zap"
)

func (i *Implementation) GetStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {

		logger.Warn("invalid HTTP method for Get Stats",
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

	stats, err := i.s.GetStats(r.Context())
	if err != nil {
		logger.Warn("Failed to get stats", zap.Error(err))

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
	json.NewEncoder(w).Encode(stats)
}

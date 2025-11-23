package testutil

import (
	"context"
	"testing"

	"github.com/Dokhoyan/avito-pr-test/internal/service"
	servicemocks "github.com/Dokhoyan/avito-pr-test/internal/service/mocks"
	"github.com/stretchr/testify/mock"
)

// NewTestTrManager creates a mock transaction manager for testing
// It executes the function and returns whatever error the function returns
func NewTestTrManager(t *testing.T) service.TrManager {
	mockTrManager := servicemocks.NewMockTrManager(t)
	// Setup default behavior: execute the function and return its error
	// This way, if the function returns an error, the transaction manager will return it
	mockTrManager.EXPECT().Do(mock.Anything, mock.AnythingOfType("func(context.Context) error")).
		RunAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}).
		Maybe()
	return mockTrManager
}

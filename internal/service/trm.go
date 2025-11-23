package service

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

// TrManager is an interface for transaction manager to allow mocking in tests
type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// ManagerAdapter adapts *manager.Manager to TrManager interface
type ManagerAdapter struct {
	*manager.Manager
}

// Do implements TrManager interface
func (a *ManagerAdapter) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return a.Manager.Do(ctx, fn)
}

// NewManagerAdapter creates an adapter from *manager.Manager to TrManager
func NewManagerAdapter(m *manager.Manager) TrManager {
	return &ManagerAdapter{Manager: m}
}

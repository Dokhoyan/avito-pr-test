package service

import (
	"context"

	"github.com/avito-tech/go-transaction-manager/trm/v2/manager"
)

// TrManager - интерфейс для менеджера транзакций, позволяющий использовать моки в тестах
type TrManager interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}

// ManagerAdapter адаптирует *manager.Manager к интерфейсу TrManager
type ManagerAdapter struct {
	*manager.Manager
}

// Do реализует интерфейс TrManager
func (a *ManagerAdapter) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return a.Manager.Do(ctx, fn)
}

// NewManagerAdapter создает адаптер из *manager.Manager в TrManager
func NewManagerAdapter(m *manager.Manager) TrManager {
	return &ManagerAdapter{Manager: m}
}

package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
)

// SyncPublisherMock implementa ports.ISyncPublisher para tests
type SyncPublisherMock struct {
	PublishInventorySyncFn func(ctx context.Context, msg ports.InventorySyncMessage) error
	// Registra las llamadas realizadas para poder verificar en tests
	Calls []ports.InventorySyncMessage
}

func (m *SyncPublisherMock) PublishInventorySync(ctx context.Context, msg ports.InventorySyncMessage) error {
	m.Calls = append(m.Calls, msg)
	if m.PublishInventorySyncFn != nil {
		return m.PublishInventorySyncFn(ctx, msg)
	}
	return nil
}

// InventoryEventPublisherMock implementa ports.IInventoryEventPublisher para tests
type InventoryEventPublisherMock struct {
	PublishInventoryEventFn func(ctx context.Context, event ports.InventoryEvent) error
	Calls                   []ports.InventoryEvent
}

func (m *InventoryEventPublisherMock) PublishInventoryEvent(ctx context.Context, event ports.InventoryEvent) error {
	if m == nil {
		return nil
	}
	m.Calls = append(m.Calls, event)
	if m.PublishInventoryEventFn != nil {
		return m.PublishInventoryEventFn(ctx, event)
	}
	return nil
}

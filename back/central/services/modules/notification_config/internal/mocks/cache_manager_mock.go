package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// CacheManagerMock - Mock del cache manager de configuraciones de notificaciones
type CacheManagerMock struct {
	WarmupCacheFn                      func(ctx context.Context) error
	CacheConfigFn                      func(ctx context.Context, config *entities.IntegrationNotificationConfig) error
	UpdateConfigInCacheFn              func(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error
	RemoveConfigFromCacheFn            func(ctx context.Context, config *entities.IntegrationNotificationConfig) error
	InvalidateConfigsByIntegrationFn   func(ctx context.Context, integrationID uint) error
	InvalidateAllFn                    func(ctx context.Context) error
}

func (m *CacheManagerMock) WarmupCache(ctx context.Context) error {
	if m.WarmupCacheFn != nil {
		return m.WarmupCacheFn(ctx)
	}
	return nil
}

func (m *CacheManagerMock) CacheConfig(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	if m.CacheConfigFn != nil {
		return m.CacheConfigFn(ctx, config)
	}
	return nil
}

func (m *CacheManagerMock) UpdateConfigInCache(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error {
	if m.UpdateConfigInCacheFn != nil {
		return m.UpdateConfigInCacheFn(ctx, oldConfig, newConfig)
	}
	return nil
}

func (m *CacheManagerMock) RemoveConfigFromCache(ctx context.Context, config *entities.IntegrationNotificationConfig) error {
	if m.RemoveConfigFromCacheFn != nil {
		return m.RemoveConfigFromCacheFn(ctx, config)
	}
	return nil
}

func (m *CacheManagerMock) InvalidateConfigsByIntegration(ctx context.Context, integrationID uint) error {
	if m.InvalidateConfigsByIntegrationFn != nil {
		return m.InvalidateConfigsByIntegrationFn(ctx, integrationID)
	}
	return nil
}

func (m *CacheManagerMock) InvalidateAll(ctx context.Context) error {
	if m.InvalidateAllFn != nil {
		return m.InvalidateAllFn(ctx)
	}
	return nil
}

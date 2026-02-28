package cache

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// ICacheManager define la interfaz para la gestión del cache de configuraciones
type ICacheManager interface {
	// WarmupCache carga todas las configuraciones activas en Redis al iniciar
	WarmupCache(ctx context.Context) error

	// CacheConfig cachea una configuración individual después de crearla
	CacheConfig(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	// UpdateConfigInCache actualiza una configuración en cache
	UpdateConfigInCache(ctx context.Context, oldConfig, newConfig *entities.IntegrationNotificationConfig) error

	// RemoveConfigFromCache elimina una configuración del cache
	RemoveConfigFromCache(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	// InvalidateConfigsByIntegration invalida todas las configs de una integración
	InvalidateConfigsByIntegration(ctx context.Context, integrationID uint) error

	// InvalidateAll invalida todo el cache
	InvalidateAll(ctx context.Context) error
}

// New crea una nueva instancia del cache manager
func New(redis redis.IRedis, repo ports.IRepository, orderStatusQuerier IOrderStatusQuerier, logger log.ILogger) ICacheManager {
	return &cacheManager{
		redis:              redis,
		repo:               repo,
		orderStatusQuerier: orderStatusQuerier,
		logger:             logger.WithModule("notification-config-cache"),
	}
}


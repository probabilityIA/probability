package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// ICacheManager define el contrato para la gestión del cache de configuraciones.
// Pertenece al dominio: es un puerto de salida (driven port) que la capa de
// infraestructura implementa.
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

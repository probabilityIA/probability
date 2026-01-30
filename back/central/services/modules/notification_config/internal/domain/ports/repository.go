package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// IRepository define el contrato del repositorio de configuraciones de notificaciones
type IRepository interface {
	// Create crea una nueva configuración
	Create(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	// Update actualiza una configuración existente
	Update(ctx context.Context, config *entities.IntegrationNotificationConfig) error

	// GetByID obtiene una configuración por su ID
	GetByID(ctx context.Context, id uint) (*entities.IntegrationNotificationConfig, error)

	// List obtiene una lista de configuraciones con filtros opcionales
	List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]entities.IntegrationNotificationConfig, error)

	// Delete elimina una configuración por su ID
	Delete(ctx context.Context, id uint) error

	// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas por integración y trigger
	// Ordenadas por prioridad descendente
	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.IntegrationNotificationConfig, error)
}

package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// IUseCase define el contrato de los casos de uso
type IUseCase interface {
	// Create crea una nueva configuración
	Create(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)

	// Update actualiza una configuración existente
	Update(ctx context.Context, id uint, dto dtos.UpdateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error)

	// GetByID obtiene una configuración por su ID
	GetByID(ctx context.Context, id uint) (*dtos.NotificationConfigResponseDTO, error)

	// List obtiene una lista de configuraciones con filtros
	List(ctx context.Context, filters dtos.FilterNotificationConfigDTO) ([]dtos.NotificationConfigResponseDTO, error)

	// Delete elimina una configuración
	Delete(ctx context.Context, id uint) error

	// ValidateConditions valida si una orden cumple las condiciones de una configuración
	ValidateConditions(config *entities.IntegrationNotificationConfig, orderStatus string, paymentMethodID uint) bool
}

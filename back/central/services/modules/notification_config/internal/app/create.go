package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

// Create crea una nueva configuración de notificación
func (uc *useCase) Create(ctx context.Context, dto dtos.CreateNotificationConfigDTO) (*dtos.NotificationConfigResponseDTO, error) {
	// Validar que la configuración no esté duplicada
	filters := dtos.FilterNotificationConfigDTO{
		IntegrationID:    &dto.IntegrationID,
		NotificationType: &dto.NotificationType,
		Trigger:          &dto.Conditions.Trigger,
	}

	existing, err := uc.repository.List(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error checking for duplicate configs")
		return nil, err
	}

	// Validar duplicados con condiciones similares
	for _, config := range existing {
		if uc.areSimilarConditions(&config.Conditions, &dto.Conditions) {
			uc.logger.Warn().
				Uint("integration_id", dto.IntegrationID).
				Str("trigger", dto.Conditions.Trigger).
				Msg("Similar notification config already exists")
			return nil, errors.ErrDuplicateConfig
		}
	}

	// Crear entidad de dominio
	entity := &entities.IntegrationNotificationConfig{
		IntegrationID:    dto.IntegrationID,
		NotificationType: dto.NotificationType,
		IsActive:         dto.IsActive,
		Conditions:       dto.Conditions,
		Config:           dto.Config,
		Description:      dto.Description,
		Priority:         dto.Priority,
	}

	// Persistir
	if err := uc.repository.Create(ctx, entity); err != nil {
		uc.logger.Error().Err(err).Msg("Error creating notification config")
		return nil, err
	}

	uc.logger.Info().
		Uint("id", entity.ID).
		Uint("integration_id", entity.IntegrationID).
		Str("type", entity.NotificationType).
		Msg("Notification config created successfully")

	return mappers.ToResponseDTO(entity), nil
}

// areSimilarConditions compara si dos configuraciones tienen condiciones similares
func (uc *useCase) areSimilarConditions(c1, c2 *entities.NotificationConditions) bool {
	// Si tienen el mismo trigger y los mismos filtros, son duplicadas
	if c1.Trigger != c2.Trigger {
		return false
	}

	// Comparar statuses
	if len(c1.Statuses) != len(c2.Statuses) {
		return false
	}
	statusMap := make(map[string]bool)
	for _, s := range c1.Statuses {
		statusMap[s] = true
	}
	for _, s := range c2.Statuses {
		if !statusMap[s] {
			return false
		}
	}

	// Comparar payment methods
	if len(c1.PaymentMethods) != len(c2.PaymentMethods) {
		return false
	}
	pmMap := make(map[uint]bool)
	for _, pm := range c1.PaymentMethods {
		pmMap[pm] = true
	}
	for _, pm := range c2.PaymentMethods {
		if !pmMap[pm] {
			return false
		}
	}

	return true
}

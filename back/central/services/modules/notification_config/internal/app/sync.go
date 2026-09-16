package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/app/mappers"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/errors"
)

func (uc *useCase) SyncByIntegration(ctx context.Context, dto dtos.SyncNotificationConfigsDTO) (*dtos.SyncNotificationConfigsResponseDTO, error) {
	uc.logger.Info().
		Uint("business_id", dto.BusinessID).
		Uint("integration_id", dto.IntegrationID).
		Int("rules_count", len(dto.Rules)).
		Msg("🔄 Sync notification configs by integration")

	businessID := dto.BusinessID
	filters := dtos.FilterNotificationConfigDTO{
		BusinessID:    &businessID,
		IntegrationID: &dto.IntegrationID,
	}

	existing, err := uc.repository.List(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error listing existing configs for sync")
		return nil, err
	}

	existingMap := make(map[uint]*entities.IntegrationNotificationConfig, len(existing))
	for i := range existing {
		existingMap[existing[i].ID] = &existing[i]
	}

	seen := make(map[string]bool)
	for _, rule := range dto.Rules {
		key := fmt.Sprintf("%d:%d", rule.NotificationTypeID, rule.NotificationEventTypeID)
		if seen[key] {
			uc.logger.Warn().
				Str("key", key).
				Msg("❌ Duplicate rule in sync request")
			return nil, domainerrors.ErrDuplicateConfig
		}
		seen[key] = true
	}

	var toCreate []*entities.IntegrationNotificationConfig
	var toUpdate []*entities.IntegrationNotificationConfig
	incomingIDs := make(map[uint]bool)

	for _, rule := range dto.Rules {
		if rule.ID != nil && *rule.ID > 0 {
			incomingIDs[*rule.ID] = true
			if _, exists := existingMap[*rule.ID]; !exists {
				uc.logger.Warn().
					Uint("rule_id", *rule.ID).
					Msg("❌ Rule ID not found for update")
				return nil, domainerrors.ErrNotificationConfigNotFound
			}
			entity := &entities.IntegrationNotificationConfig{
				ID:                      *rule.ID,
				BusinessID:              &dto.BusinessID,
				IntegrationID:           dto.IntegrationID,
				NotificationTypeID:      rule.NotificationTypeID,
				NotificationEventTypeID: rule.NotificationEventTypeID,
				Enabled:                 rule.Enabled,
				Description:             rule.Description,
				OrderStatusIDs:          rule.OrderStatusIDs,
			}
			toUpdate = append(toUpdate, entity)
		} else {
			entity := &entities.IntegrationNotificationConfig{
				BusinessID:              &dto.BusinessID,
				IntegrationID:           dto.IntegrationID,
				NotificationTypeID:      rule.NotificationTypeID,
				NotificationEventTypeID: rule.NotificationEventTypeID,
				Enabled:                 rule.Enabled,
				Description:             rule.Description,
				OrderStatusIDs:          rule.OrderStatusIDs,
			}
			toCreate = append(toCreate, entity)
		}
	}

	var toDeleteIDs []uint
	for id := range existingMap {
		if !incomingIDs[id] {
			toDeleteIDs = append(toDeleteIDs, id)
		}
	}

	uc.logger.Info().
		Int("to_create", len(toCreate)).
		Int("to_update", len(toUpdate)).
		Int("to_delete", len(toDeleteIDs)).
		Msg("📋 Sync classification complete")

	if err := uc.repository.SyncConfigs(ctx, dto.BusinessID, dto.IntegrationID, toCreate, toUpdate, toDeleteIDs); err != nil {
		uc.logger.Error().Err(err).Msg("❌ Error executing sync transaction")
		return nil, err
	}

	for i := range existing {
		if existing[i].NotificationTypeID == AssistantChannelID {
			if err := uc.cacheManager.RemoveConfigFromCache(ctx, &existing[i]); err != nil {
				uc.logger.Warn().Err(err).Uint("config_id", existing[i].ID).Msg("Error limpiando regla del asistente de la cache")
			}
		}
	}

	if err := uc.cacheManager.InvalidateConfigsByIntegration(ctx, dto.IntegrationID); err != nil {
		uc.logger.Warn().
			Err(err).
			Uint("integration_id", dto.IntegrationID).
			Msg("⚠️ Error invalidating cache after sync - cache may be stale")
	}

	updated, err := uc.repository.List(ctx, filters)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error re-fetching configs after sync")
		return nil, err
	}

	for i := range updated {
		if updated[i].Enabled {
			if cacheErr := uc.cacheManager.CacheConfig(ctx, &updated[i]); cacheErr != nil {
				uc.logger.Warn().
					Err(cacheErr).
					Uint("config_id", updated[i].ID).
					Msg("⚠️ Error re-caching config after sync")
			}
		}
	}

	responseDTOs := mappers.ToResponseDTOList(updated)

	response := &dtos.SyncNotificationConfigsResponseDTO{
		Created: len(toCreate),
		Updated: len(toUpdate),
		Deleted: len(toDeleteIDs),
		Configs: responseDTOs,
	}

	uc.logger.Info().
		Int("created", response.Created).
		Int("updated", response.Updated).
		Int("deleted", response.Deleted).
		Int("total_configs", len(response.Configs)).
		Msg("✅ Sync completed successfully")

	return response, nil
}

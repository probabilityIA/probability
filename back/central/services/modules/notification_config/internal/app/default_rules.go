package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
)

const (
	AssistantChannelID   uint = 6
	assistantChannelCode      = "assistant"
)

type DefaultRule struct {
	ChannelCode  string
	EventCode    string
	StatusCodes  []string
	BusinessWide bool
}

var DefaultNotificationRules = []DefaultRule{
	{
		ChannelCode:  assistantChannelCode,
		EventCode:    "order.status_changed",
		StatusCodes:  []string{"cancelled", "cancel_requested"},
		BusinessWide: true,
	},
}

func (uc *useCase) SetDefaultRulesQuerier(querier ports.IDefaultRulesQuerier) {
	uc.defaultRules = querier
}

func (uc *useCase) EnsureDefaultRules(ctx context.Context, businessID uint) error {
	if uc.defaultRules == nil || businessID == 0 {
		return nil
	}
	integrationID, err := uc.defaultRules.PlatformIntegrationID(ctx, businessID)
	if err != nil {
		return fmt.Errorf("buscar integracion plataforma: %w", err)
	}
	if integrationID == 0 {
		uc.logger.Warn().Uint("business_id", businessID).Msg("Negocio sin integracion Plataforma: no se crean reglas predeterminadas")
		return nil
	}

	existing, err := uc.repository.List(ctx, dtos.FilterNotificationConfigDTO{BusinessID: &businessID})
	if err != nil {
		return fmt.Errorf("listar reglas del negocio: %w", err)
	}

	var toCreate []*entities.IntegrationNotificationConfig
	for _, rule := range DefaultNotificationRules {
		config, err := uc.buildDefaultRule(ctx, businessID, integrationID, rule)
		if err != nil {
			return err
		}
		if config == nil || hasRule(existing, config, rule.BusinessWide) {
			continue
		}
		toCreate = append(toCreate, config)
	}
	if len(toCreate) == 0 {
		return nil
	}

	if err := uc.repository.SyncConfigs(ctx, businessID, integrationID, toCreate, nil, nil); err != nil {
		return fmt.Errorf("crear reglas predeterminadas: %w", err)
	}
	for _, config := range toCreate {
		if err := uc.cacheManager.CacheConfig(ctx, config); err != nil {
			uc.logger.Warn().Err(err).Uint("config_id", config.ID).Msg("Error cacheando regla predeterminada")
		}
	}
	uc.logger.Info().Uint("business_id", businessID).Int("created", len(toCreate)).Msg("Reglas de notificacion predeterminadas creadas")
	return nil
}

func (uc *useCase) buildDefaultRule(ctx context.Context, businessID, integrationID uint, rule DefaultRule) (*entities.IntegrationNotificationConfig, error) {
	channel, err := uc.notificationTypeRepo.GetByCode(ctx, rule.ChannelCode)
	if err != nil || channel == nil {
		uc.logger.Warn().Err(err).Str("channel", rule.ChannelCode).Msg("Canal de la regla predeterminada no existe")
		return nil, nil
	}
	events, err := uc.notificationEventRepo.GetByNotificationType(ctx, channel.ID)
	if err != nil {
		return nil, fmt.Errorf("listar eventos del canal %s: %w", rule.ChannelCode, err)
	}
	var eventID uint
	for _, event := range events {
		if event.EventCode == rule.EventCode {
			eventID = event.ID
			break
		}
	}
	if eventID == 0 {
		uc.logger.Warn().Str("channel", rule.ChannelCode).Str("event", rule.EventCode).Msg("Evento de la regla predeterminada no existe")
		return nil, nil
	}
	var statusIDs []uint
	if len(rule.StatusCodes) > 0 {
		statusIDs, err = uc.defaultRules.OrderStatusIDsByCodes(ctx, rule.StatusCodes)
		if err != nil {
			return nil, fmt.Errorf("resolver estados de la regla predeterminada: %w", err)
		}
		if len(statusIDs) == 0 {
			uc.logger.Warn().Strs("statuses", rule.StatusCodes).Msg("Estados de la regla predeterminada no existen")
			return nil, nil
		}
	}
	return &entities.IntegrationNotificationConfig{
		BusinessID:              &businessID,
		IntegrationID:           integrationID,
		NotificationTypeID:      channel.ID,
		NotificationEventTypeID: eventID,
		Enabled:                 true,
		Description:             "Regla predeterminada del sistema",
		OrderStatusIDs:          statusIDs,
	}, nil
}

func hasRule(existing []entities.IntegrationNotificationConfig, candidate *entities.IntegrationNotificationConfig, businessWide bool) bool {
	for _, config := range existing {
		if config.NotificationTypeID != candidate.NotificationTypeID || config.NotificationEventTypeID != candidate.NotificationEventTypeID {
			continue
		}
		if businessWide || config.IntegrationID == candidate.IntegrationID {
			return true
		}
	}
	return false
}

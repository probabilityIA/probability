package repository

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/migration/shared/models"
)

type notificationConfigRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewNotificationConfigRepository crea una nueva instancia del repositorio
func NewNotificationConfigRepository(database db.IDatabase, logger log.ILogger) ports.INotificationConfigRepository {
	return &notificationConfigRepository{
		db:     database,
		logger: logger.WithModule("notification_config_repository"),
	}
}

// GetActiveConfigsByIntegrationAndTrigger obtiene configuraciones activas
func (a *notificationConfigRepository) GetActiveConfigsByIntegrationAndTrigger(
	ctx context.Context,
	integrationID uint,
	trigger string,
) ([]dtos.NotificationConfigData, error) {
	var results []struct {
		ID               uint
		IntegrationID    uint
		NotificationType string
		IsActive         bool
		Conditions       string // JSON string
		Config           string // JSON string
		Priority         int
		Description      string
	}

	err := a.db.Conn(ctx).
		Model(&models.IntegrationNotificationConfig{}).
		Where("integration_id = ? AND is_active = ? AND conditions->>'trigger' = ?", integrationID, true, trigger).
		Order("priority DESC").
		Find(&results).Error

	if err != nil {
		a.logger.Error().Err(err).Msg("Error querying notification configs")
		return nil, err
	}

	// Convertir a DTOs parseando JSON
	configs := make([]dtos.NotificationConfigData, 0, len(results))
	for _, r := range results {
		// Parse conditions JSON
		var conditions struct {
			Trigger             string   `json:"trigger"`
			Statuses            []string `json:"statuses"`
			PaymentMethods      []uint   `json:"payment_methods"`
			SourceIntegrationID *uint    `json:"source_integration_id"`
		}
		if err := json.Unmarshal([]byte(r.Conditions), &conditions); err != nil {
			a.logger.Warn().Err(err).Uint("config_id", r.ID).Msg("Error parsing conditions JSON")
			continue
		}

		// Parse config JSON
		var config struct {
			TemplateName  string `json:"template_name"`
			RecipientType string `json:"recipient_type"`
			Language      string `json:"language"`
		}
		if err := json.Unmarshal([]byte(r.Config), &config); err != nil {
			a.logger.Warn().Err(err).Uint("config_id", r.ID).Msg("Error parsing config JSON")
			continue
		}

		configs = append(configs, dtos.NotificationConfigData{
			ID:                  r.ID,
			IntegrationID:       r.IntegrationID,
			NotificationType:    r.NotificationType,
			IsActive:            r.IsActive,
			Trigger:             conditions.Trigger,
			Statuses:            conditions.Statuses,
			PaymentMethods:      conditions.PaymentMethods,
			SourceIntegrationID: conditions.SourceIntegrationID,
			TemplateName:        config.TemplateName,
			RecipientType:       config.RecipientType,
			Language:            config.Language,
			Priority:            r.Priority,
			Description:         r.Description,
		})
	}

	return configs, nil
}

// ValidateConditions valida si una orden cumple las condiciones
func (a *notificationConfigRepository) ValidateConditions(
	config *dtos.NotificationConfigData,
	orderStatus string,
	paymentMethodID uint,
	sourceIntegrationID uint,
) bool {
	// 1. Validar source_integration_id PRIMERO (más específico)
	if config.SourceIntegrationID != nil {
		// Si la config especifica una integración origen, DEBE coincidir
		if *config.SourceIntegrationID != sourceIntegrationID {
			return false
		}
	}
	// Si config.SourceIntegrationID == nil → aplica a todas las integraciones

	// 2. Validar statuses
	if len(config.Statuses) > 0 {
		statusMatch := false
		for _, status := range config.Statuses {
			if status == orderStatus {
				statusMatch = true
				break
			}
		}
		if !statusMatch {
			return false
		}
	}

	// 3. Validar payment methods
	if len(config.PaymentMethods) > 0 {
		pmMatch := false
		for _, pm := range config.PaymentMethods {
			if pm == paymentMethodID {
				pmMatch = true
				break
			}
		}
		if !pmMatch {
			return false
		}
	}

	return true
}

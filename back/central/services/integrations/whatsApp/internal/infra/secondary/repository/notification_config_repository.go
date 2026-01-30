package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// NotificationConfigData contiene los datos de una configuración de notificación
type NotificationConfigData struct {
	ID               uint
	IntegrationID    uint
	NotificationType string
	IsActive         bool
	TemplateName     string
	Language         string
	RecipientType    string
	Trigger          string
	Statuses         []string
	PaymentMethods   []uint
	Priority         int
	Description      string
}

// INotificationConfigRepository define la interfaz para consultar configuraciones de notificación
type INotificationConfigRepository interface {
	GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]NotificationConfigData, error)
	ValidateConditions(config *NotificationConfigData, orderStatus string, paymentMethodID uint) bool
}

type notificationConfigRepository struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewNotificationConfigRepository crea una nueva instancia del repositorio
func NewNotificationConfigRepository(database db.IDatabase, logger log.ILogger) INotificationConfigRepository {
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
) ([]NotificationConfigData, error) {
	var results []struct {
		ID               uint
		IntegrationID    uint
		NotificationType string
		IsActive         bool
		Conditions       string // JSON
		Config           string // JSON
		Priority         int
		Description      string
	}

	err := a.db.Conn(ctx).
		Table("integration_notification_configs").
		Where("integration_id = ? AND is_active = ? AND conditions->>'trigger' = ?", integrationID, true, trigger).
		Order("priority DESC").
		Find(&results).Error

	if err != nil {
		a.logger.Error().Err(err).Msg("Error querying notification configs")
		return nil, err
	}

	// Convertir a DTOs simples
	configs := make([]NotificationConfigData, 0, len(results))
	for _, r := range results {
		// Parse JSON fields manualmente o usando unmarshaling simple
		// Por simplicidad, vamos a asumir que los campos vienen en formato esperado
		configs = append(configs, NotificationConfigData{
			ID:               r.ID,
			IntegrationID:    r.IntegrationID,
			NotificationType: r.NotificationType,
			IsActive:         r.IsActive,
			Trigger:          trigger,
			Priority:         r.Priority,
			Description:      r.Description,
			// Los campos JSON requieren parsing adicional
			// Por ahora los dejamos vacíos, implementar según necesidad
		})
	}

	return configs, nil
}

// ValidateConditions valida si una orden cumple las condiciones
func (a *notificationConfigRepository) ValidateConditions(
	config *NotificationConfigData,
	orderStatus string,
	paymentMethodID uint,
) bool {
	// Validar statuses
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

	// Validar payment methods
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

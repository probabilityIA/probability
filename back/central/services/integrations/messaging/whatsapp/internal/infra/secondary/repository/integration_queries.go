package repository

import (
	"context"
	"errors"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/consumer/consumerevent/request"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// integrationQueries implementa consultas de integraciones para el consumer
type integrationQueries struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewIntegrationQueries crea una nueva instancia de consultas de integraciones
func NewIntegrationQueries(database db.IDatabase, logger log.ILogger) request.IntegrationRepository {
	return &integrationQueries{
		db:     database,
		logger: logger.WithModule("integration_queries"),
	}
}

// GetWhatsAppByBusinessID obtiene la integración de WhatsApp de un business
func (a *integrationQueries) GetWhatsAppByBusinessID(ctx context.Context, businessID uint) (*request.IntegrationData, error) {
	var integration struct {
		ID         uint
		BusinessID uint
		IsActive   bool
	}

	// Buscar integración de WhatsApp
	// Asumimos que el integration_type_id de WhatsApp es conocido o podemos hacer join
	err := a.db.Conn(ctx).
		Table("integrations").
		Select("integrations.id, integrations.business_id, integrations.is_active").
		Joins("INNER JOIN integration_types ON integrations.integration_type_id = integration_types.id").
		Where("integration_types.code = ?", "whatsapp").
		Where("integrations.business_id = ?", businessID).
		Where("integrations.is_active = ?", true).
		First(&integration).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("whatsapp integration not found")
		}
		a.logger.Error().Err(err).Uint("business_id", businessID).Msg("Error getting WhatsApp integration")
		return nil, err
	}

	return &request.IntegrationData{
		ID:         integration.ID,
		BusinessID: integration.BusinessID,
		IsActive:   integration.IsActive,
	}, nil
}

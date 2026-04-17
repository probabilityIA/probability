package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
)

func ToDBWebhookLog(log *domain.WebhookLog) *models.WebhookLog {
	return &models.WebhookLog{
		ID:             log.ID,
		CreatedAt:      log.CreatedAt,
		Source:         log.Source,
		EventType:      log.EventType,
		URL:            log.URL,
		Method:         log.Method,
		Headers:        datatypes.JSON(log.Headers),
		RequestBody:    datatypes.JSON(log.RequestBody),
		RemoteIP:       log.RemoteIP,
		Status:         log.Status,
		ResponseCode:   log.ResponseCode,
		ProcessedAt:    log.ProcessedAt,
		ErrorMessage:   log.ErrorMessage,
		ShipmentID:     log.ShipmentID,
		BusinessID:     log.BusinessID,
		CorrelationID:  log.CorrelationID,
		TrackingNumber: log.TrackingNumber,
		MappedStatus:   log.MappedStatus,
		RawStatus:      log.RawStatus,
	}
}

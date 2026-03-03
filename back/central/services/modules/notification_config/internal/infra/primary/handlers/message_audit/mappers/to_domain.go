package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/infra/primary/handlers/message_audit/request"
)

// ListRequestToDomain convierte query params HTTP a DTO de dominio
func ListRequestToDomain(req *request.ListMessageAudit, businessID uint) dtos.MessageAuditFilterDTO {
	return dtos.MessageAuditFilterDTO{
		BusinessID:   businessID,
		Status:       req.Status,
		Direction:    req.Direction,
		TemplateName: req.TemplateName,
		DateFrom:     req.DateFrom,
		DateTo:       req.DateTo,
		Page:         req.Page,
		PageSize:     req.PageSize,
	}
}

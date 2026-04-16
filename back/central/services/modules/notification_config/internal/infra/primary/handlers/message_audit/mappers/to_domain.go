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

// ConversationListRequestToDomain convierte query params de conversaciones a DTO de dominio
func ConversationListRequestToDomain(req *request.ListConversations, businessID uint) dtos.ConversationListFilterDTO {
	return dtos.ConversationListFilterDTO{
		BusinessID: businessID,
		State:      req.State,
		Phone:      req.Phone,
		DateFrom:   req.DateFrom,
		DateTo:     req.DateTo,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}
}

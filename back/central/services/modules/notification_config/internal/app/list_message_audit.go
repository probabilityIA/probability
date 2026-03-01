package app

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// ListMessageAudit obtiene logs de auditoría de mensajes con filtros y paginación
func (uc *useCase) ListMessageAudit(ctx context.Context, filter dtos.MessageAuditFilterDTO) (*dtos.PaginatedMessageAuditResponseDTO, error) {
	// Validar y aplicar defaults de paginación
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 || filter.PageSize > 100 {
		filter.PageSize = 20
	}

	logs, total, err := uc.messageAuditQuerier.ListMessageLogs(ctx, filter)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error listing message audit logs")
		return nil, err
	}

	// Map to response DTOs
	data := make([]dtos.MessageAuditLogResponseDTO, len(logs))
	for i, log := range logs {
		dto := dtos.MessageAuditLogResponseDTO{
			ID:             log.ID,
			ConversationID: log.ConversationID,
			MessageID:      log.MessageID,
			Direction:      log.Direction,
			TemplateName:   log.TemplateName,
			Content:        log.Content,
			Status:         log.Status,
			CreatedAt:      log.CreatedAt.Format(time.RFC3339),
			PhoneNumber:    log.PhoneNumber,
			OrderNumber:    log.OrderNumber,
			BusinessID:     log.BusinessID,
		}
		if log.DeliveredAt != nil {
			formatted := log.DeliveredAt.Format(time.RFC3339)
			dto.DeliveredAt = &formatted
		}
		if log.ReadAt != nil {
			formatted := log.ReadAt.Format(time.RFC3339)
			dto.ReadAt = &formatted
		}
		data[i] = dto
	}

	totalPages := int(total) / filter.PageSize
	if int(total)%filter.PageSize != 0 {
		totalPages++
	}

	return &dtos.PaginatedMessageAuditResponseDTO{
		Data:       data,
		Total:      total,
		Page:       filter.Page,
		PageSize:   filter.PageSize,
		TotalPages: totalPages,
	}, nil
}

package app

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
)

// GetMessageAuditStats obtiene estadísticas agregadas de mensajes outbound
func (uc *useCase) GetMessageAuditStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*dtos.MessageAuditStatsResponseDTO, error) {
	stats, err := uc.messageAuditQuerier.GetMessageStats(ctx, businessID, dateFrom, dateTo)
	if err != nil {
		uc.logger.Error().Err(err).Msg("Error getting message audit stats")
		return nil, err
	}

	return &dtos.MessageAuditStatsResponseDTO{
		TotalSent:      stats.TotalSent,
		TotalDelivered: stats.TotalDelivered,
		TotalRead:      stats.TotalRead,
		TotalFailed:    stats.TotalFailed,
		SuccessRate:    stats.SuccessRate,
	}, nil
}

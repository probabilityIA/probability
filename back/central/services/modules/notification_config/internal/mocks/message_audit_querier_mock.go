package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

// MessageAuditQuerierMock - Mock de la interfaz IMessageAuditQuerier
type MessageAuditQuerierMock struct {
	ListMessageLogsFn  func(ctx context.Context, filter dtos.MessageAuditFilterDTO) ([]entities.MessageAuditLog, int64, error)
	GetMessageStatsFn  func(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*entities.MessageAuditStats, error)
}

func (m *MessageAuditQuerierMock) ListMessageLogs(ctx context.Context, filter dtos.MessageAuditFilterDTO) ([]entities.MessageAuditLog, int64, error) {
	if m.ListMessageLogsFn != nil {
		return m.ListMessageLogsFn(ctx, filter)
	}
	return []entities.MessageAuditLog{}, 0, nil
}

func (m *MessageAuditQuerierMock) GetMessageStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*entities.MessageAuditStats, error) {
	if m.GetMessageStatsFn != nil {
		return m.GetMessageStatsFn(ctx, businessID, dateFrom, dateTo)
	}
	return nil, nil
}

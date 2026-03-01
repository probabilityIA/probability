package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

// messageAuditQuerier consulta whatsapp_message_logs y whatsapp_conversations
// Replicado localmente para evitar compartir repositorios entre módulos
type messageAuditQuerier struct {
	db     db.IDatabase
	logger log.ILogger
}

// NewMessageAuditQuerier crea una instancia del querier de auditoría de mensajes
func NewMessageAuditQuerier(database db.IDatabase, logger log.ILogger) ports.IMessageAuditQuerier {
	return &messageAuditQuerier{
		db:     database,
		logger: logger.WithModule("message_audit_querier"),
	}
}

// messageLogRow representa una fila del resultado del JOIN
type messageLogRow struct {
	ID             uuid.UUID
	ConversationID uuid.UUID
	MessageID      string
	Direction      string
	TemplateName   string
	Content        string
	Status         string
	DeliveredAt    *time.Time
	ReadAt         *time.Time
	CreatedAt      time.Time
	PhoneNumber    string
	OrderNumber    string
	BusinessID     uint
}

// ListMessageLogs obtiene logs de mensajes con filtros y paginación
// Consulta: whatsapp_message_logs JOIN whatsapp_conversations
func (q *messageAuditQuerier) ListMessageLogs(ctx context.Context, filter dtos.MessageAuditFilterDTO) ([]entities.MessageAuditLog, int64, error) {
	baseQuery := q.db.Conn(ctx).
		Table("whatsapp_message_logs ml").
		Joins("INNER JOIN whatsapp_conversations c ON ml.conversation_id = c.id").
		Where("c.business_id = ?", filter.BusinessID)

	if filter.Status != nil && *filter.Status != "" {
		baseQuery = baseQuery.Where("ml.status = ?", *filter.Status)
	}
	if filter.Direction != nil && *filter.Direction != "" {
		baseQuery = baseQuery.Where("ml.direction = ?", *filter.Direction)
	}
	if filter.TemplateName != nil && *filter.TemplateName != "" {
		baseQuery = baseQuery.Where("ml.template_name ILIKE ?", fmt.Sprintf("%%%s%%", *filter.TemplateName))
	}
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		baseQuery = baseQuery.Where("ml.created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil && *filter.DateTo != "" {
		baseQuery = baseQuery.Where("ml.created_at < ?::date + interval '1 day'", *filter.DateTo)
	}

	// Count total
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting message logs")
		return nil, 0, err
	}

	if total == 0 {
		return []entities.MessageAuditLog{}, 0, nil
	}

	// Paginated query
	offset := (filter.Page - 1) * filter.PageSize
	var rows []messageLogRow

	err := baseQuery.
		Select(`ml.id, ml.conversation_id, ml.message_id, ml.direction,
			ml.template_name, ml.content, ml.status,
			ml.delivered_at, ml.read_at, ml.created_at,
			c.phone_number, c.order_number, c.business_id`).
		Order("ml.created_at DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&rows).Error

	if err != nil {
		q.logger.Error().Err(err).Msg("Error listing message logs")
		return nil, 0, err
	}

	// Map to domain entities
	logs := make([]entities.MessageAuditLog, len(rows))
	for i, row := range rows {
		logs[i] = entities.MessageAuditLog{
			ID:             row.ID.String(),
			ConversationID: row.ConversationID.String(),
			MessageID:      row.MessageID,
			Direction:      row.Direction,
			TemplateName:   row.TemplateName,
			Content:        row.Content,
			Status:         row.Status,
			DeliveredAt:    row.DeliveredAt,
			ReadAt:         row.ReadAt,
			CreatedAt:      row.CreatedAt,
			PhoneNumber:    row.PhoneNumber,
			OrderNumber:    row.OrderNumber,
			BusinessID:     row.BusinessID,
		}
	}

	return logs, total, nil
}

// statsResult representa el resultado de la query de estadísticas
type statsResult struct {
	TotalSent      int64
	TotalDelivered int64
	TotalRead      int64
	TotalFailed    int64
}

// GetMessageStats obtiene estadísticas agregadas de mensajes outbound
func (q *messageAuditQuerier) GetMessageStats(ctx context.Context, businessID uint, dateFrom, dateTo *string) (*entities.MessageAuditStats, error) {
	baseQuery := q.db.Conn(ctx).
		Table("whatsapp_message_logs ml").
		Joins("INNER JOIN whatsapp_conversations c ON ml.conversation_id = c.id").
		Where("c.business_id = ?", businessID).
		Where("ml.direction = ?", "outbound")

	if dateFrom != nil && *dateFrom != "" {
		baseQuery = baseQuery.Where("ml.created_at >= ?", *dateFrom)
	}
	if dateTo != nil && *dateTo != "" {
		baseQuery = baseQuery.Where("ml.created_at < ?::date + interval '1 day'", *dateTo)
	}

	var result statsResult
	err := baseQuery.Select(`
		COUNT(*) FILTER (WHERE ml.status = 'sent') AS total_sent,
		COUNT(*) FILTER (WHERE ml.status = 'delivered') AS total_delivered,
		COUNT(*) FILTER (WHERE ml.status = 'read') AS total_read,
		COUNT(*) FILTER (WHERE ml.status = 'failed') AS total_failed
	`).Scan(&result).Error

	if err != nil {
		q.logger.Error().Err(err).Msg("Error getting message stats")
		return nil, err
	}

	totalAll := result.TotalSent + result.TotalDelivered + result.TotalRead + result.TotalFailed
	var successRate float64
	if totalAll > 0 {
		successRate = float64(result.TotalSent+result.TotalDelivered+result.TotalRead) / float64(totalAll) * 100
	}

	return &entities.MessageAuditStats{
		TotalSent:      result.TotalSent,
		TotalDelivered: result.TotalDelivered,
		TotalRead:      result.TotalRead,
		TotalFailed:    result.TotalFailed,
		SuccessRate:    successRate,
	}, nil
}

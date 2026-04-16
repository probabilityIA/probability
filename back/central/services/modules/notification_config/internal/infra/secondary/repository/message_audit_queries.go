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


// emailLogRow representa una fila del resultado de la query de email_logs
type emailLogRow struct {
	BusinessID    uint
	IntegrationID uint
	ConfigID      uint
	To            string
	Subject       string
	EventType     string
	Status        string
	ErrorMessage  *string
	CreatedAt     time.Time
}

// ListEmailLogs obtiene logs de entregas de email con filtros y paginación
// Consulta: email_logs (gestionada por notification_config)
func (q *messageAuditQuerier) ListEmailLogs(ctx context.Context, businessID uint, status *string, dateFrom, dateTo *string, page, pageSize int) ([]entities.EmailDeliveryLog, int64, error) {
	baseQuery := q.db.Conn(ctx).
		Table("email_logs").
		Where("business_id = ?", businessID)

	if status != nil && *status != "" {
		baseQuery = baseQuery.Where("status = ?", *status)
	}
	if dateFrom != nil && *dateFrom != "" {
		baseQuery = baseQuery.Where("created_at >= ?", *dateFrom)
	}
	if dateTo != nil && *dateTo != "" {
		baseQuery = baseQuery.Where("created_at < ?::date + interval '1 day'", *dateTo)
	}

	// Count total
	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting email logs")
		return nil, 0, err
	}

	if total == 0 {
		return []entities.EmailDeliveryLog{}, 0, nil
	}

	// Paginated query
	offset := (page - 1) * pageSize
	var rows []emailLogRow

	err := baseQuery.
		Select("business_id, integration_id, config_id, \"to\", subject, event_type, status, error_message, created_at").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&rows).Error

	if err != nil {
		q.logger.Error().Err(err).Msg("Error listing email logs")
		return nil, 0, err
	}

	logs := make([]entities.EmailDeliveryLog, len(rows))
	for i, row := range rows {
		logs[i] = entities.EmailDeliveryLog{
			BusinessID:    row.BusinessID,
			IntegrationID: row.IntegrationID,
			ConfigID:      row.ConfigID,
			To:            row.To,
			Subject:       row.Subject,
			EventType:     row.EventType,
			Status:        row.Status,
			SentAt:        row.CreatedAt,
		}
		if row.ErrorMessage != nil {
			logs[i].ErrorMessage = *row.ErrorMessage
		}
	}

	return logs, total, nil
}


// conversationSummaryRow representa una fila del resultado de ListConversations
type conversationSummaryRow struct {
	ID                   uuid.UUID
	PhoneNumber          string
	OrderNumber          string
	BusinessID           uint
	CurrentState         string
	MessageCount         int
	LastMessageContent   string
	LastMessageDirection string
	LastMessageStatus    string
	LastActivity         time.Time
	CreatedAt            time.Time
}

// ListConversations obtiene conversaciones con resumen para la vista de lista
func (q *messageAuditQuerier) ListConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) ([]entities.ConversationSummary, int64, error) {
	// Base query: conversaciones con aggregates de mensajes
	baseQuery := q.db.Conn(ctx).
		Table("whatsapp_conversations c").
		Select(`c.id, c.phone_number, c.order_number, c.business_id, c.current_state, c.created_at,
			COALESCE(agg.message_count, 0) AS message_count,
			COALESCE(agg.last_activity, c.updated_at) AS last_activity,
			COALESCE(latest.content, '') AS last_message_content,
			COALESCE(latest.direction, '') AS last_message_direction,
			COALESCE(latest.status, '') AS last_message_status`).
		Joins(`LEFT JOIN (
			SELECT conversation_id, COUNT(*) AS message_count, MAX(created_at) AS last_activity
			FROM whatsapp_message_logs
			GROUP BY conversation_id
		) agg ON agg.conversation_id = c.id`).
		Joins(`LEFT JOIN LATERAL (
			SELECT content, direction, status
			FROM whatsapp_message_logs
			WHERE conversation_id = c.id
			ORDER BY created_at DESC
			LIMIT 1
		) latest ON true`).
		Where("c.business_id = ?", filter.BusinessID)

	// Filtros opcionales
	if filter.State != nil && *filter.State != "" {
		baseQuery = baseQuery.Where("c.current_state = ?", *filter.State)
	}
	if filter.Phone != nil && *filter.Phone != "" {
		baseQuery = baseQuery.Where("c.phone_number ILIKE ?", fmt.Sprintf("%%%s%%", *filter.Phone))
	}
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		baseQuery = baseQuery.Where("c.created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil && *filter.DateTo != "" {
		baseQuery = baseQuery.Where("c.created_at < ?::date + interval '1 day'", *filter.DateTo)
	}

	// Count total (necesitamos contar sin el select complejo)
	var total int64
	countQuery := q.db.Conn(ctx).
		Table("whatsapp_conversations c").
		Where("c.business_id = ?", filter.BusinessID)
	if filter.State != nil && *filter.State != "" {
		countQuery = countQuery.Where("c.current_state = ?", *filter.State)
	}
	if filter.Phone != nil && *filter.Phone != "" {
		countQuery = countQuery.Where("c.phone_number ILIKE ?", fmt.Sprintf("%%%s%%", *filter.Phone))
	}
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		countQuery = countQuery.Where("c.created_at >= ?", *filter.DateFrom)
	}
	if filter.DateTo != nil && *filter.DateTo != "" {
		countQuery = countQuery.Where("c.created_at < ?::date + interval '1 day'", *filter.DateTo)
	}
	if err := countQuery.Count(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting conversations")
		return nil, 0, err
	}

	if total == 0 {
		return []entities.ConversationSummary{}, 0, nil
	}

	// Paginated query
	offset := (filter.Page - 1) * filter.PageSize
	var rows []conversationSummaryRow

	err := baseQuery.
		Order("last_activity DESC").
		Offset(offset).
		Limit(filter.PageSize).
		Find(&rows).Error

	if err != nil {
		q.logger.Error().Err(err).Msg("Error listing conversations")
		return nil, 0, err
	}

	// Map to domain entities
	conversations := make([]entities.ConversationSummary, len(rows))
	for i, row := range rows {
		conversations[i] = entities.ConversationSummary{
			ID:                   row.ID.String(),
			PhoneNumber:          row.PhoneNumber,
			OrderNumber:          row.OrderNumber,
			BusinessID:           row.BusinessID,
			CurrentState:         row.CurrentState,
			MessageCount:         row.MessageCount,
			LastMessageContent:   row.LastMessageContent,
			LastMessageDirection: row.LastMessageDirection,
			LastMessageStatus:    row.LastMessageStatus,
			LastActivity:         row.LastActivity,
			CreatedAt:            row.CreatedAt,
		}
	}

	return conversations, total, nil
}

// conversationMessageRow representa una fila del resultado de GetConversationMessages
type conversationMessageRow struct {
	ID           uuid.UUID
	Direction    string
	MessageID    string
	TemplateName string
	Content      string
	Status       string
	DeliveredAt  *time.Time
	ReadAt       *time.Time
	CreatedAt    time.Time
}

// conversationMetaRow representa los metadatos de la conversación
type conversationMetaRow struct {
	ID           uuid.UUID
	PhoneNumber  string
	OrderNumber  string
	CurrentState string
}

// GetConversationMessages obtiene los mensajes de una conversación para la vista de chat
func (q *messageAuditQuerier) GetConversationMessages(ctx context.Context, conversationID string, businessID uint) (*entities.ConversationSummary, []entities.ConversationMessage, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("conversation ID inválido: %s", conversationID)
	}

	// Obtener metadatos de la conversación (y validar que pertenece al business)
	var meta conversationMetaRow
	err = q.db.Conn(ctx).
		Table("whatsapp_conversations").
		Select("id, phone_number, order_number, current_state").
		Where("id = ? AND business_id = ?", convID, businessID).
		First(&meta).Error
	if err != nil {
		return nil, nil, fmt.Errorf("conversation not found: %w", err)
	}

	// Obtener mensajes ordenados cronológicamente (ASC para chat view)
	var rows []conversationMessageRow
	err = q.db.Conn(ctx).
		Table("whatsapp_message_logs").
		Select("id, direction, message_id, template_name, content, status, delivered_at, read_at, created_at").
		Where("conversation_id = ?", convID).
		Order("created_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, nil, fmt.Errorf("error listing conversation messages: %w", err)
	}

	// Map conversación
	conv := &entities.ConversationSummary{
		ID:           meta.ID.String(),
		PhoneNumber:  meta.PhoneNumber,
		OrderNumber:  meta.OrderNumber,
		CurrentState: meta.CurrentState,
		BusinessID:   businessID,
	}

	// Map mensajes
	messages := make([]entities.ConversationMessage, len(rows))
	for i, row := range rows {
		messages[i] = entities.ConversationMessage{
			ID:           row.ID.String(),
			Direction:    row.Direction,
			MessageID:    row.MessageID,
			TemplateName: row.TemplateName,
			Content:      row.Content,
			Status:       row.Status,
			DeliveredAt:  row.DeliveredAt,
			ReadAt:       row.ReadAt,
			CreatedAt:    row.CreatedAt,
		}
	}

	return conv, messages, nil
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

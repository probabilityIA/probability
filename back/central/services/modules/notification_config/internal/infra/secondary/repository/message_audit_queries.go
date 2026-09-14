package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type messageAuditQuerier struct {
	db     db.IDatabase
	logger log.ILogger
}

func NewMessageAuditQuerier(database db.IDatabase, logger log.ILogger) ports.IMessageAuditQuerier {
	return &messageAuditQuerier{
		db:     database,
		logger: logger.WithModule("message_audit_querier"),
	}
}

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

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting message logs")
		return nil, 0, err
	}

	if total == 0 {
		return []entities.MessageAuditLog{}, 0, nil
	}

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

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting email logs")
		return nil, 0, err
	}

	if total == 0 {
		return []entities.EmailDeliveryLog{}, 0, nil
	}

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

type conversationSummaryRow struct {
	ID                   uuid.UUID
	PhoneNumber          string
	OrderNumber          string
	OrderID              string
	CampaignID           *uint
	CampaignName         string
	CustomerName         string
	UnreadCount          int
	OptedOut             bool
	ConversationType     string
	BusinessID           uint
	CurrentState         string
	MessageCount         int
	LastMessageContent   string
	LastMessageTemplate  string
	LastMessageMediaType string
	LastMessageDirection string
	LastMessageStatus    string
	LastActivity         time.Time
	CreatedAt            time.Time
}

const campaignConversationFilter = `(
	c.campaign_id = ?
	OR EXISTS (
		SELECT 1 FROM whatsapp_campaign_sends s
		WHERE s.campaign_id = ?
		  AND s.business_id = c.business_id
		  AND s.deleted_at IS NULL
		  AND right(regexp_replace(s.phone, '[^0-9]', '', 'g'), 10) = right(regexp_replace(c.phone_number, '[^0-9]', '', 'g'), 10)
	)
)`

const phoneKeyExpr = `regexp_replace(phone_number, '[^0-9]', '', 'g')`

func conversationOrderLateral(conv string) string {
	return fmt.Sprintf(`LEFT JOIN LATERAL (
		SELECT o.id::text AS id FROM orders o
		WHERE %[1]s.order_number <> ''
		  AND o.order_number = %[1]s.order_number
		  AND o.business_id = %[1]s.business_id
		  AND o.deleted_at IS NULL
		ORDER BY o.created_at DESC LIMIT 1
	) ord ON true`, conv)
}

const optOutReplyText = "dejar de recibir"

func clientNameLateral(business, phoneKey string) string {
	return fmt.Sprintf(`LEFT JOIN LATERAL (
		SELECT cl.name FROM client cl
		WHERE cl.business_id = %[1]s
		  AND cl.deleted_at IS NULL
		  AND btrim(COALESCE(cl.name, '')) <> ''
		  AND right(regexp_replace(cl.phone, '[^0-9]', '', 'g'), 10) = right(%[2]s, 10)
		ORDER BY cl.updated_at DESC LIMIT 1
	) cli ON true`, business, phoneKey)
}

func clientOptOutExists(business, phoneKey string) string {
	return fmt.Sprintf(`EXISTS (
		SELECT 1 FROM client cl
		WHERE cl.business_id = %[1]s
		  AND cl.accepts_marketing = false
		  AND cl.deleted_at IS NULL
		  AND right(regexp_replace(cl.phone, '[^0-9]', '', 'g'), 10) = right(%[2]s, 10)
	)`, business, phoneKey)
}

func conversationCampaignLateral(conv, phoneKey string) string {
	return fmt.Sprintf(`LEFT JOIN LATERAL (
		SELECT wc.id, wc.name FROM whatsapp_campaigns wc
		WHERE wc.business_id = %[1]s.business_id
		  AND wc.deleted_at IS NULL
		  AND (
			wc.id IN (
				SELECT cc.campaign_id FROM whatsapp_conversations cc
				WHERE cc.business_id = %[1]s.business_id
				  AND cc.campaign_id IS NOT NULL
				  AND regexp_replace(cc.phone_number, '[^0-9]', '', 'g') = %[2]s
			)
			OR EXISTS (
				SELECT 1 FROM whatsapp_campaign_sends s
				WHERE s.campaign_id = wc.id
				  AND s.deleted_at IS NULL
				  AND right(regexp_replace(s.phone, '[^0-9]', '', 'g'), 10) = right(%[2]s, 10)
			)
		  )
		ORDER BY wc.created_at DESC LIMIT 1
	) camp ON true`, conv, phoneKey)
}

func conversationFilters(filter dtos.ConversationListFilterDTO) (string, []any) {
	clauses := []string{"c.business_id = ?"}
	args := []any{filter.BusinessID}

	if filter.State != nil && *filter.State != "" {
		clauses = append(clauses, "c.current_state = ?")
		args = append(args, *filter.State)
	}
	if filter.Phone != nil && *filter.Phone != "" {
		clauses = append(clauses, "c.phone_number ILIKE ?")
		args = append(args, fmt.Sprintf("%%%s%%", *filter.Phone))
	}
	if filter.CampaignID != nil && *filter.CampaignID > 0 {
		clauses = append(clauses, campaignConversationFilter)
		args = append(args, *filter.CampaignID, *filter.CampaignID)
	}
	if filter.DateFrom != nil && *filter.DateFrom != "" {
		clauses = append(clauses, "c.created_at >= ?")
		args = append(args, *filter.DateFrom)
	}
	if filter.DateTo != nil && *filter.DateTo != "" {
		clauses = append(clauses, "c.created_at < ?::date + interval '1 day'")
		args = append(args, *filter.DateTo)
	}

	return strings.Join(clauses, " AND "), args
}

func (q *messageAuditQuerier) ListConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) ([]entities.ConversationSummary, int64, error) {
	where, args := conversationFilters(filter)

	var total int64
	countSQL := fmt.Sprintf(`
		SELECT COUNT(*) FROM (
			SELECT DISTINCT regexp_replace(c.phone_number, '[^0-9]', '', 'g') AS phone_key
			FROM whatsapp_conversations c
			WHERE %s
		) t`, where)

	if err := q.db.Conn(ctx).Raw(countSQL, args...).Scan(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting conversations")
		return nil, 0, err
	}

	if total == 0 {
		return []entities.ConversationSummary{}, 0, nil
	}

	offset := (filter.Page - 1) * filter.PageSize
	var rows []conversationSummaryRow

	listSQL := fmt.Sprintf(`
		WITH conv AS (
			SELECT c.id, c.phone_number, c.order_number, c.conversation_type,
			       c.business_id, c.current_state, c.created_at, c.updated_at,
			       %s AS phone_key
			FROM whatsapp_conversations c
			WHERE %s
		),
		ultima AS (
			SELECT DISTINCT ON (phone_key) * FROM conv
			ORDER BY phone_key, created_at DESC
		),
		primera AS (
			SELECT phone_key, MIN(created_at) AS created_at FROM conv GROUP BY phone_key
		),
		msgs AS (
			SELECT conv.phone_key, conv.business_id, ml.created_at, ml.content, ml.template_name,
			       ml.media_type, ml.direction, ml.status
			FROM conv
			JOIN whatsapp_message_logs ml ON ml.conversation_id = conv.id
		),
		agg AS (
			SELECT m.phone_key, COUNT(*) AS message_count, MAX(m.created_at) AS last_activity,
			       COUNT(*) FILTER (
			           WHERE m.direction = 'inbound'
			             AND m.created_at > COALESCE(rd.last_read_at, 'epoch'::timestamptz)
			       ) AS unread_count,
			       COUNT(*) FILTER (
			           WHERE m.direction = 'inbound'
			             AND lower(btrim(m.content)) = '%s'
			       ) AS opt_out_replies
			FROM msgs m
			LEFT JOIN whatsapp_conversation_reads rd
				ON rd.business_id = m.business_id AND rd.phone_key = m.phone_key
			GROUP BY m.phone_key
		),
		reciente AS (
			SELECT DISTINCT ON (phone_key) phone_key, content, template_name, media_type, direction, status
			FROM msgs
			ORDER BY phone_key, created_at DESC
		),
		pagina AS (
			SELECT u.id, u.phone_number, u.order_number, u.conversation_type, u.business_id,
			       u.current_state, u.phone_key, p.created_at,
			       COALESCE(a.message_count, 0) AS message_count,
			       COALESCE(a.last_activity, u.updated_at) AS last_activity,
			       COALESCE(a.unread_count, 0) AS unread_count,
			       COALESCE(a.opt_out_replies, 0) AS opt_out_replies,
			       COALESCE(r.content, '') AS last_message_content,
			       COALESCE(r.template_name, '') AS last_message_template,
			       COALESCE(r.media_type, '') AS last_message_media_type,
			       COALESCE(r.direction, '') AS last_message_direction,
			       COALESCE(r.status, '') AS last_message_status
			FROM ultima u
			JOIN primera p ON p.phone_key = u.phone_key
			LEFT JOIN agg a ON a.phone_key = u.phone_key
			LEFT JOIN reciente r ON r.phone_key = u.phone_key
			ORDER BY (COALESCE(a.unread_count, 0) > 0) DESC, COALESCE(a.last_activity, u.updated_at) DESC
			OFFSET ? LIMIT ?
		)
		SELECT pg.id, pg.phone_number, pg.order_number, pg.conversation_type, pg.business_id,
		       pg.current_state, pg.created_at, pg.message_count, pg.last_activity,
		       pg.last_message_content, pg.last_message_template, pg.last_message_media_type,
		       pg.last_message_direction, pg.last_message_status,
		       COALESCE(ord.id, '') AS order_id,
		       camp.id AS campaign_id,
		       COALESCE(camp.name, '') AS campaign_name,
		       COALESCE(cli.name, '') AS customer_name,
		       pg.unread_count,
		       (pg.opt_out_replies > 0 OR %s) AS opted_out
		FROM pagina pg
		%s
		%s
		%s
		ORDER BY (pg.unread_count > 0) DESC, pg.last_activity DESC`, phoneKeyExpr, where, optOutReplyText, clientOptOutExists("pg.business_id", "pg.phone_key"), conversationOrderLateral("pg"), conversationCampaignLateral("pg", "pg.phone_key"), clientNameLateral("pg.business_id", "pg.phone_key"))

	listArgs := append(append([]any{}, args...), offset, filter.PageSize)

	if err := q.db.Conn(ctx).Raw(listSQL, listArgs...).Scan(&rows).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error listing conversations")
		return nil, 0, err
	}

	previewPairs := make([][2]string, 0, len(rows))
	for _, row := range rows {
		previewPairs = append(previewPairs, [2]string{row.LastMessageTemplate, row.LastMessageContent})
	}
	previewTexts := q.loadTemplateTexts(ctx, filter.BusinessID, templateNamesToResolve(previewPairs))

	conversations := make([]entities.ConversationSummary, len(rows))
	for i, row := range rows {
		conversations[i] = entities.ConversationSummary{
			ID:                   row.ID.String(),
			PhoneNumber:          row.PhoneNumber,
			OrderNumber:          row.OrderNumber,
			OrderID:              row.OrderID,
			CampaignID:           row.CampaignID,
			CampaignName:         row.CampaignName,
			CustomerName:         row.CustomerName,
			UnreadCount:          row.UnreadCount,
			OptedOut:             row.OptedOut,
			ConversationType:     row.ConversationType,
			BusinessID:           row.BusinessID,
			CurrentState:         row.CurrentState,
			MessageCount:         row.MessageCount,
			LastMessageContent:   previewContent(resolveMessageContent(previewTexts, row.LastMessageTemplate, row.LastMessageContent), row.LastMessageMediaType),
			LastMessageDirection: row.LastMessageDirection,
			LastMessageStatus:    row.LastMessageStatus,
			LastActivity:         row.LastActivity,
			CreatedAt:            row.CreatedAt,
		}
	}

	return conversations, total, nil
}

func (q *messageAuditQuerier) CountUnreadConversations(ctx context.Context, filter dtos.ConversationListFilterDTO) (int64, error) {
	where, args := conversationFilters(filter)

	countSQL := fmt.Sprintf(`
		WITH conv AS (
			SELECT c.id, %s AS phone_key
			FROM whatsapp_conversations c
			WHERE %s
		),
		entrantes AS (
			SELECT conv.phone_key, ml.created_at
			FROM conv
			JOIN whatsapp_message_logs ml ON ml.conversation_id = conv.id
			WHERE ml.direction = 'inbound'
		)
		SELECT COUNT(DISTINCT e.phone_key)
		FROM entrantes e
		LEFT JOIN whatsapp_conversation_reads rd
			ON rd.business_id = ? AND rd.phone_key = e.phone_key
		WHERE e.created_at > COALESCE(rd.last_read_at, 'epoch'::timestamptz)`, phoneKeyExpr, where)

	countArgs := append(append([]any{}, args...), filter.BusinessID)

	var total int64
	if err := q.db.Conn(ctx).Raw(countSQL, countArgs...).Scan(&total).Error; err != nil {
		q.logger.Error().Err(err).Msg("Error counting unread conversations")
		return 0, err
	}

	return total, nil
}

func (q *messageAuditQuerier) MarkConversationRead(ctx context.Context, conversationID string, businessID uint, userID *uint) error {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return fmt.Errorf("conversation ID invalido: %s", conversationID)
	}

	result := q.db.Conn(ctx).Exec(`
		INSERT INTO whatsapp_conversation_reads (business_id, phone_key, last_read_at, read_by_id, updated_at)
		SELECT c.business_id, regexp_replace(c.phone_number, '[^0-9]', '', 'g'), NOW(), ?, NOW()
		FROM whatsapp_conversations c
		WHERE c.id = ? AND c.business_id = ?
		ON CONFLICT (business_id, phone_key) DO UPDATE SET
			last_read_at = EXCLUDED.last_read_at,
			read_by_id = EXCLUDED.read_by_id,
			updated_at = EXCLUDED.updated_at`, userID, convID, businessID)
	if result.Error != nil {
		q.logger.Error().Err(result.Error).Str("conversation_id", conversationID).Msg("Error marking conversation as read")
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("conversation not found: %s", conversationID)
	}

	return nil
}

type conversationMessageRow struct {
	ID            uuid.UUID
	Direction     string
	MessageID     string
	TemplateName  string
	Content       string
	Status        string
	DeliveredAt   *time.Time
	ReadAt        *time.Time
	CreatedAt     time.Time
	MediaType     string
	MediaKey      string
	MediaMime     string
	MediaFilename string
	MediaSize     int64
}

type conversationMetaRow struct {
	ID               uuid.UUID
	PhoneNumber      string
	OrderNumber      string
	OrderID          string
	CampaignID       *uint
	CampaignName     string
	CustomerName     string
	OptedOut         bool
	ConversationType string
	CurrentState     string
}

func (q *messageAuditQuerier) GetConversationMessages(ctx context.Context, conversationID string, businessID uint) (*entities.ConversationSummary, []entities.ConversationMessage, error) {
	convID, err := uuid.Parse(conversationID)
	if err != nil {
		return nil, nil, fmt.Errorf("conversation ID invalido: %s", conversationID)
	}

	var meta conversationMetaRow
	metaSQL := fmt.Sprintf(`
		SELECT c.id, c.phone_number, c.order_number, c.conversation_type, c.current_state,
		       COALESCE(ord.id, '') AS order_id,
		       camp.id AS campaign_id,
		       COALESCE(camp.name, '') AS campaign_name,
		       COALESCE(cli.name, '') AS customer_name,
		       (%s OR EXISTS (
		           SELECT 1 FROM whatsapp_message_logs ml
		           JOIN whatsapp_conversations c2 ON c2.id = ml.conversation_id
		           WHERE c2.business_id = c.business_id
		             AND regexp_replace(c2.phone_number, '[^0-9]', '', 'g') = regexp_replace(c.phone_number, '[^0-9]', '', 'g')
		             AND ml.direction = 'inbound'
		             AND lower(btrim(ml.content)) = '%s'
		       )) AS opted_out
		FROM whatsapp_conversations c
		%s
		%s
		%s
		WHERE c.id = ? AND c.business_id = ?
		LIMIT 1`, clientOptOutExists("c.business_id", "regexp_replace(c.phone_number, '[^0-9]', '', 'g')"), optOutReplyText, conversationOrderLateral("c"), conversationCampaignLateral("c", "regexp_replace(c.phone_number, '[^0-9]', '', 'g')"), clientNameLateral("c.business_id", "regexp_replace(c.phone_number, '[^0-9]', '', 'g')"))
	err = q.db.Conn(ctx).Raw(metaSQL, convID, businessID).Scan(&meta).Error
	if err != nil {
		return nil, nil, fmt.Errorf("conversation not found: %w", err)
	}
	if meta.ID == uuid.Nil {
		return nil, nil, fmt.Errorf("conversation not found: %s", conversationID)
	}

	var rows []conversationMessageRow
	err = q.db.Conn(ctx).
		Table("whatsapp_message_logs").
		Select("id, direction, message_id, template_name, content, status, delivered_at, read_at, created_at, COALESCE(media_type, '') AS media_type, COALESCE(media_key, '') AS media_key, COALESCE(media_mime, '') AS media_mime, COALESCE(media_filename, '') AS media_filename, COALESCE(media_size, 0) AS media_size").
		Where(fmt.Sprintf(`conversation_id IN (
			SELECT id FROM whatsapp_conversations
			WHERE business_id = ? AND %s = (
				SELECT %s FROM whatsapp_conversations WHERE id = ?
			)
		)`, phoneKeyExpr, phoneKeyExpr), businessID, convID).
		Order("created_at ASC").
		Find(&rows).Error
	if err != nil {
		return nil, nil, fmt.Errorf("error listing conversation messages: %w", err)
	}

	conv := &entities.ConversationSummary{
		ID:               meta.ID.String(),
		PhoneNumber:      meta.PhoneNumber,
		OrderNumber:      meta.OrderNumber,
		OrderID:          meta.OrderID,
		CampaignID:       meta.CampaignID,
		CampaignName:     meta.CampaignName,
		CustomerName:     meta.CustomerName,
		OptedOut:         meta.OptedOut,
		ConversationType: meta.ConversationType,
		CurrentState:     meta.CurrentState,
		BusinessID:       businessID,
	}

	messagePairs := make([][2]string, 0, len(rows))
	for _, row := range rows {
		messagePairs = append(messagePairs, [2]string{row.TemplateName, row.Content})
	}
	texts := q.loadTemplateTexts(ctx, businessID, templateNamesToResolve(messagePairs))

	messages := make([]entities.ConversationMessage, len(rows))
	for i, row := range rows {
		content, buttons := resolveMessage(texts, row.TemplateName, row.Content)
		if row.Direction != "outbound" {
			buttons = nil
		}
		messages[i] = entities.ConversationMessage{
			ID:           row.ID.String(),
			Direction:    row.Direction,
			MessageID:    row.MessageID,
			TemplateName: row.TemplateName,
			Content:      content,
			Buttons:      messageButtons(buttons),
			Media:        messageMedia(row),
			Status:       row.Status,
			DeliveredAt:  row.DeliveredAt,
			ReadAt:       row.ReadAt,
			CreatedAt:    row.CreatedAt,
		}
	}

	return conv, messages, nil
}

type statsResult struct {
	TotalSent      int64
	TotalDelivered int64
	TotalRead      int64
	TotalFailed    int64
}

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

func messageButtons(buttons []templateButton) []entities.MessageButton {
	if len(buttons) == 0 {
		return nil
	}
	out := make([]entities.MessageButton, 0, len(buttons))
	for _, button := range buttons {
		out = append(out, entities.MessageButton{Text: button.Text, Type: button.Type})
	}
	return out
}

func messageMedia(row conversationMessageRow) *entities.MessageMedia {
	if row.MediaType == "" {
		return nil
	}
	return &entities.MessageMedia{
		Type:     row.MediaType,
		Key:      row.MediaKey,
		Mime:     row.MediaMime,
		Filename: row.MediaFilename,
		Size:     row.MediaSize,
	}
}

func previewContent(content, mediaType string) string {
	if strings.TrimSpace(content) != "" || mediaType == "" {
		return content
	}
	switch mediaType {
	case "image":
		return "[Imagen]"
	case "document":
		return "[Documento]"
	case "audio":
		return "[Audio]"
	case "video":
		return "[Video]"
	case "sticker":
		return "[Sticker]"
	default:
		return "[Archivo]"
	}
}

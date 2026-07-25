package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

const whatsAppNotificationTypeID = 2

func (r *Repository) GetOrderNotificationCounters(
	ctx context.Context,
	businessID uint,
	keys []entities.OrderNotificationKey,
) (map[string][]entities.OrderNotificationCounter, error) {
	result := make(map[string][]entities.OrderNotificationCounter)
	if businessID == 0 || len(keys) == 0 {
		return result, nil
	}

	eventsByIntegration, err := r.whatsAppEventCodesByIntegration(ctx, businessID)
	if err != nil {
		return nil, err
	}
	if len(eventsByIntegration) == 0 {
		return result, nil
	}

	orderNumbers := make([]string, 0, len(keys))
	seen := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		if k.OrderNumber == "" {
			continue
		}
		if _, ok := seen[k.OrderNumber]; ok {
			continue
		}
		seen[k.OrderNumber] = struct{}{}
		orderNumbers = append(orderNumbers, k.OrderNumber)
	}

	statusByOrder, err := r.whatsAppTemplateStatusByOrder(ctx, businessID, orderNumbers)
	if err != nil {
		return nil, err
	}

	for _, k := range keys {
		eventCodes, ok := eventsByIntegration[k.IntegrationID]
		if !ok || len(eventCodes) == 0 {
			continue
		}

		templates := statusByOrder[k.OrderNumber]
		sent := 0
		failed := 0
		for _, code := range eventCodes {
			switch entities.EventStatusFromTemplates(code, templates) {
			case entities.NotificationEventStatusSent:
				sent++
			case entities.NotificationEventStatusFailed:
				failed++
			}
		}

		result[k.OrderID] = []entities.OrderNotificationCounter{
			entities.NewOrderNotificationCounter(
				entities.NotificationChannelWhatsApp,
				sent,
				len(eventCodes),
				failed,
			),
		}
	}

	return result, nil
}

func (r *Repository) whatsAppEventCodesByIntegration(ctx context.Context, businessID uint) (map[uint][]string, error) {
	type row struct {
		IntegrationID uint
		EventCode     string
	}

	var rows []row
	err := r.db.Conn(ctx).
		Table("business_notification_configs AS bnc").
		Select("DISTINCT bnc.integration_id AS integration_id, net.event_code AS event_code").
		Joins("JOIN notification_event_types net ON net.id = bnc.notification_event_type_id").
		Where("bnc.business_id = ?", businessID).
		Where("bnc.notification_type_id = ?", whatsAppNotificationTypeID).
		Where("bnc.enabled = true").
		Where("bnc.deleted_at IS NULL").
		Where("net.deleted_at IS NULL").
		Where("net.is_active = true").
		Where("bnc.integration_id IS NOT NULL").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make(map[uint][]string)
	for _, row := range rows {
		out[row.IntegrationID] = append(out[row.IntegrationID], row.EventCode)
	}
	return out, nil
}

func (r *Repository) whatsAppTemplateStatusByOrder(
	ctx context.Context,
	businessID uint,
	orderNumbers []string,
) (map[string]map[string]string, error) {
	out := make(map[string]map[string]string)
	if len(orderNumbers) == 0 {
		return out, nil
	}

	type row struct {
		OrderNumber  string
		TemplateName string
		Delivered    bool
	}

	var rows []row
	err := r.db.Conn(ctx).
		Table("whatsapp_message_logs AS ml").
		Select(`wc.order_number AS order_number, ml.template_name AS template_name,
			BOOL_OR(ml.status IN ('sent','delivered','read')) AS delivered`).
		Joins("JOIN whatsapp_conversations wc ON wc.id = ml.conversation_id").
		Where("wc.business_id = ?", businessID).
		Where("wc.order_number IN ?", orderNumbers).
		Where("ml.direction = ?", "outbound").
		Where("ml.template_name <> ''").
		Group("wc.order_number, ml.template_name").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		if out[row.OrderNumber] == nil {
			out[row.OrderNumber] = make(map[string]string)
		}
		status := entities.NotificationEventStatusFailed
		if row.Delivered {
			status = entities.NotificationEventStatusSent
		}
		out[row.OrderNumber][row.TemplateName] = status
	}
	return out, nil
}

func (r *Repository) GetOrderNotificationMessages(
	ctx context.Context,
	businessID uint,
	orderNumber string,
) ([]entities.OrderNotificationMessage, error) {
	if businessID == 0 || orderNumber == "" {
		return nil, nil
	}

	type row struct {
		ID           string
		TemplateName string
		Direction    string
		Status       string
		Content      string
		MessageID    string
		CreatedAt    time.Time
	}

	var rows []row
	err := r.db.Conn(ctx).
		Table("whatsapp_message_logs AS ml").
		Select(`ml.id, ml.template_name, ml.direction, ml.status,
			ml.content, ml.message_id, ml.created_at`).
		Joins("JOIN whatsapp_conversations wc ON wc.id = ml.conversation_id").
		Where("wc.business_id = ?", businessID).
		Where("wc.order_number = ?", orderNumber).
		Order("ml.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]entities.OrderNotificationMessage, 0, len(rows))
	for _, m := range rows {
		out = append(out, entities.OrderNotificationMessage{
			ID:           m.ID,
			Channel:      entities.NotificationChannelWhatsApp,
			TemplateName: m.TemplateName,
			Direction:    m.Direction,
			Status:       m.Status,
			Content:      m.Content,
			MessageID:    m.MessageID,
			CreatedAt:    m.CreatedAt,
		})
	}
	return out, nil
}

func (r *Repository) GetWhatsAppEventCodes(ctx context.Context, businessID, integrationID uint) ([]string, error) {
	if businessID == 0 || integrationID == 0 {
		return nil, nil
	}

	byIntegration, err := r.whatsAppEventCodesByIntegration(ctx, businessID)
	if err != nil {
		return nil, err
	}
	return byIntegration[integrationID], nil
}

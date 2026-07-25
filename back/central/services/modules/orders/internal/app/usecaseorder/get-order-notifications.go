package usecaseorder

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
)

const (
	eventOrderCreated   = "order.created"
	eventGuideGenerated = "shipment.guide_generated"
)

var resendActionByEvent = map[string]string{
	eventOrderCreated:   "request-confirmation",
	eventGuideGenerated: "send-guide-notification",
}

func (uc *UseCaseOrder) GetOrderNotifications(ctx context.Context, orderID string, businessID uint) (*dtos.OrderNotificationsResponse, error) {
	order, err := uc.repo.GetOrderByID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error getting order: %w", err)
	}
	if order == nil {
		return nil, fmt.Errorf("order not found")
	}

	if businessID == 0 && order.BusinessID != nil {
		businessID = *order.BusinessID
	}
	if order.BusinessID == nil || *order.BusinessID != businessID {
		return nil, fmt.Errorf("la orden no pertenece al negocio")
	}

	messages, err := uc.repo.GetOrderNotificationMessages(ctx, businessID, order.OrderNumber)
	if err != nil {
		return nil, fmt.Errorf("error getting notification messages: %w", err)
	}

	eventCodes, err := uc.repo.GetWhatsAppEventCodes(ctx, businessID, order.IntegrationID)
	if err != nil {
		return nil, fmt.Errorf("error getting notification config: %w", err)
	}

	statusByTemplate := make(map[string]string, len(messages))
	for _, m := range messages {
		if m.Direction != "outbound" || m.TemplateName == "" {
			continue
		}
		if isDeliveredStatus(m.Status) {
			statusByTemplate[m.TemplateName] = entities.NotificationEventStatusSent
			continue
		}
		if _, ok := statusByTemplate[m.TemplateName]; !ok {
			statusByTemplate[m.TemplateName] = entities.NotificationEventStatusFailed
		}
	}

	events := make([]dtos.OrderNotificationEvent, 0, len(eventCodes))
	sent := 0
	failed := 0
	for _, code := range eventCodes {
		status := entities.EventStatusFromTemplates(code, statusByTemplate)

		switch status {
		case entities.NotificationEventStatusSent:
			sent++
		case entities.NotificationEventStatusFailed:
			failed++
		}

		events = append(events, dtos.OrderNotificationEvent{
			EventCode:    code,
			Status:       status,
			ResendAction: resendActionByEvent[code],
			CanResend:    resendActionByEvent[code] != "" && status != entities.NotificationEventStatusSent,
		})
	}

	counter := entities.NewOrderNotificationCounter(
		entities.NotificationChannelWhatsApp,
		sent,
		len(eventCodes),
		failed,
	)

	return &dtos.OrderNotificationsResponse{
		OrderID:     order.ID,
		OrderNumber: order.OrderNumber,
		Channel:     entities.NotificationChannelWhatsApp,
		Sent:        counter.Sent,
		Expected:    counter.Expected,
		Failed:      counter.Failed,
		State:       counter.State,
		Events:      events,
		Messages:    toNotificationMessageDTOs(messages),
	}, nil
}

func isDeliveredStatus(status string) bool {
	return status == "sent" || status == "delivered" || status == "read"
}

func toNotificationMessageDTOs(messages []entities.OrderNotificationMessage) []dtos.OrderNotificationMessage {
	out := make([]dtos.OrderNotificationMessage, 0, len(messages))
	for _, m := range messages {
		out = append(out, dtos.OrderNotificationMessage{
			ID:           m.ID,
			Channel:      m.Channel,
			TemplateName: m.TemplateName,
			Direction:    m.Direction,
			Status:       m.Status,
			Content:      m.Content,
			CreatedAt:    m.CreatedAt,
		})
	}
	return out
}

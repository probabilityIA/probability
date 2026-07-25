package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/infra/primary/handlers/response"
)

func OrderNotificationsToResponse(dto *dtos.OrderNotificationsResponse) *response.OrderNotifications {
	if dto == nil {
		return nil
	}

	events := make([]response.OrderNotificationEvent, 0, len(dto.Events))
	for _, e := range dto.Events {
		events = append(events, response.OrderNotificationEvent{
			EventCode:    e.EventCode,
			Status:       e.Status,
			ResendAction: e.ResendAction,
			CanResend:    e.CanResend,
		})
	}

	messages := make([]response.OrderNotificationMessage, 0, len(dto.Messages))
	for _, m := range dto.Messages {
		messages = append(messages, response.OrderNotificationMessage{
			ID:           m.ID,
			Channel:      m.Channel,
			TemplateName: m.TemplateName,
			Direction:    m.Direction,
			Status:       m.Status,
			Content:      m.Content,
			CreatedAt:    m.CreatedAt,
		})
	}

	return &response.OrderNotifications{
		OrderID:     dto.OrderID,
		OrderNumber: dto.OrderNumber,
		Channel:     dto.Channel,
		Sent:        dto.Sent,
		Expected:    dto.Expected,
		Failed:      dto.Failed,
		State:       dto.State,
		Events:      events,
		Messages:    messages,
	}
}

package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

// WhatsAppMock implementa ports.IWhatsApp para tests unitarios
type WhatsAppMock struct {
	SendMessageFn func(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error)
}

func (m *WhatsAppMock) SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error) {
	if m.SendMessageFn != nil {
		return m.SendMessageFn(ctx, phoneNumberID, msg, accessToken)
	}
	return "wamid.mock123", nil
}

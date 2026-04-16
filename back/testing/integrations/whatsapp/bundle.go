package whatsapp

import (
	"github.com/secamc93/probability/back/testing/integrations/whatsapp/internal/app/usecases"
	"github.com/secamc93/probability/back/testing/integrations/whatsapp/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// New inicializa el módulo de WhatsApp para pruebas de integración
func New(config env.IConfig, logger log.ILogger) *WhatsAppIntegration {
	webhookClient := usecases.NewWebhookClient(config, logger)
	conversationSimulator := usecases.NewConversationSimulator(webhookClient, config, logger)

	return &WhatsAppIntegration{
		conversationSimulator: conversationSimulator,
		logger:                logger,
	}
}

// WhatsAppIntegration representa el módulo de integración de WhatsApp
type WhatsAppIntegration struct {
	conversationSimulator *usecases.ConversationSimulator
	logger                log.ILogger
}

// SimulateUserResponse simula una respuesta manual de usuario
func (w *WhatsAppIntegration) SimulateUserResponse(phoneNumber string, response string) error {
	return w.conversationSimulator.SimulateUserResponse(phoneNumber, response)
}

// SimulateAutoResponse simula una respuesta automática basada en el template
func (w *WhatsAppIntegration) SimulateAutoResponse(phoneNumber, templateName string) error {
	return w.conversationSimulator.SimulateAutoResponse(phoneNumber, templateName)
}

// GetAllConversations retorna todas las conversaciones almacenadas
func (w *WhatsAppIntegration) GetAllConversations() []*domain.Conversation {
	return w.conversationSimulator.GetAllConversations()
}

// GetConversation obtiene una conversación por número de teléfono
func (w *WhatsAppIntegration) GetConversation(phoneNumber string) (*domain.Conversation, bool) {
	return w.conversationSimulator.GetConversation(phoneNumber)
}

// GetMessages obtiene todos los mensajes de una conversación
func (w *WhatsAppIntegration) GetMessages(conversationID string) []*domain.MessageLog {
	return w.conversationSimulator.GetMessages(conversationID)
}

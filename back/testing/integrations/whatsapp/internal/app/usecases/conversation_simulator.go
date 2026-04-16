package usecases

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/testing/integrations/whatsapp/internal/domain"
	"github.com/secamc93/probability/back/testing/shared/env"
	"github.com/secamc93/probability/back/testing/shared/log"
)

// ConversationSimulator simula conversaciones de WhatsApp
type ConversationSimulator struct {
	webhookClient         domain.IWebhookClient
	config                env.IConfig
	logger                log.ILogger
	conversationRepository *domain.ConversationRepository
	autoReplyDelay        time.Duration
}

// SimulateUserResponse simula una respuesta de usuario para un template específico
func (s *ConversationSimulator) SimulateUserResponse(phoneNumber string, response string) error {
	s.logger.Info().
		Str("phone_number", phoneNumber).
		Str("response", response).
		Msg("Simulando respuesta de usuario")

	// Verificar que existe una conversación
	conv, exists := s.conversationRepository.GetConversation(phoneNumber)
	if !exists {
		return fmt.Errorf("no existe conversación para el número %s", phoneNumber)
	}

	// Esperar delay configurado
	time.Sleep(s.autoReplyDelay)

	// Enviar webhooks de estado primero
	messageID := fmt.Sprintf("wamid.HBg%s", uuid.New().String()[:20])
	s.sendStatusWebhook(messageID, phoneNumber, "delivered")
	time.Sleep(500 * time.Millisecond)
	s.sendStatusWebhook(messageID, phoneNumber, "read")
	time.Sleep(500 * time.Millisecond)

	// Enviar webhook de respuesta de usuario
	s.sendMessageWebhook(phoneNumber, response)

	// Guardar mensaje en el log
	msg := &domain.MessageLog{
		ID:             uuid.New().String(),
		ConversationID: conv.ID,
		Direction:      "INBOUND",
		MessageType:    "button",
		Content:        response,
		Status:         "RECEIVED",
		CreatedAt:      time.Now(),
	}
	s.conversationRepository.SaveMessage(msg)

	return nil
}

// SimulateAutoResponse simula una respuesta automática basada en el template
func (s *ConversationSimulator) SimulateAutoResponse(phoneNumber, templateName string) error {
	response := s.getUserResponseForTemplate(templateName)
	if response == "" {
		s.logger.Info().
			Str("template", templateName).
			Msg("Template no espera respuesta automática")
		return nil
	}

	return s.SimulateUserResponse(phoneNumber, response)
}

// GetAllConversations retorna todas las conversaciones almacenadas
func (s *ConversationSimulator) GetAllConversations() []*domain.Conversation {
	return s.conversationRepository.GetAllConversations()
}

// GetConversation obtiene una conversación por número de teléfono
func (s *ConversationSimulator) GetConversation(phoneNumber string) (*domain.Conversation, bool) {
	return s.conversationRepository.GetConversation(phoneNumber)
}

// GetMessages obtiene todos los mensajes de una conversación
func (s *ConversationSimulator) GetMessages(conversationID string) []*domain.MessageLog {
	return s.conversationRepository.GetMessages(conversationID)
}

// sendStatusWebhook envía un webhook de estado de mensaje
func (s *ConversationSimulator) sendStatusWebhook(messageID, phoneNumber, status string) {
	payload := domain.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []domain.WebhookEntry{{
			ID: "123456789",
			Changes: []domain.WebhookChange{{
				Value: domain.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata: domain.WebhookMetadata{
						DisplayPhoneNumber: "+15551234567",
						PhoneNumberID:      "123456789012345",
					},
					Statuses: []domain.WebhookStatus{{
						ID:          messageID,
						Status:      status,
						Timestamp:   fmt.Sprintf("%d", time.Now().Unix()),
						RecipientID: phoneNumber,
						Conversation: &domain.ConversationInfo{
							ID: "conversation_" + messageID,
							Origin: domain.ConversationOrigin{Type: "business_initiated"},
						},
					}},
				},
				Field: "messages",
			}},
		}},
	}

	err := s.webhookClient.SendWebhook(payload)
	if err != nil {
		s.logger.Error().Err(err).Str("status", status).Msg("Error al enviar webhook de estado")
	} else {
		s.logger.Info().Str("status", status).Str("message_id", messageID).Msg("Webhook de estado enviado")
	}
}

// sendMessageWebhook envía un webhook de mensaje de usuario
func (s *ConversationSimulator) sendMessageWebhook(phoneNumber, buttonText string) {
	newMessageID := fmt.Sprintf("wamid.HBg%s", uuid.New().String()[:20])

	payload := domain.WebhookPayload{
		Object: "whatsapp_business_account",
		Entry: []domain.WebhookEntry{{
			ID: "123456789",
			Changes: []domain.WebhookChange{{
				Value: domain.WebhookValue{
					MessagingProduct: "whatsapp",
					Metadata: domain.WebhookMetadata{
						DisplayPhoneNumber: "+15551234567",
						PhoneNumberID:      "123456789012345",
					},
					Contacts: []domain.WebhookContact{{
						Profile: domain.WebhookProfile{Name: "Test User"},
						WaID:    phoneNumber,
					}},
					Messages: []domain.WebhookMessage{{
						From:      phoneNumber,
						ID:        newMessageID,
						Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
						Type:      "button",
						Button: &domain.ButtonResponse{
							Payload: buttonText,
							Text:    buttonText,
						},
					}},
				},
				Field: "messages",
			}},
		}},
	}

	err := s.webhookClient.SendWebhook(payload)
	if err != nil {
		s.logger.Error().Err(err).Str("button_text", buttonText).Msg("Error al enviar webhook de mensaje")
	} else {
		s.logger.Info().Str("button_text", buttonText).Str("phone_number", phoneNumber).Msg("Webhook de mensaje enviado")
	}
}

// getUserResponseForTemplate retorna la respuesta del usuario según el template
func (s *ConversationSimulator) getUserResponseForTemplate(templateName string) string {
	responses := map[string]string{
		"confirmacion_pedido_contraentrega": "Confirmar pedido",
		"menu_no_confirmacion":              "Presentar novedad",
		"solicitud_novedad":                 "Otro",
		"novedad_otro":                      "El producto llegó dañado",
		"pedido_cancelado":                  "", // No espera respuesta
		"pedido_confirmado":                 "", // No espera respuesta
		"solicitud_cancelacion":             "Sí, cancelar",
		"motivo_cancelacion":                "Ya no lo necesito",
		"confirmacion_cancelacion_pago":     "Sí, he cancelado el pago",
	}

	if response, ok := responses[templateName]; ok {
		return response
	}

	return "Asesor humano" // Default: solicitar asesor
}

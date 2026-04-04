package usecasemessaging

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

// SendManualReply envía un mensaje de texto libre desde el dashboard del agente.
// Requiere que el cliente haya enviado un mensaje en las últimas 24h (ventana de servicio WhatsApp).
// conversationID: UUID de la conversación en BD
// phoneNumber: número destino (formato internacional, ej: "573001234567")
// businessID: negocio propietario de las credenciales WhatsApp
// text: cuerpo del mensaje
// sentBy: user_id del agente que envía (para auditoría)
func (u *usecases) SendManualReply(
	ctx context.Context,
	conversationID string,
	phoneNumber string,
	businessID uint,
	text string,
	sentBy string,
) (string, error) {
	u.log.Info(ctx).
		Str("conversation_id", conversationID).
		Str("phone_number", phoneNumber).
		Uint("business_id", businessID).
		Str("sent_by", sentBy).
		Msg("[WhatsApp UseCase] - enviando reply manual del agente")

	// 1. Normalizar número
	phoneNumber = NormalizePhoneNumber(phoneNumber)
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return "", fmt.Errorf("número de teléfono inválido: %w", err)
	}

	// 2. Obtener credenciales WhatsApp del negocio
	whatsappConfig, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Uint("business_id", businessID).
			Msg("[WhatsApp UseCase] - error obteniendo config WhatsApp para reply manual")
		return "", fmt.Errorf("error obteniendo configuración de WhatsApp: %w", err)
	}

	// 3. Seleccionar cliente (con URL específica del negocio si está configurada)
	waClient := u.whatsApp
	if whatsappConfig.WhatsAppURL != "" && u.clientFactory != nil {
		waClient = u.clientFactory(whatsappConfig.WhatsAppURL)
	}

	// 4. Enviar mensaje de texto libre
	messageID, err := waClient.SendTextMessage(ctx, whatsappConfig.PhoneNumberID, phoneNumber, text, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp UseCase] - error enviando reply manual")
		return "", fmt.Errorf("error al enviar mensaje: %w", err)
	}

	// 5. Activar sesión humana en Redis para que las respuestas del cliente lleguen al dashboard
	if hsErr := u.conversationCache.ActivateHumanSession(ctx, phoneNumber, conversationID, businessID); hsErr != nil {
		u.log.Error(ctx).Err(hsErr).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp UseCase] - error activando human session")
		// No retornamos error: el mensaje ya fue enviado
	}

	// 6. Persistir en message_logs (async via RabbitMQ)
	messageLog := &entities.MessageLog{
		ConversationID: conversationID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   "",
		Content:        text,
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
	}

	if err := u.persistPublisher.PublishMessageLogCreated(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Msg("[WhatsApp UseCase] - error publicando message log de reply manual")
		// No retornamos error: el mensaje ya fue enviado
	}

	u.log.Info(ctx).
		Str("message_id", messageID).
		Str("conversation_id", conversationID).
		Str("sent_by", sentBy).
		Msg("[WhatsApp UseCase] - reply manual enviado exitosamente")

	return messageID, nil
}

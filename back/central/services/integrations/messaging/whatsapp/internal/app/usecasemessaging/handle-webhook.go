package usecasemessaging

import (
	"context"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
)

// HandleIncomingMessage procesa mensajes entrantes del usuario
func (u *usecases) HandleIncomingMessage(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error {
	u.log.Info(ctx).Msg("[WhatsApp Webhook] - procesando mensaje entrante")

	// Extraer mensajes del webhook
	for _, entry := range whPayload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			for _, message := range change.Value.Messages {
				if err := u.processIncomingMessage(ctx, message, change.Value.Metadata); err != nil {
					u.log.Error(ctx).Err(err).
						Str("message_id", message.ID).
						Str("from", message.From).
						Msg("[WhatsApp Webhook] - error procesando mensaje")
					// No retornamos error para no bloquear otros mensajes
					continue
				}
			}
		}
	}

	return nil
}

// processIncomingMessage procesa un mensaje individual
func (u *usecases) processIncomingMessage(ctx context.Context, message dtos.WebhookMessageDTO, metadata dtos.WebhookMetadataDTO) error {
	phoneNumber := message.From
	messageText := message.GetMessageText()

	u.log.Info(ctx).
		Str("from", phoneNumber).
		Str("message_id", message.ID).
		Str("text", messageText).
		Str("type", message.Type).
		Msg("[WhatsApp Webhook] - procesando mensaje del usuario")

	// 1. Buscar conversación activa del usuario
	conversation, err := u.conversationRepo.GetActiveByPhone(ctx, phoneNumber)
	if err != nil {
		u.log.Debug(ctx).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp Webhook] - no hay conversación activa para este usuario")
		// No hay conversación activa, ignorar mensaje
		return nil
	}

	// 2. Verificar que no ha expirado
	if conversation.IsExpired() {
		u.log.Warn(ctx).
			Str("conversation_id", conversation.ID).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp Webhook] - conversación expirada")
		return &errors.ErrConversationExpired{ConversationID: conversation.ID}
	}

	// 3. Registrar mensaje entrante en log
	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionInbound,
		MessageID:      message.ID,
		Content:        messageText,
		Status:         entities.MessageStatusDelivered, // Los mensajes entrantes ya están entregados
		CreatedAt:      time.Now(),
	}

	if err := u.messageLogRepo.Create(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", message.ID).
			Msg("[WhatsApp Webhook] - error registrando mensaje entrante")
		// Continuamos aunque falle el log
	}

	// 4. Procesar según el estado actual de la conversación
	if err := u.processConversationFlow(ctx, conversation, messageText); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Str("state", string(conversation.CurrentState)).
			Msg("[WhatsApp Webhook] - error procesando flujo conversacional")
		return err
	}

	return nil
}

// processConversationFlow maneja el flujo de la conversación según el estado actual
func (u *usecases) processConversationFlow(ctx context.Context, conversation *entities.Conversation, userResponse string) error {
	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("current_state", string(conversation.CurrentState)).
		Str("user_response", userResponse).
		Msg("[WhatsApp Webhook] - evaluando transición de estado")

	// 1. Usar el ConversationManager para determinar la transición
	transition, err := u.TransitionState(ctx, conversation, userResponse)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("state", string(conversation.CurrentState)).
			Str("response", userResponse).
			Msg("[WhatsApp Webhook] - transición inválida")
		return err
	}

	// 2. Guardar metadata del evento si existe
	if transition.EventMetadata != nil {
		if conversation.Metadata == nil {
			conversation.Metadata = make(map[string]interface{})
		}
		for key, value := range transition.EventMetadata {
			conversation.Metadata[key] = value
		}
	}

	// 3. Enviar siguiente plantilla
	_, err = u.SendTemplateWithConversation(
		ctx,
		transition.TemplateName,
		conversation.PhoneNumber,
		transition.Variables,
		conversation.ID,
	)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("template", transition.TemplateName).
			Msg("[WhatsApp Webhook] - error enviando plantilla")
		return err
	}

	// 4. Actualizar estado de la conversación
	conversation.CurrentState = transition.NextState
	conversation.UpdatedAt = time.Now()

	if err := u.conversationRepo.Update(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error actualizando conversación")
		return err
	}

	// 5. Publicar evento de negocio si aplica
	if transition.PublishEvent {
		if err := u.publishBusinessEvent(ctx, transition.EventType, conversation); err != nil {
			u.log.Error(ctx).Err(err).
				Str("event_type", transition.EventType).
				Msg("[WhatsApp Webhook] - error publicando evento de negocio")
			// No retornamos error
		}
	}

	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("new_state", string(transition.NextState)).
		Str("template_sent", transition.TemplateName).
		Bool("event_published", transition.PublishEvent).
		Msg("[WhatsApp Webhook] - transición de estado completada")

	return nil
}

// publishBusinessEvent publica eventos de negocio según el tipo
func (u *usecases) publishBusinessEvent(ctx context.Context, eventType string, conversation *entities.Conversation) error {
	switch eventType {
	case "confirmed":
		return u.publisher.PublishOrderConfirmed(
			ctx,
			conversation.OrderNumber,
			conversation.PhoneNumber,
			conversation.BusinessID,
		)
	case "cancelled":
		reason := ""
		if r, ok := conversation.Metadata["cancellation_reason"].(string); ok {
			reason = r
		}
		return u.publisher.PublishOrderCancelled(
			ctx,
			conversation.OrderNumber,
			reason,
			conversation.PhoneNumber,
			conversation.BusinessID,
		)
	case "novelty":
		noveltyType := ""
		if n, ok := conversation.Metadata["novelty_type"].(string); ok {
			noveltyType = n
		}
		return u.publisher.PublishNoveltyRequested(
			ctx,
			conversation.OrderNumber,
			noveltyType,
			conversation.PhoneNumber,
			conversation.BusinessID,
		)
	case "handoff":
		return u.publisher.PublishHandoffRequested(
			ctx,
			conversation.OrderNumber,
			conversation.PhoneNumber,
			conversation.BusinessID,
			conversation.ID,
		)
	}
	return nil
}

// HandleMessageStatus procesa cambios de estado de mensajes (delivered, read)
func (u *usecases) HandleMessageStatus(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error {
	u.log.Info(ctx).Msg("[WhatsApp Webhook] - procesando cambios de estado de mensajes")

	for _, entry := range whPayload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			for _, status := range change.Value.Statuses {
				if err := u.processMessageStatus(ctx, status); err != nil {
					u.log.Error(ctx).Err(err).
						Str("message_id", status.ID).
						Str("status", status.Status).
						Msg("[WhatsApp Webhook] - error procesando estado de mensaje")
					// Continuamos con otros estados
					continue
				}
			}
		}
	}

	return nil
}

// processMessageStatus procesa un cambio de estado individual
func (u *usecases) processMessageStatus(ctx context.Context, status dtos.WebhookStatusDTO) error {
	u.log.Info(ctx).
		Str("message_id", status.ID).
		Str("status", status.Status).
		Msg("[WhatsApp Webhook] - actualizando estado de mensaje")

	// Convertir status string a MessageStatus
	var messageStatus entities.MessageStatus
	switch status.Status {
	case "sent":
		messageStatus = entities.MessageStatusSent
	case "delivered":
		messageStatus = entities.MessageStatusDelivered
	case "read":
		messageStatus = entities.MessageStatusRead
	case "failed":
		messageStatus = entities.MessageStatusFailed
	default:
		u.log.Warn(ctx).
			Str("status", status.Status).
			Msg("[WhatsApp Webhook] - estado de mensaje desconocido")
		return nil
	}

	// Preparar timestamps
	timestamps := make(map[string]time.Time)
	timestampInt, _ := strconv.ParseInt(status.Timestamp, 10, 64)
	timestamp := time.Unix(timestampInt, 0)

	if messageStatus == entities.MessageStatusDelivered {
		timestamps["delivered_at"] = timestamp
	} else if messageStatus == entities.MessageStatusRead {
		timestamps["read_at"] = timestamp
	}

	// Actualizar en el repositorio
	if err := u.messageLogRepo.UpdateStatus(ctx, status.ID, messageStatus, timestamps); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", status.ID).
			Msg("[WhatsApp Webhook] - error actualizando estado en BD")
		return err
	}

	u.log.Info(ctx).
		Str("message_id", status.ID).
		Str("status", string(messageStatus)).
		Msg("[WhatsApp Webhook] - estado de mensaje actualizado exitosamente")

	return nil
}

package usecasemessaging

import (
	"context"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func (u *usecases) HandleIncomingMessage(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error {
	u.log.Info(ctx).Msg("[WhatsApp Webhook] - procesando mensaje entrante")

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
					continue
				}
			}
		}
	}

	return nil
}

func (u *usecases) processIncomingMessage(ctx context.Context, message dtos.WebhookMessageDTO, metadata dtos.WebhookMetadataDTO) error {
	phoneNumber := message.From
	messageText := message.GetMessageText()

	u.log.Info(ctx).
		Str("from", phoneNumber).
		Str("message_id", message.ID).
		Str("text", messageText).
		Str("type", message.Type).
		Msg("[WhatsApp Webhook] - procesando mensaje del usuario")

	conversation, err := u.conversationCache.GetActiveByPhone(ctx, phoneNumber)
	if err != nil {
		u.log.Debug(ctx).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp Webhook] - no hay conversación activa para este usuario")

		if humanSession, hsErr := u.conversationCache.GetHumanSession(ctx, phoneNumber); hsErr == nil && humanSession != nil {
			u.log.Info(ctx).
				Str("phone_number", phoneNumber).
				Str("conversation_id", humanSession.ConversationID).
				Msg("[WhatsApp Webhook] - mensaje enrutado a sesión humana (dashboard)")

			humanLog := &entities.MessageLog{
				ConversationID: humanSession.ConversationID,
				Direction:      entities.MessageDirectionInbound,
				MessageID:      message.ID,
				Content:        messageText,
				Status:         entities.MessageStatusDelivered,
				CreatedAt:      time.Now(),
			}
			if logErr := u.persistPublisher.PublishMessageLogCreated(ctx, humanLog); logErr != nil {
				u.log.Error(ctx).Err(logErr).
					Str("conversation_id", humanSession.ConversationID).
					Msg("[WhatsApp Webhook] - error persistiendo mensaje de sesión humana")
			}

			if sseErr := u.ssePublisher.PublishMessageReceived(
				ctx,
				humanSession.BusinessID,
				humanSession.ConversationID,
				phoneNumber,
				message.ID,
				messageText,
			); sseErr != nil {
				u.log.Error(ctx).Err(sseErr).
					Str("conversation_id", humanSession.ConversationID).
					Msg("[WhatsApp Webhook] - error publicando SSE de sesión humana")
			}

			return nil
		}

		if handled, routeErr := u.routeToOwnNumberBusiness(ctx, message, metadata, messageText); routeErr != nil {
			return routeErr
		} else if handled {
			return nil
		}

		if u.aiForwarder != nil {
			if fwdErr := u.aiForwarder.ForwardToAI(ctx, phoneNumber, messageText, message.ID, message.Type); fwdErr != nil {
				u.log.Error(ctx).Err(fwdErr).
					Str("phone_number", phoneNumber).
					Msg("[WhatsApp Webhook] - error reenviando a AI Sales")
			}
		}

		return nil
	}

	if conversation.IsExpired() {
		u.log.Warn(ctx).
			Str("conversation_id", conversation.ID).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp Webhook] - conversación expirada")
		return &errors.ErrConversationExpired{ConversationID: conversation.ID}
	}

	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionInbound,
		MessageID:      message.ID,
		Content:        messageText,
		Status:         entities.MessageStatusDelivered,
		CreatedAt:      time.Now(),
	}

	if err := u.persistPublisher.PublishMessageLogCreated(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", message.ID).
			Msg("[WhatsApp Webhook] - error publicando mensaje entrante en log")
	}

	if sseErr := u.ssePublisher.PublishMessageReceived(
		ctx,
		conversation.BusinessID,
		conversation.ID,
		phoneNumber,
		message.ID,
		messageText,
	); sseErr != nil {
		u.log.Error(ctx).Err(sseErr).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error publicando SSE message_received")
	}

	if err := u.processConversationFlow(ctx, conversation, messageText); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Str("state", string(conversation.CurrentState)).
			Msg("[WhatsApp Webhook] - error procesando flujo conversacional")
		return err
	}

	return nil
}

func (u *usecases) resolveBusinessByPhoneNumberID(ctx context.Context, phoneNumberID string) *ports.IntegrationOwner {
	if phoneNumberID == "" || u.credentialsCache == nil {
		return nil
	}

	owner, err := u.credentialsCache.ResolveByPhoneNumberID(ctx, phoneNumberID)
	if err != nil {
		u.log.Warn(ctx).Err(err).
			Str("phone_number_id", phoneNumberID).
			Msg("[WhatsApp Webhook] - no se pudo resolver el negocio por phone_number_id")
		return nil
	}
	if owner == nil || owner.BusinessID == 0 {
		return nil
	}

	return owner
}

func (u *usecases) routeToOwnNumberBusiness(
	ctx context.Context,
	message dtos.WebhookMessageDTO,
	metadata dtos.WebhookMetadataDTO,
	messageText string,
) (bool, error) {
	owner := u.resolveBusinessByPhoneNumberID(ctx, metadata.PhoneNumberID)
	if owner == nil {
		return false, nil
	}

	phoneNumber := message.From

	conversation, err := u.getOrCreateInboundConversation(ctx, phoneNumber, owner.BusinessID)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Uint("business_id", owner.BusinessID).
			Msg("[WhatsApp Webhook] - error creando conversación entrante para número propio")
		return true, err
	}

	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionInbound,
		MessageID:      message.ID,
		Content:        messageText,
		Status:         entities.MessageStatusDelivered,
		CreatedAt:      time.Now(),
	}
	if logErr := u.persistPublisher.PublishMessageLogCreated(ctx, messageLog); logErr != nil {
		u.log.Error(ctx).Err(logErr).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error persistiendo mensaje entrante de número propio")
	}

	if sseErr := u.ssePublisher.PublishMessageReceived(
		ctx,
		owner.BusinessID,
		conversation.ID,
		phoneNumber,
		message.ID,
		messageText,
	); sseErr != nil {
		u.log.Error(ctx).Err(sseErr).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error publicando SSE de mensaje entrante")
	}

	u.log.Info(ctx).
		Str("phone_number", phoneNumber).
		Str("phone_number_id", metadata.PhoneNumberID).
		Uint("business_id", owner.BusinessID).
		Uint("integration_id", owner.IntegrationID).
		Str("conversation_id", conversation.ID).
		Msg("[WhatsApp Webhook] - mensaje enrutado al negocio dueño del número")

	return true, nil
}

func (u *usecases) getOrCreateInboundConversation(
	ctx context.Context,
	phoneNumber string,
	businessID uint,
) (*entities.Conversation, error) {
	if existing, err := u.conversationCache.GetActiveByPhone(ctx, phoneNumber); err == nil && existing != nil {
		return existing, nil
	}

	conversation := &entities.Conversation{
		PhoneNumber:      phoneNumber,
		ConversationType: entities.ConversationTypeInbound,
		BusinessID:       businessID,
		CurrentState:     entities.StateHandoffToHuman,
		Metadata:         make(map[string]interface{}),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	if err := u.conversationCache.Save(ctx, conversation); err != nil {
		return nil, err
	}

	if err := u.persistPublisher.PublishConversationCreated(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error publicando creación de conversación entrante")
	}

	if err := u.ssePublisher.PublishConversationStarted(ctx, businessID, conversation.ID, phoneNumber); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error publicando SSE conversation_started")
	}

	return conversation, nil
}

func (u *usecases) processConversationFlow(ctx context.Context, conversation *entities.Conversation, userResponse string) error {
	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("current_state", string(conversation.CurrentState)).
		Str("user_response", userResponse).
		Msg("[WhatsApp Webhook] - evaluando transición de estado")

	transition, err := u.TransitionState(ctx, conversation, userResponse)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("state", string(conversation.CurrentState)).
			Str("response", userResponse).
			Msg("[WhatsApp Webhook] - transición inválida")
		return err
	}

	if transition.EventMetadata != nil {
		if conversation.Metadata == nil {
			conversation.Metadata = make(map[string]interface{})
		}
		for key, value := range transition.EventMetadata {
			conversation.Metadata[key] = value
		}
	}

	var templateSendErr error
	_, templateSendErr = u.SendTemplateWithConversation(
		ctx,
		transition.TemplateName,
		conversation.PhoneNumber,
		transition.Variables,
		conversation.ID,
	)
	if templateSendErr != nil {
		u.log.Error(ctx).Err(templateSendErr).
			Str("template", transition.TemplateName).
			Msg("[WhatsApp Webhook] - error enviando plantilla (continuando con evento de negocio)")
	}

	conversation.CurrentState = transition.NextState
	conversation.UpdatedAt = time.Now()

	if err := u.conversationCache.Save(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error actualizando conversación en cache")
	}

	if err := u.persistPublisher.PublishConversationUpdated(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp Webhook] - error publicando actualización de conversación")
	}

	if transition.PublishEvent {
		if err := u.publishBusinessEvent(ctx, transition.EventType, conversation); err != nil {
			u.log.Error(ctx).Err(err).
				Str("event_type", transition.EventType).
				Msg("[WhatsApp Webhook] - error publicando evento de negocio")
		}
	}

	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("new_state", string(transition.NextState)).
		Str("template_sent", transition.TemplateName).
		Bool("event_published", transition.PublishEvent).
		Bool("template_sent_ok", templateSendErr == nil).
		Msg("[WhatsApp Webhook] - transición de estado completada")

	return templateSendErr
}

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

func (u *usecases) HandleMessageStatus(ctx context.Context, whPayload dtos.WebhookPayloadDTO) error {
	u.log.Info(ctx).Msg("[WhatsApp Webhook] - procesando cambios de estado de mensajes")

	for _, entry := range whPayload.Entry {
		for _, change := range entry.Changes {
			if change.Field != "messages" {
				continue
			}

			for _, status := range change.Value.Statuses {
				if err := u.processMessageStatus(ctx, status, change.Value.Metadata); err != nil {
					u.log.Error(ctx).Err(err).
						Str("message_id", status.ID).
						Str("status", status.Status).
						Msg("[WhatsApp Webhook] - error procesando estado de mensaje")
					continue
				}
			}
		}
	}

	return nil
}

func (u *usecases) processMessageStatus(ctx context.Context, status dtos.WebhookStatusDTO, metadata dtos.WebhookMetadataDTO) error {
	u.log.Info(ctx).
		Str("message_id", status.ID).
		Str("status", status.Status).
		Msg("[WhatsApp Webhook] - actualizando estado de mensaje")

	if status.Status == "failed" && len(status.Errors) > 0 {
		for _, wErr := range status.Errors {
			u.log.Error(ctx).
				Str("message_id", status.ID).
				Int("error_code", wErr.Code).
				Str("error_title", wErr.Title).
				Str("error_message", wErr.Message).
				Str("error_details", wErr.Details).
				Msg("[WhatsApp Webhook] - Meta rechazo el mensaje")
		}
	}

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

	timestamps := make(map[string]time.Time)
	timestampInt, _ := strconv.ParseInt(status.Timestamp, 10, 64)
	timestamp := time.Unix(timestampInt, 0)

	if messageStatus == entities.MessageStatusDelivered {
		timestamps["delivered_at"] = timestamp
	} else if messageStatus == entities.MessageStatusRead {
		timestamps["read_at"] = timestamp
	}

	if err := u.persistPublisher.PublishMessageStatusUpdated(ctx, status.ID, messageStatus, timestamps); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", status.ID).
			Msg("[WhatsApp Webhook] - error publicando actualización de estado")
		return err
	}

	businessID := uint(0)
	if conv, convErr := u.conversationCache.GetActiveByPhone(ctx, status.RecipientID); convErr == nil && conv != nil {
		businessID = conv.BusinessID
	} else if owner := u.resolveBusinessByPhoneNumberID(ctx, metadata.PhoneNumberID); owner != nil {
		businessID = owner.BusinessID
	}

	if businessID > 0 {
		if sseErr := u.ssePublisher.PublishMessageStatusUpdated(ctx, businessID, status.ID, string(messageStatus)); sseErr != nil {
			u.log.Error(ctx).Err(sseErr).
				Str("message_id", status.ID).
				Msg("[WhatsApp Webhook] - error publicando SSE message_status_updated")
		}
	}

	u.log.Info(ctx).
		Str("message_id", status.ID).
		Str("status", string(messageStatus)).
		Msg("[WhatsApp Webhook] - estado de mensaje publicado para persistencia")

	return nil
}

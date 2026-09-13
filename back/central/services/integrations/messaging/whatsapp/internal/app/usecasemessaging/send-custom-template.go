package usecasemessaging

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
)

func (u *usecases) SendCustomTemplate(
	ctx context.Context,
	templateName string,
	language string,
	phoneNumber string,
	parameters []string,
	headerImageURL string,
	businessID uint,
) (string, error) {
	templateName = strings.TrimSpace(templateName)
	if templateName == "" {
		return "", fmt.Errorf("el nombre de la plantilla es obligatorio")
	}

	if strings.TrimSpace(language) == "" {
		language = "es"
	}

	phoneNumber = NormalizePhoneNumber(phoneNumber)
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		return "", fmt.Errorf("numero de telefono invalido: %w", err)
	}

	whatsappConfig, err := u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		return "", fmt.Errorf("error obteniendo configuracion de WhatsApp: %w", err)
	}

	components := []entities.TemplateComponent{}

	if headerImage := strings.TrimSpace(headerImageURL); headerImage != "" {
		components = append(components, entities.TemplateComponent{
			Type: "header",
			Parameters: []entities.TemplateParameter{
				{Type: "image", ImageLink: headerImage},
			},
		})
	}

	if len(parameters) > 0 {
		bodyParams := make([]entities.TemplateParameter, 0, len(parameters))
		for _, value := range parameters {
			bodyParams = append(bodyParams, entities.TemplateParameter{
				Type: "text",
				Text: value,
			})
		}
		components = append(components, entities.TemplateComponent{
			Type:       "body",
			Parameters: bodyParams,
		})
	}

	msg := entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: entities.TemplateData{
			Name:       templateName,
			Language:   entities.TemplateLanguage{Code: language},
			Components: components,
		},
	}

	waClient := u.whatsApp
	if whatsappConfig.WhatsAppURL != "" && u.clientFactory != nil {
		waClient = u.clientFactory(whatsappConfig.WhatsAppURL)
	}

	messageID, err := waClient.SendMessage(ctx, whatsappConfig.PhoneNumberID, msg, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("template_name", templateName).
			Str("phone_number", phoneNumber).
			Uint("business_id", businessID).
			Msg("[WhatsApp UseCase] - error enviando plantilla propia del negocio")
		return "", err
	}

	u.recordCustomTemplateSend(ctx, templateName, phoneNumber, messageID, parameters, businessID)

	u.log.Info(ctx).
		Str("template_name", templateName).
		Str("phone_number", phoneNumber).
		Str("message_id", messageID).
		Msg("[WhatsApp UseCase] - plantilla propia del negocio enviada")

	return messageID, nil
}

func (u *usecases) recordCustomTemplateSend(
	ctx context.Context,
	templateName string,
	phoneNumber string,
	messageID string,
	parameters []string,
	businessID uint,
) {
	if u.persistPublisher == nil || u.conversationCache == nil {
		return
	}

	conversation, err := u.getOrCreateOutboundConversation(ctx, phoneNumber, businessID)
	if err != nil || conversation == nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Str("template_name", templateName).
			Msg("[WhatsApp UseCase] - no se pudo registrar el envio: el encadenado de botones no va a resolver")
		return
	}

	content := templateName
	if len(parameters) > 0 {
		content = fmt.Sprintf("%s: %s", templateName, strings.Join(parameters, " | "))
	}

	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        content,
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
	}

	if err := u.persistPublisher.PublishMessageLogCreated(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Msg("[WhatsApp UseCase] - error publicando el envio en el log")
	}

	conversation.LastMessageID = messageID
	conversation.LastTemplateID = templateName
	conversation.UpdatedAt = time.Now()

	if err := u.conversationCache.Save(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - error actualizando la conversacion en cache")
	}

	if err := u.persistPublisher.PublishConversationUpdated(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - error publicando la actualizacion de la conversacion")
	}
}

func (u *usecases) getOrCreateOutboundConversation(
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
			Msg("[WhatsApp UseCase] - error publicando la creacion de la conversacion saliente")
	}

	if u.ssePublisher != nil {
		if err := u.ssePublisher.PublishConversationStarted(ctx, businessID, conversation.ID, phoneNumber); err != nil {
			u.log.Error(ctx).Err(err).
				Str("conversation_id", conversation.ID).
				Msg("[WhatsApp UseCase] - error publicando SSE conversation_started")
		}
	}

	return conversation, nil
}

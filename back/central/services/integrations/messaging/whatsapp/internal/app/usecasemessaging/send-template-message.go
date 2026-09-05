package usecasemessaging

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func (u *usecases) SendTemplate(
	ctx context.Context,
	templateName string,
	phoneNumber string,
	variables map[string]string,
	orderNumber string,
	businessID uint,
) (string, error) {
	return u.sendTemplate(ctx, templateName, phoneNumber, variables, orderNumber, businessID, false)
}

func (u *usecases) SendPlatformTemplate(
	ctx context.Context,
	templateName string,
	phoneNumber string,
	variables map[string]string,
	orderNumber string,
	businessID uint,
) (string, error) {
	return u.sendTemplate(ctx, templateName, phoneNumber, variables, orderNumber, businessID, true)
}

func (u *usecases) sendTemplate(
	ctx context.Context,
	templateName string,
	phoneNumber string,
	variables map[string]string,
	orderNumber string,
	businessID uint,
	forcePlatformNumber bool,
) (string, error) {
	u.log.Info(ctx).
		Str("template_name", templateName).
		Str("phone_number", phoneNumber).
		Str("order_number", orderNumber).
		Bool("platform_number", forcePlatformNumber).
		Msg("[WhatsApp UseCase] - enviando plantilla")

	templateDef, exists := entities.GetTemplateDefinition(templateName)
	if !exists {
		u.log.Error(ctx).
			Str("template_name", templateName).
			Msg("[WhatsApp UseCase] - plantilla no encontrada")
		return "", &errors.ErrTemplateNotFound{TemplateName: templateName}
	}

	variables = SanitizeTemplateVariables(variables)

	if err := entities.ValidateTemplateVariables(templateName, variables); err != nil {
		u.log.Error(ctx).Err(err).
			Str("template_name", templateName).
			Msg("[WhatsApp UseCase] - variables faltantes")
		return "", err
	}

	phoneNumber = NormalizePhoneNumber(phoneNumber)
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp UseCase] - número de teléfono inválido")
		return "", fmt.Errorf("número de teléfono inválido: %w", err)
	}

	var whatsappConfig *ports.WhatsAppConfig
	var err error
	if forcePlatformNumber {
		whatsappConfig, err = u.credentialsCache.GetWhatsAppDefaultConfig(ctx)
	} else {
		whatsappConfig, err = u.credentialsCache.GetWhatsAppConfig(ctx, businessID)
	}
	if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error obteniendo configuración de WhatsApp")
		return "", fmt.Errorf("error obteniendo configuración de WhatsApp: %w", err)
	}

	msg := u.buildTemplateMessage(templateName, phoneNumber, variables, templateDef)

	conversation, err := u.getOrCreateConversation(ctx, phoneNumber, orderNumber, businessID)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Str("order_number", orderNumber).
			Msg("[WhatsApp UseCase] - error obteniendo/creando conversación")
		return "", err
	}

	if conversation.Metadata == nil {
		conversation.Metadata = make(map[string]interface{})
	}
	for i, varName := range templateDef.Variables {
		varKey := strconv.Itoa(i + 1)
		if val, ok := variables[varKey]; ok && val != "" {
			conversation.Metadata[varName] = val
		}
	}

	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("template_name", templateName).
		Uint("phone_number_id", whatsappConfig.PhoneNumberID).
		Str("whatsapp_url", whatsappConfig.WhatsAppURL).
		Msg("[WhatsApp UseCase] - enviando mensaje a WhatsApp API")

	waClient := u.whatsApp
	if whatsappConfig.WhatsAppURL != "" && u.clientFactory != nil {
		waClient = u.clientFactory(whatsappConfig.WhatsAppURL)
	}

	messageID, err := waClient.SendMessage(ctx, whatsappConfig.PhoneNumberID, msg, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("template_name", templateName).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp UseCase] - error enviando mensaje")
		return "", fmt.Errorf("error al enviar mensaje de WhatsApp: %w", err)
	}

	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        buildTemplateContent(templateName, variables),
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
	}

	if err := u.persistPublisher.PublishMessageLogCreated(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Msg("[WhatsApp UseCase] - error publicando mensaje en log")
	}

	conversation.LastMessageID = messageID
	conversation.LastTemplateID = templateName
	conversation.UpdatedAt = time.Now()

	if err := u.conversationCache.Save(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - error actualizando conversación en cache")
	}

	if err := u.persistPublisher.PublishConversationUpdated(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - error publicando actualización de conversación")
	}

	u.log.Info(ctx).
		Str("message_id", messageID).
		Str("conversation_id", conversation.ID).
		Str("template_name", templateName).
		Msg("[WhatsApp UseCase] - mensaje enviado exitosamente")

	return messageID, nil
}

func (u *usecases) SendTemplateWithConversation(
	ctx context.Context,
	templateName string,
	phoneNumber string,
	variables map[string]string,
	conversationID string,
) (string, error) {
	u.log.Info(ctx).
		Str("template_name", templateName).
		Str("conversation_id", conversationID).
		Msg("[WhatsApp UseCase] - enviando plantilla con conversación existente")

	templateDef, exists := entities.GetTemplateDefinition(templateName)
	if !exists {
		return "", &errors.ErrTemplateNotFound{TemplateName: templateName}
	}

	variables = SanitizeTemplateVariables(variables)

	if err := entities.ValidateTemplateVariables(templateName, variables); err != nil {
		return "", err
	}

	phoneNumber = NormalizePhoneNumber(phoneNumber)

	conversation, err := u.conversationCache.GetByID(ctx, conversationID)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp UseCase] - conversación no encontrada")
		return "", err
	}

	if conversation.IsExpired() {
		u.log.Error(ctx).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp UseCase] - conversación expirada")
		return "", &errors.ErrConversationExpired{ConversationID: conversationID}
	}

	whatsappConfig, err := u.credentialsCache.GetWhatsAppConfig(ctx, conversation.BusinessID)
	if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error obteniendo configuración de WhatsApp")
		return "", fmt.Errorf("error obteniendo configuración de WhatsApp: %w", err)
	}

	msg := u.buildTemplateMessage(templateName, phoneNumber, variables, templateDef)

	waClient := u.whatsApp
	if whatsappConfig.WhatsAppURL != "" && u.clientFactory != nil {
		waClient = u.clientFactory(whatsappConfig.WhatsAppURL)
	}

	messageID, err := waClient.SendMessage(ctx, whatsappConfig.PhoneNumberID, msg, whatsappConfig.AccessToken)
	if err != nil {
		return "", err
	}

	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        buildTemplateContent(templateName, nil),
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
	}
	u.persistPublisher.PublishMessageLogCreated(ctx, messageLog)

	conversation.LastMessageID = messageID
	conversation.LastTemplateID = templateName
	conversation.UpdatedAt = time.Now()
	u.conversationCache.Save(ctx, conversation)
	u.persistPublisher.PublishConversationUpdated(ctx, conversation)

	return messageID, nil
}

func (u *usecases) buildTemplateMessage(
	templateName string,
	phoneNumber string,
	variables map[string]string,
	templateDef entities.TemplateDefinition,
) entities.TemplateMessage {
	components := []entities.TemplateComponent{}

	if len(templateDef.Variables) > 0 {
		bodyParams := []entities.TemplateParameter{}
		for i := range templateDef.Variables {
			varKey := strconv.Itoa(i + 1)
			bodyParams = append(bodyParams, entities.TemplateParameter{
				Type: "text",
				Text: variables[varKey],
			})
		}
		components = append(components, entities.TemplateComponent{
			Type:       "body",
			Parameters: bodyParams,
		})
	}

	return entities.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: entities.TemplateData{
			Name:       templateName,
			Language:   entities.TemplateLanguage{Code: templateDef.Language},
			Components: components,
		},
	}
}

func (u *usecases) getOrCreateConversation(
	ctx context.Context,
	phoneNumber string,
	orderNumber string,
	businessID uint,
) (*entities.Conversation, error) {
	conversation, err := u.conversationCache.GetByPhoneAndOrder(ctx, phoneNumber, orderNumber)
	if err == nil {
		if conversation.IsActive() {
			return conversation, nil
		}
		u.log.Info(ctx).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - conversación expirada, creando nueva")
	}

	conversationType := entities.ConversationTypeOrder
	initialState := entities.StateAwaitingConfirmation
	if orderNumber == "" {
		conversationType = entities.ConversationTypeSystemAlert
		initialState = entities.StateStart
	}

	newConversation := &entities.Conversation{
		PhoneNumber:      phoneNumber,
		OrderNumber:      orderNumber,
		ConversationType: conversationType,
		BusinessID:       businessID,
		CurrentState:     initialState,
		Metadata:         make(map[string]interface{}),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
		ExpiresAt:        time.Now().Add(24 * time.Hour),
	}

	if err := u.conversationCache.Save(ctx, newConversation); err != nil {
		return nil, err
	}

	if err := u.persistPublisher.PublishConversationCreated(ctx, newConversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", newConversation.ID).
			Msg("[WhatsApp UseCase] - error publicando creación de conversación")
	}

	u.log.Info(ctx).
		Str("conversation_id", newConversation.ID).
		Str("phone_number", phoneNumber).
		Str("order_number", orderNumber).
		Msg("[WhatsApp UseCase] - nueva conversación creada")

	return newConversation, nil
}

func buildTemplateContent(templateName string, variables map[string]string) string {
	if rendered := entities.RenderTemplateBody(templateName, variables); rendered != "" {
		return rendered
	}
	if len(variables) == 0 {
		return fmt.Sprintf("Plantilla: %s", templateName)
	}
	return fmt.Sprintf("Plantilla: %s (variables: %v)", templateName, variables)
}

package usecasemessaging

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain/errors"
)

// SendTemplate envía una plantilla de WhatsApp y crea/actualiza la conversación
func (u *Usecases) SendTemplate(
	ctx context.Context,
	templateName string,
	phoneNumber string,
	variables map[string]string,
	orderNumber string,
	businessID uint,
) (string, error) {
	u.log.Info(ctx).
		Str("template_name", templateName).
		Str("phone_number", phoneNumber).
		Str("order_number", orderNumber).
		Msg("[WhatsApp UseCase] - enviando plantilla")

	// 1. Validar que la plantilla existe
	templateDef, exists := entities.GetTemplateDefinition(templateName)
	if !exists {
		u.log.Error(ctx).
			Str("template_name", templateName).
			Msg("[WhatsApp UseCase] - plantilla no encontrada")
		return "", &errors.ErrTemplateNotFound{TemplateName: templateName}
	}

	// 2. Validar que se proveen todas las variables requeridas
	if err := entities.ValidateTemplateVariables(templateName, variables); err != nil {
		u.log.Error(ctx).Err(err).
			Str("template_name", templateName).
			Msg("[WhatsApp UseCase] - variables faltantes")
		return "", err
	}

	// 3. Validar número de teléfono
	if err := ValidatePhoneNumber(phoneNumber); err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp UseCase] - número de teléfono inválido")
		return "", fmt.Errorf("número de teléfono inválido: %w", err)
	}

	// 4. Obtener configuración de WhatsApp (phone_number_id + access_token) desde el repositorio
	whatsappConfig, err := u.integrationRepo.GetWhatsAppConfig(ctx, businessID)
	if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error obteniendo configuración de WhatsApp")
		return "", fmt.Errorf("error obteniendo configuración de WhatsApp: %w", err)
	}

	// 5. Construir mensaje con botones si aplica
	msg := u.buildTemplateMessage(templateName, phoneNumber, variables, templateDef)

	// 6. Buscar o crear conversación
	conversation, err := u.getOrCreateConversation(ctx, phoneNumber, orderNumber, businessID)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("phone_number", phoneNumber).
			Str("order_number", orderNumber).
			Msg("[WhatsApp UseCase] - error obteniendo/creando conversación")
		return "", err
	}

	// 7. Enviar mensaje
	u.log.Info(ctx).
		Str("conversation_id", conversation.ID).
		Str("template_name", templateName).
		Uint("phone_number_id", whatsappConfig.PhoneNumberID).
		Msg("[WhatsApp UseCase] - enviando mensaje a WhatsApp API")

	messageID, err := u.whatsApp.SendMessage(ctx, whatsappConfig.PhoneNumberID, msg, whatsappConfig.AccessToken)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("template_name", templateName).
			Str("phone_number", phoneNumber).
			Msg("[WhatsApp UseCase] - error enviando mensaje")
		return "", fmt.Errorf("error al enviar mensaje de WhatsApp: %w", err)
	}

	// 8. Registrar en message_log
	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        fmt.Sprintf("Template: %s, Variables: %v", templateName, variables),
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
	}

	if err := u.messageLogRepo.Create(ctx, messageLog); err != nil {
		u.log.Error(ctx).Err(err).
			Str("message_id", messageID).
			Msg("[WhatsApp UseCase] - error registrando mensaje en log")
		// No retornamos error porque el mensaje ya fue enviado
	}

	// 9. Actualizar conversación
	conversation.LastMessageID = messageID
	conversation.LastTemplateID = templateName
	conversation.UpdatedAt = time.Now()

	if err := u.conversationRepo.Update(ctx, conversation); err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - error actualizando conversación")
		// No retornamos error porque el mensaje ya fue enviado
	}

	u.log.Info(ctx).
		Str("message_id", messageID).
		Str("conversation_id", conversation.ID).
		Str("template_name", templateName).
		Msg("[WhatsApp UseCase] - mensaje enviado exitosamente")

	return messageID, nil
}

// SendTemplateWithConversation envía una plantilla usando una conversación existente
func (u *Usecases) SendTemplateWithConversation(
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

	// 1. Validar plantilla y variables
	templateDef, exists := entities.GetTemplateDefinition(templateName)
	if !exists {
		return "", &errors.ErrTemplateNotFound{TemplateName: templateName}
	}

	if err := entities.ValidateTemplateVariables(templateName, variables); err != nil {
		return "", err
	}

	// 2. Obtener conversación existente
	conversation, err := u.conversationRepo.GetByID(ctx, conversationID)
	if err != nil {
		u.log.Error(ctx).Err(err).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp UseCase] - conversación no encontrada")
		return "", err
	}

	// 3. Verificar que la conversación no ha expirado
	if conversation.IsExpired() {
		u.log.Error(ctx).
			Str("conversation_id", conversationID).
			Msg("[WhatsApp UseCase] - conversación expirada")
		return "", &errors.ErrConversationExpired{ConversationID: conversationID}
	}

	// 4. Obtener configuración de WhatsApp desde el repositorio
	whatsappConfig, err := u.integrationRepo.GetWhatsAppConfig(ctx, conversation.BusinessID)
	if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error obteniendo configuración de WhatsApp")
		return "", fmt.Errorf("error obteniendo configuración de WhatsApp: %w", err)
	}

	// 5. Construir y enviar mensaje
	msg := u.buildTemplateMessage(templateName, phoneNumber, variables, templateDef)
	messageID, err := u.whatsApp.SendMessage(ctx, whatsappConfig.PhoneNumberID, msg, whatsappConfig.AccessToken)
	if err != nil {
		return "", err
	}

	// 6. Registrar en log
	messageLog := &entities.MessageLog{
		ConversationID: conversation.ID,
		Direction:      entities.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        fmt.Sprintf("Template: %s", templateName),
		Status:         entities.MessageStatusSent,
		CreatedAt:      time.Now(),
	}
	u.messageLogRepo.Create(ctx, messageLog)

	// 7. Actualizar conversación
	conversation.LastMessageID = messageID
	conversation.LastTemplateID = templateName
	conversation.UpdatedAt = time.Now()
	u.conversationRepo.Update(ctx, conversation)

	return messageID, nil
}

// buildTemplateMessage construye el mensaje de plantilla con todos sus componentes
func (u *Usecases) buildTemplateMessage(
	templateName string,
	phoneNumber string,
	variables map[string]string,
	templateDef entities.TemplateDefinition,
) entities.TemplateMessage {
	// Construir componentes
	components := []entities.TemplateComponent{}

	// Agregar componente body con variables si existen
	if len(templateDef.Variables) > 0 {
		bodyParams := []entities.TemplateParameter{}
		for i := range templateDef.Variables {
			varKey := string(rune('1' + i))
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

	// NOTA: Los botones de tipo "quick_reply" NO se envían en el payload.
	// Estos botones ya están definidos en la plantilla en Meta y se
	// renderizan automáticamente. Solo enviamos parámetros del body/header.

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

// getOrCreateConversation obtiene una conversación existente o crea una nueva
func (u *Usecases) getOrCreateConversation(
	ctx context.Context,
	phoneNumber string,
	orderNumber string,
	businessID uint,
) (*entities.Conversation, error) {
	// Intentar obtener conversación existente
	conversation, err := u.conversationRepo.GetByPhoneAndOrder(ctx, phoneNumber, orderNumber)
	if err == nil {
		// Conversación encontrada
		if conversation.IsActive() {
			return conversation, nil
		}
		// Conversación expirada, crear una nueva
		u.log.Info(ctx).
			Str("conversation_id", conversation.ID).
			Msg("[WhatsApp UseCase] - conversación expirada, creando nueva")
	}

	// Crear nueva conversación
	newConversation := &entities.Conversation{
		PhoneNumber:  phoneNumber,
		OrderNumber:  orderNumber,
		BusinessID:   businessID,
		CurrentState: entities.StateStart,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour), // Ventana de 24h
	}

	if err := u.conversationRepo.Create(ctx, newConversation); err != nil {
		return nil, err
	}

	u.log.Info(ctx).
		Str("conversation_id", newConversation.ID).
		Str("phone_number", phoneNumber).
		Str("order_number", orderNumber).
		Msg("[WhatsApp UseCase] - nueva conversación creada")

	return newConversation, nil
}

package app

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/whatsApp/internal/domain"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
	"gorm.io/gorm"
)

// ISendTemplateMessageUseCase define la interfaz para el caso de uso de envío de plantillas
type ISendTemplateMessageUseCase interface {
	SendTemplate(ctx context.Context, templateName, phoneNumber string, variables map[string]string, orderNumber string, businessID uint) (string, error)
	SendTemplateWithConversation(ctx context.Context, templateName, phoneNumber string, variables map[string]string, conversationID string) (string, error)
}

// SendTemplateMessageUseCase implementa el caso de uso de envío de plantillas
type SendTemplateMessageUseCase struct {
	whatsApp         domain.IWhatsApp
	conversationRepo domain.IConversationRepository
	messageLogRepo   domain.IMessageLogRepository
	db               *gorm.DB
	log              log.ILogger
	config           env.IConfig
	encryptionKey    []byte
}

// NewSendTemplateMessage crea una nueva instancia del usecase
func NewSendTemplateMessage(
	whatsApp domain.IWhatsApp,
	conversationRepo domain.IConversationRepository,
	messageLogRepo domain.IMessageLogRepository,
	db *gorm.DB,
	logger log.ILogger,
	config env.IConfig,
) ISendTemplateMessageUseCase {
	// Obtener encryption key desde env
	encryptionKeyStr := config.Get("ENCRYPTION_KEY")
	var encryptionKey []byte

	// Intentar decodificar como base64
	decoded, err := base64.StdEncoding.DecodeString(encryptionKeyStr)
	if err == nil && len(decoded) == 32 {
		encryptionKey = decoded
	} else {
		encryptionKey = []byte(encryptionKeyStr)
	}

	return &SendTemplateMessageUseCase{
		whatsApp:         whatsApp,
		conversationRepo: conversationRepo,
		messageLogRepo:   messageLogRepo,
		db:               db,
		log:              logger,
		config:           config,
		encryptionKey:    encryptionKey,
	}
}

// SendTemplate envía una plantilla de WhatsApp y crea/actualiza la conversación
func (u *SendTemplateMessageUseCase) SendTemplate(
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
	templateDef, exists := domain.GetTemplateDefinition(templateName)
	if !exists {
		u.log.Error(ctx).
			Str("template_name", templateName).
			Msg("[WhatsApp UseCase] - plantilla no encontrada")
		return "", &domain.ErrTemplateNotFound{TemplateName: templateName}
	}

	// 2. Validar que se proveen todas las variables requeridas
	if err := domain.ValidateTemplateVariables(templateName, variables); err != nil {
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

	// 4. Obtener configuración de WhatsApp (phone_number_id + access_token) desde la base de datos
	whatsappConfig, err := u.getWhatsAppConfig(ctx, businessID)
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
	messageLog := &domain.MessageLog{
		ConversationID: conversation.ID,
		Direction:      domain.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        fmt.Sprintf("Template: %s, Variables: %v", templateName, variables),
		Status:         domain.MessageStatusSent,
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
func (u *SendTemplateMessageUseCase) SendTemplateWithConversation(
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
	templateDef, exists := domain.GetTemplateDefinition(templateName)
	if !exists {
		return "", &domain.ErrTemplateNotFound{TemplateName: templateName}
	}

	if err := domain.ValidateTemplateVariables(templateName, variables); err != nil {
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
		return "", &domain.ErrConversationExpired{ConversationID: conversationID}
	}

	// 4. Obtener configuración de WhatsApp desde la base de datos
	whatsappConfig, err := u.getWhatsAppConfig(ctx, conversation.BusinessID)
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
	messageLog := &domain.MessageLog{
		ConversationID: conversation.ID,
		Direction:      domain.MessageDirectionOutbound,
		MessageID:      messageID,
		TemplateName:   templateName,
		Content:        fmt.Sprintf("Template: %s", templateName),
		Status:         domain.MessageStatusSent,
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
func (u *SendTemplateMessageUseCase) buildTemplateMessage(
	templateName string,
	phoneNumber string,
	variables map[string]string,
	templateDef domain.TemplateDefinition,
) domain.TemplateMessage {
	// Construir componentes
	components := []domain.TemplateComponent{}

	// Agregar componente body con variables si existen
	if len(templateDef.Variables) > 0 {
		bodyParams := []domain.TemplateParameter{}
		for i := range templateDef.Variables {
			varKey := string(rune('1' + i))
			bodyParams = append(bodyParams, domain.TemplateParameter{
				Type: "text",
				Text: variables[varKey],
			})
		}
		components = append(components, domain.TemplateComponent{
			Type:       "body",
			Parameters: bodyParams,
		})
	}

	// NOTA: Los botones de tipo "quick_reply" NO se envían en el payload.
	// Estos botones ya están definidos en la plantilla en Meta y se
	// renderizan automáticamente. Solo enviamos parámetros del body/header.

	return domain.TemplateMessage{
		MessagingProduct: "whatsapp",
		RecipientType:    "individual",
		To:               phoneNumber,
		Type:             "template",
		Template: domain.TemplateData{
			Name:       templateName,
			Language:   domain.TemplateLanguage{Code: templateDef.Language},
			Components: components,
		},
	}
}

// getOrCreateConversation obtiene una conversación existente o crea una nueva
func (u *SendTemplateMessageUseCase) getOrCreateConversation(
	ctx context.Context,
	phoneNumber string,
	orderNumber string,
	businessID uint,
) (*domain.Conversation, error) {
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
	newConversation := &domain.Conversation{
		PhoneNumber:  phoneNumber,
		OrderNumber:  orderNumber,
		BusinessID:   businessID,
		CurrentState: domain.StateStart,
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

// WhatsAppConfig contiene la configuración de WhatsApp obtenida desde la base de datos
type WhatsAppConfig struct {
	PhoneNumberID uint
	AccessToken   string
	IntegrationID uint
}

// getWhatsAppConfig obtiene phone_number_id y access_token desde la base de datos
func (u *SendTemplateMessageUseCase) getWhatsAppConfig(ctx context.Context, businessID uint) (*WhatsAppConfig, error) {
	// Estructura para almacenar la integración
	type Integration struct {
		ID          uint            `json:"id"`
		Config      json.RawMessage `json:"config"`
		Credentials json.RawMessage `json:"credentials"`
	}

	var integration Integration

	// Primero intentar obtener la integración del business específico
	err := u.db.WithContext(ctx).
		Table("integrations").
		Select("id, config, credentials").
		Where("integration_type_id = ?", 2).
		Where("business_id = ?", businessID).
		First(&integration).Error

	if err == gorm.ErrRecordNotFound {
		// Si no existe, usar la integración global (business_id IS NULL)
		u.log.Info(ctx).
			Uint("business_id", businessID).
			Msg("[WhatsApp UseCase] - no se encontró integración específica, usando global")

		err = u.db.WithContext(ctx).
			Table("integrations").
			Select("id, config, credentials").
			Where("integration_type_id = ?", 2).
			Where("business_id IS NULL").
			First(&integration).Error

		if err != nil {
			u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - no se encontró integración de WhatsApp")
			return nil, fmt.Errorf("no se encontró integración de WhatsApp")
		}
	} else if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error consultando integración")
		return nil, fmt.Errorf("error consultando integración: %w", err)
	}

	// 1. Parsear el config JSON para obtener phone_number_id
	var config map[string]interface{}
	if err := json.Unmarshal(integration.Config, &config); err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error parseando config")
		return nil, fmt.Errorf("error parseando config: %w", err)
	}

	// Extraer phone_number_id
	phoneNumberIDValue, exists := config["phone_number_id"]
	if !exists {
		u.log.Error(ctx).Msg("[WhatsApp UseCase] - phone_number_id no encontrado en config")
		return nil, fmt.Errorf("phone_number_id no encontrado en configuración")
	}

	phoneNumberIDStr, ok := phoneNumberIDValue.(string)
	if !ok {
		u.log.Error(ctx).Msg("[WhatsApp UseCase] - phone_number_id no es string")
		return nil, fmt.Errorf("phone_number_id debe ser string")
	}

	phoneNumberID, err := strconv.ParseUint(phoneNumberIDStr, 10, 64)
	if err != nil {
		u.log.Error(ctx).Err(err).Str("phone_number_id", phoneNumberIDStr).Msg("[WhatsApp UseCase] - error parseando phone_number_id")
		return nil, fmt.Errorf("error parseando phone_number_id: %w", err)
	}

	// 2. Parsear credentials para obtener el wrapper encriptado
	var credentialsWrapper map[string]interface{}
	if err := json.Unmarshal(integration.Credentials, &credentialsWrapper); err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error parseando credentials wrapper")
		return nil, fmt.Errorf("error parseando credentials: %w", err)
	}

	// Extraer el valor encriptado del wrapper
	encryptedValue, exists := credentialsWrapper["encrypted"]
	if !exists {
		u.log.Error(ctx).Msg("[WhatsApp UseCase] - credentials no tienen campo 'encrypted'")
		return nil, fmt.Errorf("credentials inválidas: falta campo 'encrypted'")
	}

	encryptedStr, ok := encryptedValue.(string)
	if !ok {
		u.log.Error(ctx).Msg("[WhatsApp UseCase] - encrypted no es string")
		return nil, fmt.Errorf("credentials inválidas: 'encrypted' debe ser string")
	}

	// Decodificar base64
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error decodificando credentials desde base64")
		return nil, fmt.Errorf("error decodificando credentials: %w", err)
	}

	// Desencriptar usando AES-GCM
	decryptedCredentials, err := u.decryptCredentials(ctx, encryptedBytes)
	if err != nil {
		u.log.Error(ctx).Err(err).Msg("[WhatsApp UseCase] - error desencriptando credentials")
		return nil, fmt.Errorf("error desencriptando credentials: %w", err)
	}

	// Extraer access_token de las credenciales desencriptadas
	accessTokenValue, exists := decryptedCredentials["access_token"]
	if !exists {
		u.log.Error(ctx).Msg("[WhatsApp UseCase] - access_token no encontrado en credentials")
		return nil, fmt.Errorf("access_token no encontrado en credentials")
	}

	accessToken, ok := accessTokenValue.(string)
	if !ok {
		u.log.Error(ctx).Msg("[WhatsApp UseCase] - access_token no es string")
		return nil, fmt.Errorf("access_token debe ser string")
	}

	u.log.Info(ctx).
		Uint("integration_id", integration.ID).
		Uint("phone_number_id", uint(phoneNumberID)).
		Msg("[WhatsApp UseCase] - Configuración de WhatsApp obtenida desde DB")

	return &WhatsAppConfig{
		PhoneNumberID: uint(phoneNumberID),
		AccessToken:   accessToken,
		IntegrationID: integration.ID,
	}, nil
}

// decryptCredentials desencripta credenciales usando AES-256-GCM
func (u *SendTemplateMessageUseCase) decryptCredentials(ctx context.Context, ciphertext []byte) (map[string]interface{}, error) {
	block, err := aes.NewCipher(u.encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("error al crear cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("error al crear GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext demasiado corto")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error al desencriptar: %w", err)
	}

	// Convertir JSON a mapa
	var credentials map[string]interface{}
	if err := json.Unmarshal(plaintext, &credentials); err != nil {
		return nil, fmt.Errorf("error al deserializar credenciales: %w", err)
	}

	return credentials, nil
}

// getPhoneNumberID mantener por compatibilidad - ahora usa getWhatsAppConfig
func (u *SendTemplateMessageUseCase) getPhoneNumberID(ctx context.Context, businessID uint) (uint, error) {
	config, err := u.getWhatsAppConfig(ctx, businessID)
	if err != nil {
		return 0, err
	}
	return config.PhoneNumberID, nil
}

package usecasemessaging

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/mocks"
)

// ---------------------------------------------------------------------------
// Helpers para construir payloads de webhook
// ---------------------------------------------------------------------------

func buildWebhookWithTextMessage(phoneNumber, text, messageID string) dtos.WebhookPayloadDTO {
	return dtos.WebhookPayloadDTO{
		Object: "whatsapp_business_account",
		Entry: []dtos.WebhookEntryDTO{
			{
				ID: "entry-001",
				Changes: []dtos.WebhookChangeDTO{
					{
						Field: "messages",
						Value: dtos.WebhookValueDTO{
							Messages: []dtos.WebhookMessageDTO{
								{
									ID:   messageID,
									From: phoneNumber,
									Type: "text",
									Text: &dtos.TextContentDTO{Body: text},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildWebhookWithButtonMessage(phoneNumber, buttonText, messageID string) dtos.WebhookPayloadDTO {
	return dtos.WebhookPayloadDTO{
		Object: "whatsapp_business_account",
		Entry: []dtos.WebhookEntryDTO{
			{
				Changes: []dtos.WebhookChangeDTO{
					{
						Field: "messages",
						Value: dtos.WebhookValueDTO{
							Messages: []dtos.WebhookMessageDTO{
								{
									ID:     messageID,
									From:   phoneNumber,
									Type:   "button",
									Button: &dtos.ButtonResponseDTO{Text: buttonText},
								},
							},
						},
					},
				},
			},
		},
	}
}

func buildWebhookWithStatus(messageID, status, timestamp string) dtos.WebhookPayloadDTO {
	return dtos.WebhookPayloadDTO{
		Entry: []dtos.WebhookEntryDTO{
			{
				Changes: []dtos.WebhookChangeDTO{
					{
						Field: "messages",
						Value: dtos.WebhookValueDTO{
							Statuses: []dtos.WebhookStatusDTO{
								{
									ID:        messageID,
									Status:    status,
									Timestamp: timestamp,
								},
							},
						},
					},
				},
			},
		},
	}
}

// ---------------------------------------------------------------------------
// HandleIncomingMessage
// ---------------------------------------------------------------------------

func TestHandleIncomingMessage_SinConversacionActiva_IgnoraMensaje(t *testing.T) {
	// Si no hay conversación activa, el mensaje se ignora (no es error)
	convRepoMock := &mocks.ConversationRepositoryMock{
		GetActiveByPhoneFn: func(_ context.Context, _ string) (*entities.Conversation, error) {
			return nil, errors.New("sin conversación activa")
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	payload := buildWebhookWithTextMessage("+573001234567", "Hola", "msg-001")
	err := uc.HandleIncomingMessage(context.Background(), payload)

	if err != nil {
		t.Errorf("HandleIncomingMessage() no debería retornar error cuando no hay conversación activa, obtuvo: %v", err)
	}
}

func TestHandleIncomingMessage_ConversacionExpirada_RetornaError(t *testing.T) {
	expiradaConv := &entities.Conversation{
		ID:          "conv-exp-001",
		PhoneNumber: "+573001234567",
		ExpiresAt:   time.Now().Add(-2 * time.Hour), // expirada
	}

	convRepoMock := &mocks.ConversationRepositoryMock{
		GetActiveByPhoneFn: func(_ context.Context, _ string) (*entities.Conversation, error) {
			return expiradaConv, nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	// Los errores por mensaje individual no bloquean el handler (continue),
	// así que HandleIncomingMessage retorna nil aunque haya error interno.
	payload := buildWebhookWithTextMessage("+573001234567", "Confirmar pedido", "msg-002")
	err := uc.HandleIncomingMessage(context.Background(), payload)

	// El diseño del use case usa "continue" en el loop de mensajes,
	// por lo que el error se registra pero no se propaga.
	if err != nil {
		t.Errorf("HandleIncomingMessage() retornó error inesperado: %v", err)
	}
}

func TestHandleIncomingMessage_CambioNoMessages_Ignorado(t *testing.T) {
	// Un change con field distinto a "messages" no debe procesarse
	payload := dtos.WebhookPayloadDTO{
		Entry: []dtos.WebhookEntryDTO{
			{
				Changes: []dtos.WebhookChangeDTO{
					{Field: "message_template_status_update", Value: dtos.WebhookValueDTO{}},
				},
			},
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	err := uc.HandleIncomingMessage(context.Background(), payload)
	if err != nil {
		t.Errorf("HandleIncomingMessage() no debería retornar error con change.Field != messages: %v", err)
	}
}

func TestHandleIncomingMessage_PayloadVacio_NoError(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	err := uc.HandleIncomingMessage(context.Background(), dtos.WebhookPayloadDTO{})
	if err != nil {
		t.Errorf("HandleIncomingMessage() no debería retornar error con payload vacío: %v", err)
	}
}

// ---------------------------------------------------------------------------
// processIncomingMessage — flujo completo con transición de estado
// ---------------------------------------------------------------------------

func TestHandleIncomingMessage_FlujoBotonesConfirmarPedido(t *testing.T) {
	// Conversación activa en estado AWAITING_CONFIRMATION
	conv := &entities.Conversation{
		ID:           "conv-flujo-001",
		PhoneNumber:  "+573001234567",
		OrderNumber:  "ORD-555",
		BusinessID:   1,
		CurrentState: entities.StateAwaitingConfirmation,
		Metadata:     make(map[string]interface{}),
		ExpiresAt:    time.Now().Add(12 * time.Hour),
	}

	var updateCalled bool
	convRepoMock := &mocks.ConversationRepositoryMock{
		GetActiveByPhoneFn: func(_ context.Context, _ string) (*entities.Conversation, error) {
			return conv, nil
		},
		GetByIDFn: func(_ context.Context, _ string) (*entities.Conversation, error) {
			// Requerido internamente por SendTemplateWithConversation
			return conv, nil
		},
		UpdateFn: func(_ context.Context, _ *entities.Conversation) error {
			updateCalled = true
			return nil
		},
	}

	integRepo := &mocks.IntegrationRepositoryMock{
		GetWhatsAppConfigFn: func(_ context.Context, _ uint) (*ports.WhatsAppConfig, error) {
			return &ports.WhatsAppConfig{PhoneNumberID: 111, AccessToken: "tok"}, nil
		},
	}

	var publishConfirmedCalled bool
	publisherMock := &mocks.EventPublisherMock{
		PublishOrderConfirmedFn: func(_ context.Context, _, _ string, _ uint) error {
			publishConfirmedCalled = true
			return nil
		},
	}

	waClient := &mocks.WhatsAppMock{} // retorna "wamid.mock123" por defecto

	uc := newUsecasesForTest(
		waClient,
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		integRepo,
		publisherMock,
		&mocks.ConfigMock{},
	)

	// El usuario presiona el botón "Confirmar pedido" (tipo button)
	payload := buildWebhookWithButtonMessage("+573001234567", "Confirmar pedido", "msg-btn-001")
	_ = uc.HandleIncomingMessage(context.Background(), payload)

	if !updateCalled {
		t.Error("se esperaba que Update() fuera llamado para actualizar el estado de la conversación")
	}
	if !publishConfirmedCalled {
		t.Error("se esperaba que PublishOrderConfirmed() fuera llamado al confirmar el pedido")
	}
}

// ---------------------------------------------------------------------------
// HandleMessageStatus
// ---------------------------------------------------------------------------

func TestHandleMessageStatus_Entregado(t *testing.T) {
	var capturedStatus entities.MessageStatus

	msgLogMock := &mocks.MessageLogRepositoryMock{
		UpdateStatusFn: func(_ context.Context, _ string, status entities.MessageStatus, _ map[string]time.Time) error {
			capturedStatus = status
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		msgLogMock,
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	payload := buildWebhookWithStatus("wamid.abc123", "delivered", "1700000000")
	err := uc.HandleMessageStatus(context.Background(), payload)

	if err != nil {
		t.Fatalf("HandleMessageStatus() error inesperado: %v", err)
	}
	if capturedStatus != entities.MessageStatusDelivered {
		t.Errorf("capturedStatus = %q, quería %q", capturedStatus, entities.MessageStatusDelivered)
	}
}

func TestHandleMessageStatus_Leido(t *testing.T) {
	var capturedStatus entities.MessageStatus

	msgLogMock := &mocks.MessageLogRepositoryMock{
		UpdateStatusFn: func(_ context.Context, _ string, status entities.MessageStatus, _ map[string]time.Time) error {
			capturedStatus = status
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		msgLogMock,
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	payload := buildWebhookWithStatus("wamid.abc123", "read", "1700000100")
	err := uc.HandleMessageStatus(context.Background(), payload)

	if err != nil {
		t.Fatalf("HandleMessageStatus() error inesperado: %v", err)
	}
	if capturedStatus != entities.MessageStatusRead {
		t.Errorf("capturedStatus = %q, quería %q", capturedStatus, entities.MessageStatusRead)
	}
}

func TestHandleMessageStatus_Enviado(t *testing.T) {
	var capturedStatus entities.MessageStatus

	msgLogMock := &mocks.MessageLogRepositoryMock{
		UpdateStatusFn: func(_ context.Context, _ string, status entities.MessageStatus, _ map[string]time.Time) error {
			capturedStatus = status
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		msgLogMock,
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	payload := buildWebhookWithStatus("wamid.sent001", "sent", "1700000200")
	err := uc.HandleMessageStatus(context.Background(), payload)

	if err != nil {
		t.Fatalf("HandleMessageStatus() error inesperado: %v", err)
	}
	if capturedStatus != entities.MessageStatusSent {
		t.Errorf("capturedStatus = %q, quería %q", capturedStatus, entities.MessageStatusSent)
	}
}

func TestHandleMessageStatus_Fallido(t *testing.T) {
	var capturedStatus entities.MessageStatus

	msgLogMock := &mocks.MessageLogRepositoryMock{
		UpdateStatusFn: func(_ context.Context, _ string, status entities.MessageStatus, _ map[string]time.Time) error {
			capturedStatus = status
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		msgLogMock,
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	payload := buildWebhookWithStatus("wamid.fail001", "failed", "1700000300")
	err := uc.HandleMessageStatus(context.Background(), payload)

	if err != nil {
		t.Fatalf("HandleMessageStatus() error inesperado: %v", err)
	}
	if capturedStatus != entities.MessageStatusFailed {
		t.Errorf("capturedStatus = %q, quería %q", capturedStatus, entities.MessageStatusFailed)
	}
}

func TestHandleMessageStatus_EstadoDesconocido_NoError(t *testing.T) {
	var updateCalled bool

	msgLogMock := &mocks.MessageLogRepositoryMock{
		UpdateStatusFn: func(_ context.Context, _ string, _ entities.MessageStatus, _ map[string]time.Time) error {
			updateCalled = true
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		msgLogMock,
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	payload := buildWebhookWithStatus("wamid.unknown", "pending_magic", "1700000400")
	err := uc.HandleMessageStatus(context.Background(), payload)

	if err != nil {
		t.Errorf("HandleMessageStatus() no debería retornar error con estado desconocido: %v", err)
	}
	// Con estado desconocido, no se llama a UpdateStatus
	if updateCalled {
		t.Error("UpdateStatus no debería ser llamado con un estado desconocido")
	}
}

func TestHandleMessageStatus_ErrorRepositorio_ContinuaConOtros(t *testing.T) {
	// Si UpdateStatus falla para un mensaje, el loop continúa con los demás
	var callCount int

	msgLogMock := &mocks.MessageLogRepositoryMock{
		UpdateStatusFn: func(_ context.Context, messageID string, _ entities.MessageStatus, _ map[string]time.Time) error {
			callCount++
			if messageID == "wamid.falla" {
				return errors.New("error de base de datos")
			}
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		msgLogMock,
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	// Payload con dos statuses: uno fallará, el otro no
	payload := dtos.WebhookPayloadDTO{
		Entry: []dtos.WebhookEntryDTO{
			{
				Changes: []dtos.WebhookChangeDTO{
					{
						Field: "messages",
						Value: dtos.WebhookValueDTO{
							Statuses: []dtos.WebhookStatusDTO{
								{ID: "wamid.falla", Status: "delivered", Timestamp: "1700000000"},
								{ID: "wamid.ok", Status: "read", Timestamp: "1700000001"},
							},
						},
					},
				},
			},
		},
	}

	err := uc.HandleMessageStatus(context.Background(), payload)

	// El handler siempre retorna nil (usa "continue" en el loop)
	if err != nil {
		t.Errorf("HandleMessageStatus() retornó error, debería retornar nil siempre: %v", err)
	}
	// Ambos statuses fueron procesados (aunque uno falló)
	if callCount != 2 {
		t.Errorf("UpdateStatus fue llamado %d veces, quería 2", callCount)
	}
}

func TestHandleMessageStatus_PayloadVacio_NoError(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	err := uc.HandleMessageStatus(context.Background(), dtos.WebhookPayloadDTO{})
	if err != nil {
		t.Errorf("HandleMessageStatus() no debería retornar error con payload vacío: %v", err)
	}
}

// ---------------------------------------------------------------------------
// publishBusinessEvent
// ---------------------------------------------------------------------------

func TestPublishBusinessEvent_Confirmado(t *testing.T) {
	var called bool

	publisherMock := &mocks.EventPublisherMock{
		PublishOrderConfirmedFn: func(_ context.Context, orderNumber, phoneNumber string, businessID uint) error {
			called = true
			if orderNumber != "ORD-111" {
				t.Errorf("orderNumber = %q, quería %q", orderNumber, "ORD-111")
			}
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		publisherMock,
		&mocks.ConfigMock{},
	)

	conv := &entities.Conversation{
		OrderNumber: "ORD-111",
		PhoneNumber: "+573001234567",
		BusinessID:  1,
		Metadata:    make(map[string]interface{}),
	}

	err := uc.publishBusinessEvent(context.Background(), "confirmed", conv)
	if err != nil {
		t.Errorf("publishBusinessEvent(confirmed) error inesperado: %v", err)
	}
	if !called {
		t.Error("PublishOrderConfirmed no fue llamado")
	}
}

func TestPublishBusinessEvent_Cancelado_ConRazon(t *testing.T) {
	var capturedReason string

	publisherMock := &mocks.EventPublisherMock{
		PublishOrderCancelledFn: func(_ context.Context, _, reason, _ string, _ uint) error {
			capturedReason = reason
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		publisherMock,
		&mocks.ConfigMock{},
	)

	conv := &entities.Conversation{
		OrderNumber: "ORD-222",
		PhoneNumber: "+573001234567",
		BusinessID:  1,
		Metadata: map[string]interface{}{
			"cancellation_reason": "El cliente cambió de opinión",
		},
	}

	err := uc.publishBusinessEvent(context.Background(), "cancelled", conv)
	if err != nil {
		t.Errorf("publishBusinessEvent(cancelled) error inesperado: %v", err)
	}
	if capturedReason != "El cliente cambió de opinión" {
		t.Errorf("capturedReason = %q, quería %q", capturedReason, "El cliente cambió de opinión")
	}
}

func TestPublishBusinessEvent_Novedad(t *testing.T) {
	var capturedNoveltyType string

	publisherMock := &mocks.EventPublisherMock{
		PublishNoveltyRequestedFn: func(_ context.Context, _, noveltyType, _ string, _ uint) error {
			capturedNoveltyType = noveltyType
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		publisherMock,
		&mocks.ConfigMock{},
	)

	conv := &entities.Conversation{
		OrderNumber: "ORD-333",
		PhoneNumber: "+573001234567",
		BusinessID:  1,
		Metadata: map[string]interface{}{
			"novelty_type": "change_address",
		},
	}

	err := uc.publishBusinessEvent(context.Background(), "novelty", conv)
	if err != nil {
		t.Errorf("publishBusinessEvent(novelty) error inesperado: %v", err)
	}
	if capturedNoveltyType != "change_address" {
		t.Errorf("capturedNoveltyType = %q, quería %q", capturedNoveltyType, "change_address")
	}
}

func TestPublishBusinessEvent_Handoff(t *testing.T) {
	var called bool

	publisherMock := &mocks.EventPublisherMock{
		PublishHandoffRequestedFn: func(_ context.Context, _, _ string, _ uint, convID string) error {
			called = true
			if convID != "conv-handoff-001" {
				t.Errorf("conversationID = %q, quería %q", convID, "conv-handoff-001")
			}
			return nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		publisherMock,
		&mocks.ConfigMock{},
	)

	conv := &entities.Conversation{
		ID:          "conv-handoff-001",
		OrderNumber: "ORD-444",
		PhoneNumber: "+573001234567",
		BusinessID:  1,
		Metadata:    make(map[string]interface{}),
	}

	err := uc.publishBusinessEvent(context.Background(), "handoff", conv)
	if err != nil {
		t.Errorf("publishBusinessEvent(handoff) error inesperado: %v", err)
	}
	if !called {
		t.Error("PublishHandoffRequested no fue llamado")
	}
}

func TestPublishBusinessEvent_TipoDesconocido_NoError(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := &entities.Conversation{Metadata: make(map[string]interface{})}
	err := uc.publishBusinessEvent(context.Background(), "tipo_desconocido", conv)

	if err != nil {
		t.Errorf("publishBusinessEvent(tipo_desconocido) no debería retornar error: %v", err)
	}
}

// ---------------------------------------------------------------------------
// WebhookMessageDTO.GetMessageText
// ---------------------------------------------------------------------------

func TestGetMessageText_TextoSimple(t *testing.T) {
	msg := &dtos.WebhookMessageDTO{
		Type: "text",
		Text: &dtos.TextContentDTO{Body: "Hola mundo"},
	}
	got := msg.GetMessageText()
	if got != "Hola mundo" {
		t.Errorf("GetMessageText() = %q, quería %q", got, "Hola mundo")
	}
}

func TestGetMessageText_BotonQuickReply(t *testing.T) {
	msg := &dtos.WebhookMessageDTO{
		Type:   "button",
		Button: &dtos.ButtonResponseDTO{Text: "Confirmar pedido"},
	}
	got := msg.GetMessageText()
	if got != "Confirmar pedido" {
		t.Errorf("GetMessageText() = %q, quería %q", got, "Confirmar pedido")
	}
}

func TestGetMessageText_Interactivo_ButtonReply(t *testing.T) {
	msg := &dtos.WebhookMessageDTO{
		Type: "interactive",
		Interactive: &dtos.InteractiveResponseDTO{
			Type:        "button_reply",
			ButtonReply: &dtos.ButtonReplyDataDTO{Title: "Cancelar pedido"},
		},
	}
	got := msg.GetMessageText()
	if got != "Cancelar pedido" {
		t.Errorf("GetMessageText() = %q, quería %q", got, "Cancelar pedido")
	}
}

func TestGetMessageText_Interactivo_ListReply(t *testing.T) {
	msg := &dtos.WebhookMessageDTO{
		Type: "interactive",
		Interactive: &dtos.InteractiveResponseDTO{
			Type:      "list_reply",
			ListReply: &dtos.ListReplyDataDTO{Title: "Opción 1"},
		},
	}
	got := msg.GetMessageText()
	if got != "Opción 1" {
		t.Errorf("GetMessageText() = %q, quería %q", got, "Opción 1")
	}
}

func TestGetMessageText_TipoDesconocido_Vacio(t *testing.T) {
	msg := &dtos.WebhookMessageDTO{Type: "image"}
	got := msg.GetMessageText()
	if got != "" {
		t.Errorf("GetMessageText() = %q, quería string vacío para tipo desconocido", got)
	}
}

func TestIsButtonResponse(t *testing.T) {
	tests := []struct {
		name string
		msg  dtos.WebhookMessageDTO
		want bool
	}{
		{
			name: "tipo button es respuesta de botón",
			msg:  dtos.WebhookMessageDTO{Type: "button", Button: &dtos.ButtonResponseDTO{}},
			want: true,
		},
		{
			name: "tipo interactive con ButtonReply es respuesta de botón",
			msg: dtos.WebhookMessageDTO{
				Type: "interactive",
				Interactive: &dtos.InteractiveResponseDTO{
					ButtonReply: &dtos.ButtonReplyDataDTO{},
				},
			},
			want: true,
		},
		{
			name: "tipo text NO es respuesta de botón",
			msg:  dtos.WebhookMessageDTO{Type: "text", Text: &dtos.TextContentDTO{Body: "texto"}},
			want: false,
		},
		{
			name: "tipo interactive sin ButtonReply NO es respuesta de botón",
			msg: dtos.WebhookMessageDTO{
				Type:        "interactive",
				Interactive: &dtos.InteractiveResponseDTO{ListReply: &dtos.ListReplyDataDTO{}},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.msg.IsButtonResponse()
			if got != tt.want {
				t.Errorf("IsButtonResponse() = %v, quería %v", got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Errores del dominio
// ---------------------------------------------------------------------------

func TestErrTemplateNotFound_Error(t *testing.T) {
	err := &domainerrors.ErrTemplateNotFound{TemplateName: "mi_plantilla"}
	got := err.Error()
	if !searchSubstring(got, "mi_plantilla") {
		t.Errorf("ErrTemplateNotFound.Error() = %q, no contiene el nombre de la plantilla", got)
	}
}

func TestErrConversationExpired_Error(t *testing.T) {
	err := &domainerrors.ErrConversationExpired{ConversationID: "conv-exp-999"}
	got := err.Error()
	if !searchSubstring(got, "conv-exp-999") {
		t.Errorf("ErrConversationExpired.Error() = %q, no contiene el ID de conversación", got)
	}
}

func TestErrInvalidStateTransition_Error(t *testing.T) {
	err := &domainerrors.ErrInvalidStateTransition{
		CurrentState: "AWAITING_CONFIRMATION",
		UserResponse: "respuesta_rara",
	}
	got := err.Error()
	if !searchSubstring(got, "AWAITING_CONFIRMATION") {
		t.Errorf("ErrInvalidStateTransition.Error() = %q, no contiene el estado actual", got)
	}
}

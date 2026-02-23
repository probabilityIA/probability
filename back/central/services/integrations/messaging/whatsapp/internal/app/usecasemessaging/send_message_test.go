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
// SendMessage (legacy — usa variables de entorno)
// ---------------------------------------------------------------------------

func TestSendMessage_ExitoConVariablesDeEntorno(t *testing.T) {
	// Arrange
	const expectedMessageID = "wamid.success001"

	waMock := &mocks.WhatsAppMock{
		SendMessageFn: func(_ context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error) {
			if msg.To != "+573001234567" {
				t.Errorf("To = %q, quería %q", msg.To, "+573001234567")
			}
			if msg.Template.Name != "order_status_9" {
				t.Errorf("Template.Name = %q, quería %q", msg.Template.Name, "order_status_9")
			}
			return expectedMessageID, nil
		},
	}

	cfg := &mocks.ConfigMock{
		Values: map[string]string{
			"WHATSAPP_PHONE_NUMBER_ID": "123456",
			"WHATSAPP_TOKEN":           "token-test",
		},
	}

	uc := newUsecasesForTest(
		waMock,
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		cfg,
	)

	// Act
	messageID, err := uc.SendMessage(context.Background(), dtos.SendMessageRequest{
		PhoneNumber: "+573001234567",
		OrderNumber: "ORD-001",
	})

	// Assert
	if err != nil {
		t.Fatalf("SendMessage() error inesperado: %v", err)
	}
	if messageID != expectedMessageID {
		t.Errorf("SendMessage() messageID = %q, quería %q", messageID, expectedMessageID)
	}
}

func TestSendMessage_ErrorNumeroTelefonoInvalido(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, err := uc.SendMessage(context.Background(), dtos.SendMessageRequest{
		PhoneNumber: "numero_invalido_sin_codigo_pais",
		OrderNumber: "ORD-001",
	})

	if err == nil {
		t.Fatal("SendMessage() esperaba error con número inválido, no obtuvo ninguno")
	}
}

func TestSendMessage_ErrorPhoneNumberIDNoConfigurado(t *testing.T) {
	cfg := &mocks.ConfigMock{
		Values: map[string]string{
			"WHATSAPP_PHONE_NUMBER_ID": "", // vacío
			"WHATSAPP_TOKEN":           "token-test",
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		cfg,
	)

	_, err := uc.SendMessage(context.Background(), dtos.SendMessageRequest{
		PhoneNumber: "+573001234567",
		OrderNumber: "ORD-001",
	})

	if err == nil {
		t.Fatal("SendMessage() esperaba error cuando WHATSAPP_PHONE_NUMBER_ID está vacío")
	}
}

func TestSendMessage_ErrorTokenNoConfigurado(t *testing.T) {
	cfg := &mocks.ConfigMock{
		Values: map[string]string{
			"WHATSAPP_PHONE_NUMBER_ID": "123456",
			"WHATSAPP_TOKEN":           "", // vacío
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		cfg,
	)

	_, err := uc.SendMessage(context.Background(), dtos.SendMessageRequest{
		PhoneNumber: "+573001234567",
		OrderNumber: "ORD-001",
	})

	if err == nil {
		t.Fatal("SendMessage() esperaba error cuando WHATSAPP_TOKEN está vacío")
	}
}

func TestSendMessage_ErrorPhoneNumberIDInvalido(t *testing.T) {
	cfg := &mocks.ConfigMock{
		Values: map[string]string{
			"WHATSAPP_PHONE_NUMBER_ID": "no_es_un_numero",
			"WHATSAPP_TOKEN":           "token-test",
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		cfg,
	)

	_, err := uc.SendMessage(context.Background(), dtos.SendMessageRequest{
		PhoneNumber: "+573001234567",
		OrderNumber: "ORD-001",
	})

	if err == nil {
		t.Fatal("SendMessage() esperaba error con WHATSAPP_PHONE_NUMBER_ID no numérico")
	}
}

func TestSendMessage_ErrorDelClienteWhatsApp(t *testing.T) {
	expectedErr := errors.New("fallo de red al enviar mensaje")

	waMock := &mocks.WhatsAppMock{
		SendMessageFn: func(_ context.Context, _ uint, _ entities.TemplateMessage, _ string) (string, error) {
			return "", expectedErr
		},
	}

	cfg := &mocks.ConfigMock{
		Values: map[string]string{
			"WHATSAPP_PHONE_NUMBER_ID": "123456",
			"WHATSAPP_TOKEN":           "token-test",
		},
	}

	uc := newUsecasesForTest(
		waMock,
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		cfg,
	)

	_, err := uc.SendMessage(context.Background(), dtos.SendMessageRequest{
		PhoneNumber: "+573001234567",
		OrderNumber: "ORD-001",
	})

	if err == nil {
		t.Fatal("SendMessage() esperaba error cuando el cliente WhatsApp falla")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("SendMessage() error = %v, quería wrapping de %v", err, expectedErr)
	}
}

// ---------------------------------------------------------------------------
// SendTemplate
// ---------------------------------------------------------------------------

func TestSendTemplate_Exito(t *testing.T) {
	const expectedMsgID = "wamid.template001"

	waMock := &mocks.WhatsAppMock{
		SendMessageFn: func(_ context.Context, phoneNumberID uint, msg entities.TemplateMessage, _ string) (string, error) {
			if msg.Template.Name != "pedido_confirmado" {
				t.Errorf("Template.Name = %q, quería %q", msg.Template.Name, "pedido_confirmado")
			}
			return expectedMsgID, nil
		},
	}

	convRepoMock := &mocks.ConversationRepositoryMock{
		GetByPhoneAndOrderFn: func(_ context.Context, _, _ string) (*entities.Conversation, error) {
			// Sin conversación previa
			return nil, errors.New("not found")
		},
		CreateFn: func(_ context.Context, conv *entities.Conversation) error {
			conv.ID = "conv-new-001"
			return nil
		},
		UpdateFn: func(_ context.Context, _ *entities.Conversation) error {
			return nil
		},
	}

	integRepoMock := &mocks.IntegrationRepositoryMock{
		GetWhatsAppConfigFn: func(_ context.Context, _ uint) (*ports.WhatsAppConfig, error) {
			return &ports.WhatsAppConfig{
				PhoneNumberID: 123456,
				AccessToken:   "access-token",
			}, nil
		},
	}

	uc := newUsecasesForTest(
		waMock,
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		integRepoMock,
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	messageID, err := uc.SendTemplate(
		context.Background(),
		"pedido_confirmado",
		"+573001234567",
		map[string]string{"1": "ORD-001"},
		"ORD-001",
		1,
	)

	if err != nil {
		t.Fatalf("SendTemplate() error inesperado: %v", err)
	}
	if messageID != expectedMsgID {
		t.Errorf("SendTemplate() messageID = %q, quería %q", messageID, expectedMsgID)
	}
}

func TestSendTemplate_ErrorPlantillaNoExiste(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, err := uc.SendTemplate(
		context.Background(),
		"plantilla_que_no_existe",
		"+573001234567",
		map[string]string{},
		"ORD-001",
		1,
	)

	if err == nil {
		t.Fatal("SendTemplate() esperaba error con plantilla inexistente")
	}

	var errTmpl *domainerrors.ErrTemplateNotFound
	if !errors.As(err, &errTmpl) {
		t.Errorf("SendTemplate() error = %T, quería *errors.ErrTemplateNotFound", err)
	}
}

func TestSendTemplate_ErrorVariableFaltante(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	// pedido_confirmado requiere variable "1" (numero_pedido)
	_, err := uc.SendTemplate(
		context.Background(),
		"pedido_confirmado",
		"+573001234567",
		map[string]string{}, // faltan variables
		"ORD-001",
		1,
	)

	if err == nil {
		t.Fatal("SendTemplate() esperaba error con variables faltantes")
	}

	var errMissing *domainerrors.ErrMissingVariable
	if !errors.As(err, &errMissing) {
		t.Errorf("SendTemplate() error = %T, quería *errors.ErrMissingVariable", err)
	}
}

func TestSendTemplate_ErrorNumeroTelefonoInvalido(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, err := uc.SendTemplate(
		context.Background(),
		"pedido_confirmado",
		"numero_invalido",
		map[string]string{"1": "ORD-001"},
		"ORD-001",
		1,
	)

	if err == nil {
		t.Fatal("SendTemplate() esperaba error con número de teléfono inválido")
	}
}

func TestSendTemplate_ErrorObtenendoConfigWhatsApp(t *testing.T) {
	expectedErr := errors.New("integración no encontrada")

	integRepoMock := &mocks.IntegrationRepositoryMock{
		GetWhatsAppConfigFn: func(_ context.Context, _ uint) (*ports.WhatsAppConfig, error) {
			return nil, expectedErr
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		integRepoMock,
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, err := uc.SendTemplate(
		context.Background(),
		"pedido_confirmado",
		"+573001234567",
		map[string]string{"1": "ORD-001"},
		"ORD-001",
		1,
	)

	if err == nil {
		t.Fatal("SendTemplate() esperaba error cuando falla la obtención de config de WhatsApp")
	}
}

func TestSendTemplate_ErrorEnvioClienteWhatsApp(t *testing.T) {
	expectedErr := errors.New("API de WhatsApp rechazó el mensaje")

	waMock := &mocks.WhatsAppMock{
		SendMessageFn: func(_ context.Context, _ uint, _ entities.TemplateMessage, _ string) (string, error) {
			return "", expectedErr
		},
	}

	convRepoMock := &mocks.ConversationRepositoryMock{
		GetByPhoneAndOrderFn: func(_ context.Context, _, _ string) (*entities.Conversation, error) {
			return nil, errors.New("not found")
		},
		CreateFn: func(_ context.Context, conv *entities.Conversation) error {
			conv.ID = "conv-001"
			return nil
		},
	}

	integRepoMock := &mocks.IntegrationRepositoryMock{
		GetWhatsAppConfigFn: func(_ context.Context, _ uint) (*ports.WhatsAppConfig, error) {
			return &ports.WhatsAppConfig{PhoneNumberID: 123456, AccessToken: "tok"}, nil
		},
	}

	uc := newUsecasesForTest(
		waMock,
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		integRepoMock,
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, err := uc.SendTemplate(
		context.Background(),
		"pedido_confirmado",
		"+573001234567",
		map[string]string{"1": "ORD-001"},
		"ORD-001",
		1,
	)

	if err == nil {
		t.Fatal("SendTemplate() esperaba error cuando el cliente WhatsApp falla")
	}
}

func TestSendTemplate_ConversacionActivaExistenteEsReutilizada(t *testing.T) {
	existingConv := &entities.Conversation{
		ID:           "conv-existente-001",
		PhoneNumber:  "+573001234567",
		OrderNumber:  "ORD-001",
		BusinessID:   1,
		CurrentState: entities.StateAwaitingConfirmation,
		Metadata:     make(map[string]interface{}),
		ExpiresAt:    time.Now().Add(12 * time.Hour), // activa
	}

	var createCalled bool
	convRepoMock := &mocks.ConversationRepositoryMock{
		GetByPhoneAndOrderFn: func(_ context.Context, _, _ string) (*entities.Conversation, error) {
			return existingConv, nil
		},
		CreateFn: func(_ context.Context, _ *entities.Conversation) error {
			createCalled = true
			return nil
		},
		UpdateFn: func(_ context.Context, _ *entities.Conversation) error {
			return nil
		},
	}

	integRepoMock := &mocks.IntegrationRepositoryMock{
		GetWhatsAppConfigFn: func(_ context.Context, _ uint) (*ports.WhatsAppConfig, error) {
			return &ports.WhatsAppConfig{PhoneNumberID: 111, AccessToken: "tok"}, nil
		},
	}

	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		integRepoMock,
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, _ = uc.SendTemplate(
		context.Background(),
		"pedido_confirmado",
		"+573001234567",
		map[string]string{"1": "ORD-001"},
		"ORD-001",
		1,
	)

	if createCalled {
		t.Error("no se debería crear una nueva conversación cuando ya existe una activa")
	}
}

// ---------------------------------------------------------------------------
// SendTemplateWithConversation
// ---------------------------------------------------------------------------

func TestSendTemplateWithConversation_Exito(t *testing.T) {
	const expectedMsgID = "wamid.withconv001"

	existingConv := &entities.Conversation{
		ID:          "conv-abc",
		PhoneNumber: "+573001234567",
		BusinessID:  1,
		ExpiresAt:   time.Now().Add(12 * time.Hour),
	}

	convRepoMock := &mocks.ConversationRepositoryMock{
		GetByIDFn: func(_ context.Context, id string) (*entities.Conversation, error) {
			if id != "conv-abc" {
				return nil, errors.New("not found")
			}
			return existingConv, nil
		},
		UpdateFn: func(_ context.Context, _ *entities.Conversation) error {
			return nil
		},
	}

	integRepoMock := &mocks.IntegrationRepositoryMock{
		GetWhatsAppConfigFn: func(_ context.Context, _ uint) (*ports.WhatsAppConfig, error) {
			return &ports.WhatsAppConfig{PhoneNumberID: 999, AccessToken: "tok"}, nil
		},
	}

	waMock := &mocks.WhatsAppMock{
		SendMessageFn: func(_ context.Context, _ uint, _ entities.TemplateMessage, _ string) (string, error) {
			return expectedMsgID, nil
		},
	}

	uc := newUsecasesForTest(
		waMock,
		convRepoMock,
		&mocks.MessageLogRepositoryMock{},
		integRepoMock,
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	// tipo_novedad_pedido no tiene variables requeridas
	msgID, err := uc.SendTemplateWithConversation(
		context.Background(),
		"tipo_novedad_pedido",
		"+573001234567",
		map[string]string{},
		"conv-abc",
	)

	if err != nil {
		t.Fatalf("SendTemplateWithConversation() error inesperado: %v", err)
	}
	if msgID != expectedMsgID {
		t.Errorf("SendTemplateWithConversation() msgID = %q, quería %q", msgID, expectedMsgID)
	}
}

func TestSendTemplateWithConversation_ErrorConversacionExpirada(t *testing.T) {
	expiradaConv := &entities.Conversation{
		ID:        "conv-expirada",
		ExpiresAt: time.Now().Add(-1 * time.Hour), // ya expiró
	}

	convRepoMock := &mocks.ConversationRepositoryMock{
		GetByIDFn: func(_ context.Context, _ string) (*entities.Conversation, error) {
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

	_, err := uc.SendTemplateWithConversation(
		context.Background(),
		"tipo_novedad_pedido",
		"+573001234567",
		map[string]string{},
		"conv-expirada",
	)

	if err == nil {
		t.Fatal("SendTemplateWithConversation() esperaba error con conversación expirada")
	}

	var errExpired *domainerrors.ErrConversationExpired
	if !errors.As(err, &errExpired) {
		t.Errorf("SendTemplateWithConversation() error = %T, quería *errors.ErrConversationExpired", err)
	}
}

func TestSendTemplateWithConversation_ErrorPlantillaNoExiste(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	_, err := uc.SendTemplateWithConversation(
		context.Background(),
		"plantilla_inexistente",
		"+573001234567",
		map[string]string{},
		"conv-abc",
	)

	if err == nil {
		t.Fatal("SendTemplateWithConversation() esperaba error con plantilla inexistente")
	}

	var errTmpl *domainerrors.ErrTemplateNotFound
	if !errors.As(err, &errTmpl) {
		t.Errorf("SendTemplateWithConversation() error = %T, quería *errors.ErrTemplateNotFound", err)
	}
}

func TestSendTemplateWithConversation_ErrorConversacionNoEncontrada(t *testing.T) {
	convRepoMock := &mocks.ConversationRepositoryMock{
		GetByIDFn: func(_ context.Context, _ string) (*entities.Conversation, error) {
			return nil, errors.New("conversación no encontrada en BD")
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

	_, err := uc.SendTemplateWithConversation(
		context.Background(),
		"tipo_novedad_pedido",
		"+573001234567",
		map[string]string{},
		"conv-inexistente",
	)

	if err == nil {
		t.Fatal("SendTemplateWithConversation() esperaba error cuando la conversación no existe en repositorio")
	}
}

// ---------------------------------------------------------------------------
// buildTemplateMessage (función interna — table-driven)
// ---------------------------------------------------------------------------

func TestBuildTemplateMessage_ConVariables(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	templateDef := entities.TemplateDefinition{
		Name:      "pedido_confirmado",
		Language:  "es",
		Variables: []string{"numero_pedido"},
	}

	msg := uc.buildTemplateMessage(
		"pedido_confirmado",
		"+573001234567",
		map[string]string{"1": "ORD-999"},
		templateDef,
	)

	if msg.MessagingProduct != "whatsapp" {
		t.Errorf("MessagingProduct = %q, quería %q", msg.MessagingProduct, "whatsapp")
	}
	if msg.To != "+573001234567" {
		t.Errorf("To = %q, quería %q", msg.To, "+573001234567")
	}
	if msg.Template.Name != "pedido_confirmado" {
		t.Errorf("Template.Name = %q, quería %q", msg.Template.Name, "pedido_confirmado")
	}
	if msg.Template.Language.Code != "es" {
		t.Errorf("Language.Code = %q, quería %q", msg.Template.Language.Code, "es")
	}
	if len(msg.Template.Components) != 1 {
		t.Fatalf("Components tiene %d elementos, quería 1", len(msg.Template.Components))
	}
	if len(msg.Template.Components[0].Parameters) != 1 {
		t.Fatalf("Body Parameters tiene %d elementos, quería 1", len(msg.Template.Components[0].Parameters))
	}
	if msg.Template.Components[0].Parameters[0].Text != "ORD-999" {
		t.Errorf("Parameter[0].Text = %q, quería %q", msg.Template.Components[0].Parameters[0].Text, "ORD-999")
	}
}

func TestBuildTemplateMessage_SinVariables(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	templateDef := entities.TemplateDefinition{
		Name:      "handoff_asesor",
		Language:  "es",
		Variables: []string{}, // sin variables
	}

	msg := uc.buildTemplateMessage("handoff_asesor", "+573001234567", map[string]string{}, templateDef)

	if len(msg.Template.Components) != 0 {
		t.Errorf("Components tiene %d elementos, quería 0 (sin variables)", len(msg.Template.Components))
	}
}

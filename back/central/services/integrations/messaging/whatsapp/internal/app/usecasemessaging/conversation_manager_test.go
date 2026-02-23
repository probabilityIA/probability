package usecasemessaging

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/mocks"
)

// newUsecasesForTest construye la instancia de Usecases con todos los mocks.
// Centraliza el setup para que cada test solo sobreescriba lo que necesita.
func newUsecasesForTest(
	waClient *mocks.WhatsAppMock,
	convRepo *mocks.ConversationRepositoryMock,
	msgRepo *mocks.MessageLogRepositoryMock,
	integRepo *mocks.IntegrationRepositoryMock,
	publisher *mocks.EventPublisherMock,
	cfg *mocks.ConfigMock,
) *Usecases {
	return &Usecases{
		whatsApp:         waClient,
		conversationRepo: convRepo,
		messageLogRepo:   msgRepo,
		integrationRepo:  integRepo,
		publisher:        publisher,
		log:              &mocks.LoggerMock{},
		config:           cfg,
	}
}

// newActiveConversation construye una conversación activa de prueba
func newActiveConversation(state entities.ConversationState) *entities.Conversation {
	return &entities.Conversation{
		ID:           "conv-test-001",
		PhoneNumber:  "+573001234567",
		OrderNumber:  "ORD-9999",
		BusinessID:   1,
		CurrentState: state,
		Metadata:     make(map[string]interface{}),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}
}

// ---------------------------------------------------------------------------
// GetInitialState
// ---------------------------------------------------------------------------

func TestGetInitialState(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	got := uc.GetInitialState()
	want := entities.StateAwaitingConfirmation

	if got != want {
		t.Errorf("GetInitialState() = %q, quería %q", got, want)
	}
}

// ---------------------------------------------------------------------------
// IsTerminalState
// ---------------------------------------------------------------------------

func TestIsTerminalState(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	tests := []struct {
		name  string
		state entities.ConversationState
		want  bool
	}{
		{
			name:  "StateCompleted es terminal",
			state: entities.StateCompleted,
			want:  true,
		},
		{
			name:  "StateHandoffToHuman es terminal",
			state: entities.StateHandoffToHuman,
			want:  true,
		},
		{
			name:  "StateStart NO es terminal",
			state: entities.StateStart,
			want:  false,
		},
		{
			name:  "StateAwaitingConfirmation NO es terminal",
			state: entities.StateAwaitingConfirmation,
			want:  false,
		},
		{
			name:  "StateAwaitingMenuSelection NO es terminal",
			state: entities.StateAwaitingMenuSelection,
			want:  false,
		},
		{
			name:  "StateAwaitingNoveltyType NO es terminal",
			state: entities.StateAwaitingNoveltyType,
			want:  false,
		},
		{
			name:  "StateAwaitingCancelConfirm NO es terminal",
			state: entities.StateAwaitingCancelConfirm,
			want:  false,
		},
		{
			name:  "StateAwaitingCancelReason NO es terminal",
			state: entities.StateAwaitingCancelReason,
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := uc.IsTerminalState(tt.state)
			if got != tt.want {
				t.Errorf("IsTerminalState(%q) = %v, quería %v", tt.state, got, tt.want)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// TransitionState — AWAITING_CONFIRMATION
// ---------------------------------------------------------------------------

func TestTransitionState_AwaitingConfirmation_ConfirmarPedido(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingConfirmation)

	transition, err := uc.TransitionState(context.Background(), conv, "Confirmar pedido")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateCompleted {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateCompleted)
	}
	if transition.TemplateName != "pedido_confirmado" {
		t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, "pedido_confirmado")
	}
	if !transition.PublishEvent {
		t.Error("PublishEvent debería ser true")
	}
	if transition.EventType != "confirmed" {
		t.Errorf("EventType = %q, quería %q", transition.EventType, "confirmed")
	}
}

func TestTransitionState_AwaitingConfirmation_NoConfirmar(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingConfirmation)

	transition, err := uc.TransitionState(context.Background(), conv, "No confirmar")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateAwaitingMenuSelection {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateAwaitingMenuSelection)
	}
	if transition.TemplateName != "menu_no_confirmacion" {
		t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, "menu_no_confirmacion")
	}
	if transition.PublishEvent {
		t.Error("PublishEvent debería ser false")
	}
}

func TestTransitionState_AwaitingConfirmation_RespuestaInvalida(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingConfirmation)

	_, err := uc.TransitionState(context.Background(), conv, "respuesta_desconocida")

	if err == nil {
		t.Fatal("TransitionState() esperaba error con respuesta inválida, no obtuvo ninguno")
	}
}

// ---------------------------------------------------------------------------
// TransitionState — AWAITING_MENU_SELECTION
// ---------------------------------------------------------------------------

func TestTransitionState_AwaitingMenuSelection_PresentarNovedad(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingMenuSelection)

	transition, err := uc.TransitionState(context.Background(), conv, "Presentar novedad")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateAwaitingNoveltyType {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateAwaitingNoveltyType)
	}
	if transition.TemplateName != "tipo_novedad_pedido" {
		t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, "tipo_novedad_pedido")
	}
}

func TestTransitionState_AwaitingMenuSelection_CancelarPedido(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingMenuSelection)

	transition, err := uc.TransitionState(context.Background(), conv, "Cancelar pedido")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateAwaitingCancelConfirm {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateAwaitingCancelConfirm)
	}
	if transition.TemplateName != "confirmar_cancelacion_pedido" {
		t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, "confirmar_cancelacion_pedido")
	}
}

func TestTransitionState_AwaitingMenuSelection_Asesor(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingMenuSelection)

	transition, err := uc.TransitionState(context.Background(), conv, "Asesor")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateHandoffToHuman {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateHandoffToHuman)
	}
	if !transition.PublishEvent || transition.EventType != "handoff" {
		t.Errorf("se esperaba PublishEvent=true y EventType=handoff, obtuvo PublishEvent=%v EventType=%q",
			transition.PublishEvent, transition.EventType)
	}
}

// ---------------------------------------------------------------------------
// TransitionState — AWAITING_NOVELTY_TYPE
// ---------------------------------------------------------------------------

func TestTransitionState_AwaitingNoveltyType(t *testing.T) {
	tests := []struct {
		userResponse     string
		wantTemplateName string
		wantNoveltyType  string
	}{
		{"Cambio de dirección", "novedad_cambio_direccion", "change_address"},
		{"Cambio de productos", "novedad_cambio_productos", "change_products"},
		{"Cambio medio de pago", "novedad_cambio_medio_pago", "change_payment"},
	}

	for _, tt := range tests {
		t.Run("novedad_"+tt.wantNoveltyType, func(t *testing.T) {
			uc := newUsecasesForTest(
				&mocks.WhatsAppMock{},
				&mocks.ConversationRepositoryMock{},
				&mocks.MessageLogRepositoryMock{},
				&mocks.IntegrationRepositoryMock{},
				&mocks.EventPublisherMock{},
				&mocks.ConfigMock{},
			)

			conv := newActiveConversation(entities.StateAwaitingNoveltyType)
			transition, err := uc.TransitionState(context.Background(), conv, tt.userResponse)

			if err != nil {
				t.Fatalf("TransitionState() error inesperado: %v", err)
			}
			if transition.NextState != entities.StateCompleted {
				t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateCompleted)
			}
			if transition.TemplateName != tt.wantTemplateName {
				t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, tt.wantTemplateName)
			}
			if !transition.PublishEvent || transition.EventType != "novelty" {
				t.Errorf("PublishEvent=%v EventType=%q, se esperaba PublishEvent=true EventType=novelty",
					transition.PublishEvent, transition.EventType)
			}
			if v, ok := transition.EventMetadata["novelty_type"]; !ok || v != tt.wantNoveltyType {
				t.Errorf("EventMetadata[novelty_type] = %v, quería %q", v, tt.wantNoveltyType)
			}
		})
	}
}

func TestTransitionState_AwaitingNoveltyType_RespuestaInvalida(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingNoveltyType)
	_, err := uc.TransitionState(context.Background(), conv, "opción_desconocida")

	if err == nil {
		t.Fatal("TransitionState() esperaba error con respuesta inválida, no obtuvo ninguno")
	}
}

// ---------------------------------------------------------------------------
// TransitionState — AWAITING_CANCEL_CONFIRM
// ---------------------------------------------------------------------------

func TestTransitionState_AwaitingCancelConfirm_SiCancelar(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingCancelConfirm)
	transition, err := uc.TransitionState(context.Background(), conv, "Sí, cancelar")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateAwaitingCancelReason {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateAwaitingCancelReason)
	}
	if transition.TemplateName != "motivo_cancelacion_pedido" {
		t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, "motivo_cancelacion_pedido")
	}
}

func TestTransitionState_AwaitingCancelConfirm_NoVolver(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingCancelConfirm)
	transition, err := uc.TransitionState(context.Background(), conv, "No, volver")

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateAwaitingMenuSelection {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateAwaitingMenuSelection)
	}
}

// ---------------------------------------------------------------------------
// TransitionState — AWAITING_CANCEL_REASON (texto libre)
// ---------------------------------------------------------------------------

func TestTransitionState_AwaitingCancelReason_CualquierTexto(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation(entities.StateAwaitingCancelReason)
	motivo := "No me llegó a tiempo el pedido"

	transition, err := uc.TransitionState(context.Background(), conv, motivo)

	if err != nil {
		t.Fatalf("TransitionState() error inesperado: %v", err)
	}
	if transition.NextState != entities.StateCompleted {
		t.Errorf("NextState = %q, quería %q", transition.NextState, entities.StateCompleted)
	}
	if transition.TemplateName != "pedido_cancelado" {
		t.Errorf("TemplateName = %q, quería %q", transition.TemplateName, "pedido_cancelado")
	}
	if !transition.PublishEvent || transition.EventType != "cancelled" {
		t.Errorf("PublishEvent=%v EventType=%q, se esperaba PublishEvent=true EventType=cancelled",
			transition.PublishEvent, transition.EventType)
	}
	if reason, ok := transition.EventMetadata["cancellation_reason"]; !ok || reason != motivo {
		t.Errorf("EventMetadata[cancellation_reason] = %v, quería %q", reason, motivo)
	}
}

// ---------------------------------------------------------------------------
// TransitionState — estado no reconocido
// ---------------------------------------------------------------------------

func TestTransitionState_EstadoNoReconocido(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	conv := newActiveConversation("ESTADO_FANTASMA")

	_, err := uc.TransitionState(context.Background(), conv, "cualquier cosa")

	if err == nil {
		t.Fatal("TransitionState() esperaba error con estado no reconocido, no obtuvo ninguno")
	}
}

// ---------------------------------------------------------------------------
// ValidateTransition
// ---------------------------------------------------------------------------

func TestValidateTransition_TransicionesValidas(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	tests := []struct {
		from entities.ConversationState
		to   entities.ConversationState
	}{
		{entities.StateStart, entities.StateAwaitingConfirmation},
		{entities.StateAwaitingConfirmation, entities.StateCompleted},
		{entities.StateAwaitingConfirmation, entities.StateAwaitingMenuSelection},
		{entities.StateAwaitingMenuSelection, entities.StateAwaitingNoveltyType},
		{entities.StateAwaitingMenuSelection, entities.StateAwaitingCancelConfirm},
		{entities.StateAwaitingMenuSelection, entities.StateHandoffToHuman},
		{entities.StateAwaitingNoveltyType, entities.StateCompleted},
		{entities.StateAwaitingCancelConfirm, entities.StateAwaitingCancelReason},
		{entities.StateAwaitingCancelConfirm, entities.StateAwaitingMenuSelection},
		{entities.StateAwaitingCancelReason, entities.StateCompleted},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to), func(t *testing.T) {
			err := uc.ValidateTransition(context.Background(), tt.from, tt.to)
			if err != nil {
				t.Errorf("ValidateTransition(%q, %q) error inesperado: %v", tt.from, tt.to, err)
			}
		})
	}
}

func TestValidateTransition_TransicionesInvalidas(t *testing.T) {
	uc := newUsecasesForTest(
		&mocks.WhatsAppMock{},
		&mocks.ConversationRepositoryMock{},
		&mocks.MessageLogRepositoryMock{},
		&mocks.IntegrationRepositoryMock{},
		&mocks.EventPublisherMock{},
		&mocks.ConfigMock{},
	)

	tests := []struct {
		from entities.ConversationState
		to   entities.ConversationState
	}{
		// Retrocesos no permitidos
		{entities.StateAwaitingConfirmation, entities.StateStart},
		{entities.StateCompleted, entities.StateStart},
		// Saltos no permitidos
		{entities.StateStart, entities.StateCompleted},
		{entities.StateAwaitingConfirmation, entities.StateAwaitingNoveltyType},
	}

	for _, tt := range tests {
		t.Run(string(tt.from)+"->"+string(tt.to)+"_invalida", func(t *testing.T) {
			err := uc.ValidateTransition(context.Background(), tt.from, tt.to)
			if err == nil {
				t.Errorf("ValidateTransition(%q, %q) esperaba error, no obtuvo ninguno", tt.from, tt.to)
			}
		})
	}
}

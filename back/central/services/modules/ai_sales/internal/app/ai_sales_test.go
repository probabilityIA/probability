package app

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"strings"
	"testing"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type suite struct {
	ai       *mocks.AIProviderMock
	cache    *mocks.SessionCacheMock
	products *mocks.ProductRepositoryMock
	clientes *mocks.CustomerRepositoryMock
	resp     *mocks.ResponsePublisherMock
	orders   *mocks.OrderPublisherMock
	config   *mocks.ConfigProviderMock
	persist  *mocks.PersistencePublisherMock
	pause    *mocks.PauseCheckerMock
	uc       IUseCase
}

func nuevaSuite() *suite {
	s := &suite{
		ai:       &mocks.AIProviderMock{},
		cache:    &mocks.SessionCacheMock{},
		products: &mocks.ProductRepositoryMock{},
		clientes: &mocks.CustomerRepositoryMock{},
		resp:     &mocks.ResponsePublisherMock{},
		orders:   &mocks.OrderPublisherMock{},
		config:   &mocks.ConfigProviderMock{},
		persist:  &mocks.PersistencePublisherMock{},
		pause:    &mocks.PauseCheckerMock{},
	}
	s.uc = New(s.ai, s.cache, s.products, s.clientes, s.resp, s.orders,
		s.config, s.persist, s.pause, mocks.NewSilentLogger())
	return s
}

func mensaje() domain.IncomingMessageDTO {
	return domain.IncomingMessageDTO{PhoneNumber: "573001234567", MessageText: "hola", BusinessID: 26}
}

func depsDe(s *suite, businessID uint) *toolDeps {
	return &toolDeps{
		productRepo:    s.products,
		customerRepo:   s.clientes,
		orderPublisher: s.orders,
		businessID:     businessID,
	}
}

func TestHandleIncoming_ConversacionPausada_IgnoraElMensaje(t *testing.T) {
	s := nuevaSuite()
	s.pause.IsAIPausedFn = func(ctx context.Context, phone string) bool { return true }

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Empty(t, s.ai.Calls(), "si un humano tomo el chat la IA no debe responder")
	assert.Empty(t, s.resp.Published)
	assert.Empty(t, s.persist.MessageLogs)
}

func TestHandleIncoming_Deshabilitado_IgnoraElMensaje(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return &domain.AIConfig{Enabled: false, DemoBusinessID: 1}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Empty(t, s.ai.Calls())
	assert.Empty(t, s.resp.Published)
}

func TestHandleIncoming_ErrorDeConfig_RespondeMensajeDeDisculpa(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return nil, stderrors.New("dynamo caido")
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Empty(t, s.ai.Calls())
	require.Len(t, s.resp.Published, 1,
		"el cliente debe recibir algo, no quedarse esperando en silencio")
	assert.Contains(t, s.resp.Published[0].Text, "no puedo procesar")
	assert.EqualValues(t, 1, s.resp.Published[0].BusinessID,
		"sin config se cae al business por defecto")
}

func TestHandleIncoming_ElBusinessDelMensajeGanaSobreElDeConfig(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return &domain.AIConfig{Enabled: true, DemoBusinessID: 1, MaxToolIterations: 5}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	require.Len(t, s.resp.Published, 1)
	assert.EqualValues(t, 26, s.resp.Published[0].BusinessID)
}

func TestHandleIncoming_SinBusinessEnElMensaje_UsaElDemo(t *testing.T) {
	s := nuevaSuite()
	dto := mensaje()
	dto.BusinessID = 0

	err := s.uc.HandleIncoming(context.Background(), dto)

	require.NoError(t, err)
	require.Len(t, s.resp.Published, 1)
	assert.EqualValues(t, 1, s.resp.Published[0].BusinessID)
}

func TestHandleIncoming_SesionNueva_CreaConversacionYRegistraEntrante(t *testing.T) {
	s := nuevaSuite()

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	require.Len(t, s.persist.Upserts, 1, "una sesion nueva crea la conversacion")
	assert.Equal(t, "573001234567", s.persist.Upserts[0].PhoneNumber)
	assert.NotEmpty(t, s.persist.Upserts[0].ID)

	require.Len(t, s.persist.MessageLogs, 2)
	assert.Equal(t, "inbound", s.persist.MessageLogs[0].Direction)
	assert.Equal(t, "hola", s.persist.MessageLogs[0].Content)
	assert.Equal(t, "outbound", s.persist.MessageLogs[1].Direction)
}

func TestHandleIncoming_SesionExistente_NoRecreaLaConversacion(t *testing.T) {
	s := nuevaSuite()
	s.cache.GetFn = func(ctx context.Context, phone string) (*domain.AISession, error) {
		return &domain.AISession{ID: "sesion-vieja", PhoneNumber: phone, BusinessID: 26}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Empty(t, s.persist.Upserts, "la conversacion ya existia")
	require.NotEmpty(t, s.persist.MessageLogs)
	assert.Equal(t, "sesion-vieja", s.persist.MessageLogs[0].ConversationID)
}

func TestHandleIncoming_FalloAlLeerLaSesion_ArrancaUnaNueva(t *testing.T) {
	s := nuevaSuite()
	s.cache.GetFn = func(ctx context.Context, phone string) (*domain.AISession, error) {
		return nil, stderrors.New("redis caido")
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err, "si Redis falla se arranca conversacion nueva, no se pierde el mensaje")
	assert.Len(t, s.persist.Upserts, 1)
	assert.Len(t, s.resp.Published, 1)
}

func TestHandleIncoming_RecortaElHistorialAlMaximo(t *testing.T) {
	s := nuevaSuite()
	viejos := make([]domain.AIMessage, 50)
	for i := range viejos {
		viejos[i] = domain.AIMessage{
			Role:    domain.RoleUser,
			Content: []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "viejo"}},
		}
	}
	s.cache.GetFn = func(ctx context.Context, phone string) (*domain.AISession, error) {
		return &domain.AISession{ID: "s1", PhoneNumber: phone, Messages: viejos}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	require.NotEmpty(t, s.ai.Calls())
	assert.Len(t, s.ai.Calls()[0].Messages, maxHistoryMessages,
		"sin recorte el prompt crece sin limite y dispara el costo por token")
}

func TestHandleIncoming_ElUltimoMensajeDelHistorialEsElDelUsuario(t *testing.T) {
	s := nuevaSuite()

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	msgs := s.ai.Calls()[0].Messages
	ultimo := msgs[len(msgs)-1]
	assert.Equal(t, domain.RoleUser, ultimo.Role)
	require.Len(t, ultimo.Content, 1)
	assert.Equal(t, "hola", ultimo.Content[0].Text)
}

func TestHandleIncoming_MandaElPromptYLasHerramientas(t *testing.T) {
	s := nuevaSuite()

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	call := s.ai.Calls()[0]
	assert.NotEmpty(t, call.SystemPrompt)
	assert.Len(t, call.Tools, len(GetToolDefinitions()))
}

func TestHandleIncoming_ErrorDelProveedor_RespondeSinRomper(t *testing.T) {
	s := nuevaSuite()
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		return nil, stderrors.New("bedrock caido")
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	require.Len(t, s.resp.Published, 1)
	assert.Contains(t, s.resp.Published[0].Text, "tuve un problema")
}

func TestHandleIncoming_RespuestaVacia_MandaUnFallback(t *testing.T) {
	s := nuevaSuite()
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		return &domain.AIResponse{StopReason: domain.StopReasonEndTurn, Content: nil}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	require.Len(t, s.resp.Published, 1)
	assert.Contains(t, s.resp.Published[0].Text, "No pude generar una respuesta",
		"nunca se envia un mensaje vacio al cliente")
}

func TestHandleIncoming_MaxTokens_TambienCierraElTurno(t *testing.T) {
	s := nuevaSuite()
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		return &domain.AIResponse{
			StopReason: domain.StopReasonMaxToken,
			Content:    []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "respuesta cortada"}},
		}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Len(t, s.ai.Calls(), 1, "no se reintenta")
	require.Len(t, s.resp.Published, 1)
	assert.Equal(t, "respuesta cortada", s.resp.Published[0].Text)
}

func TestHandleIncoming_EjecutaLasHerramientasYVuelveAPreguntar(t *testing.T) {
	s := nuevaSuite()
	llamadas := 0
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		llamadas++
		if llamadas == 1 {
			return &domain.AIResponse{
				StopReason: domain.StopReasonToolUse,
				Content: []domain.ContentBlock{{
					Type: domain.ContentTypeToolUse, ToolUseID: "t1",
					ToolName: "SearchProducts", Input: `{"query":"proteina"}`,
				}},
			}, nil
		}
		return &domain.AIResponse{
			StopReason: domain.StopReasonEndTurn,
			Content:    []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "Tenemos proteina"}},
		}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Equal(t, 2, llamadas)
	assert.Equal(t, []string{"proteina"}, s.products.SearchQuery)
	require.Len(t, s.resp.Published, 1)
	assert.Equal(t, "Tenemos proteina", s.resp.Published[0].Text)
}

func TestHandleIncoming_ElResultadoDeLaHerramientaVuelveComoMensajeDeUsuario(t *testing.T) {
	s := nuevaSuite()
	llamadas := 0
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		llamadas++
		if llamadas == 1 {
			return &domain.AIResponse{
				StopReason: domain.StopReasonToolUse,
				Content: []domain.ContentBlock{{
					Type: domain.ContentTypeToolUse, ToolUseID: "t1",
					ToolName: "SearchProducts", Input: `{"query":"x"}`,
				}},
			}, nil
		}
		return &domain.AIResponse{
			StopReason: domain.StopReasonEndTurn,
			Content:    []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "ok"}},
		}, nil
	}

	require.NoError(t, s.uc.HandleIncoming(context.Background(), mensaje()))

	segunda := s.ai.Calls()[1].Messages
	ultimo := segunda[len(segunda)-1]
	assert.Equal(t, domain.RoleUser, ultimo.Role)
	require.Len(t, ultimo.Content, 1)
	assert.Equal(t, domain.ContentTypeToolResult, ultimo.Content[0].Type)
	assert.Equal(t, "t1", ultimo.Content[0].ToolUseID,
		"el tool_use_id debe coincidir o el proveedor rechaza la conversacion")
}

func TestHandleIncoming_HerramientaQueFalla_DevuelveElErrorAlModelo(t *testing.T) {
	s := nuevaSuite()
	s.products.SearchProductsFn = func(ctx context.Context, b uint, q string, l int) ([]domain.ProductSearchResult, error) {
		return nil, stderrors.New("postgres caido")
	}
	llamadas := 0
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		llamadas++
		if llamadas == 1 {
			return &domain.AIResponse{
				StopReason: domain.StopReasonToolUse,
				Content: []domain.ContentBlock{{
					Type: domain.ContentTypeToolUse, ToolUseID: "t1",
					ToolName: "SearchProducts", Input: `{"query":"x"}`,
				}},
			}, nil
		}
		return &domain.AIResponse{
			StopReason: domain.StopReasonEndTurn,
			Content:    []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "disculpa"}},
		}, nil
	}

	require.NoError(t, s.uc.HandleIncoming(context.Background(), mensaje()))

	segunda := s.ai.Calls()[1].Messages
	ultimo := segunda[len(segunda)-1]
	assert.Contains(t, ultimo.Content[0].Content, "error",
		"el modelo debe enterarse del fallo para poder disculparse, no quedarse colgado")
}

func TestHandleIncoming_IgnoraLosBloquesQueNoSonToolUse(t *testing.T) {
	s := nuevaSuite()
	llamadas := 0
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		llamadas++
		if llamadas == 1 {
			return &domain.AIResponse{
				StopReason: domain.StopReasonToolUse,
				Content: []domain.ContentBlock{
					{Type: domain.ContentTypeText, Text: "dejame buscar"},
					{Type: domain.ContentTypeToolUse, ToolUseID: "t1", ToolName: "SearchProducts", Input: `{"query":"x"}`},
				},
			}, nil
		}
		return &domain.AIResponse{
			StopReason: domain.StopReasonEndTurn,
			Content:    []domain.ContentBlock{{Type: domain.ContentTypeText, Text: "ok"}},
		}, nil
	}

	require.NoError(t, s.uc.HandleIncoming(context.Background(), mensaje()))

	segunda := s.ai.Calls()[1].Messages
	ultimo := segunda[len(segunda)-1]
	assert.Len(t, ultimo.Content, 1, "solo el bloque de tool_use produce resultado")
}

func TestHandleIncoming_LlegaAlMaximoDeIteraciones_RespondeYCorta(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return &domain.AIConfig{Enabled: true, DemoBusinessID: 1, MaxToolIterations: 3}, nil
	}
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		return &domain.AIResponse{
			StopReason: domain.StopReasonToolUse,
			Content: []domain.ContentBlock{{
				Type: domain.ContentTypeToolUse, ToolUseID: "t1",
				ToolName: "SearchProducts", Input: `{"query":"x"}`,
			}},
		}, nil
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Len(t, s.ai.Calls(), 3, "el tope corta el ciclo: sin el, se gasta plata indefinidamente")
	require.Len(t, s.resp.Published, 1)
	assert.Contains(t, s.resp.Published[0].Text, "no pude completar tu solicitud")
}

func TestHandleIncoming_SinMaximoConfigurado_UsaCinco(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return &domain.AIConfig{Enabled: true, DemoBusinessID: 1, MaxToolIterations: 0}, nil
	}
	s.ai.ConverseFn = func(ctx context.Context, m []domain.AIMessage, p string, tl []domain.ToolDefinition) (*domain.AIResponse, error) {
		return &domain.AIResponse{
			StopReason: domain.StopReasonToolUse,
			Content: []domain.ContentBlock{{
				Type: domain.ContentTypeToolUse, ToolUseID: "t1",
				ToolName: "SearchProducts", Input: `{"query":"x"}`,
			}},
		}, nil
	}

	require.NoError(t, s.uc.HandleIncoming(context.Background(), mensaje()))

	assert.Len(t, s.ai.Calls(), 5)
}

func TestHandleIncoming_GuardaLaSesionConVencimiento(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return &domain.AIConfig{Enabled: true, DemoBusinessID: 1, MaxToolIterations: 5, SessionTTLMinutes: 30}, nil
	}

	require.NoError(t, s.uc.HandleIncoming(context.Background(), mensaje()))

	require.Len(t, s.cache.Saved, 1)
	guardada := s.cache.Saved[0]
	assert.False(t, guardada.ExpiresAt.IsZero())
	minutos := guardada.ExpiresAt.Sub(guardada.UpdatedAt).Minutes()
	assert.InDelta(t, 30.0, minutos, 1.0)
}

func TestHandleIncoming_SinTTLConfigurado_UsaVeinteMinutos(t *testing.T) {
	s := nuevaSuite()
	s.config.GetAIConfigFn = func(ctx context.Context) (*domain.AIConfig, error) {
		return &domain.AIConfig{Enabled: true, DemoBusinessID: 1, MaxToolIterations: 5, SessionTTLMinutes: 0}, nil
	}

	require.NoError(t, s.uc.HandleIncoming(context.Background(), mensaje()))

	require.Len(t, s.cache.Saved, 1)
	minutos := s.cache.Saved[0].ExpiresAt.Sub(s.cache.Saved[0].UpdatedAt).Minutes()
	assert.InDelta(t, 20.0, minutos, 1.0)
}

func TestHandleIncoming_FalloAlGuardarLaSesion_NoRompeLaRespuesta(t *testing.T) {
	s := nuevaSuite()
	s.cache.SaveFn = func(ctx context.Context, session *domain.AISession) error {
		return stderrors.New("redis caido")
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err, "el cliente recibe su respuesta aunque no se pueda cachear la sesion")
	assert.Len(t, s.resp.Published, 1)
}

func TestHandleIncoming_FallosDePersistencia_NoRompenLaRespuesta(t *testing.T) {
	s := nuevaSuite()
	s.persist.PublishConversationUpsertFn = func(ctx context.Context, session *domain.AISession) error {
		return stderrors.New("cola caida")
	}
	s.persist.PublishMessageLogFn = func(ctx context.Context, c, p, d, ct string) error {
		return stderrors.New("cola caida")
	}

	err := s.uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Len(t, s.resp.Published, 1)
}

func TestHandleIncoming_SinPersistencia_NoRompe(t *testing.T) {
	s := nuevaSuite()
	uc := New(s.ai, s.cache, s.products, s.clientes, s.resp, s.orders,
		s.config, nil, s.pause, mocks.NewSilentLogger())

	err := uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Len(t, s.resp.Published, 1)
}

func TestHandleIncoming_SinPauseChecker_Procesa(t *testing.T) {
	s := nuevaSuite()
	uc := New(s.ai, s.cache, s.products, s.clientes, s.resp, s.orders,
		s.config, s.persist, nil, mocks.NewSilentLogger())

	err := uc.HandleIncoming(context.Background(), mensaje())

	require.NoError(t, err)
	assert.Len(t, s.resp.Published, 1)
}

func TestExtractTextFromContent_LimpiaElBloqueDePensamiento(t *testing.T) {
	casos := []struct {
		nombre  string
		entrada string
		want    string
	}{
		{"sin thinking", "Hola, como estas", "Hola, como estas"},
		{"con thinking", "<thinking>el usuario quiere X</thinking>Hola", "Hola"},
		{"thinking al final", "Hola<thinking>razonamiento</thinking>", "Hola"},
		{"thinking multilinea", "<thinking>linea1\nlinea2</thinking>\nRespuesta", "Respuesta"},
		{"varios thinking", "<thinking>a</thinking>Hola<thinking>b</thinking>", "Hola"},
		{"con espacios", "  Hola  ", "Hola"},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			got := extractTextFromContent([]domain.ContentBlock{
				{Type: domain.ContentTypeText, Text: tc.entrada},
			})
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestExtractTextFromContent_SoloThinking_DevuelveVacio(t *testing.T) {
	got := extractTextFromContent([]domain.ContentBlock{
		{Type: domain.ContentTypeText, Text: "<thinking>solo razonamiento</thinking>"},
	})

	assert.Empty(t, got, "un mensaje que solo trae razonamiento no se envia al cliente")
}

func TestExtractTextFromContent_IgnoraBloquesNoTexto(t *testing.T) {
	got := extractTextFromContent([]domain.ContentBlock{
		{Type: domain.ContentTypeToolUse, ToolName: "SearchProducts"},
		{Type: domain.ContentTypeText, Text: "Aqui esta"},
	})

	assert.Equal(t, "Aqui esta", got)
}

func TestExtractTextFromContent_SinBloques_DevuelveVacio(t *testing.T) {
	assert.Empty(t, extractTextFromContent(nil))
	assert.Empty(t, extractTextFromContent([]domain.ContentBlock{}))
}

func TestExtractTextFromContent_DevuelveElPrimerTextoUtil(t *testing.T) {
	got := extractTextFromContent([]domain.ContentBlock{
		{Type: domain.ContentTypeText, Text: "<thinking>x</thinking>"},
		{Type: domain.ContentTypeText, Text: "segundo"},
	})

	assert.Equal(t, "segundo", got)
}

func TestBuildSystemPrompt_IncluyeElNegocioYLasReglasClave(t *testing.T) {
	got := BuildSystemPrompt("Mi Tienda")

	assert.Contains(t, got, "Mi Tienda")
	assert.Contains(t, got, "SearchProducts",
		"el prompt debe obligar a buscar antes de recomendar")
	assert.Contains(t, got, "NUNCA inventes productos ni precios")
	assert.Contains(t, got, "NUNCA crees un pedido sin nombre del cliente y sin direccion")
}

func TestGetToolDefinitions_ExponeLasCuatroHerramientas(t *testing.T) {
	got := GetToolDefinitions()

	nombres := make([]string, 0, len(got))
	for _, d := range got {
		nombres = append(nombres, d.Name)
		assert.NotEmpty(t, d.Description, "la herramienta %s no tiene descripcion", d.Name)
		var schema map[string]any
		require.NoError(t, json.Unmarshal([]byte(d.InputSchema), &schema),
			"el schema de %s no es JSON valido", d.Name)
	}
	assert.ElementsMatch(t,
		[]string{"SearchProducts", "CreateOrder", "SearchCustomer", "GetCustomerLastAddress"},
		nombres)
}

func TestDispatchTool_HerramientaDesconocida_Rechaza(t *testing.T) {
	s := nuevaSuite()

	_, err := DispatchTool(context.Background(), "BorrarTodo", "{}", depsDe(s, 26))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unknown tool")
}

func TestDispatchTool_RuteaCadaHerramientaDeclarada(t *testing.T) {
	s := nuevaSuite()
	entradas := map[string]string{
		"SearchProducts":         `{"query":"x"}`,
		"SearchCustomer":         `{"query":"x"}`,
		"GetCustomerLastAddress": `{"customer_id":1}`,
		"CreateOrder":            `{"items":[]}`,
	}

	for _, def := range GetToolDefinitions() {
		_, err := DispatchTool(context.Background(), def.Name, entradas[def.Name], depsDe(s, 26))
		assert.NoError(t, err, "la herramienta declarada %s no esta ruteada", def.Name)
	}
}

func TestSearchProducts_SinLimite_PideCinco(t *testing.T) {
	s := nuevaSuite()

	_, err := executeSearchProducts(context.Background(), `{"query":"proteina"}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Equal(t, []int{5}, s.products.SearchLimits)
}

func TestSearchProducts_ConLimite_LoRespeta(t *testing.T) {
	s := nuevaSuite()

	_, err := executeSearchProducts(context.Background(), `{"query":"x","limit":12}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Equal(t, []int{12}, s.products.SearchLimits)
}

func TestSearchProducts_SinResultados_DevuelveMensajeExplicito(t *testing.T) {
	s := nuevaSuite()

	got, err := executeSearchProducts(context.Background(), `{"query":"nada"}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "No se encontraron productos")
}

func TestSearchProducts_DisponibilidadSegunInventario(t *testing.T) {
	casos := []struct {
		nombre    string
		track     bool
		stock     int
		wantAvail bool
	}{
		{"no trackea inventario siempre disponible", false, 0, true},
		{"trackea con stock", true, 5, true},
		{"trackea sin stock", true, 0, false},
		{"trackea con stock negativo", true, -3, false},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			s := nuevaSuite()
			s.products.SearchProductsFn = func(ctx context.Context, b uint, q string, l int) ([]domain.ProductSearchResult, error) {
				return []domain.ProductSearchResult{{
					SKU: "S-1", Name: "P", TrackInventory: tc.track, StockQuantity: tc.stock,
				}}, nil
			}

			got, err := executeSearchProducts(context.Background(), `{"query":"x"}`, depsDe(s, 26))

			require.NoError(t, err)
			var out map[string]any
			require.NoError(t, json.Unmarshal([]byte(got), &out))
			results := out["results"].([]any)
			require.Len(t, results, 1)
			assert.Equal(t, tc.wantAvail, results[0].(map[string]any)["available"])
		})
	}
}

func TestSearchProducts_ErroresDeEntradaYRepo(t *testing.T) {
	s := nuevaSuite()

	_, err := executeSearchProducts(context.Background(), `no es json`, depsDe(s, 26))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing")

	dbErr := stderrors.New("postgres caido")
	s2 := nuevaSuite()
	s2.products.SearchProductsFn = func(ctx context.Context, b uint, q string, l int) ([]domain.ProductSearchResult, error) {
		return nil, dbErr
	}
	_, err = executeSearchProducts(context.Background(), `{"query":"x"}`, depsDe(s2, 26))
	assert.ErrorIs(t, err, dbErr)
}

func TestSearchCustomer_QueryVacio_NoConsulta(t *testing.T) {
	s := nuevaSuite()
	consultado := false
	s.clientes.SearchCustomersFn = func(ctx context.Context, b uint, q string) ([]domain.CustomerSearchResult, error) {
		consultado = true
		return nil, nil
	}

	got, err := executeSearchCustomer(context.Background(), `{"query":""}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "termino de busqueda")
	assert.False(t, consultado)
}

func TestSearchCustomer_SinResultados_DevuelveMensaje(t *testing.T) {
	s := nuevaSuite()

	got, err := executeSearchCustomer(context.Background(), `{"query":"3001234567"}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "No se encontro ningun cliente")
}

func TestSearchCustomer_BuscaEnElNegocioDeLaSesion(t *testing.T) {
	s := nuevaSuite()
	var visto uint
	s.clientes.SearchCustomersFn = func(ctx context.Context, b uint, q string) ([]domain.CustomerSearchResult, error) {
		visto = b
		return []domain.CustomerSearchResult{{ID: 1, Name: "Ana"}}, nil
	}

	got, err := executeSearchCustomer(context.Background(), `{"query":"ana"}`, depsDe(s, 99))

	require.NoError(t, err)
	assert.Equal(t, uint(99), visto,
		"la IA no debe poder consultar clientes de otro negocio")
	assert.Contains(t, got, "Ana")
}

func TestSearchCustomer_ErroresDeEntradaYRepo(t *testing.T) {
	s := nuevaSuite()

	_, err := executeSearchCustomer(context.Background(), `no es json`, depsDe(s, 26))
	require.Error(t, err)

	dbErr := stderrors.New("timeout")
	s2 := nuevaSuite()
	s2.clientes.SearchCustomersFn = func(ctx context.Context, b uint, q string) ([]domain.CustomerSearchResult, error) {
		return nil, dbErr
	}
	_, err = executeSearchCustomer(context.Background(), `{"query":"x"}`, depsDe(s2, 26))
	assert.ErrorIs(t, err, dbErr)
}

func TestCreateOrder_SinItems_NoPublica(t *testing.T) {
	s := nuevaSuite()

	got, err := executeCreateOrder(context.Background(), `{"items":[]}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "No se proporcionaron productos")
	assert.Empty(t, s.orders.Published)
}

func TestCreateOrder_CantidadInvalida_NoPublica(t *testing.T) {
	for _, cantidad := range []int{0, -1} {
		s := nuevaSuite()
		input := `{"items":[{"product_sku":"S-1","quantity":` + itoa(cantidad) + `}]}`

		got, err := executeCreateOrder(context.Background(), input, depsDe(s, 26))

		require.NoError(t, err)
		assert.Contains(t, got, "Cantidad invalida")
		assert.Empty(t, s.orders.Published)
	}
}

func TestCreateOrder_ProductoInexistente_NoPublica(t *testing.T) {
	s := nuevaSuite()
	s.products.GetProductBySKUFn = func(ctx context.Context, b uint, sku string) (*domain.ProductSearchResult, error) {
		return nil, stderrors.New("no existe")
	}

	got, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"FANTASMA","quantity":1}]}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "Producto no encontrado: FANTASMA",
		"la IA puede alucinar un SKU: el backend debe rechazarlo")
	assert.Empty(t, s.orders.Published)
}

func TestCreateOrder_StockInsuficiente_NoPublica(t *testing.T) {
	s := nuevaSuite()
	s.products.GetProductBySKUFn = func(ctx context.Context, b uint, sku string) (*domain.ProductSearchResult, error) {
		return &domain.ProductSearchResult{
			SKU: sku, Name: "Proteina", TrackInventory: true, StockQuantity: 2, Price: 1000,
		}, nil
	}

	got, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"S-1","quantity":5}]}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "Stock insuficiente")
	assert.Empty(t, s.orders.Published)
}

func TestCreateOrder_SinTrackeoDeInventario_NoValidaStock(t *testing.T) {
	s := nuevaSuite()
	s.products.GetProductBySKUFn = func(ctx context.Context, b uint, sku string) (*domain.ProductSearchResult, error) {
		return &domain.ProductSearchResult{
			SKU: sku, Name: "Servicio", TrackInventory: false, StockQuantity: 0,
			Price: 1000, Currency: "COP",
		}, nil
	}

	_, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"S-1","quantity":100}]}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Len(t, s.orders.Published, 1)
}

func TestCreateOrder_ElPrecioSaleDelCatalogo(t *testing.T) {
	s := nuevaSuite()
	s.products.GetProductBySKUFn = func(ctx context.Context, b uint, sku string) (*domain.ProductSearchResult, error) {
		return &domain.ProductSearchResult{
			ID: "P-1", SKU: sku, Name: "Proteina", Price: 85000, Currency: "COP",
		}, nil
	}

	got, err := executeCreateOrder(context.Background(),
		`{"customer_name":"Ana","customer_phone":"300","shipping_address":"Cra 1","shipping_city":"Bogota","items":[{"product_sku":"S-1","quantity":2}]}`,
		depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, `"total":170000`)

	require.Len(t, s.orders.Published, 1)
	var orden map[string]any
	require.NoError(t, json.Unmarshal(s.orders.Published[0], &orden))
	assert.InDelta(t, 170000.0, orden["total_amount"], 0.001,
		"el total lo calcula el backend con el precio real, no lo que diga el modelo")
	assert.InDelta(t, 170000.0, orden["subtotal"], 0.001)
}

func TestCreateOrder_SinIntegracionDeWhatsApp_NoPublica(t *testing.T) {
	s := nuevaSuite()
	s.clientes.GetWhatsAppIntegrationIDFn = func(ctx context.Context, b uint) (uint, error) {
		return 0, stderrors.New("no configurada")
	}

	got, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"S-1","quantity":1}]}`, depsDe(s, 26))

	require.NoError(t, err)
	assert.Contains(t, got, "No se encontro integracion de WhatsApp")
	assert.Empty(t, s.orders.Published)
}

func TestCreateOrder_MarcaElOrigenYNaceEnPendiente(t *testing.T) {
	s := nuevaSuite()

	_, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"S-1","quantity":1}]}`, depsDe(s, 26))

	require.NoError(t, err)
	require.Len(t, s.orders.Published, 1)
	var orden map[string]any
	require.NoError(t, json.Unmarshal(s.orders.Published[0], &orden))
	assert.Equal(t, "whatsapp_ai", orden["integration_type"])
	assert.Equal(t, "whatsapp_ai", orden["platform"])
	assert.Equal(t, "pending", orden["status"])
	assert.Equal(t, "pending", orden["original_status"])
	assert.EqualValues(t, 26, orden["business_id"])
	assert.EqualValues(t, 9, orden["integration_id"])
}

func TestCreateOrder_CopiaLosDatosDelCliente(t *testing.T) {
	s := nuevaSuite()

	_, err := executeCreateOrder(context.Background(),
		`{"customer_name":"Ana Perez","customer_phone":"573001234567","shipping_address":"Cra 1 #2-3","shipping_city":"Bogota","items":[{"product_sku":"S-1","quantity":1}]}`,
		depsDe(s, 26))

	require.NoError(t, err)
	var orden map[string]any
	require.NoError(t, json.Unmarshal(s.orders.Published[0], &orden))
	assert.Equal(t, "Ana Perez", orden["customer_name"])
	assert.Equal(t, "573001234567", orden["customer_phone"])
	assert.Equal(t, "Cra 1 #2-3", orden["shipping_address"])
	assert.Equal(t, "Bogota", orden["shipping_city"])
}

func TestCreateOrder_GeneraIdentificadoresUnicos(t *testing.T) {
	externos := map[string]bool{}
	numeros := map[string]bool{}

	for i := 0; i < 20; i++ {
		s := nuevaSuite()
		_, err := executeCreateOrder(context.Background(),
			`{"items":[{"product_sku":"S-1","quantity":1}]}`, depsDe(s, 26))
		require.NoError(t, err)

		var orden map[string]any
		require.NoError(t, json.Unmarshal(s.orders.Published[0], &orden))
		ext := orden["external_id"].(string)
		num := orden["order_number"].(string)
		require.False(t, externos[ext], "external_id repetido: %q", ext)
		externos[ext] = true
		numeros[num] = true
		assert.True(t, strings.HasPrefix(num, "AI-"))
	}
	assert.Len(t, externos, 20)
}

func TestCreateOrder_VariosItems_SumaYDetalla(t *testing.T) {
	precios := map[string]float64{"S-1": 50000, "S-2": 1000}
	s := nuevaSuite()
	s.products.GetProductBySKUFn = func(ctx context.Context, b uint, sku string) (*domain.ProductSearchResult, error) {
		return &domain.ProductSearchResult{
			ID: "P-" + sku, SKU: sku, Name: "P " + sku, Price: precios[sku], Currency: "COP",
		}, nil
	}

	_, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"S-1","quantity":2},{"product_sku":"S-2","quantity":3}]}`,
		depsDe(s, 26))

	require.NoError(t, err)
	var orden map[string]any
	require.NoError(t, json.Unmarshal(s.orders.Published[0], &orden))
	assert.InDelta(t, 103000.0, orden["total_amount"], 0.001)
	items := orden["order_items"].([]any)
	require.Len(t, items, 2)
	assert.InDelta(t, 100000.0, items[0].(map[string]any)["total_price"], 0.001)
	assert.InDelta(t, 3000.0, items[1].(map[string]any)["total_price"], 0.001)
}

func TestCreateOrder_FallaLaPublicacion_SePropaga(t *testing.T) {
	pubErr := stderrors.New("rabbit caido")
	s := nuevaSuite()
	s.orders.PublishOrderFn = func(ctx context.Context, payload []byte) error { return pubErr }

	_, err := executeCreateOrder(context.Background(),
		`{"items":[{"product_sku":"S-1","quantity":1}]}`, depsDe(s, 26))

	assert.ErrorIs(t, err, pubErr)
}

func TestCreateOrder_EntradaInvalida_Rechaza(t *testing.T) {
	s := nuevaSuite()

	_, err := executeCreateOrder(context.Background(), `no es json`, depsDe(s, 26))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing")
	assert.Empty(t, s.orders.Published)
}

func TestGetCustomerLastAddress_EntradaInvalida_Rechaza(t *testing.T) {
	s := nuevaSuite()

	_, err := executeGetCustomerLastAddress(context.Background(), `no es json`, depsDe(s, 26))

	require.Error(t, err)
}

func TestGetCustomerLastAddress_BuscaEnElNegocioDeLaSesion(t *testing.T) {
	s := nuevaSuite()
	var vistoNegocio, vistoCliente uint
	s.clientes.GetCustomerLastAddressFn = func(ctx context.Context, b, c uint) (*domain.CustomerLastAddress, error) {
		vistoNegocio, vistoCliente = b, c
		return &domain.CustomerLastAddress{Street: "Cra 1", City: "Bogota"}, nil
	}

	got, err := executeGetCustomerLastAddress(context.Background(),
		`{"customer_id":7}`, depsDe(s, 99))

	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio)
	assert.Equal(t, uint(7), vistoCliente)
	assert.Contains(t, got, "Cra 1")
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}
	neg := v < 0
	if neg {
		v = -v
	}
	var b []byte
	for v > 0 {
		b = append([]byte{byte('0' + v%10)}, b...)
		v /= 10
	}
	if neg {
		return "-" + string(b)
	}
	return string(b)
}

package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/errors"
)

func ordersCatalog() *entities.NavigationCatalog {
	return &entities.NavigationCatalog{Allowed: []entities.Destination{
		{Key: "orders", Label: "Ordenes", Route: "/orders"},
		{Key: "shipments", Label: "Envios", Route: "/shipments"},
	}}
}

func toolCallReply(name string, input map[string]any) *dtos.ModelReply {
	return &dtos.ModelReply{ToolCalls: []dtos.ToolCall{{ID: "t1", Name: name, Input: input}}, InputTokens: 100, OutputTokens: 10}
}

func TestChat_ConsultaUnaOrdenDelNegocioDelToken(t *testing.T) {
	reader := &readerFake{orders: []entities.OrderInfo{{Number: "DEM-0048", StatusName: "Pendiente", Total: 300000}}}
	model := &modelFake{replies: []*dtos.ModelReply{
		toolCallReply(toolFindOrder, map[string]any{"numero": "DEM-0048"}),
		{Message: "La orden DEM-0048 esta Pendiente.", DestinationKey: "orders", InputTokens: 200, OutputTokens: 20},
	}}
	uc := newDataUseCase(model, reader, ordersCatalog())

	reply, err := uc.Chat(context.Background(), dtos.ChatInput{
		Scope:    dtos.AccessScope{UserID: 7, TokenBusinessID: 26, RequestedBusinessID: 99},
		Messages: []entities.ChatMessage{userMessage("en que va la DEM-0048")},
	})

	require.NoError(t, err)
	assert.Equal(t, "La orden DEM-0048 esta Pendiente.", reply.Message)
	assert.Equal(t, []uint{26}, reader.businessIDs, "un usuario normal consulta siempre su negocio del token, nunca el enviado")
	require.Len(t, model.requests, 2)
	assert.NotEmpty(t, model.requests[0].Tools)
	assert.False(t, model.requests[0].ForceReply)

	last := model.requests[1].Messages
	require.GreaterOrEqual(t, len(last), 3)
	result := last[len(last)-1].ToolResults
	require.Len(t, result, 1)
	assert.Equal(t, "t1", result[0].ToolCallID)
	assert.Equal(t, 1, result[0].Content["encontradas"])
	assert.Contains(t, model.requests[0].SystemPrompt, "14 de septiembre de 2026")
}

func TestChat_SuperAdminSinNegocioNoTieneHerramientas(t *testing.T) {
	reader := &readerFake{}
	model := &modelFake{reply: &dtos.ModelReply{Message: "Elige un negocio primero.", DestinationKey: "none"}}
	uc := newDataUseCase(model, reader, ordersCatalog())

	_, err := uc.Chat(context.Background(), dtos.ChatInput{
		Scope:    dtos.AccessScope{UserID: 1},
		Messages: []entities.ChatMessage{userMessage("cuantas ordenes tengo")},
	})

	require.NoError(t, err)
	require.Len(t, model.requests, 1)
	assert.Empty(t, model.requests[0].Tools)
	assert.True(t, model.requests[0].ForceReply)
	assert.Contains(t, model.requests[0].SystemPrompt, "No hay un negocio seleccionado")
	assert.Empty(t, reader.businessIDs)
}

func TestChat_SuperAdminConsultaElNegocioSeleccionado(t *testing.T) {
	reader := &readerFake{}
	model := &modelFake{replies: []*dtos.ModelReply{
		toolCallReply(toolOrdersSummary, map[string]any{}),
		{Message: "Esta semana no hay ordenes.", DestinationKey: "none"},
	}}
	uc := newDataUseCase(model, reader, ordersCatalog())

	_, err := uc.Chat(context.Background(), dtos.ChatInput{
		Scope:    dtos.AccessScope{UserID: 1, RequestedBusinessID: 26},
		Messages: []entities.ChatMessage{userMessage("resumen de la semana")},
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{26}, reader.businessIDs)
}

func TestChat_SinPermisoDeOrdenesLaHerramientaSeRechaza(t *testing.T) {
	reader := &readerFake{}
	catalog := &entities.NavigationCatalog{Allowed: []entities.Destination{{Key: "shipments", Label: "Envios", Route: "/shipments"}}}
	model := &modelFake{replies: []*dtos.ModelReply{
		toolCallReply(toolFindOrder, map[string]any{"numero": "1"}),
		{Message: "No puedo ver ordenes.", DestinationKey: "none"},
	}}
	uc := newDataUseCase(model, reader, catalog)

	_, err := uc.Chat(context.Background(), dtos.ChatInput{
		Scope:    dtos.AccessScope{UserID: 7, TokenBusinessID: 26},
		Messages: []entities.ChatMessage{userMessage("orden 1")},
	})

	require.NoError(t, err)
	assert.Empty(t, reader.businessIDs, "consultar_orden no debe ejecutarse sin permiso de ordenes")
	for _, tool := range model.requests[0].Tools {
		assert.NotEqual(t, toolFindOrder, tool.Name)
	}
	results := model.requests[1].Messages[len(model.requests[1].Messages)-1].ToolResults
	require.Len(t, results, 1)
	assert.Contains(t, results[0].Content, "error")
}

func TestChat_CortaElBucleDeHerramientas(t *testing.T) {
	reader := &readerFake{}
	model := &modelFake{reply: toolCallReply(toolListOrders, map[string]any{})}
	uc := newDataUseCase(model, reader, ordersCatalog())

	_, err := uc.Chat(context.Background(), dtos.ChatInput{
		Scope:    dtos.AccessScope{UserID: 7, TokenBusinessID: 26},
		Messages: []entities.ChatMessage{userMessage("lista")},
	})

	assert.ErrorIs(t, err, domainerrors.ErrModelUnavailable)
	require.Len(t, model.requests, MaxToolRounds+1)
	assert.True(t, model.requests[MaxToolRounds].ForceReply)
}

func TestDateRangeArgs_HastaEsInclusivo(t *testing.T) {
	from, to := dateRangeArgs(map[string]any{"desde": "2026-09-08", "hasta": "2026-09-14"})
	require.NotNil(t, from)
	require.NotNil(t, to)
	assert.Equal(t, "2026-09-15", to.Format(dateLayout))
}

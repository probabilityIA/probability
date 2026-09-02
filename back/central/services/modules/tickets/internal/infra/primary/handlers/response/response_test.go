package response

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromTicket_TrasladaTodosLosCampos(t *testing.T) {
	creado := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	vence := creado.Add(48 * time.Hour)
	negocio := uint(26)
	responsable := uint(7)
	sprint := uint(9)

	got := FromTicket(&entities.Ticket{
		ID: 1, Code: "TKT-000001",
		BusinessID: &negocio, BusinessName: "Demo",
		CreatedByID: 42, CreatedByName: "Ana", CreatedByAvatarURL: "s3://ana.png",
		AssignedToID: &responsable, AssignedToName: "Luis", AssignedToAvatarURL: "s3://luis.png",
		SprintID: &sprint, SprintName: "Sprint 1",
		Title: "No carga", Description: "queda en blanco",
		Type: "bug", Category: "ordenes", Priority: "high", Status: "open",
		Source: "business", Severity: "high", Area: "soporte",
		EscalatedToDev: true, EscalatedAt: &creado,
		DueDate: &vence, ResolvedAt: &creado, ClosedAt: &creado,
		CreatedAt: creado, UpdatedAt: creado,
		CommentsCount: 3, AttachmentsCount: 2,
	})

	assert.Equal(t, uint(1), got.ID)
	assert.Equal(t, "TKT-000001", got.Code)
	require.NotNil(t, got.BusinessID)
	assert.Equal(t, uint(26), *got.BusinessID)
	assert.Equal(t, "Demo", got.BusinessName)
	assert.Equal(t, uint(42), got.CreatedByID)
	assert.Equal(t, "Ana", got.CreatedByName)
	assert.Equal(t, "s3://ana.png", got.CreatedByAvatarURL)
	require.NotNil(t, got.AssignedToID)
	assert.Equal(t, uint(7), *got.AssignedToID)
	assert.Equal(t, "Luis", got.AssignedToName)
	assert.Equal(t, "s3://luis.png", got.AssignedToAvatarURL)
	require.NotNil(t, got.SprintID)
	assert.Equal(t, uint(9), *got.SprintID)
	assert.Equal(t, "Sprint 1", got.SprintName)
	assert.Equal(t, "No carga", got.Title)
	assert.Equal(t, "queda en blanco", got.Description)
	assert.Equal(t, "bug", got.Type)
	assert.Equal(t, "ordenes", got.Category)
	assert.Equal(t, "high", got.Priority)
	assert.Equal(t, "open", got.Status)
	assert.Equal(t, "business", got.Source)
	assert.Equal(t, "high", got.Severity)
	assert.Equal(t, "soporte", got.Area)
	assert.True(t, got.EscalatedToDev)
	assert.Equal(t, &creado, got.EscalatedAt)
	assert.Equal(t, &vence, got.DueDate)
	assert.Equal(t, &creado, got.ResolvedAt)
	assert.Equal(t, &creado, got.ClosedAt)
	assert.Equal(t, creado, got.CreatedAt)
	assert.Equal(t, creado, got.UpdatedAt)
	assert.Equal(t, int64(3), got.CommentsCount)
	assert.Equal(t, int64(2), got.AttachmentsCount)
}

func TestFromTicket_NoExponeLaDescripcionDelDominioNiRelaciones(t *testing.T) {
	crudo, err := json.Marshal(FromTicket(&entities.Ticket{
		ID: 1, Title: "t",
		Comments:    []entities.TicketComment{{ID: 1, Body: "interno"}},
		Attachments: []entities.TicketAttachment{{ID: 1}},
		History:     []entities.TicketStatusHistory{{ID: 1}},
	}))

	require.NoError(t, err)
	assert.NotContains(t, string(crudo), "interno",
		"la respuesta de ticket no arrastra comentarios: se piden por su propio endpoint")
	assert.NotContains(t, string(crudo), "history")
}

func TestFromTicket_CamposVaciosSeOmitenEnElJSON(t *testing.T) {
	crudo, err := json.Marshal(FromTicket(&entities.Ticket{ID: 1, Title: "t", Status: "open"}))

	require.NoError(t, err)
	var salida map[string]any
	require.NoError(t, json.Unmarshal(crudo, &salida))

	for _, omitido := range []string{"business_name", "severity", "area", "due_date", "resolved_at", "closed_at"} {
		assert.NotContains(t, salida, omitido, "campo vacio %q no debe viajar al front", omitido)
	}
	assert.Contains(t, salida, "business_id", "el id de negocio viaja siempre, aunque sea null")
	assert.Contains(t, salida, "sprint_name")
	assert.Contains(t, salida, "status")
}

func TestFromComment_IncluyeSusAdjuntos(t *testing.T) {
	creado := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)

	got := FromComment(&entities.TicketComment{
		ID: 3, TicketID: 1, UserID: 42, UserName: "Ana",
		Body: "ya lo revise", IsInternal: true, CreatedAt: creado,
		Attachments: []entities.TicketAttachment{
			{ID: 10, FileName: "a.png"},
			{ID: 11, FileName: "b.pdf"},
		},
	})

	assert.Equal(t, uint(3), got.ID)
	assert.Equal(t, uint(1), got.TicketID)
	assert.Equal(t, uint(42), got.UserID)
	assert.Equal(t, "Ana", got.UserName)
	assert.Equal(t, "ya lo revise", got.Body)
	assert.True(t, got.IsInternal)
	assert.Equal(t, creado, got.CreatedAt)
	require.Len(t, got.Attachments, 2)
	assert.Equal(t, "a.png", got.Attachments[0].FileName)
	assert.Equal(t, "b.pdf", got.Attachments[1].FileName)
}

func TestFromComment_SinAdjuntos_DejaLaListaNula(t *testing.T) {
	got := FromComment(&entities.TicketComment{ID: 3, Body: "hola"})

	assert.Nil(t, got.Attachments)

	crudo, err := json.Marshal(got)
	require.NoError(t, err)
	assert.NotContains(t, string(crudo), "attachments")
}

func TestFromAttachment_TrasladaTodosLosCampos(t *testing.T) {
	creado := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	comentario := uint(3)

	got := FromAttachment(&entities.TicketAttachment{
		ID: 10, TicketID: 1, CommentID: &comentario,
		UploadedByID: 42, UploadedByName: "Ana",
		FileURL: "s3://bucket/a.png", FileName: "a.png",
		MimeType: "image/png", Size: 1234, CreatedAt: creado,
	})

	assert.Equal(t, uint(10), got.ID)
	assert.Equal(t, uint(1), got.TicketID)
	require.NotNil(t, got.CommentID)
	assert.Equal(t, uint(3), *got.CommentID)
	assert.Equal(t, uint(42), got.UploadedByID)
	assert.Equal(t, "Ana", got.UploadedByName)
	assert.Equal(t, "s3://bucket/a.png", got.FileURL)
	assert.Equal(t, "a.png", got.FileName)
	assert.Equal(t, "image/png", got.MimeType)
	assert.Equal(t, int64(1234), got.Size)
	assert.Equal(t, creado, got.CreatedAt)
}

func TestFromHistory_TrasladaAmbosTiposDeCambio(t *testing.T) {
	creado := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)

	estado := FromHistory(&entities.TicketStatusHistory{
		ID: 1, TicketID: 5, ChangeType: "status",
		FromStatus: "open", ToStatus: "in_review",
		ChangedByID: 42, ChangedByName: "Ana", Note: "empezamos", CreatedAt: creado,
	})
	assert.Equal(t, uint(1), estado.ID)
	assert.Equal(t, uint(5), estado.TicketID)
	assert.Equal(t, "status", estado.ChangeType)
	assert.Equal(t, "open", estado.FromStatus)
	assert.Equal(t, "in_review", estado.ToStatus)
	assert.Empty(t, estado.FromArea)
	assert.Equal(t, uint(42), estado.ChangedByID)
	assert.Equal(t, "Ana", estado.ChangedByName)
	assert.Equal(t, "empezamos", estado.Note)
	assert.Equal(t, creado, estado.CreatedAt)

	area := FromHistory(&entities.TicketStatusHistory{
		ID: 2, ChangeType: "area", FromArea: "soporte", ToArea: "desarrollo",
	})
	assert.Equal(t, "area", area.ChangeType)
	assert.Equal(t, "soporte", area.FromArea)
	assert.Equal(t, "desarrollo", area.ToArea)
	assert.Empty(t, area.FromStatus)
}

package repository

import (
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestEntityToModel_TrasladaLosCamposPersistibles(t *testing.T) {
	negocio, responsable, sprint := uint(26), uint(7), uint(9)
	escalado := time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)
	vence := escalado.Add(48 * time.Hour)

	m := entityToModel(&entities.Ticket{
		ID: 99, Code: "TKT-000001",
		BusinessID: &negocio, CreatedByID: 42, AssignedToID: &responsable, SprintID: &sprint,
		Title: "No carga", Description: "queda en blanco",
		Type: "bug", Category: "ordenes", Priority: "high", Status: "open",
		Source: "business", Severity: "high", Area: "soporte",
		EscalatedToDev: true, EscalatedAt: &escalado,
		DueDate: &vence, ResolvedAt: &escalado, ClosedAt: &escalado,
		CreatedAt: escalado,
	})

	assert.Equal(t, "TKT-000001", m.Code)
	require.NotNil(t, m.BusinessID)
	assert.Equal(t, uint(26), *m.BusinessID)
	assert.Equal(t, uint(42), m.CreatedByID)
	require.NotNil(t, m.AssignedToID)
	assert.Equal(t, uint(7), *m.AssignedToID)
	require.NotNil(t, m.SprintID)
	assert.Equal(t, uint(9), *m.SprintID)
	assert.Equal(t, "No carga", m.Title)
	assert.Equal(t, "queda en blanco", m.Description)
	assert.Equal(t, "bug", m.Type)
	assert.Equal(t, "ordenes", m.Category)
	assert.Equal(t, "high", m.Priority)
	assert.Equal(t, "open", m.Status)
	assert.Equal(t, "business", m.Source)
	assert.Equal(t, "high", m.Severity)
	assert.Equal(t, "soporte", m.Area)
	assert.True(t, m.EscalatedToDev)
	assert.Equal(t, &escalado, m.EscalatedAt)
	assert.Equal(t, &vence, m.DueDate)
	assert.Equal(t, &escalado, m.ResolvedAt)
	assert.Equal(t, &escalado, m.ClosedAt)

	assert.Zero(t, m.ID, "el id y las fechas las pone la base, no la entidad")
	assert.True(t, m.CreatedAt.IsZero())
}

func TestModelToEntity_TrasladaLosCamposYLasRelacionesCargadas(t *testing.T) {
	negocio, responsable, sprint := uint(26), uint(7), uint(9)
	creado := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)

	got := modelToEntity(&models.Ticket{
		Model:      gorm.Model{ID: 99, CreatedAt: creado, UpdatedAt: creado},
		Code:       "TKT-000099",
		BusinessID: &negocio, Business: &models.Business{Name: "Demo"},
		CreatedByID: 42, CreatedBy: models.User{Name: "Ana", AvatarURL: "s3://ana.png"},
		AssignedToID: &responsable, AssignedTo: &models.User{Name: "Luis", AvatarURL: "s3://luis.png"},
		SprintID: &sprint, Sprint: &models.Sprint{Name: "Sprint 1"},
		Title: "No carga", Description: "queda en blanco",
		Type: "bug", Category: "ordenes", Priority: "high", Status: "open",
		Source: "business", Severity: "high", Area: "soporte",
		EscalatedToDev: true, EscalatedAt: &creado,
		DueDate: &creado, ResolvedAt: &creado, ClosedAt: &creado,
	})

	assert.Equal(t, uint(99), got.ID)
	assert.Equal(t, "TKT-000099", got.Code)
	assert.Equal(t, "Demo", got.BusinessName)
	assert.Equal(t, "Ana", got.CreatedByName)
	assert.Equal(t, "s3://ana.png", got.CreatedByAvatarURL)
	assert.Equal(t, "Luis", got.AssignedToName)
	assert.Equal(t, "s3://luis.png", got.AssignedToAvatarURL)
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
	assert.Equal(t, creado, got.CreatedAt)
	assert.Equal(t, creado, got.UpdatedAt)
}

func TestModelToEntity_SinRelacionesCargadas_NoRevienta(t *testing.T) {
	got := modelToEntity(&models.Ticket{
		Model: gorm.Model{ID: 1}, Code: "TKT-000001", Title: "t",
	})

	require.NotNil(t, got)
	assert.Empty(t, got.BusinessName, "un ticket interno no tiene negocio")
	assert.Empty(t, got.AssignedToName, "un ticket sin responsable no muestra nombre")
	assert.Empty(t, got.SprintName)
	assert.Empty(t, got.CreatedByName)
	assert.Nil(t, got.BusinessID)
	assert.Nil(t, got.AssignedToID)
	assert.Nil(t, got.SprintID)
}

func TestCommentToEntity_IncluyeAutorYAdjuntos(t *testing.T) {
	creado := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)

	got := commentToEntity(&models.TicketComment{
		Model:    gorm.Model{ID: 31, CreatedAt: creado},
		TicketID: 5, UserID: 42, User: models.User{Name: "Ana"},
		Body: "ya lo revise", IsInternal: true,
		Attachments: []models.TicketAttachment{
			{Model: gorm.Model{ID: 1}, FileName: "a.png"},
			{Model: gorm.Model{ID: 2}, FileName: "b.pdf"},
		},
	})

	assert.Equal(t, uint(31), got.ID)
	assert.Equal(t, uint(5), got.TicketID)
	assert.Equal(t, uint(42), got.UserID)
	assert.Equal(t, "Ana", got.UserName)
	assert.Equal(t, "ya lo revise", got.Body)
	assert.True(t, got.IsInternal)
	assert.Equal(t, creado, got.CreatedAt)
	require.Len(t, got.Attachments, 2)
	assert.Equal(t, "a.png", got.Attachments[0].FileName)
	assert.Equal(t, "b.pdf", got.Attachments[1].FileName)
}

func TestCommentToEntity_SinAdjuntos_DejaLaListaNula(t *testing.T) {
	got := commentToEntity(&models.TicketComment{Model: gorm.Model{ID: 31}, Body: "hola"})

	assert.Nil(t, got.Attachments)
}

func TestAttachmentToEntity_TrasladaTodosLosCampos(t *testing.T) {
	creado := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)
	comentario := uint(31)

	got := attachmentToEntity(&models.TicketAttachment{
		Model:    gorm.Model{ID: 77, CreatedAt: creado},
		TicketID: 5, CommentID: &comentario,
		UploadedByID: 42, UploadedBy: models.User{Name: "Ana"},
		FileURL: "s3://bucket/a.png", FileName: "a.png",
		MimeType: "image/png", Size: 1234,
	})

	assert.Equal(t, uint(77), got.ID)
	assert.Equal(t, uint(5), got.TicketID)
	require.NotNil(t, got.CommentID)
	assert.Equal(t, uint(31), *got.CommentID)
	assert.Equal(t, uint(42), got.UploadedByID)
	assert.Equal(t, "Ana", got.UploadedByName)
	assert.Equal(t, "s3://bucket/a.png", got.FileURL)
	assert.Equal(t, "a.png", got.FileName)
	assert.Equal(t, "image/png", got.MimeType)
	assert.Equal(t, int64(1234), got.Size)
	assert.Equal(t, creado, got.CreatedAt)
}

func TestHistoryToEntity_TrasladaAmbosTiposDeCambio(t *testing.T) {
	creado := time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC)

	estado := historyToEntity(&models.TicketStatusHistory{
		Model:    gorm.Model{ID: 1, CreatedAt: creado},
		TicketID: 5, ChangeType: "status", FromStatus: "open", ToStatus: "closed",
		ChangedByID: 42, ChangedBy: models.User{Name: "Ana"}, Note: "duplicado",
	})
	assert.Equal(t, uint(1), estado.ID)
	assert.Equal(t, uint(5), estado.TicketID)
	assert.Equal(t, "status", estado.ChangeType)
	assert.Equal(t, "open", estado.FromStatus)
	assert.Equal(t, "closed", estado.ToStatus)
	assert.Empty(t, estado.FromArea)
	assert.Empty(t, estado.ToArea)
	assert.Equal(t, uint(42), estado.ChangedByID)
	assert.Equal(t, "Ana", estado.ChangedByName)
	assert.Equal(t, "duplicado", estado.Note)
	assert.Equal(t, creado, estado.CreatedAt)

	area := historyToEntity(&models.TicketStatusHistory{
		Model: gorm.Model{ID: 2}, ChangeType: "area",
		FromArea: "soporte", ToArea: "desarrollo",
	})
	assert.Equal(t, "area", area.ChangeType)
	assert.Equal(t, "soporte", area.FromArea)
	assert.Equal(t, "desarrollo", area.ToArea)
	assert.Empty(t, area.FromStatus)
	assert.Empty(t, area.ToStatus)
}

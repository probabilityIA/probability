package repository

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestAddComment_GuardaYReleeElComentario(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_comments"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 5, 0, "ya lo revise", true).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(31))
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_comments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "user_id", "body", "is_internal"}).
			AddRow(31, 5, 0, "ya lo revise", true))

	got, err := New(base).AddComment(context.Background(), dtos.CreateCommentDTO{
		TicketID: 5, Body: "ya lo revise", IsInternal: true,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(31), got.ID)
	assert.Equal(t, uint(5), got.TicketID)
	assert.Equal(t, "ya lo revise", got.Body)
	assert.True(t, got.IsInternal)
	base.SinPendientes(t)
}

func TestAddComment_ErrorAlInsertar_NoRelee(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_comments"`).WillReturnError(stderrors.New("fk violation"))

	got, err := New(base).AddComment(context.Background(), dtos.CreateCommentDTO{TicketID: 404, Body: "hola"})

	require.Error(t, err)
	assert.Nil(t, got)
	base.SinPendientes(t)
}

func TestAddComment_ErrorAlReleer_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_comments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(31))
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_comments"`).WillReturnError(gorm.ErrRecordNotFound)

	got, err := New(base).AddComment(context.Background(), dtos.CreateCommentDTO{TicketID: 5, Body: "hola"})

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestListComments_AlClienteLeOcultaLosInternos(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_comments" WHERE ticket_id = \$1 AND is_internal = \$2`).
		WithArgs(5, false).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "body"}).AddRow(1, 5, "publico"))
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := New(base).ListComments(context.Background(), 5, false)

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, "publico", items[0].Body)
	base.SinPendientes(t)
}

func TestListComments_AlEquipoLeMuestraTodos(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_comments" WHERE ticket_id = \$1 AND "ticket_comments"."deleted_at" IS NULL ORDER BY created_at ASC`).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "body", "is_internal"}).
			AddRow(1, 5, "publico", false).
			AddRow(2, 5, "interno", true))
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := New(base).ListComments(context.Background(), 5, true)

	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.True(t, items[1].IsInternal)
	base.SinPendientes(t)
}

func TestListComments_SinComentarios_DevuelveListaVaciaNoNula(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_comments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := New(base).ListComments(context.Background(), 5, true)

	require.NoError(t, err)
	assert.NotNil(t, items)
	assert.Empty(t, items)
}

func TestListComments_ErrorDeLaBase_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_comments"`).WillReturnError(stderrors.New("timeout"))

	items, err := New(base).ListComments(context.Background(), 5, true)

	require.Error(t, err)
	assert.Nil(t, items)
}

func TestAddAttachment_GuardaYReleeElAdjunto(t *testing.T) {
	comentario := uint(31)
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(77))
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "comment_id", "file_url", "file_name", "mime_type", "size"}).
			AddRow(77, 5, 31, "s3://bucket/a.png", "a.png", "image/png", 1234))

	got, err := New(base).AddAttachment(context.Background(), dtos.CreateAttachmentDTO{
		TicketID: 5, CommentID: &comentario, FileURL: "s3://bucket/a.png",
		FileName: "a.png", MimeType: "image/png", Size: 1234,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(77), got.ID)
	assert.Equal(t, uint(5), got.TicketID)
	require.NotNil(t, got.CommentID)
	assert.Equal(t, uint(31), *got.CommentID)
	assert.Equal(t, "s3://bucket/a.png", got.FileURL)
	assert.Equal(t, int64(1234), got.Size)
	base.SinPendientes(t)
}

func TestAddAttachment_ErrorAlInsertar_NoRelee(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_attachments"`).WillReturnError(stderrors.New("fk violation"))

	got, err := New(base).AddAttachment(context.Background(), dtos.CreateAttachmentDTO{TicketID: 404})

	require.Error(t, err)
	assert.Nil(t, got)
	base.SinPendientes(t)
}

func TestAddAttachment_ErrorAlReleer_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(77))
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).WillReturnError(gorm.ErrRecordNotFound)

	got, err := New(base).AddAttachment(context.Background(), dtos.CreateAttachmentDTO{TicketID: 5})

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestGetAttachment_Existente(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments" WHERE id = \$1`).
		WithArgs(77, 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "uploaded_by_id", "file_url"}).
			AddRow(77, 5, 42, "s3://bucket/a.png"))
	base.Mock.ExpectQuery(`SELECT \* FROM "users" WHERE "users"."id" = \$1`).
		WithArgs(42).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(42, "Ana"))

	got, err := New(base).GetAttachment(context.Background(), 77)

	require.NoError(t, err)
	assert.Equal(t, uint(77), got.ID)
	assert.Equal(t, uint(42), got.UploadedByID,
		"el caso de uso usa este id para decidir si el solicitante puede borrarlo")
	assert.Equal(t, "Ana", got.UploadedByName, "el nombre viene del preload de usuario")
	base.SinPendientes(t)
}

func TestGetAttachment_Inexistente_DevuelveErrorDeDominio(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).WillReturnError(gorm.ErrRecordNotFound)

	got, err := New(base).GetAttachment(context.Background(), 404)

	assert.ErrorIs(t, err, dom.ErrAttachmentNotFound)
	assert.Nil(t, got)
}

func TestGetAttachment_ErrorRealDeLaBase_NoSeDisfraza(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).WillReturnError(stderrors.New("timeout"))

	got, err := New(base).GetAttachment(context.Background(), 77)

	require.Error(t, err)
	assert.NotErrorIs(t, err, dom.ErrAttachmentNotFound)
	assert.Nil(t, got)
}

func TestDeleteAttachment_BorraDefinitivamente(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`DELETE FROM "ticket_attachments" WHERE "ticket_attachments"."id" = \$1`).
		WithArgs(77).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := New(base).DeleteAttachment(context.Background(), 77)

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestDeleteAttachment_ErrorDeLaBase_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`DELETE FROM "ticket_attachments"`).WillReturnError(stderrors.New("timeout"))

	err := New(base).DeleteAttachment(context.Background(), 77)

	require.Error(t, err)
}

func TestListAttachments_OrdenadosPorFecha(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments" WHERE ticket_id = \$1 AND "ticket_attachments"."deleted_at" IS NULL ORDER BY created_at ASC`).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "file_name"}).
			AddRow(1, 5, "a.png").
			AddRow(2, 5, "b.pdf"))

	items, err := New(base).ListAttachments(context.Background(), 5)

	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "a.png", items[0].FileName)
	assert.Equal(t, "b.pdf", items[1].FileName)
	base.SinPendientes(t)
}

func TestListAttachments_SinAdjuntos_DevuelveListaVaciaNoNula(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := New(base).ListAttachments(context.Background(), 5)

	require.NoError(t, err)
	assert.NotNil(t, items)
	assert.Empty(t, items)
}

func TestListAttachments_ErrorDeLaBase_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_attachments"`).WillReturnError(stderrors.New("timeout"))

	items, err := New(base).ListAttachments(context.Background(), 5)

	require.Error(t, err)
	assert.Nil(t, items)
}

func TestAddHistory_GuardaElCambioDeEstado(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_status_history"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			5, "status", "open", "closed", "", "", 42, "duplicado").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := New(base).AddHistory(context.Background(), 5, "open", "closed", 42, "duplicado")

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestAddAreaHistory_GuardaElCambioDeArea(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "ticket_status_history"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			5, "area", "", "", "soporte", "desarrollo", 42, "es un bug").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := New(base).AddAreaHistory(context.Background(), 5, "soporte", "desarrollo", 42, "es un bug")

	require.NoError(t, err,
		"el cambio de area viaja en las columnas de area, no en las de estado")
	base.SinPendientes(t)
}

func TestAddHistory_ErroresDeLaBase_SePropagan(t *testing.T) {
	t.Run("estado", func(t *testing.T) {
		base := testkit.NewDB(t)
		base.Mock.ExpectQuery(`INSERT INTO "ticket_status_history"`).WillReturnError(stderrors.New("fk violation"))

		err := New(base).AddHistory(context.Background(), 404, "open", "closed", 42, "")

		require.Error(t, err)
	})

	t.Run("area", func(t *testing.T) {
		base := testkit.NewDB(t)
		base.Mock.ExpectQuery(`INSERT INTO "ticket_status_history"`).WillReturnError(stderrors.New("fk violation"))

		err := New(base).AddAreaHistory(context.Background(), 404, "soporte", "desarrollo", 42, "")

		require.Error(t, err)
	})
}

func TestListHistory_OrdenadoCronologicamente(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_status_history" WHERE ticket_id = \$1 AND "ticket_status_history"."deleted_at" IS NULL ORDER BY created_at ASC`).
		WithArgs(5).
		WillReturnRows(sqlmock.NewRows([]string{"id", "ticket_id", "change_type", "from_status", "to_status"}).
			AddRow(1, 5, "status", "", "open").
			AddRow(2, 5, "status", "open", "closed"))

	items, err := New(base).ListHistory(context.Background(), 5)

	require.NoError(t, err)
	require.Len(t, items, 2)
	assert.Equal(t, "open", items[0].ToStatus)
	assert.Equal(t, "closed", items[1].ToStatus)
	base.SinPendientes(t)
}

func TestListHistory_SinHistorial_DevuelveListaVaciaNoNula(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_status_history"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, err := New(base).ListHistory(context.Background(), 5)

	require.NoError(t, err)
	assert.NotNil(t, items)
	assert.Empty(t, items)
}

func TestListHistory_ErrorDeLaBase_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "ticket_status_history"`).WillReturnError(stderrors.New("timeout"))

	items, err := New(base).ListHistory(context.Background(), 5)

	require.Error(t, err)
	assert.Nil(t, items)
}

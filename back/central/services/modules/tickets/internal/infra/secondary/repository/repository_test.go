package repository

import (
	"context"
	stderrors "errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/entities"
	dom "github.com/secamc93/probability/back/central/services/modules/tickets/internal/domain/errors"
	"github.com/secamc93/probability/back/central/shared/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func esperarConteosDeGetByID(base *testkit.DBMock, comentarios, adjuntos int64) {
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "ticket_comments"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(comentarios))
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "ticket_attachments"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(adjuntos))
}

func filaDeTicket(id uint, estado string) *sqlmock.Rows {
	return sqlmock.NewRows([]string{"id", "code", "title", "status", "area", "created_by_id"}).
		AddRow(id, "TKT-000001", "No carga", estado, "soporte", 0)
}

func TestNextCode_CuentaDesdeElMaximoIDIncluyendoBorrados(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT COALESCE\(MAX\(id\), 0\) FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(41))

	code, err := New(base).NextCode(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "TKT-000042", code,
		"el consecutivo sale del maximo id sin filtrar borrados, para no reutilizar codigos")
	base.SinPendientes(t)
}

func TestNextCode_BaseVacia_EmpiezaEnUno(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT COALESCE\(MAX\(id\), 0\) FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"coalesce"}).AddRow(0))

	code, err := New(base).NextCode(context.Background())

	require.NoError(t, err)
	assert.Equal(t, "TKT-000001", code)
}

func TestNextCode_ErrorDeLaBase_SePropaga(t *testing.T) {
	fallo := stderrors.New("connection refused")
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT COALESCE\(MAX\(id\), 0\) FROM "tickets"`).WillReturnError(fallo)

	code, err := New(base).NextCode(context.Background())

	require.Error(t, err)
	assert.Empty(t, code)
}

func TestUserExists_IgnoraLosUsuariosBorrados(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "user" WHERE id = \$1 AND deleted_at IS NULL`).
		WithArgs(7).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	existe, err := New(base).UserExists(context.Background(), 7)

	require.NoError(t, err)
	assert.True(t, existe)
	base.SinPendientes(t)
}

func TestUserExists_SinFilas_DevuelveFalso(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "user"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	existe, err := New(base).UserExists(context.Background(), 404)

	require.NoError(t, err)
	assert.False(t, existe)
}

func TestUserExists_ErrorDeLaBase_NoAfirmaQueExiste(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "user"`).WillReturnError(stderrors.New("timeout"))

	existe, err := New(base).UserExists(context.Background(), 7)

	require.Error(t, err)
	assert.False(t, existe)
}

func TestCreate_DevuelveElIDYLasFechasQueAsignoLaBase(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(77))

	ticket := &entities.Ticket{Code: "TKT-000077", Title: "t", Description: "d", Status: "open"}
	got, err := New(base).Create(context.Background(), ticket)

	require.NoError(t, err)
	assert.Equal(t, uint(77), got.ID)
	assert.Equal(t, uint(77), ticket.ID, "el id vuelve sobre la entidad recibida")
	assert.False(t, got.CreatedAt.IsZero())
	base.SinPendientes(t)
}

func TestCreate_ErrorDeLaBase_NoDevuelveTicket(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`INSERT INTO "tickets"`).WillReturnError(stderrors.New("unique violation"))

	got, err := New(base).Create(context.Background(), &entities.Ticket{Code: "TKT-000077"})

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestGetByID_TraeElTicketConSusConteos(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE id = \$1`).
		WithArgs(5, 1).
		WillReturnRows(filaDeTicket(5, "open"))
	esperarConteosDeGetByID(base, 3, 2)

	got, err := New(base).GetByID(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), got.ID)
	assert.Equal(t, "TKT-000001", got.Code)
	assert.Equal(t, "open", got.Status)
	assert.Equal(t, int64(3), got.CommentsCount)
	assert.Equal(t, int64(2), got.AttachmentsCount)
	base.SinPendientes(t)
}

func TestGetByID_SinFila_DevuelveErrorDeDominio(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets"`).WillReturnError(gorm.ErrRecordNotFound)

	got, err := New(base).GetByID(context.Background(), 404)

	assert.ErrorIs(t, err, dom.ErrTicketNotFound,
		"el caso de uso espera el error del dominio, no el de gorm")
	assert.Nil(t, got)
}

func TestGetByID_ErrorRealDeLaBase_NoSeDisfrazaDeNoEncontrado(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets"`).WillReturnError(stderrors.New("connection refused"))

	got, err := New(base).GetByID(context.Background(), 5)

	require.Error(t, err)
	assert.NotErrorIs(t, err, dom.ErrTicketNotFound)
	assert.Nil(t, got)
}

func TestList_UsuarioSinNegocio_NoVeNada(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE 1 = 0`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE 1 = 0`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	items, total, err := New(base).List(context.Background(),
		dtos.ListTicketsParams{Page: 1, PageSize: 10, IsSuperAdmin: false, BusinessID: nil})

	require.NoError(t, err)
	assert.Empty(t, items)
	assert.Zero(t, total)
	base.SinPendientes(t)
}

func TestList_UsuarioDeNegocio_FiltraPorSuNegocio(t *testing.T) {
	negocio := uint(26)
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE business_id = \$1`).
		WithArgs(26).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE business_id = \$1`).
		WithArgs(26, 10).
		WillReturnRows(filaDeTicket(1, "open"))

	items, total, err := New(base).List(context.Background(),
		dtos.ListTicketsParams{Page: 1, PageSize: 10, BusinessID: &negocio})

	require.NoError(t, err)
	require.Len(t, items, 1)
	assert.Equal(t, uint(1), items[0].ID)
	assert.Equal(t, int64(1), total)
	base.SinPendientes(t)
}

func TestList_SuperAdminSinNegocio_VeTodo(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE "tickets"."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE "tickets"."deleted_at" IS NULL`).
		WillReturnRows(filaDeTicket(1, "open").AddRow(2, "TKT-000002", "otro", "closed", "soporte", 0))

	items, total, err := New(base).List(context.Background(),
		dtos.ListTicketsParams{Page: 1, PageSize: 10, IsSuperAdmin: true})

	require.NoError(t, err)
	assert.Len(t, items, 2)
	assert.Equal(t, int64(2), total)
	base.SinPendientes(t)
}

func TestList_AplicaTodosLosFiltros(t *testing.T) {
	negocio, creador, responsable, sprint := uint(26), uint(42), uint(7), uint(9)
	base := testkit.NewDB(t)
	patron := `SELECT count\(\*\) FROM "tickets" WHERE business_id = \$1 AND created_by_id = \$2 ` +
		`AND assigned_to_id = \$3 AND sprint_id = \$4 AND \(created_by_id = \$5 OR assigned_to_id = \$6\) ` +
		`AND status IN \(\$7\) AND priority IN \(\$8\) AND type IN \(\$9\) AND area IN \(\$10\) ` +
		`AND source = \$11 AND escalated_to_dev = \$12 AND \(title ILIKE`
	base.Mock.ExpectQuery(patron).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
		Page: 1, PageSize: 10, IsSuperAdmin: true,
		BusinessID: &negocio, CreatedByID: &creador, AssignedToID: &responsable, SprintID: &sprint,
		OnlyMine: true, UserID: 42,
		Status: []string{"open"}, Priority: []string{"high"}, Type: []string{"bug"}, Area: []string{"soporte"},
		Source: "business", EscalatedOnly: true, Search: "  factura  ",
	})

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestList_SprintNoneGanaSobreElSprintConcreto(t *testing.T) {
	sprint := uint(9)
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE sprint_id IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE sprint_id IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
		Page: 1, PageSize: 10, IsSuperAdmin: true, SprintNone: true, SprintID: &sprint,
	})

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestList_OnlyMineSinUsuario_NoFiltra(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE "tickets"."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
		Page: 1, PageSize: 10, IsSuperAdmin: true, OnlyMine: true, UserID: 0,
	})

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestList_BusquedaEnBlanco_NoAgregaElLike(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets" WHERE "tickets"."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
		Page: 1, PageSize: 10, IsSuperAdmin: true, Search: "   ",
	})

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestList_OrdenamientoPermitido(t *testing.T) {
	casos := []struct {
		sortBy, sortOrder string
		esperado          string
	}{
		{"", "", `ORDER BY created_at desc`},
		{"inventado", "asc", `ORDER BY created_at asc`},
		{"  CODE  ", "  ASC  ", `ORDER BY code asc`},
		{"priority", "descendente", `ORDER BY priority desc`},
		{"due_date", "asc", `ORDER BY due_date asc`},
	}

	for _, tc := range casos {
		base := testkit.NewDB(t)
		base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		base.Mock.ExpectQuery(tc.esperado).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
			Page: 1, PageSize: 10, IsSuperAdmin: true, SortBy: tc.sortBy, SortOrder: tc.sortOrder,
		})

		require.NoError(t, err, "sort_by %q", tc.sortBy)
		base.SinPendientes(t)
	}
}

func TestList_ColumnaDeOrdenNoPermitida_NoLlegaAlSQL(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`ORDER BY created_at asc`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
		Page: 1, PageSize: 10, IsSuperAdmin: true,
		SortBy: "id; DROP TABLE tickets", SortOrder: "asc",
	})

	require.NoError(t, err, "una columna fuera de la lista blanca cae al orden por defecto")
	base.SinPendientes(t)
}

func TestList_CalculaElOffsetDeLaPagina(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
	base.Mock.ExpectQuery(`LIMIT \$1 OFFSET \$2`).
		WithArgs(25, 50).
		WillReturnRows(sqlmock.NewRows([]string{"id"}))

	_, _, err := New(base).List(context.Background(), dtos.ListTicketsParams{
		Page: 3, PageSize: 25, IsSuperAdmin: true,
	})

	require.NoError(t, err)
	base.SinPendientes(t)
}

func TestList_ErroresDeLaBase_SePropagan(t *testing.T) {
	fallo := stderrors.New("timeout")

	t.Run("al contar", func(t *testing.T) {
		base := testkit.NewDB(t)
		base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).WillReturnError(fallo)

		items, total, err := New(base).List(context.Background(),
			dtos.ListTicketsParams{Page: 1, PageSize: 10, IsSuperAdmin: true})

		require.Error(t, err)
		assert.Nil(t, items)
		assert.Zero(t, total)
	})

	t.Run("al traer la pagina", func(t *testing.T) {
		base := testkit.NewDB(t)
		base.Mock.ExpectQuery(`SELECT count\(\*\) FROM "tickets"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(9))
		base.Mock.ExpectQuery(`SELECT \* FROM "tickets"`).WillReturnError(fallo)

		items, total, err := New(base).List(context.Background(),
			dtos.ListTicketsParams{Page: 1, PageSize: 10, IsSuperAdmin: true})

		require.Error(t, err)
		assert.Nil(t, items)
		assert.Zero(t, total, "si la pagina fallo no se devuelve el total a medias")
	})
}

func TestUpdate_EscribeYReleeElTicket(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`UPDATE "tickets" SET`).
		WillReturnResult(sqlmock.NewResult(0, 1))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE id = \$1`).
		WithArgs(5, 1).
		WillReturnRows(filaDeTicket(5, "closed"))
	esperarConteosDeGetByID(base, 0, 0)

	got, err := New(base).Update(context.Background(), 5, map[string]any{"status": "closed"})

	require.NoError(t, err)
	assert.Equal(t, "closed", got.Status)
	base.SinPendientes(t)
}

func TestUpdate_SinFilasAfectadasPeroElTicketExiste_NoEsError(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`UPDATE "tickets" SET`).WillReturnResult(sqlmock.NewResult(0, 0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE id = \$1`).
		WillReturnRows(filaDeTicket(5, "open"))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE id = \$1`).
		WillReturnRows(filaDeTicket(5, "open"))
	esperarConteosDeGetByID(base, 0, 0)

	got, err := New(base).Update(context.Background(), 5, map[string]any{"status": "open"})

	require.NoError(t, err,
		"actualizar con el mismo valor no afecta filas y no debe parecer un ticket inexistente")
	assert.Equal(t, uint(5), got.ID)
	base.SinPendientes(t)
}

func TestUpdate_SinFilasAfectadasYTicketInexistente_DevuelveErrorDeDominio(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`UPDATE "tickets" SET`).WillReturnResult(sqlmock.NewResult(0, 0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE id = \$1`).
		WillReturnError(gorm.ErrRecordNotFound)

	got, err := New(base).Update(context.Background(), 404, map[string]any{"status": "open"})

	assert.ErrorIs(t, err, dom.ErrTicketNotFound)
	assert.Nil(t, got)
	base.SinPendientes(t)
}

func TestUpdate_SinFilasAfectadasYErrorReal_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`UPDATE "tickets" SET`).WillReturnResult(sqlmock.NewResult(0, 0))
	base.Mock.ExpectQuery(`SELECT \* FROM "tickets" WHERE id = \$1`).
		WillReturnError(stderrors.New("connection refused"))

	got, err := New(base).Update(context.Background(), 5, map[string]any{"status": "open"})

	require.Error(t, err)
	assert.NotErrorIs(t, err, dom.ErrTicketNotFound)
	assert.Nil(t, got)
}

func TestUpdate_ErrorAlEscribir_NoRelee(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`UPDATE "tickets" SET`).WillReturnError(stderrors.New("check violation"))

	got, err := New(base).Update(context.Background(), 5, map[string]any{"status": "raro"})

	require.Error(t, err)
	assert.Nil(t, got)
	base.SinPendientes(t)
}

func TestDelete_BorraDefinitivamente(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`DELETE FROM "tickets" WHERE "tickets"."id" = \$1`).
		WithArgs(5).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err := New(base).Delete(context.Background(), 5)

	require.NoError(t, err,
		"el borrado de tickets es fisico (Unscoped), no un soft delete")
	base.SinPendientes(t)
}

func TestDelete_ErrorDeLaBase_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectExec(`DELETE FROM "tickets"`).WillReturnError(stderrors.New("fk violation"))

	err := New(base).Delete(context.Background(), 5)

	require.Error(t, err)
}

func TestFindSprintName_Existente(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT "name" FROM "sprints" WHERE id = \$1`).
		WithArgs(9, 1).
		WillReturnRows(sqlmock.NewRows([]string{"name"}).AddRow("Sprint de octubre"))

	nombre, encontrado, err := New(base).FindSprintName(context.Background(), 9)

	require.NoError(t, err)
	assert.True(t, encontrado)
	assert.Equal(t, "Sprint de octubre", nombre)
	base.SinPendientes(t)
}

func TestFindSprintName_Inexistente_NoEsError(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT "name" FROM "sprints"`).WillReturnError(gorm.ErrRecordNotFound)

	nombre, encontrado, err := New(base).FindSprintName(context.Background(), 404)

	require.NoError(t, err, "que el sprint no exista lo decide el caso de uso, no es un fallo tecnico")
	assert.False(t, encontrado)
	assert.Empty(t, nombre)
}

func TestFindSprintName_ErrorRealDeLaBase_SePropaga(t *testing.T) {
	base := testkit.NewDB(t)
	base.Mock.ExpectQuery(`SELECT "name" FROM "sprints"`).WillReturnError(stderrors.New("timeout"))

	nombre, encontrado, err := New(base).FindSprintName(context.Background(), 9)

	require.Error(t, err)
	assert.False(t, encontrado)
	assert.Empty(t, nombre)
}

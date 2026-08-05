package app

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/syncruns/internal/mocks"
)

func nuevoCaso(repo *mocks.RepositoryMock) domain.IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func corridaValida() *domain.SyncRun {
	return &domain.SyncRun{
		BusinessID:    26,
		IntegrationID: 197,
		Kind:          domain.KindInventory,
		Status:        domain.StatusCompleted,
		Total:         10,
		Updated:       8,
	}
}

func TestRecord_GuardaLaCorridaCuandoLaIntegracionEsDelNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	err := nuevoCaso(repo).Record(context.Background(), corridaValida())

	require.NoError(t, err)
	require.Len(t, repo.ConsultasDueno, 1)
	assert.Equal(t, mocks.ConsultaDueno{IntegrationID: 197, BusinessID: 26}, repo.ConsultasDueno[0],
		"la propiedad se valida antes de escribir: un integration_id de otro negocio no puede ensuciar su historial de sincronizaciones")
	assert.Len(t, repo.Guardados, 1)
}

func TestRecord_RechazaUnaIntegracionDeOtroNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		IntegrationBelongsToBusinessFn: func(ctx context.Context, integrationID, businessID uint) (bool, error) {
			return false, nil
		},
	}

	err := nuevoCaso(repo).Record(context.Background(), corridaValida())

	assert.ErrorIs(t, err, domain.ErrIntegrationNotOwned)
	assert.Empty(t, repo.Guardados)
}

func TestRecord_RechazaCorridasSinNegocioOSinIntegracion(t *testing.T) {
	casos := map[string]*domain.SyncRun{
		"corrida nula":    nil,
		"sin business":    {BusinessID: 0, IntegrationID: 197},
		"sin integracion": {BusinessID: 26, IntegrationID: 0},
		"ambos en cero":   {},
	}

	for nombre, run := range casos {
		t.Run(nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}

			err := nuevoCaso(repo).Record(context.Background(), run)

			assert.ErrorIs(t, err, domain.ErrIntegrationNotOwned)
			assert.Empty(t, repo.ConsultasDueno, "ni siquiera se consulta la base")
			assert.Empty(t, repo.Guardados)
		})
	}
}

func TestRecord_PropagaElErrorDeLaValidacionDePropiedad(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		IntegrationBelongsToBusinessFn: func(ctx context.Context, integrationID, businessID uint) (bool, error) {
			return false, fallo
		},
	}

	err := nuevoCaso(repo).Record(context.Background(), corridaValida())

	assert.ErrorIs(t, err, fallo,
		"si no se puede comprobar la propiedad no se guarda: falla cerrado")
	assert.Empty(t, repo.Guardados)
}

func TestRecord_PropagaElErrorDeGuardado(t *testing.T) {
	fallo := errors.New("insert fallido")
	repo := &mocks.RepositoryMock{
		SaveFn: func(ctx context.Context, run *domain.SyncRun) error { return fallo },
	}

	err := nuevoCaso(repo).Record(context.Background(), corridaValida())

	assert.ErrorIs(t, err, fallo)
}

func TestRecord_NormalizaLaCorridaAntesDeGuardarla(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	run := &domain.SyncRun{BusinessID: 26, IntegrationID: 197, Kind: "loquesea", Status: "loquesea"}

	require.NoError(t, nuevoCaso(repo).Record(context.Background(), run))

	guardada := repo.Guardados[0]
	assert.Equal(t, domain.KindInventory, guardada.Kind)
	assert.Equal(t, domain.StatusCompleted, guardada.Status)
	assert.False(t, guardada.StartedAt.IsZero())
	assert.NotNil(t, guardada.FinishedAt)
}

func TestListLast_DevuelveElHistorialDelNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListLastByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.SyncRun, error) {
			return []domain.SyncRun{{ID: 1, BusinessID: businessID}}, nil
		},
	}

	corridas, err := nuevoCaso(repo).ListLast(context.Background(), 26)

	require.NoError(t, err)
	require.Len(t, corridas, 1)
	assert.EqualValues(t, 26, repo.ConsultasLast[0])
}

func TestListLast_NoValidaElNegocioYDelegaElFiltroAlRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := nuevoCaso(repo).ListLast(context.Background(), 0)

	require.NoError(t, err)
	assert.Equal(t, []uint{0}, repo.ConsultasLast,
		"business_id 0 (super admin) llega tal cual al repositorio: el aislamiento depende de que el WHERE del query no devuelva nada, no de una guarda aqui")
}

func TestListLast_PropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListLastByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.SyncRun, error) {
			return nil, fallo
		},
	}

	_, err := nuevoCaso(repo).ListLast(context.Background(), 26)

	assert.ErrorIs(t, err, fallo)
}

func TestListDetail_ValidaLaPropiedadAntesDeConsultar(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := nuevoCaso(repo).ListDetail(context.Background(), domain.DetailQuery{
		BusinessID: 26, IntegrationID: 197,
	})

	require.NoError(t, err)
	assert.Len(t, repo.ConsultasDueno, 1)
	assert.Len(t, repo.ConsultasDetail, 1)
}

func TestListDetail_RechazaUnaIntegracionAjena(t *testing.T) {
	repo := &mocks.RepositoryMock{
		IntegrationBelongsToBusinessFn: func(ctx context.Context, integrationID, businessID uint) (bool, error) {
			return false, nil
		},
	}

	pagina, err := nuevoCaso(repo).ListDetail(context.Background(), domain.DetailQuery{
		BusinessID: 26, IntegrationID: 999,
	})

	assert.ErrorIs(t, err, domain.ErrIntegrationNotOwned)
	assert.Nil(t, pagina)
	assert.Empty(t, repo.ConsultasDetail,
		"el detalle expone SKUs del catalogo, por eso no se consulta sin confirmar el dueno")
}

func TestListDetail_RechazaConsultasSinNegocioOSinIntegracion(t *testing.T) {
	casos := map[string]domain.DetailQuery{
		"sin business":    {IntegrationID: 197},
		"sin integracion": {BusinessID: 26},
		"vacia":           {},
	}

	for nombre, query := range casos {
		t.Run(nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}

			_, err := nuevoCaso(repo).ListDetail(context.Background(), query)

			assert.ErrorIs(t, err, domain.ErrIntegrationNotOwned)
			assert.Empty(t, repo.ConsultasDueno)
		})
	}
}

func TestListDetail_PropagaElErrorDeLaValidacionDePropiedad(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		IntegrationBelongsToBusinessFn: func(ctx context.Context, integrationID, businessID uint) (bool, error) {
			return true, fallo
		},
	}

	_, err := nuevoCaso(repo).ListDetail(context.Background(), domain.DetailQuery{
		BusinessID: 26, IntegrationID: 197,
	})

	assert.ErrorIs(t, err, fallo,
		"aunque el repositorio diga true, el error manda: falla cerrado")
	assert.Empty(t, repo.ConsultasDetail)
}

func TestListDetail_NormalizaLaPaginacionAntesDeConsultar(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := nuevoCaso(repo).ListDetail(context.Background(), domain.DetailQuery{
		BusinessID: 26, IntegrationID: 197, Page: 0, PageSize: 5000, Kind: "loquesea",
	})

	require.NoError(t, err)
	consulta := repo.ConsultasDetail[0]
	assert.Equal(t, 1, consulta.Page)
	assert.Equal(t, domain.MaxPageSize, consulta.PageSize,
		"el tope de 200 protege la memoria del backend: el detalle puede traer miles de SKUs")
	assert.Equal(t, domain.KindInventory, consulta.Kind)
}

func TestListDetail_PropagaElErrorDelRepositorio(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListDetailFn: func(ctx context.Context, query domain.DetailQuery) (*domain.DetailPage, error) {
			return nil, fallo
		},
	}

	_, err := nuevoCaso(repo).ListDetail(context.Background(), domain.DetailQuery{
		BusinessID: 26, IntegrationID: 197,
	})

	assert.ErrorIs(t, err, fallo)
}

func TestSyncRunNormalize_SoloAceptaProductsComoKindAlternativo(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{domain.KindProducts, domain.KindProducts},
		{domain.KindInventory, domain.KindInventory},
		{"", domain.KindInventory},
		{"PRODUCTS", domain.KindInventory},
		{"pedidos", domain.KindInventory},
	}

	for _, caso := range casos {
		run := &domain.SyncRun{Kind: caso.entrada}
		run.Normalize()
		assert.Equal(t, caso.esperado, run.Kind,
			"todo lo que no sea exactamente 'products' cae a inventory, incluida la mayuscula")
	}
}

func TestSyncRunNormalize_SoloAceptaFailedComoEstadoAlternativo(t *testing.T) {
	casos := []struct {
		entrada  string
		esperado string
	}{
		{domain.StatusFailed, domain.StatusFailed},
		{domain.StatusCompleted, domain.StatusCompleted},
		{"", domain.StatusCompleted},
		{"error", domain.StatusCompleted},
	}

	for _, caso := range casos {
		run := &domain.SyncRun{Status: caso.entrada}
		run.Normalize()
		assert.Equal(t, caso.esperado, run.Status,
			"un estado escrito mal se guarda como completed: una corrida que fallo puede quedar reportada como exitosa")
	}
}

func TestSyncRunNormalize_RellenaLasFechasQueFalten(t *testing.T) {
	antes := time.Now()
	run := &domain.SyncRun{}

	run.Normalize()

	assert.False(t, run.StartedAt.IsZero())
	assert.True(t, run.StartedAt.After(antes) || run.StartedAt.Equal(antes))
	require.NotNil(t, run.FinishedAt)
}

func TestSyncRunNormalize_RespetaLasFechasQueYaVienen(t *testing.T) {
	inicio := time.Date(2026, 8, 1, 10, 0, 0, 0, time.UTC)
	fin := inicio.Add(time.Minute)
	run := &domain.SyncRun{StartedAt: inicio, FinishedAt: &fin}

	run.Normalize()

	assert.Equal(t, inicio, run.StartedAt)
	assert.Equal(t, fin, *run.FinishedAt)
}

func TestSyncRunNormalize_RecortaElMensajeA500Caracteres(t *testing.T) {
	run := &domain.SyncRun{Message: strings.Repeat("x", 600)}

	run.Normalize()

	assert.Len(t, run.Message, 500,
		"el mensaje suele ser el error crudo del canal; se recorta para que quepa en la columna")
}

func TestSyncRunNormalize_NoTocaUnMensajeCorto(t *testing.T) {
	run := &domain.SyncRun{Message: "todo bien"}

	run.Normalize()

	assert.Equal(t, "todo bien", run.Message)
}

func TestDetailQueryNormalize_AjustaPaginaYTamano(t *testing.T) {
	casos := []struct {
		nombre           string
		page, size       int
		pageEsp, sizeEsp int
	}{
		{"pagina cero", 0, 50, 1, 50},
		{"pagina negativa", -3, 50, 1, 50},
		{"tamano cero", 1, 0, 1, domain.DefaultPageSize},
		{"tamano negativo", 1, -1, 1, domain.DefaultPageSize},
		{"tamano sobre el tope", 1, 201, 1, domain.MaxPageSize},
		{"tamano en el tope", 1, 200, 1, 200},
		{"valores validos", 4, 25, 4, 25},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			q := domain.DetailQuery{Page: caso.page, PageSize: caso.size}
			q.Normalize()
			assert.Equal(t, caso.pageEsp, q.Page)
			assert.Equal(t, caso.sizeEsp, q.PageSize)
		})
	}
}

func TestDetailQueryOffset_CuentaDesdeLaPrimeraPagina(t *testing.T) {
	casos := []struct {
		page, size, esperado int
	}{
		{1, 50, 0},
		{2, 50, 50},
		{3, 25, 50},
	}

	for _, caso := range casos {
		q := domain.DetailQuery{Page: caso.page, PageSize: caso.size}
		assert.Equal(t, caso.esperado, q.Offset(),
			"la pagina 1 arranca en offset 0, no en 50")
	}
}

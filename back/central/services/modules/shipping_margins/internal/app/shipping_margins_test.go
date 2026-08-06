package app

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/shipping_margins/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMarginsUseCase(repo *mocks.RepositoryMock, cache *mocks.CacheWriterMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	if cache == nil {
		return New(repo, nil)
	}
	return New(repo, cache)
}

func margenValido() dtos.CreateShippingMarginDTO {
	return dtos.CreateShippingMarginDTO{
		BusinessID: 26, CarrierCode: "servientrega", CarrierName: "Servientrega",
		MarginAmount: 2000, InsuranceMargin: 500, CODMarginPercent: 3, IsActive: true,
	}
}

func TestCreate_CodigoDeTransportadoraVacio_Rechaza(t *testing.T) {
	for _, code := range []string{"", "   ", "\t", "\n"} {
		repo := &mocks.RepositoryMock{}
		dto := margenValido()
		dto.CarrierCode = code

		got, err := newMarginsUseCase(repo, nil).Create(context.Background(), dto)

		assert.ErrorIs(t, err, domainerrors.ErrInvalidCarrierCode, "code=%q", code)
		assert.Nil(t, got)
		assert.Empty(t, repo.CreatedCarriers)
	}
}

func TestCreate_NormalizaElCodigoAMinusculasYSinEspacios(t *testing.T) {
	casos := []struct {
		entrada string
		want    string
	}{
		{"servientrega", "servientrega"},
		{"SERVIENTREGA", "servientrega"},
		{"  Servientrega  ", "servientrega"},
		{"InterRapidisimo", "interrapidisimo"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			var vistoExists string
			repo := &mocks.RepositoryMock{
				ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
					vistoExists = carrierCode
					return false, nil
				},
			}
			dto := margenValido()
			dto.CarrierCode = tc.entrada

			got, err := newMarginsUseCase(repo, nil).Create(context.Background(), dto)

			require.NoError(t, err)
			assert.Equal(t, tc.want, vistoExists,
				"la busqueda de duplicados usa el codigo normalizado")
			assert.Equal(t, tc.want, got.CarrierCode)
		})
	}
}

func TestCreate_MargenesInvalidos_Rechaza(t *testing.T) {
	casos := []struct {
		nombre string
		ajuste func(*dtos.CreateShippingMarginDTO)
	}{
		{"margen negativo", func(d *dtos.CreateShippingMarginDTO) { d.MarginAmount = -1 }},
		{"seguro negativo", func(d *dtos.CreateShippingMarginDTO) { d.InsuranceMargin = -0.01 }},
		{"cod negativo", func(d *dtos.CreateShippingMarginDTO) { d.CODMarginPercent = -1 }},
		{"cod sobre cien", func(d *dtos.CreateShippingMarginDTO) { d.CODMarginPercent = 100.1 }},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{}
			dto := margenValido()
			tc.ajuste(&dto)

			got, err := newMarginsUseCase(repo, nil).Create(context.Background(), dto)

			assert.ErrorIs(t, err, domainerrors.ErrInvalidMargin)
			assert.Nil(t, got)
			assert.Empty(t, repo.CreatedCarriers)
		})
	}
}

func TestCreate_LosLimitesExactosSeAceptan(t *testing.T) {
	casos := []dtos.CreateShippingMarginDTO{
		{BusinessID: 26, CarrierCode: "x", MarginAmount: 0, InsuranceMargin: 0, CODMarginPercent: 0},
		{BusinessID: 26, CarrierCode: "x", MarginAmount: 0, InsuranceMargin: 0, CODMarginPercent: 100},
	}

	for i, dto := range casos {
		got, err := newMarginsUseCase(nil, nil).Create(context.Background(), dto)
		require.NoError(t, err, "caso %d", i)
		assert.NotNil(t, got)
	}
}

func TestCreate_TransportadoraDuplicada_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	got, err := newMarginsUseCase(repo, nil).Create(context.Background(), margenValido())

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateCarrier)
	assert.Nil(t, got)
	assert.Empty(t, repo.CreatedCarriers,
		"dos margenes para la misma transportadora harian ambiguo el cobro de la guia")
}

func TestCreate_LaUnicidadEsPorNegocio(t *testing.T) {
	var vistoNegocio uint
	repo := &mocks.RepositoryMock{
		ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
			vistoNegocio = businessID
			return false, nil
		},
	}
	dto := margenValido()
	dto.BusinessID = 99

	_, err := newMarginsUseCase(repo, nil).Create(context.Background(), dto)

	require.NoError(t, err)
	assert.Equal(t, uint(99), vistoNegocio,
		"cada negocio tiene su propio margen para la misma transportadora")
}

func TestCreate_ErrorAlVerificarDuplicado_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}

	_, err := newMarginsUseCase(repo, nil).Create(context.Background(), margenValido())

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.CreatedCarriers)
}

func TestCreate_ErrorAlPersistir_NoTocaElCache(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	cache := &mocks.CacheWriterMock{}
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			return nil, dbErr
		},
	}

	got, err := newMarginsUseCase(repo, cache).Create(context.Background(), margenValido())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Empty(t, cache.UpsertCalls, "el cache no debe quedar con un margen que no existe en la base")
}

func TestCreate_Exitoso_CalientaElCache(t *testing.T) {
	cache := &mocks.CacheWriterMock{}

	_, err := newMarginsUseCase(nil, cache).Create(context.Background(), margenValido())

	require.NoError(t, err)
	require.Len(t, cache.UpsertCalls, 1)
	assert.Equal(t, uint(26), cache.UpsertCalls[0].BusinessID)
	assert.Equal(t, "servientrega", cache.UpsertCalls[0].CarrierCode)
}

func TestCreate_FalloDelCache_NoRompeLaCreacion(t *testing.T) {
	cache := &mocks.CacheWriterMock{
		UpsertFn: func(ctx context.Context, m *entities.ShippingMargin) error {
			return stderrors.New("redis caido")
		},
	}

	got, err := newMarginsUseCase(nil, cache).Create(context.Background(), margenValido())

	require.NoError(t, err, "el cache es un acelerador, no la fuente de verdad")
	assert.NotNil(t, got)
}

func TestCreate_SinCache_NoFalla(t *testing.T) {
	got, err := newMarginsUseCase(nil, nil).Create(context.Background(), margenValido())

	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestCreate_MapeaTodosLosCampos(t *testing.T) {
	var creado *entities.ShippingMargin
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			creado = m
			m.ID = 9
			return m, nil
		},
	}

	got, err := newMarginsUseCase(repo, nil).Create(context.Background(), margenValido())

	require.NoError(t, err)
	require.NotNil(t, creado)
	assert.Equal(t, uint(26), creado.BusinessID)
	assert.Equal(t, "Servientrega", creado.CarrierName)
	assert.InDelta(t, 2000.0, creado.MarginAmount, 0.001)
	assert.InDelta(t, 500.0, creado.InsuranceMargin, 0.001)
	assert.InDelta(t, 3.0, creado.CODMarginPercent, 0.001)
	assert.True(t, creado.IsActive)
	assert.Equal(t, uint(9), got.ID)
}

func TestUpdate_MargenesInvalidos_NoLeeNiEscribe(t *testing.T) {
	casos := []struct {
		nombre string
		dto    dtos.UpdateShippingMarginDTO
	}{
		{"margen negativo", dtos.UpdateShippingMarginDTO{ID: 1, BusinessID: 26, MarginAmount: -1}},
		{"seguro negativo", dtos.UpdateShippingMarginDTO{ID: 1, BusinessID: 26, InsuranceMargin: -1}},
		{"cod negativo", dtos.UpdateShippingMarginDTO{ID: 1, BusinessID: 26, CODMarginPercent: -1}},
		{"cod sobre cien", dtos.UpdateShippingMarginDTO{ID: 1, BusinessID: 26, CODMarginPercent: 101}},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			leido := false
			repo := &mocks.RepositoryMock{
				GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
					leido = true
					return &entities.ShippingMargin{ID: id}, nil
				},
			}

			got, err := newMarginsUseCase(repo, nil).Update(context.Background(), tc.dto)

			assert.ErrorIs(t, err, domainerrors.ErrInvalidMargin)
			assert.Nil(t, got)
			assert.False(t, leido, "se valida antes de tocar la base")
		})
	}
}

func TestUpdate_NoPermiteCambiarNegocioNiCodigoDeTransportadora(t *testing.T) {
	original := &entities.ShippingMargin{
		ID: 5, BusinessID: 26, CarrierCode: "servientrega", CarrierName: "Vieja",
		MarginAmount: 1000, IsActive: true,
	}
	var guardado *entities.ShippingMargin
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			return original, nil
		},
		UpdateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			guardado = m
			return m, nil
		},
	}

	_, err := newMarginsUseCase(repo, nil).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 5, BusinessID: 26, CarrierName: "Nueva",
		MarginAmount: 3000, InsuranceMargin: 700, CODMarginPercent: 5, IsActive: true,
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, uint(26), guardado.BusinessID)
	assert.Equal(t, "servientrega", guardado.CarrierCode,
		"el codigo de transportadora es la llave del margen, no se cambia por update")
	assert.Equal(t, "Nueva", guardado.CarrierName)
	assert.InDelta(t, 3000.0, guardado.MarginAmount, 0.001)
	assert.InDelta(t, 5.0, guardado.CODMarginPercent, 0.001)
}

func TestUpdate_LeePorNegocioYPorID(t *testing.T) {
	var vistoNegocio, vistoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			vistoNegocio, vistoID = businessID, id
			return &entities.ShippingMargin{ID: id, BusinessID: businessID}, nil
		},
	}

	_, err := newMarginsUseCase(repo, nil).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 5, BusinessID: 26, IsActive: true,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"el negocio va en la lectura: un negocio no puede editar el margen de otro")
	assert.Equal(t, uint(5), vistoID)
}

func TestUpdate_MargenInexistente_SePropaga(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			return nil, domainerrors.ErrShippingMarginNotFound
		},
		UpdateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			actualizado = true
			return m, nil
		},
	}

	_, err := newMarginsUseCase(repo, nil).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 404, BusinessID: 26, IsActive: true,
	})

	assert.ErrorIs(t, err, domainerrors.ErrShippingMarginNotFound)
	assert.False(t, actualizado)
}

func TestUpdate_ErrorAlPersistir_NoTocaElCache(t *testing.T) {
	dbErr := stderrors.New("deadlock")
	cache := &mocks.CacheWriterMock{}
	repo := &mocks.RepositoryMock{
		UpdateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			return nil, dbErr
		},
	}

	got, err := newMarginsUseCase(repo, cache).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 5, BusinessID: 26, IsActive: true,
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Empty(t, cache.UpsertCalls)
	assert.Empty(t, cache.DeleteCalls)
}

func TestUpdate_QuedaActivo_RefrescaElCache(t *testing.T) {
	cache := &mocks.CacheWriterMock{}
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			return &entities.ShippingMargin{ID: id, BusinessID: businessID, CarrierCode: "servientrega"}, nil
		},
	}

	_, err := newMarginsUseCase(repo, cache).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 5, BusinessID: 26, IsActive: true,
	})

	require.NoError(t, err)
	require.Len(t, cache.UpsertCalls, 1)
	assert.Equal(t, "servientrega", cache.UpsertCalls[0].CarrierCode)
	assert.Empty(t, cache.DeleteCalls)
}

func TestUpdate_QuedaInactivo_LoSacaDelCache(t *testing.T) {
	cache := &mocks.CacheWriterMock{}
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			return &entities.ShippingMargin{ID: id, BusinessID: businessID, CarrierCode: "servientrega"}, nil
		},
	}

	_, err := newMarginsUseCase(repo, cache).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 5, BusinessID: 26, IsActive: false,
	})

	require.NoError(t, err)
	require.Len(t, cache.DeleteCalls, 1)
	assert.Equal(t, uint(26), cache.DeleteCalls[0].BusinessID)
	assert.Equal(t, "servientrega", cache.DeleteCalls[0].CarrierCode,
		"desactivar debe sacar el margen del cache o se seguiria cobrando")
	assert.Empty(t, cache.UpsertCalls)
}

func TestUpdate_FalloDelCache_NoRompeLaActualizacion(t *testing.T) {
	cache := &mocks.CacheWriterMock{
		UpsertFn: func(ctx context.Context, m *entities.ShippingMargin) error {
			return stderrors.New("redis caido")
		},
		DeleteFn: func(ctx context.Context, businessID uint, carrierCode string) error {
			return stderrors.New("redis caido")
		},
	}

	for _, activo := range []bool{true, false} {
		got, err := newMarginsUseCase(nil, cache).Update(context.Background(), dtos.UpdateShippingMarginDTO{
			ID: 5, BusinessID: 26, IsActive: activo,
		})
		require.NoError(t, err, "activo=%v", activo)
		assert.NotNil(t, got)
	}
}

func TestUpdate_SinCache_NoFalla(t *testing.T) {
	got, err := newMarginsUseCase(nil, nil).Update(context.Background(), dtos.UpdateShippingMarginDTO{
		ID: 5, BusinessID: 26, IsActive: true,
	})

	require.NoError(t, err)
	assert.NotNil(t, got)
}

func TestGet_PasaNegocioYID(t *testing.T) {
	var vistoNegocio, vistoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			vistoNegocio, vistoID = businessID, id
			return &entities.ShippingMargin{ID: id, BusinessID: businessID}, nil
		},
	}

	got, err := newMarginsUseCase(repo, nil).Get(context.Background(), 26, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, uint(5), vistoID)
	assert.Equal(t, uint(5), got.ID)
}

func TestGet_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, id uint) (*entities.ShippingMargin, error) {
			return nil, domainerrors.ErrShippingMarginNotFound
		},
	}

	got, err := newMarginsUseCase(repo, nil).Get(context.Background(), 26, 404)

	assert.ErrorIs(t, err, domainerrors.ErrShippingMarginNotFound)
	assert.Nil(t, got)
}

func TestList_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		nombre       string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{"page cero", 0, 20, 1, 20},
		{"page negativo", -5, 20, 1, 20},
		{"pageSize cero", 1, 0, 1, 20},
		{"pageSize negativo", 1, -1, 1, 20},
		{"pageSize sobre el maximo cae al default", 1, 101, 1, 20},
		{"pageSize en el maximo se respeta", 1, 100, 1, 100},
		{"valores validos", 3, 50, 3, 50},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var visto dtos.ListShippingMarginsParams
			repo := &mocks.RepositoryMock{
				ListFn: func(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error) {
					visto = params
					return nil, 0, nil
				},
			}

			_, _, err := newMarginsUseCase(repo, nil).List(context.Background(), dtos.ListShippingMarginsParams{
				BusinessID: 26, Page: tc.page, PageSize: tc.pageSize,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, visto.Page)
			assert.Equal(t, tc.wantPageSize, visto.PageSize)
		})
	}
}

func TestListParams_OffsetSeCalculaDesdeLaPaginaNormalizada(t *testing.T) {
	casos := []struct {
		page     int
		pageSize int
		want     int
	}{
		{1, 20, 0},
		{2, 20, 20},
		{3, 20, 40},
		{1, 100, 0},
		{5, 10, 40},
	}

	for _, tc := range casos {
		p := dtos.ListShippingMarginsParams{Page: tc.page, PageSize: tc.pageSize}
		assert.Equal(t, tc.want, p.Offset(), "page=%d pageSize=%d", tc.page, tc.pageSize)
	}
}

func TestList_PasaNegocioYFiltroDeTransportadora(t *testing.T) {
	var visto dtos.ListShippingMarginsParams
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error) {
			visto = params
			return []entities.ShippingMargin{{ID: 1}}, 1, nil
		},
	}

	got, total, err := newMarginsUseCase(repo, nil).List(context.Background(), dtos.ListShippingMarginsParams{
		BusinessID: 26, CarrierCode: "servientrega", Page: 1, PageSize: 20,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto.BusinessID)
	assert.Equal(t, "servientrega", visto.CarrierCode)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestList_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, params dtos.ListShippingMarginsParams) ([]entities.ShippingMargin, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, total, err := newMarginsUseCase(repo, nil).List(context.Background(), dtos.ListShippingMarginsParams{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.Zero(t, total)
}

func TestEnsureDefaults_NegocioCero_NoHaceNada(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 0)

	require.NoError(t, err)
	assert.Empty(t, repo.CreatedCarriers,
		"business_id cero es super admin, no tiene negocio propio al que sembrarle margenes")
}

func TestEnsureDefaults_SiembraTodasLasTransportadorasDelCatalogo(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 26)

	require.NoError(t, err)
	require.Len(t, repo.CreatedCarriers, len(entities.DefaultCarriers))
	for _, c := range entities.DefaultCarriers {
		assert.Contains(t, repo.CreatedCarriers, c.Code)
	}
}

func TestEnsureDefaults_EsIdempotente(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 26)

	require.NoError(t, err)
	assert.Empty(t, repo.CreatedCarriers,
		"correr el sembrado dos veces no debe duplicar margenes")
}

func TestEnsureDefaults_SoloCreaLasQueFaltan(t *testing.T) {
	yaExisten := map[string]bool{"servientrega": true, "coordinadora": true}
	repo := &mocks.RepositoryMock{
		ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
			return yaExisten[carrierCode], nil
		},
	}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 26)

	require.NoError(t, err)
	assert.Len(t, repo.CreatedCarriers, len(entities.DefaultCarriers)-2)
	assert.NotContains(t, repo.CreatedCarriers, "servientrega")
	assert.NotContains(t, repo.CreatedCarriers, "coordinadora")
	assert.Contains(t, repo.CreatedCarriers, "tcc")
}

func TestEnsureDefaults_NacenEnCeroYActivas(t *testing.T) {
	var creados []*entities.ShippingMargin
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			creados = append(creados, m)
			return m, nil
		},
	}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 26)

	require.NoError(t, err)
	require.NotEmpty(t, creados)
	for _, m := range creados {
		assert.Zero(t, m.MarginAmount, "%s no debe nacer cobrando", m.CarrierCode)
		assert.Zero(t, m.InsuranceMargin)
		assert.True(t, m.IsActive)
		assert.Equal(t, uint(26), m.BusinessID)
		assert.NotEmpty(t, m.CarrierName)
	}
}

func TestEnsureDefaults_CalientaElCachePorCadaUna(t *testing.T) {
	cache := &mocks.CacheWriterMock{}

	err := newMarginsUseCase(nil, cache).EnsureDefaultsForBusiness(context.Background(), 26)

	require.NoError(t, err)
	assert.Len(t, cache.UpsertCalls, len(entities.DefaultCarriers))
}

func TestEnsureDefaults_ErrorAlVerificar_Aborta(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ExistsByCarrierFn: func(ctx context.Context, businessID uint, carrierCode string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, repo.CreatedCarriers)
}

func TestEnsureDefaults_ErrorAlCrear_Aborta(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, m *entities.ShippingMargin) (*entities.ShippingMargin, error) {
			return nil, dbErr
		},
	}

	err := newMarginsUseCase(repo, nil).EnsureDefaultsForBusiness(context.Background(), 26)

	assert.ErrorIs(t, err, dbErr)
	assert.Len(t, repo.CreatedCarriers, 1, "se corta en la primera que falla")
}

func TestEnsureDefaults_FalloDelCache_NoAborta(t *testing.T) {
	cache := &mocks.CacheWriterMock{
		UpsertFn: func(ctx context.Context, m *entities.ShippingMargin) error {
			return stderrors.New("redis caido")
		},
	}
	repo := &mocks.RepositoryMock{}

	err := newMarginsUseCase(repo, cache).EnsureDefaultsForBusiness(context.Background(), 26)

	require.NoError(t, err)
	assert.Len(t, repo.CreatedCarriers, len(entities.DefaultCarriers))
}

func TestDefaultCarriers_NoTieneCodigosRepetidos(t *testing.T) {
	vistos := make(map[string]bool, len(entities.DefaultCarriers))
	for _, c := range entities.DefaultCarriers {
		require.False(t, vistos[c.Code], "codigo repetido en el catalogo: %q", c.Code)
		vistos[c.Code] = true
		assert.Equal(t, c.Code, strings.ToLower(c.Code),
			"los codigos del catalogo deben estar en minusculas para que el Create normalizado los encuentre")
		assert.NotEmpty(t, c.Name)
	}
}

func TestProfitReport_DelegaYPropaga(t *testing.T) {
	esperado := &dtos.ProfitReportResponse{}
	var visto dtos.ProfitReportParams
	repo := &mocks.RepositoryMock{
		ProfitReportFn: func(ctx context.Context, params dtos.ProfitReportParams) (*dtos.ProfitReportResponse, error) {
			visto = params
			return esperado, nil
		},
	}
	params := dtos.ProfitReportParams{BusinessID: 26}

	got, err := newMarginsUseCase(repo, nil).ProfitReport(context.Background(), params)

	require.NoError(t, err)
	assert.Equal(t, params, visto)
	assert.Same(t, esperado, got)

	dbErr := stderrors.New("timeout")
	repoErr := &mocks.RepositoryMock{
		ProfitReportFn: func(ctx context.Context, params dtos.ProfitReportParams) (*dtos.ProfitReportResponse, error) {
			return nil, dbErr
		},
	}
	_, err = newMarginsUseCase(repoErr, nil).ProfitReport(context.Background(), params)
	assert.ErrorIs(t, err, dbErr)
}

func TestProfitReportDetail_DelegaYPropaga(t *testing.T) {
	esperado := &dtos.ProfitReportDetailResponse{}
	var visto dtos.ProfitReportDetailParams
	repo := &mocks.RepositoryMock{
		ProfitReportDetailFn: func(ctx context.Context, params dtos.ProfitReportDetailParams) (*dtos.ProfitReportDetailResponse, error) {
			visto = params
			return esperado, nil
		},
	}
	params := dtos.ProfitReportDetailParams{BusinessID: 26}

	got, err := newMarginsUseCase(repo, nil).ProfitReportDetail(context.Background(), params)

	require.NoError(t, err)
	assert.Equal(t, params, visto)
	assert.Same(t, esperado, got)

	dbErr := stderrors.New("timeout")
	repoErr := &mocks.RepositoryMock{
		ProfitReportDetailFn: func(ctx context.Context, params dtos.ProfitReportDetailParams) (*dtos.ProfitReportDetailResponse, error) {
			return nil, dbErr
		},
	}
	_, err = newMarginsUseCase(repoErr, nil).ProfitReportDetail(context.Background(), params)
	assert.ErrorIs(t, err, dbErr)
}

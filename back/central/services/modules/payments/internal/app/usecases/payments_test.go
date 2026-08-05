package usecases

import (
	"context"
	stderrors "errors"
	"sort"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/entities"
	payerrs "github.com/secamc93/probability/back/central/services/modules/payments/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/payments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newPaymentsUseCase(repo *mocks.RepositoryMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo)
}

func boolPtr(v bool) *bool    { return &v }
func strPtr(v string) *string { return &v }

func TestCreatePaymentMethod_CodigoDuplicado_RechazaSinCrear(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		PaymentMethodExistsFn: func(ctx context.Context, code string) (bool, error) { return true, nil },
		CreatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error {
			creado = true
			return nil
		},
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMethod(context.Background(), &dtos.CreatePaymentMethod{Code: "nequi"})

	assert.ErrorIs(t, err, payerrs.ErrPaymentMethodCodeAlreadyExists)
	assert.Nil(t, got)
	assert.False(t, creado)
}

func TestCreatePaymentMethod_ErrorAlVerificarDuplicado_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	creado := false
	repo := &mocks.RepositoryMock{
		PaymentMethodExistsFn: func(ctx context.Context, code string) (bool, error) { return false, dbErr },
		CreatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error {
			creado = true
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).CreatePaymentMethod(context.Background(), &dtos.CreatePaymentMethod{Code: "x"})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado)
}

func TestCreatePaymentMethod_NaceActivoPorDefecto(t *testing.T) {
	var creado *entities.PaymentMethod
	repo := &mocks.RepositoryMock{
		CreatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error {
			creado = method
			method.ID = 9
			return nil
		},
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMethod(context.Background(), &dtos.CreatePaymentMethod{
		Code: "nequi", Name: "Nequi", Category: "digital_wallet", Provider: "Bancolombia",
		Icon: "nequi.svg", Color: "#FF00A0", Description: "billetera",
	})

	require.NoError(t, err)
	require.NotNil(t, creado)
	assert.True(t, creado.IsActive, "un metodo nuevo nace activo")
	assert.Equal(t, "nequi", creado.Code)
	assert.Equal(t, "digital_wallet", creado.Category)
	assert.Equal(t, uint(9), got.ID)
	assert.Equal(t, "#FF00A0", got.Color)
}

func TestCreatePaymentMethod_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error { return dbErr },
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMethod(context.Background(), &dtos.CreatePaymentMethod{Code: "x"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdatePaymentMethod_NoPermiteCambiarCodigoNiEstadoActivo(t *testing.T) {
	original := &entities.PaymentMethod{
		ID: 5, Code: "nequi", Name: "Nequi", IsActive: true,
	}
	var guardado *entities.PaymentMethod
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return original, nil
		},
		UpdatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error {
			guardado = method
			return nil
		},
	}

	got, err := newPaymentsUseCase(repo).UpdatePaymentMethod(context.Background(), 5, &dtos.UpdatePaymentMethod{
		Name: "Nequi Empresas", Description: "nueva", Category: "digital_wallet",
		Provider: "Bancolombia", Icon: "x.svg", Color: "#000000",
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "nequi", guardado.Code, "el codigo es la llave de negocio, no se toca en update")
	assert.True(t, guardado.IsActive, "el estado activo se cambia solo por el toggle")
	assert.Equal(t, "Nequi Empresas", guardado.Name)
	assert.Equal(t, "Nequi Empresas", got.Name)
}

func TestUpdatePaymentMethod_MetodoInexistente_SePropaga(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return nil, payerrs.ErrPaymentMethodNotFound
		},
		UpdatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error {
			actualizado = true
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).UpdatePaymentMethod(context.Background(), 404, &dtos.UpdatePaymentMethod{})

	assert.ErrorIs(t, err, payerrs.ErrPaymentMethodNotFound)
	assert.False(t, actualizado)
}

func TestUpdatePaymentMethod_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("deadlock")
	repo := &mocks.RepositoryMock{
		UpdatePaymentMethodFn: func(ctx context.Context, method *entities.PaymentMethod) error { return dbErr },
	}

	got, err := newPaymentsUseCase(repo).UpdatePaymentMethod(context.Background(), 1, &dtos.UpdatePaymentMethod{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDeletePaymentMethod_ConMapeosActivos_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		PaymentMethodHasActiveMapsFn: func(ctx context.Context, id uint) (bool, error) { return true, nil },
		DeletePaymentMethodFn: func(ctx context.Context, id uint) error {
			borrado = true
			return nil
		},
	}

	err := newPaymentsUseCase(repo).DeletePaymentMethod(context.Background(), 5)

	assert.ErrorIs(t, err, payerrs.ErrPaymentMethodHasActiveMappings)
	assert.False(t, borrado, "borrar un metodo con mapeos activos dejaria ordenes sin metodo de pago")
}

func TestDeletePaymentMethod_ErrorAlVerificarMapeos_NoBorra(t *testing.T) {
	dbErr := stderrors.New("timeout")
	borrado := false
	repo := &mocks.RepositoryMock{
		PaymentMethodHasActiveMapsFn: func(ctx context.Context, id uint) (bool, error) { return false, dbErr },
		DeletePaymentMethodFn: func(ctx context.Context, id uint) error {
			borrado = true
			return nil
		},
	}

	err := newPaymentsUseCase(repo).DeletePaymentMethod(context.Background(), 5)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, borrado)
}

func TestDeletePaymentMethod_SinMapeosActivos_Borra(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		DeletePaymentMethodFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	err := newPaymentsUseCase(repo).DeletePaymentMethod(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), borradoID)
}

func TestDeletePaymentMethod_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation")
	repo := &mocks.RepositoryMock{
		DeletePaymentMethodFn: func(ctx context.Context, id uint) error { return dbErr },
	}

	err := newPaymentsUseCase(repo).DeletePaymentMethod(context.Background(), 5)

	assert.ErrorIs(t, err, dbErr)
}

func TestListPaymentMethods_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		nombre       string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{"page cero", 0, 10, 1, 10},
		{"page negativo", -3, 10, 1, 10},
		{"pageSize cero", 1, 0, 1, 10},
		{"pageSize negativo", 1, -1, 1, 10},
		{"pageSize sobre el maximo cae al default", 1, 101, 1, 10},
		{"pageSize en el maximo se respeta", 1, 100, 1, 100},
		{"valores validos", 3, 25, 3, 25},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var vistoPage, vistoSize int
			repo := &mocks.RepositoryMock{
				ListPaymentMethodsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error) {
					vistoPage, vistoSize = page, pageSize
					return nil, 0, nil
				},
			}

			got, err := newPaymentsUseCase(repo).ListPaymentMethods(context.Background(), tc.page, tc.pageSize, nil)

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, vistoPage)
			assert.Equal(t, tc.wantPageSize, vistoSize)
			assert.Equal(t, tc.wantPage, got.Page)
			assert.Equal(t, tc.wantPageSize, got.PageSize)
		})
	}
}

func TestListPaymentMethods_CalculaTotalPagesConRedondeoHaciaArriba(t *testing.T) {
	casos := []struct {
		total    int64
		pageSize int
		want     int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{25, 10, 3},
	}

	for _, tc := range casos {
		repo := &mocks.RepositoryMock{
			ListPaymentMethodsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error) {
				return nil, tc.total, nil
			},
		}

		got, err := newPaymentsUseCase(repo).ListPaymentMethods(context.Background(), 1, tc.pageSize, nil)

		require.NoError(t, err)
		assert.Equal(t, tc.want, got.TotalPages, "total=%d pageSize=%d", tc.total, tc.pageSize)
	}
}

func TestListPaymentMethods_PasaLosFiltrosTalCual(t *testing.T) {
	esperados := map[string]interface{}{"is_active": true, "category": "card"}
	var recibidos map[string]interface{}
	repo := &mocks.RepositoryMock{
		ListPaymentMethodsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error) {
			recibidos = filters
			return nil, 0, nil
		},
	}

	_, err := newPaymentsUseCase(repo).ListPaymentMethods(context.Background(), 1, 10, esperados)

	require.NoError(t, err)
	assert.Equal(t, esperados, recibidos)
}

func TestListPaymentMethods_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListPaymentMethodsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).ListPaymentMethods(context.Background(), 1, 10, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListPaymentMethods_MapeaCadaEntidad(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListPaymentMethodsFn: func(ctx context.Context, page, pageSize int, filters map[string]interface{}) ([]entities.PaymentMethod, int64, error) {
			return []entities.PaymentMethod{
				{ID: 1, Code: "nequi", Name: "Nequi", IsActive: true},
				{ID: 2, Code: "cash", Name: "Efectivo", IsActive: false},
			}, 2, nil
		},
	}

	got, err := newPaymentsUseCase(repo).ListPaymentMethods(context.Background(), 1, 10, nil)

	require.NoError(t, err)
	require.Len(t, got.Data, 2)
	assert.Equal(t, "nequi", got.Data[0].Code)
	assert.True(t, got.Data[0].IsActive)
	assert.Equal(t, "cash", got.Data[1].Code)
	assert.False(t, got.Data[1].IsActive)
	assert.Equal(t, int64(2), got.Total)
}

func TestGetPaymentMethodByID_MapeaYPropagaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return &entities.PaymentMethod{ID: id, Code: "nequi"}, nil
		},
	}
	got, err := newPaymentsUseCase(repo).GetPaymentMethodByID(context.Background(), 3)
	require.NoError(t, err)
	assert.Equal(t, uint(3), got.ID)
	assert.Equal(t, "nequi", got.Code)

	dbErr := stderrors.New("no encontrado")
	repoErr := &mocks.RepositoryMock{
		GetPaymentMethodByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return nil, dbErr
		},
	}
	_, err = newPaymentsUseCase(repoErr).GetPaymentMethodByID(context.Background(), 3)
	assert.ErrorIs(t, err, dbErr)
}

func TestGetPaymentMethodByCode_PasaElCodigoYMapea(t *testing.T) {
	var visto string
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByCodeFn: func(ctx context.Context, code string) (*entities.PaymentMethod, error) {
			visto = code
			return &entities.PaymentMethod{ID: 1, Code: code}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).GetPaymentMethodByCode(context.Background(), "bancolombia_qr")

	require.NoError(t, err)
	assert.Equal(t, "bancolombia_qr", visto)
	assert.Equal(t, "bancolombia_qr", got.Code)
}

func TestGetPaymentMethodByCode_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no encontrado")
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByCodeFn: func(ctx context.Context, code string) (*entities.PaymentMethod, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).GetPaymentMethodByCode(context.Background(), "x")

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestTogglePaymentMethodActive_DevuelveElEstadoNuevo(t *testing.T) {
	var visto uint
	repo := &mocks.RepositoryMock{
		TogglePaymentMethodActiveFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			visto = id
			return &entities.PaymentMethod{ID: id, Code: "nequi", IsActive: false}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).TogglePaymentMethodActive(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), visto)
	assert.False(t, got.IsActive)
}

func TestTogglePaymentMethodActive_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no encontrado")
	repo := &mocks.RepositoryMock{
		TogglePaymentMethodActiveFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).TogglePaymentMethodActive(context.Background(), 5)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreatePaymentMapping_MapeoDuplicado_Rechaza(t *testing.T) {
	creado := false
	var vistoTipo, vistoOriginal string
	repo := &mocks.RepositoryMock{
		PaymentMappingExistsFn: func(ctx context.Context, integrationType, originalMethod string) (bool, error) {
			vistoTipo, vistoOriginal = integrationType, originalMethod
			return true, nil
		},
		CreatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			creado = true
			return nil
		},
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMapping(context.Background(), &dtos.CreatePaymentMapping{
		IntegrationType: "shopify", OriginalMethod: "Cash on Delivery",
	})

	assert.ErrorIs(t, err, payerrs.ErrPaymentMappingAlreadyExists)
	assert.Nil(t, got)
	assert.False(t, creado)
	assert.Equal(t, "shopify", vistoTipo)
	assert.Equal(t, "Cash on Delivery", vistoOriginal)
}

func TestCreatePaymentMapping_ErrorAlVerificarDuplicado_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	creado := false
	repo := &mocks.RepositoryMock{
		PaymentMappingExistsFn: func(ctx context.Context, integrationType, originalMethod string) (bool, error) {
			return false, dbErr
		},
		CreatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			creado = true
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).CreatePaymentMapping(context.Background(), &dtos.CreatePaymentMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado)
}

func TestCreatePaymentMapping_MetodoDestinoInexistente_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return nil, stderrors.New("record not found")
		},
		CreatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			creado = true
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).CreatePaymentMapping(context.Background(), &dtos.CreatePaymentMapping{
		PaymentMethodID: 999,
	})

	assert.ErrorIs(t, err, payerrs.ErrPaymentMethodNotFound)
	assert.False(t, creado, "un mapeo apuntando a un metodo inexistente dejaria ordenes sin resolver")
}

func TestCreatePaymentMapping_NaceActivoYRecargaConElMetodo(t *testing.T) {
	var creado *entities.PaymentMethodMapping
	repo := &mocks.RepositoryMock{
		CreatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			creado = mapping
			mapping.ID = 33
			return nil
		},
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return &entities.PaymentMethodMapping{
				ID: id, IntegrationType: "shopify", OriginalMethod: "COD",
				PaymentMethodID: 2, IsActive: true, Priority: 5,
				PaymentMethod: entities.PaymentMethod{ID: 2, Code: "cash", Name: "Efectivo"},
			}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMapping(context.Background(), &dtos.CreatePaymentMapping{
		IntegrationType: "shopify", OriginalMethod: "COD", PaymentMethodID: 2, Priority: 5,
	})

	require.NoError(t, err)
	require.NotNil(t, creado)
	assert.True(t, creado.IsActive, "un mapeo nuevo nace activo")
	assert.Equal(t, 5, creado.Priority)
	assert.Equal(t, uint(33), got.ID)
	assert.Equal(t, "cash", got.PaymentMethod.Code, "la respuesta trae el metodo destino ya cargado")
}

func TestCreatePaymentMapping_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			return dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMapping(context.Background(), &dtos.CreatePaymentMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreatePaymentMapping_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout al recargar")
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).CreatePaymentMapping(context.Background(), &dtos.CreatePaymentMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdatePaymentMapping_NoPermiteCambiarElTipoDeIntegracion(t *testing.T) {
	original := &entities.PaymentMethodMapping{
		ID: 7, IntegrationType: "shopify", OriginalMethod: "COD", PaymentMethodID: 1, Priority: 1, IsActive: true,
	}
	var guardado *entities.PaymentMethodMapping
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return original, nil
		},
		UpdatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			guardado = mapping
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).UpdatePaymentMapping(context.Background(), 7, &dtos.UpdatePaymentMapping{
		OriginalMethod: "Contra entrega", PaymentMethodID: 2, Priority: 9,
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "shopify", guardado.IntegrationType, "el canal de origen no se cambia por update")
	assert.True(t, guardado.IsActive, "el estado activo se cambia solo por el toggle")
	assert.Equal(t, "Contra entrega", guardado.OriginalMethod)
	assert.Equal(t, uint(2), guardado.PaymentMethodID)
	assert.Equal(t, 9, guardado.Priority)
}

func TestUpdatePaymentMapping_MapeoInexistente_SePropaga(t *testing.T) {
	dbErr := stderrors.New("record not found")
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return nil, dbErr
		},
		UpdatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			actualizado = true
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).UpdatePaymentMapping(context.Background(), 404, &dtos.UpdatePaymentMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, actualizado)
}

func TestUpdatePaymentMapping_MetodoDestinoInexistente_Rechaza(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetPaymentMethodByIDFn: func(ctx context.Context, id uint) (*entities.PaymentMethod, error) {
			return nil, stderrors.New("record not found")
		},
		UpdatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			actualizado = true
			return nil
		},
	}

	_, err := newPaymentsUseCase(repo).UpdatePaymentMapping(context.Background(), 7, &dtos.UpdatePaymentMapping{
		PaymentMethodID: 999,
	})

	assert.ErrorIs(t, err, payerrs.ErrPaymentMethodNotFound)
	assert.False(t, actualizado)
}

func TestUpdatePaymentMapping_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("deadlock")
	repo := &mocks.RepositoryMock{
		UpdatePaymentMappingFn: func(ctx context.Context, mapping *entities.PaymentMethodMapping) error {
			return dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).UpdatePaymentMapping(context.Background(), 7, &dtos.UpdatePaymentMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdatePaymentMapping_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			llamadas++
			return nil, dbErr
		},
	}

	_, err := newPaymentsUseCase(repo).UpdatePaymentMapping(context.Background(), 7, &dtos.UpdatePaymentMapping{})

	assert.ErrorIs(t, err, dbErr)
	assert.Equal(t, 1, llamadas)
}

func TestDeletePaymentMapping_DelegaAlRepositorio(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		DeletePaymentMappingFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}

	err := newPaymentsUseCase(repo).DeletePaymentMapping(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, uint(7), borradoID)
}

func TestDeletePaymentMapping_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation")
	repo := &mocks.RepositoryMock{
		DeletePaymentMappingFn: func(ctx context.Context, id uint) error { return dbErr },
	}

	err := newPaymentsUseCase(repo).DeletePaymentMapping(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
}

func TestTogglePaymentMappingActive_RecargaConElMetodoDestino(t *testing.T) {
	var toggled uint
	repo := &mocks.RepositoryMock{
		TogglePaymentMappingActiveFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			toggled = id
			return &entities.PaymentMethodMapping{ID: id, IsActive: false}, nil
		},
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return &entities.PaymentMethodMapping{
				ID: id, IsActive: false,
				PaymentMethod: entities.PaymentMethod{ID: 2, Code: "cash"},
			}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).TogglePaymentMappingActive(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, uint(7), toggled)
	assert.False(t, got.IsActive)
	assert.Equal(t, "cash", got.PaymentMethod.Code)
}

func TestTogglePaymentMappingActive_ErrorEnToggle_NoRecarga(t *testing.T) {
	dbErr := stderrors.New("no encontrado")
	recargado := false
	repo := &mocks.RepositoryMock{
		TogglePaymentMappingActiveFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return nil, dbErr
		},
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			recargado = true
			return nil, nil
		},
	}

	_, err := newPaymentsUseCase(repo).TogglePaymentMappingActive(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, recargado)
}

func TestTogglePaymentMappingActive_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).TogglePaymentMappingActive(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetPaymentMappingByID_MapeaConElMetodoDestino(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return &entities.PaymentMethodMapping{
				ID: id, IntegrationType: "meli", OriginalMethod: "account_money",
				PaymentMethod: entities.PaymentMethod{ID: 3, Code: "meli_wallet", Name: "Mercado Pago"},
			}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).GetPaymentMappingByID(context.Background(), 7)

	require.NoError(t, err)
	assert.Equal(t, "meli", got.IntegrationType)
	assert.Equal(t, "Mercado Pago", got.PaymentMethod.Name)
}

func TestGetPaymentMappingByID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("no encontrado")
	repo := &mocks.RepositoryMock{
		GetPaymentMappingByIDWithMethodFn: func(ctx context.Context, id uint) (*entities.PaymentMethodMapping, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).GetPaymentMappingByID(context.Background(), 7)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetPaymentMappingsByIntegrationType_PasaElTipoYMapea(t *testing.T) {
	var visto string
	repo := &mocks.RepositoryMock{
		GetPaymentMappingsByIntegrationTypeWithMethFn: func(ctx context.Context, integrationType string) ([]entities.PaymentMethodMapping, error) {
			visto = integrationType
			return []entities.PaymentMethodMapping{
				{ID: 1, IntegrationType: integrationType, OriginalMethod: "COD"},
				{ID: 2, IntegrationType: integrationType, OriginalMethod: "Card"},
			}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).GetPaymentMappingsByIntegrationType(context.Background(), "shopify")

	require.NoError(t, err)
	assert.Equal(t, "shopify", visto)
	require.Len(t, got, 2)
	assert.Equal(t, "COD", got[0].OriginalMethod)
}

func TestGetPaymentMappingsByIntegrationType_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetPaymentMappingsByIntegrationTypeWithMethFn: func(ctx context.Context, integrationType string) ([]entities.PaymentMethodMapping, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).GetPaymentMappingsByIntegrationType(context.Background(), "shopify")

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetAllPaymentMappingsGroupedByIntegration_AgrupaPorCanal(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListPaymentMappingsWithMethodsFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
			return []entities.PaymentMethodMapping{
				{ID: 1, IntegrationType: "shopify", OriginalMethod: "COD"},
				{ID: 2, IntegrationType: "meli", OriginalMethod: "account_money"},
				{ID: 3, IntegrationType: "shopify", OriginalMethod: "Card"},
			}, 3, nil
		},
	}

	got, err := newPaymentsUseCase(repo).GetAllPaymentMappingsGroupedByIntegration(context.Background())

	require.NoError(t, err)
	require.Len(t, got, 2)

	sort.Slice(got, func(i, j int) bool { return got[i].IntegrationType < got[j].IntegrationType })
	assert.Equal(t, "meli", got[0].IntegrationType)
	assert.Len(t, got[0].Mappings, 1)
	assert.Equal(t, "shopify", got[1].IntegrationType)
	assert.Len(t, got[1].Mappings, 2)
}

func TestGetAllPaymentMappingsGroupedByIntegration_NoConfundeMapeosEntreCanales(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListPaymentMappingsWithMethodsFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
			return []entities.PaymentMethodMapping{
				{ID: 1, IntegrationType: "shopify", OriginalMethod: "COD", PaymentMethodID: 10},
				{ID: 2, IntegrationType: "meli", OriginalMethod: "COD", PaymentMethodID: 20},
			}, 2, nil
		},
	}

	got, err := newPaymentsUseCase(repo).GetAllPaymentMappingsGroupedByIntegration(context.Background())

	require.NoError(t, err)
	porCanal := map[string]uint{}
	for _, g := range got {
		require.Len(t, g.Mappings, 1)
		porCanal[g.IntegrationType] = g.Mappings[0].PaymentMethodID
	}
	assert.Equal(t, uint(10), porCanal["shopify"])
	assert.Equal(t, uint(20), porCanal["meli"],
		"el mismo OriginalMethod en dos canales debe resolver a metodos distintos")
}

func TestGetAllPaymentMappingsGroupedByIntegration_SinMapeos_DevuelveVacioNoNil(t *testing.T) {
	got, err := newPaymentsUseCase(nil).GetAllPaymentMappingsGroupedByIntegration(context.Background())

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestGetAllPaymentMappingsGroupedByIntegration_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListPaymentMappingsWithMethodsFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).GetAllPaymentMappingsGroupedByIntegration(context.Background())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListPaymentMappings_PasaFiltrosYDevuelveTotal(t *testing.T) {
	esperados := map[string]interface{}{"integration_type": "shopify"}
	var recibidos map[string]interface{}
	repo := &mocks.RepositoryMock{
		ListPaymentMappingsWithMethodsFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
			recibidos = filters
			return []entities.PaymentMethodMapping{{ID: 1}}, 1, nil
		},
	}

	got, err := newPaymentsUseCase(repo).ListPaymentMappings(context.Background(), esperados)

	require.NoError(t, err)
	assert.Equal(t, esperados, recibidos)
	assert.Equal(t, int64(1), got.Total)
	assert.Len(t, got.Data, 1)
}

func TestListPaymentMappings_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListPaymentMappingsWithMethodsFn: func(ctx context.Context, filters map[string]interface{}) ([]entities.PaymentMethodMapping, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).ListPaymentMappings(context.Background(), nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListPaymentStatuses_MapeaTodosLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListPaymentStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error) {
			return []entities.PaymentStatus{
				{ID: 1, Code: "paid", Name: "Pagado", Description: "ok", Category: "final", Color: "#0F0", IsActive: true},
			}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).ListPaymentStatuses(context.Background(), nil)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, dtos.PaymentStatusInfo{
		ID: 1, Code: "paid", Name: "Pagado", Description: "ok", Category: "final", Color: "#0F0", IsActive: true,
	}, got[0])
}

func TestListPaymentStatuses_PasaElFiltroIsActive(t *testing.T) {
	for _, filtro := range []*bool{nil, boolPtr(true), boolPtr(false)} {
		var visto *bool
		repo := &mocks.RepositoryMock{
			ListPaymentStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error) {
				visto = isActive
				return nil, nil
			},
		}

		_, err := newPaymentsUseCase(repo).ListPaymentStatuses(context.Background(), filtro)

		require.NoError(t, err)
		assert.Equal(t, filtro, visto)
	}
}

func TestListPaymentStatuses_SinResultados_DevuelveVacioNoNil(t *testing.T) {
	got, err := newPaymentsUseCase(nil).ListPaymentStatuses(context.Background(), nil)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestListPaymentStatuses_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListPaymentStatusesFn: func(ctx context.Context, isActive *bool) ([]entities.PaymentStatus, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).ListPaymentStatuses(context.Background(), nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListChannelPaymentMethods_MapeaTodosLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListChannelPaymentMethodsFn: func(ctx context.Context, integrationType *string, isActive *bool) ([]entities.ChannelPaymentMethod, error) {
			return []entities.ChannelPaymentMethod{
				{ID: 1, IntegrationType: "shopify", Code: "cod", Name: "Contra entrega", Description: "d", IsActive: true, DisplayOrder: 3},
			}, nil
		},
	}

	got, err := newPaymentsUseCase(repo).ListChannelPaymentMethods(context.Background(), nil, nil)

	require.NoError(t, err)
	require.Len(t, got, 1)
	assert.Equal(t, dtos.ChannelPaymentMethodInfo{
		ID: 1, IntegrationType: "shopify", Code: "cod", Name: "Contra entrega",
		Description: "d", IsActive: true, DisplayOrder: 3,
	}, got[0])
}

func TestListChannelPaymentMethods_PasaAmbosFiltros(t *testing.T) {
	var vistoTipo *string
	var vistoActivo *bool
	repo := &mocks.RepositoryMock{
		ListChannelPaymentMethodsFn: func(ctx context.Context, integrationType *string, isActive *bool) ([]entities.ChannelPaymentMethod, error) {
			vistoTipo, vistoActivo = integrationType, isActive
			return nil, nil
		},
	}

	_, err := newPaymentsUseCase(repo).ListChannelPaymentMethods(context.Background(), strPtr("meli"), boolPtr(true))

	require.NoError(t, err)
	require.NotNil(t, vistoTipo)
	require.NotNil(t, vistoActivo)
	assert.Equal(t, "meli", *vistoTipo)
	assert.True(t, *vistoActivo)
}

func TestListChannelPaymentMethods_SinResultados_DevuelveVacioNoNil(t *testing.T) {
	got, err := newPaymentsUseCase(nil).ListChannelPaymentMethods(context.Background(), nil, nil)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestListChannelPaymentMethods_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListChannelPaymentMethodsFn: func(ctx context.Context, integrationType *string, isActive *bool) ([]entities.ChannelPaymentMethod, error) {
			return nil, dbErr
		},
	}

	got, err := newPaymentsUseCase(repo).ListChannelPaymentMethods(context.Background(), nil, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

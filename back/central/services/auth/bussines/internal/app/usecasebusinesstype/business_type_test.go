package usecasebusinesstype

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/bussines/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTypeUseCase(repo *mocks.BusinessRepositoryMock) IUseCaseBusinessType {
	if repo == nil {
		repo = &mocks.BusinessRepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func TestGenerateBusinessTypeCodeFromName_NormalizaYAgregaSufijo(t *testing.T) {
	casos := []struct {
		entrada    string
		wantPrefix string
	}{
		{"Retail", "retail_"},
		{"RETAIL", "retail_"},
		{"Comida Rapida", "comida_rapida_"},
		{"Comida-Rapida", "comida_rapida_"},
		{"Tipo#1", "tipo1_"},
	}

	for _, tc := range casos {
		t.Run(tc.entrada, func(t *testing.T) {
			got := generateBusinessTypeCodeFromName(tc.entrada)
			assert.True(t, strings.HasPrefix(got, tc.wantPrefix),
				"esperaba prefijo %q, obtuve %q", tc.wantPrefix, got)
		})
	}
}

func TestGenerateBusinessTypeCodeFromName_TruncaLaBaseAVeinte(t *testing.T) {
	got := generateBusinessTypeCodeFromName("TipoDeNegocioConNombreLarguisimoQueNoCabe")

	base := got[:strings.LastIndex(got, "_")]
	assert.LessOrEqual(t, len(base), 20)
}

func TestGenerateBusinessTypeCodeFromName_NoRepiteCodigos(t *testing.T) {
	vistos := make(map[string]bool, 200)
	for i := 0; i < 200; i++ {
		got := generateBusinessTypeCodeFromName("Retail")
		require.False(t, vistos[got], "codigo repetido en la iteracion %d: %q", i, got)
		vistos[got] = true
	}
}

func TestCreateBusinessType_NombreYaEnUso_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			return &domain.BusinessType{ID: 1, Name: name}, nil
		},
		CreateBusinessTypeFn: func(ctx context.Context, bt domain.BusinessType) (string, error) {
			creado = true
			return "", nil
		},
	}

	got, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{Name: "Retail"})

	assert.ErrorIs(t, err, domain.ErrBusinessTypeNameAlreadyExists)
	assert.Nil(t, got)
	assert.False(t, creado)
}

func TestCreateBusinessType_NotFoundAlBuscarNombre_NoEsBloqueante(t *testing.T) {
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return nil, domain.ErrBusinessTypeNotFound
			}
			return &domain.BusinessType{ID: 9, Name: name}, nil
		},
	}

	got, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{Name: "Retail"})

	require.NoError(t, err, "que el nombre no exista es justamente lo que buscamos")
	assert.Equal(t, uint(9), got.ID)
}

func TestCreateBusinessType_ErrorRealAlBuscarNombre_Aborta(t *testing.T) {
	dbErr := errors.New("connection refused")
	creado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			return nil, dbErr
		},
		CreateBusinessTypeFn: func(ctx context.Context, bt domain.BusinessType) (string, error) {
			creado = true
			return "", nil
		},
	}

	_, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{Name: "Retail"})

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado)
}

func TestCreateBusinessType_CodigoVacio_SeGeneraDesdeElNombre(t *testing.T) {
	var creado domain.BusinessType
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return nil, nil
			}
			return &domain.BusinessType{ID: 9, Name: name}, nil
		},
		CreateBusinessTypeFn: func(ctx context.Context, bt domain.BusinessType) (string, error) {
			creado = bt
			return "ok", nil
		},
	}

	_, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{Name: "Comida Rapida"})

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(creado.Code, "comida_rapida_"), "codigo generado: %q", creado.Code)
}

func TestCreateBusinessType_CodigoExplicito_SeRespeta(t *testing.T) {
	var creado domain.BusinessType
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return nil, nil
			}
			return &domain.BusinessType{ID: 9}, nil
		},
		CreateBusinessTypeFn: func(ctx context.Context, bt domain.BusinessType) (string, error) {
			creado = bt
			return "ok", nil
		},
	}

	_, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{
		Name: "Retail", Code: "RETAIL_FIJO",
	})

	require.NoError(t, err)
	assert.Equal(t, "RETAIL_FIJO", creado.Code)
}

func TestCreateBusinessType_PropagaDescripcionIconoYEstado(t *testing.T) {
	var creado domain.BusinessType
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return nil, nil
			}
			return &domain.BusinessType{ID: 9}, nil
		},
		CreateBusinessTypeFn: func(ctx context.Context, bt domain.BusinessType) (string, error) {
			creado = bt
			return "ok", nil
		},
	}

	_, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{
		Name: "Retail", Description: "tiendas", Icon: "shop.svg", IsActive: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "tiendas", creado.Description)
	assert.Equal(t, "shop.svg", creado.Icon)
	assert.True(t, creado.IsActive)
}

func TestCreateBusinessType_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := errors.New("constraint violation")
	repo := &mocks.BusinessRepositoryMock{
		CreateBusinessTypeFn: func(ctx context.Context, bt domain.BusinessType) (string, error) { return "", dbErr },
	}

	got, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{Name: "X"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateBusinessType_ErrorAlRecargarElCreado_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByNameFn: func(ctx context.Context, name string) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return nil, nil
			}
			return nil, dbErr
		},
	}

	got, err := newTypeUseCase(repo).CreateBusinessType(context.Background(), domain.BusinessTypeRequest{Name: "X"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetBusinessTypes_MapeaTodosLosCampos(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypesFn: func(ctx context.Context) ([]domain.BusinessType, error) {
			return []domain.BusinessType{
				{ID: 1, Name: "Retail", Code: "retail", Description: "d", Icon: "i", IsActive: true},
				{ID: 2, Name: "Comida", Code: "food", IsActive: false},
			}, nil
		},
	}

	got, err := newTypeUseCase(repo).GetBusinessTypes(context.Background())

	require.NoError(t, err)
	require.Len(t, got, 2)
	assert.Equal(t, "Retail", got[0].Name)
	assert.Equal(t, "d", got[0].Description)
	assert.Equal(t, "i", got[0].Icon)
	assert.True(t, got[0].IsActive)
	assert.False(t, got[1].IsActive)
}

func TestGetBusinessTypes_SinResultados_DevuelveVacioNoNil(t *testing.T) {
	got, err := newTypeUseCase(nil).GetBusinessTypes(context.Background())

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Empty(t, got)
}

func TestGetBusinessTypes_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypesFn: func(ctx context.Context) ([]domain.BusinessType, error) { return nil, dbErr },
	}

	got, err := newTypeUseCase(repo).GetBusinessTypes(context.Background())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetBusinessTypeByID_Inexistente_RetornaError(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) { return nil, nil },
	}

	got, err := newTypeUseCase(repo).GetBusinessTypeByID(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "tipo de negocio no encontrado")
}

func TestGetBusinessTypeByID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) { return nil, dbErr },
	}

	_, err := newTypeUseCase(repo).GetBusinessTypeByID(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetBusinessTypeByID_MapeaTodosLosCampos(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return &domain.BusinessType{
				ID: id, Name: "Retail", Code: "retail", Description: "d", Icon: "i", IsActive: true,
			}, nil
		},
	}

	got, err := newTypeUseCase(repo).GetBusinessTypeByID(context.Background(), 4)

	require.NoError(t, err)
	assert.Equal(t, uint(4), got.ID)
	assert.Equal(t, "Retail", got.Name)
	assert.Equal(t, "retail", got.Code)
	assert.Equal(t, "i", got.Icon)
	assert.True(t, got.IsActive)
}

func TestUpdateBusinessType_Inexistente_NoActualiza(t *testing.T) {
	actualizado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) { return nil, nil },
		UpdateBusinessTypeFn: func(ctx context.Context, id uint, bt domain.BusinessType) (string, error) {
			actualizado = true
			return "", nil
		},
	}

	_, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 404, domain.BusinessTypeRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tipo de negocio no encontrado")
	assert.False(t, actualizado)
}

func TestUpdateBusinessType_ErrorAlBuscar_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) { return nil, dbErr },
	}

	_, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateBusinessType_CodigoNuevoDuplicado_Rechaza(t *testing.T) {
	actualizado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return &domain.BusinessType{ID: id, Code: "viejo"}, nil
		},
		GetBusinessTypeByCodeFn: func(ctx context.Context, code string) (*domain.BusinessType, error) {
			return &domain.BusinessType{ID: 99, Code: code}, nil
		},
		UpdateBusinessTypeFn: func(ctx context.Context, id uint, bt domain.BusinessType) (string, error) {
			actualizado = true
			return "", nil
		},
	}

	_, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{Code: "ocupado"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "el código 'ocupado' ya existe")
	assert.False(t, actualizado)
}

func TestUpdateBusinessType_MismoCodigoNoDisparaChequeo(t *testing.T) {
	chequeado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return &domain.BusinessType{ID: id, Code: "igual"}, nil
		},
		GetBusinessTypeByCodeFn: func(ctx context.Context, code string) (*domain.BusinessType, error) {
			chequeado = true
			return nil, nil
		},
	}

	_, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{Code: "igual"})

	require.NoError(t, err)
	assert.False(t, chequeado)
}

func TestUpdateBusinessType_NotFoundAlChequearCodigo_NoEsBloqueante(t *testing.T) {
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return &domain.BusinessType{ID: id, Code: "viejo"}, nil
		},
		GetBusinessTypeByCodeFn: func(ctx context.Context, code string) (*domain.BusinessType, error) {
			return nil, errors.New("tipo de negocio no encontrado")
		},
	}

	_, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{Code: "libre"})

	require.NoError(t, err)
}

func TestUpdateBusinessType_ErrorRealAlChequearCodigo_Aborta(t *testing.T) {
	dbErr := errors.New("connection refused")
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			return &domain.BusinessType{ID: id, Code: "viejo"}, nil
		},
		GetBusinessTypeByCodeFn: func(ctx context.Context, code string) (*domain.BusinessType, error) {
			return nil, dbErr
		},
	}

	_, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{Code: "nuevo"})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateBusinessType_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := errors.New("deadlock")
	repo := &mocks.BusinessRepositoryMock{
		UpdateBusinessTypeFn: func(ctx context.Context, id uint, bt domain.BusinessType) (string, error) {
			return "", dbErr
		},
	}

	got, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateBusinessType_Exitoso_ReleeYDevuelveElActualizado(t *testing.T) {
	var guardado domain.BusinessType
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return &domain.BusinessType{ID: id, Code: "viejo", Name: "Viejo"}, nil
			}
			return &domain.BusinessType{ID: id, Code: "nuevo", Name: "Nuevo", IsActive: true}, nil
		},
		UpdateBusinessTypeFn: func(ctx context.Context, id uint, bt domain.BusinessType) (string, error) {
			guardado = bt
			return "ok", nil
		},
	}

	got, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{
		Name: "Nuevo", Code: "nuevo", Description: "d", Icon: "i", IsActive: true,
	})

	require.NoError(t, err)
	assert.Equal(t, "Nuevo", guardado.Name)
	assert.Equal(t, "nuevo", guardado.Code)
	assert.Equal(t, "Nuevo", got.Name, "la respuesta viene de la relectura")
	assert.True(t, got.IsActive)
}

func TestUpdateBusinessType_ErrorAlReleer_SePropaga(t *testing.T) {
	dbErr := errors.New("timeout")
	llamadas := 0
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) {
			llamadas++
			if llamadas == 1 {
				return &domain.BusinessType{ID: id, Code: "igual"}, nil
			}
			return nil, dbErr
		},
	}

	got, err := newTypeUseCase(repo).UpdateBusinessType(context.Background(), 1, domain.BusinessTypeRequest{Code: "igual"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDeleteBusinessType_Inexistente_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) { return nil, nil },
		DeleteBusinessTypeFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	err := newTypeUseCase(repo).DeleteBusinessType(context.Background(), 404)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "tipo de negocio no encontrado")
	assert.False(t, borrado)
}

func TestDeleteBusinessType_ErrorAlBuscar_NoBorra(t *testing.T) {
	dbErr := errors.New("timeout")
	borrado := false
	repo := &mocks.BusinessRepositoryMock{
		GetBusinessTypeByIDFn: func(ctx context.Context, id uint) (*domain.BusinessType, error) { return nil, dbErr },
		DeleteBusinessTypeFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	err := newTypeUseCase(repo).DeleteBusinessType(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, borrado)
}

func TestDeleteBusinessType_Exitoso(t *testing.T) {
	var borradoID uint
	repo := &mocks.BusinessRepositoryMock{
		DeleteBusinessTypeFn: func(ctx context.Context, id uint) (string, error) {
			borradoID = id
			return "ok", nil
		},
	}

	err := newTypeUseCase(repo).DeleteBusinessType(context.Background(), 4)

	require.NoError(t, err)
	assert.Equal(t, uint(4), borradoID)
}

func TestDeleteBusinessType_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := errors.New("fk violation: hay negocios de este tipo")
	repo := &mocks.BusinessRepositoryMock{
		DeleteBusinessTypeFn: func(ctx context.Context, id uint) (string, error) { return "", dbErr },
	}

	err := newTypeUseCase(repo).DeleteBusinessType(context.Background(), 4)

	assert.ErrorIs(t, err, dbErr)
}

package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func lpnActiva(id uint) *entities.LicensePlate {
	return &entities.LicensePlate{ID: id, Code: "LPN-001", LpnType: "pallet", Status: "active"}
}

func TestCreateLPN_CodigoDuplicado_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		LPNExistsByCodeFn: func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateLPN(context.Background(), request.CreateLPNDTO{BusinessID: 10, Code: "LPN-001"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateLPNCode)
	assert.Nil(t, got)
}

func TestCreateLPN_FallaVerificar_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		LPNExistsByCodeFn: func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateLPN(context.Background(), request.CreateLPNDTO{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateLPN_TipoPorDefectoEsPallet(t *testing.T) {
	var guardada *entities.LicensePlate
	repo := &mocks.RepositoryMock{
		LPNExistsByCodeFn: func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
			return false, nil
		},
		CreateLPNFn: func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
			guardada = lpn
			return lpn, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateLPN(context.Background(), request.CreateLPNDTO{BusinessID: 10, Code: "LPN-001"})

	require.NoError(t, err)
	require.NotNil(t, guardada)
	assert.Equal(t, "pallet", guardada.LpnType)
	assert.Equal(t, "active", guardada.Status)
	assert.Equal(t, uint(10), guardada.BusinessID)
}

func TestCreateLPN_TipoExplicitoSeRespeta(t *testing.T) {
	var guardada *entities.LicensePlate
	repo := &mocks.RepositoryMock{
		LPNExistsByCodeFn: func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
			return false, nil
		},
		CreateLPNFn: func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
			guardada = lpn
			return lpn, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateLPN(context.Background(), request.CreateLPNDTO{
		BusinessID: 10,
		Code:       "LPN-001",
		LpnType:    "caja",
	})

	require.NoError(t, err)
	assert.Equal(t, "caja", guardada.LpnType)
}

func TestListLPNs_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 300, wantPage: 1, wantPageSize: 10},
		{page: 2, pageSize: 50, wantPage: 2, wantPageSize: 50},
	}

	for _, tc := range cases {
		var recibido dtos.ListLPNParams
		repo := &mocks.RepositoryMock{
			ListLPNsFn: func(ctx context.Context, params dtos.ListLPNParams) ([]entities.LicensePlate, int64, error) {
				recibido = params
				return nil, 0, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, _, err := uc.ListLPNs(context.Background(), dtos.ListLPNParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, recibido.Page)
		assert.Equal(t, tc.wantPageSize, recibido.PageSize)
	}
}

func TestGetLPN_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GetLPN(context.Background(), 10, 4)

	require.NoError(t, err)
	assert.Equal(t, uint(4), got.ID)
}

func TestUpdateLPN_CodigoDuplicado_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
		LPNExistsByCodeFn: func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdateLPN(context.Background(), request.UpdateLPNDTO{ID: 1, BusinessID: 10, Code: "LPN-002"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateLPNCode)
	assert.Nil(t, got)
}

func TestUpdateLPN_ActualizaCampos(t *testing.T) {
	ubicacion := uint(15)
	var guardada *entities.LicensePlate
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
		LPNExistsByCodeFn: func(ctx context.Context, businessID uint, code string, excludeID *uint) (bool, error) {
			return false, nil
		},
		UpdateLPNFn: func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
			guardada = lpn
			return lpn, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateLPN(context.Background(), request.UpdateLPNDTO{
		ID:         1,
		BusinessID: 10,
		Code:       "LPN-002",
		LpnType:    "caja",
		LocationID: &ubicacion,
		Status:     "inactive",
	})

	require.NoError(t, err)
	assert.Equal(t, "LPN-002", guardada.Code)
	assert.Equal(t, "caja", guardada.LpnType)
	assert.Equal(t, uint(15), *guardada.CurrentLocationID)
	assert.Equal(t, "inactive", guardada.Status)
}

func TestUpdateLPN_CamposVacios_NoBorranLosExistentes(t *testing.T) {
	var guardada *entities.LicensePlate
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
		UpdateLPNFn: func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
			guardada = lpn
			return lpn, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateLPN(context.Background(), request.UpdateLPNDTO{ID: 1, BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, "LPN-001", guardada.Code)
	assert.Equal(t, "pallet", guardada.LpnType)
	assert.Equal(t, "active", guardada.Status)
}

func TestDeleteLPN_DelegaEnRepositorio(t *testing.T) {
	var borrado uint
	repo := &mocks.RepositoryMock{
		DeleteLPNFn: func(ctx context.Context, businessID, id uint) error {
			borrado = id
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	require.NoError(t, uc.DeleteLPN(context.Background(), 10, 4))
	assert.Equal(t, uint(4), borrado)
}

func TestAddToLPN_CantidadInvalida_RetornaError(t *testing.T) {
	uc := buildUseCase(&mocks.RepositoryMock{}, nil, nil)

	for _, qty := range []int{0, -5} {
		got, err := uc.AddToLPN(context.Background(), request.AddToLPNDTO{BusinessID: 10, LpnID: 1, Qty: qty})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
		assert.Nil(t, got)
	}
}

func TestAddToLPN_LPNDisuelta_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			lpn := lpnActiva(id)
			lpn.Status = "dissolved"
			return lpn, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.AddToLPN(context.Background(), request.AddToLPNDTO{BusinessID: 10, LpnID: 1, Qty: 5})

	assert.ErrorIs(t, err, domainerrors.ErrLPNDissolved)
	assert.Nil(t, got)
}

func TestAddToLPN_Exitoso_CreaLinea(t *testing.T) {
	var linea *entities.LicensePlateLine
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
		AddLPNLineFn: func(ctx context.Context, line *entities.LicensePlateLine) (*entities.LicensePlateLine, error) {
			linea = line
			return line, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.AddToLPN(context.Background(), request.AddToLPNDTO{
		BusinessID: 10,
		LpnID:      1,
		ProductID:  "prod-1",
		Qty:        5,
	})

	require.NoError(t, err)
	require.NotNil(t, linea)
	assert.Equal(t, uint(1), linea.LpnID)
	assert.Equal(t, "prod-1", linea.ProductID)
	assert.Equal(t, 5, linea.Qty)
	assert.Equal(t, uint(10), linea.BusinessID)
}

func TestAddToLPN_LPNInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.AddToLPN(context.Background(), request.AddToLPNDTO{BusinessID: 10, LpnID: 1, Qty: 5})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestMoveLPN_CambiaLaUbicacion(t *testing.T) {
	var guardada *entities.LicensePlate
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
		UpdateLPNFn: func(ctx context.Context, lpn *entities.LicensePlate) (*entities.LicensePlate, error) {
			guardada = lpn
			return lpn, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.MoveLPN(context.Background(), request.MoveLPNDTO{
		BusinessID:    10,
		LpnID:         1,
		NewLocationID: 42,
	})

	require.NoError(t, err)
	require.NotNil(t, guardada.CurrentLocationID)
	assert.Equal(t, uint(42), *guardada.CurrentLocationID)
}

func TestMoveLPN_LPNInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.MoveLPN(context.Background(), request.MoveLPNDTO{BusinessID: 10, LpnID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDissolveLPN_YaDisuelta_Rechaza(t *testing.T) {
	disuelta := false
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			lpn := lpnActiva(id)
			lpn.Status = "dissolved"
			return lpn, nil
		},
		DissolveLPNFn: func(ctx context.Context, businessID, id uint) error {
			disuelta = true
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	err := uc.DissolveLPN(context.Background(), request.DissolveLPNDTO{BusinessID: 10, LpnID: 1})

	assert.ErrorIs(t, err, domainerrors.ErrLPNDissolved)
	assert.False(t, disuelta, "no debe disolverse dos veces")
}

func TestDissolveLPN_Activa_SeDisuelve(t *testing.T) {
	var disueltaID uint
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return lpnActiva(id), nil
		},
		DissolveLPNFn: func(ctx context.Context, businessID, id uint) error {
			disueltaID = id
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	require.NoError(t, uc.DissolveLPN(context.Background(), request.DissolveLPNDTO{BusinessID: 10, LpnID: 1}))
	assert.Equal(t, uint(1), disueltaID)
}

func TestDissolveLPN_LPNInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	assert.ErrorIs(t, uc.DissolveLPN(context.Background(), request.DissolveLPNDTO{BusinessID: 10, LpnID: 1}), dbErr)
}

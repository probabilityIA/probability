package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repoLotBase() *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		LotExistsByCodeFn: func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
			return false, nil
		},
		CreateLotFn: func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
			lot.ID = 1
			return lot, nil
		},
	}
}

func TestCreateLot_ProductoInexistente_RetornaError(t *testing.T) {
	repo := repoLotBase()
	repo.GetProductByIDFn = func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
		return "", "", false, errors.New("not found")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateLot(context.Background(), request.CreateLotDTO{ProductID: "prod-1", BusinessID: 10})

	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
	assert.Nil(t, got)
}

func TestCreateLot_CodigoDuplicado_RetornaError(t *testing.T) {
	repo := repoLotBase()
	repo.LotExistsByCodeFn = func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
		return true, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateLot(context.Background(), request.CreateLotDTO{
		ProductID:  "prod-1",
		BusinessID: 10,
		LotCode:    "L-001",
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateLotCode)
	assert.Nil(t, got)
}

func TestCreateLot_FallaVerificarDuplicado_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := repoLotBase()
	repo.LotExistsByCodeFn = func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
		return false, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateLot(context.Background(), request.CreateLotDTO{ProductID: "prod-1", BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateLot_EstadoPorDefectoEsActive(t *testing.T) {
	var guardado *entities.InventoryLot
	repo := repoLotBase()
	repo.CreateLotFn = func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
		guardado = lot
		return lot, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateLot(context.Background(), request.CreateLotDTO{
		ProductID:  "prod-1",
		BusinessID: 10,
		LotCode:    "L-001",
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "active", guardado.Status)
}

func TestCreateLot_CopiaTodosLosCampos(t *testing.T) {
	fabricacion := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	vencimiento := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	recibido := time.Date(2026, 2, 1, 0, 0, 0, 0, time.UTC)
	proveedor := uint(7)

	var guardado *entities.InventoryLot
	repo := repoLotBase()
	repo.CreateLotFn = func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
		guardado = lot
		return lot, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateLot(context.Background(), request.CreateLotDTO{
		ProductID:       "prod-1",
		BusinessID:      10,
		LotCode:         "L-001",
		ManufactureDate: &fabricacion,
		ExpirationDate:  &vencimiento,
		ReceivedAt:      &recibido,
		SupplierID:      &proveedor,
		Status:          "quarantine",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(10), guardado.BusinessID)
	assert.Equal(t, "prod-1", guardado.ProductID)
	assert.Equal(t, "L-001", guardado.LotCode)
	assert.Equal(t, "quarantine", guardado.Status)
	assert.True(t, guardado.ExpirationDate.Equal(vencimiento))
	assert.Equal(t, uint(7), *guardado.SupplierID)
}

func TestListLots_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "pagina cero", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "page size cero", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{name: "page size sobre el maximo", page: 1, pageSize: 500, wantPage: 1, wantPageSize: 10},
		{name: "page size en el limite", page: 2, pageSize: 100, wantPage: 2, wantPageSize: 100},
		{name: "valores validos", page: 3, pageSize: 25, wantPage: 3, wantPageSize: 25},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var recibido dtos.ListLotsParams
			repo := &mocks.RepositoryMock{
				ListLotsFn: func(ctx context.Context, params dtos.ListLotsParams) ([]entities.InventoryLot, int64, error) {
					recibido = params
					return nil, 0, nil
				},
			}
			uc := buildUseCase(repo, nil, nil)

			_, _, err := uc.ListLots(context.Background(), dtos.ListLotsParams{
				Page:     tc.page,
				PageSize: tc.pageSize,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, recibido.Page)
			assert.Equal(t, tc.wantPageSize, recibido.PageSize)
		})
	}
}

func TestGetLot_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLotByIDFn: func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
			return &entities.InventoryLot{ID: lotID, BusinessID: businessID}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GetLot(context.Background(), 10, 3)

	require.NoError(t, err)
	assert.Equal(t, uint(3), got.ID)
	assert.Equal(t, uint(10), got.BusinessID)
}

func TestUpdateLot_LoteInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetLotByIDFn: func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdateLot(context.Background(), request.UpdateLotDTO{ID: 1, BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateLot_CodigoDuplicado_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetLotByIDFn: func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
			return &entities.InventoryLot{ID: lotID, LotCode: "L-001", ProductID: "prod-1"}, nil
		},
		LotExistsByCodeFn: func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdateLot(context.Background(), request.UpdateLotDTO{
		ID:         1,
		BusinessID: 10,
		LotCode:    "L-002",
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateLotCode)
	assert.Nil(t, got)
}

func TestUpdateLot_MismoCodigo_NoVerificaDuplicado(t *testing.T) {
	verificado := false
	repo := &mocks.RepositoryMock{
		GetLotByIDFn: func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
			return &entities.InventoryLot{ID: lotID, LotCode: "L-001"}, nil
		},
		LotExistsByCodeFn: func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
			verificado = true
			return false, nil
		},
		UpdateLotFn: func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
			return lot, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateLot(context.Background(), request.UpdateLotDTO{
		ID:         1,
		BusinessID: 10,
		LotCode:    "L-001",
	})

	require.NoError(t, err)
	assert.False(t, verificado)
}

func TestUpdateLot_ActualizaCamposOpcionales(t *testing.T) {
	vencimiento := time.Date(2027, 6, 1, 0, 0, 0, 0, time.UTC)
	proveedor := uint(9)

	var guardado *entities.InventoryLot
	repo := &mocks.RepositoryMock{
		GetLotByIDFn: func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
			return &entities.InventoryLot{ID: lotID, LotCode: "L-001", Status: "active"}, nil
		},
		LotExistsByCodeFn: func(ctx context.Context, businessID uint, productID, code string, excludeID *uint) (bool, error) {
			return false, nil
		},
		UpdateLotFn: func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
			guardado = lot
			return lot, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateLot(context.Background(), request.UpdateLotDTO{
		ID:             1,
		BusinessID:     10,
		LotCode:        "L-002",
		ExpirationDate: &vencimiento,
		SupplierID:     &proveedor,
		Status:         "quarantine",
	})

	require.NoError(t, err)
	assert.Equal(t, "L-002", guardado.LotCode)
	assert.True(t, guardado.ExpirationDate.Equal(vencimiento))
	assert.Equal(t, uint(9), *guardado.SupplierID)
	assert.Equal(t, "quarantine", guardado.Status)
}

func TestUpdateLot_CamposVacios_NoBorranLosExistentes(t *testing.T) {
	original := time.Date(2027, 1, 1, 0, 0, 0, 0, time.UTC)
	proveedor := uint(7)

	var guardado *entities.InventoryLot
	repo := &mocks.RepositoryMock{
		GetLotByIDFn: func(ctx context.Context, businessID, lotID uint) (*entities.InventoryLot, error) {
			return &entities.InventoryLot{
				ID:             lotID,
				LotCode:        "L-001",
				Status:         "active",
				ExpirationDate: &original,
				SupplierID:     &proveedor,
			}, nil
		},
		UpdateLotFn: func(ctx context.Context, lot *entities.InventoryLot) (*entities.InventoryLot, error) {
			guardado = lot
			return lot, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateLot(context.Background(), request.UpdateLotDTO{ID: 1, BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, "L-001", guardado.LotCode)
	assert.Equal(t, "active", guardado.Status)
	assert.True(t, guardado.ExpirationDate.Equal(original))
	assert.Equal(t, uint(7), *guardado.SupplierID)
}

func TestDeleteLot_DelegaEnRepositorio(t *testing.T) {
	var borradoBusiness, borradoLot uint
	repo := &mocks.RepositoryMock{
		DeleteLotFn: func(ctx context.Context, businessID, lotID uint) error {
			borradoBusiness, borradoLot = businessID, lotID
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	err := uc.DeleteLot(context.Background(), 10, 3)

	require.NoError(t, err)
	assert.Equal(t, uint(10), borradoBusiness)
	assert.Equal(t, uint(3), borradoLot)
}

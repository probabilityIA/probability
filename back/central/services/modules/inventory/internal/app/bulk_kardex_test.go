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

func intPtrInv(v int) *int {
	return &v
}

func repoBulkBase(rastrea bool) *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		WarehouseExistsFn: func(ctx context.Context, warehouseID, businessID uint) (bool, error) {
			return true, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductBySKUFn: func(ctx context.Context, sku string, businessID uint) (string, string, bool, error) {
			return "prod-1", "Producto", rastrea, nil
		},
		EnableProductTrackInventoryFn: func(ctx context.Context, productID string) error {
			return nil
		},
		AdjustStockTxFn: func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
			return &dtos.AdjustStockTxResult{
				Movement:    &entities.StockMovement{PreviousQty: 0, NewQty: params.Quantity},
				NewQuantity: params.Quantity,
			}, nil
		},
	}
}

func bulkDTO(items ...request.BulkLoadItem) request.BulkLoadDTO {
	return request.BulkLoadDTO{
		BusinessID:  10,
		WarehouseID: 5,
		Items:       items,
	}
}

func TestBulkLoad_SinItems_RetornaError(t *testing.T) {
	uc := buildUseCase(repoBulkBase(true), nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO())

	assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
	assert.Nil(t, got)
}

func TestBulkLoad_BodegaInexistente_RetornaError(t *testing.T) {
	repo := repoBulkBase(true)
	repo.WarehouseExistsFn = func(ctx context.Context, warehouseID, businessID uint) (bool, error) {
		return false, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	assert.ErrorIs(t, err, domainerrors.ErrWarehouseNotFound)
	assert.Nil(t, got)
}

func TestBulkLoad_FallaVerificarBodega_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := repoBulkBase(true)
	repo.WarehouseExistsFn = func(ctx context.Context, warehouseID, businessID uint) (bool, error) {
		return false, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestBulkLoad_SinTipoDeMovimiento_RetornaError(t *testing.T) {
	repo := repoBulkBase(true)
	repo.GetMovementTypeIDByCodeFn = func(ctx context.Context, code string) (uint, error) {
		return 0, errors.New("catalogo vacio")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.Error(t, err)
	assert.Contains(t, err.Error(), "inbound movement type")
	assert.Nil(t, got)
}

func TestBulkLoad_CantidadInvalida_FallaSoloEseItem(t *testing.T) {
	uc := buildUseCase(repoBulkBase(true), nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(
		request.BulkLoadItem{SKU: "SKU-1", Quantity: 0},
		request.BulkLoadItem{SKU: "SKU-2", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 2, got.TotalItems)
	assert.Equal(t, 1, got.SuccessCount)
	assert.Equal(t, 1, got.FailureCount)
	assert.Equal(t, "quantity must be positive", got.Items[0].Error)
}

func TestBulkLoad_SKUInexistente_FallaSoloEseItem(t *testing.T) {
	repo := repoBulkBase(true)
	repo.GetProductBySKUFn = func(ctx context.Context, sku string, businessID uint) (string, string, bool, error) {
		return "", "", false, errors.New("not found")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-X", Quantity: 10}))

	require.NoError(t, err)
	assert.Equal(t, 1, got.FailureCount)
	assert.Contains(t, got.Items[0].Error, "product not found for SKU SKU-X")
}

func TestBulkLoad_ProductoSinRastreo_LoActiva(t *testing.T) {
	activado := ""
	repo := repoBulkBase(false)
	repo.EnableProductTrackInventoryFn = func(ctx context.Context, productID string) error {
		activado = productID
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.NoError(t, err)
	assert.Equal(t, "prod-1", activado, "la carga masiva activa el rastreo de inventario del producto")
	assert.Equal(t, 1, got.SuccessCount)
}

func TestBulkLoad_FallaActivarRastreo_FallaElItem(t *testing.T) {
	repo := repoBulkBase(false)
	repo.EnableProductTrackInventoryFn = func(ctx context.Context, productID string) error {
		return errors.New("db down")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.NoError(t, err)
	assert.Equal(t, 1, got.FailureCount)
	assert.Contains(t, got.Items[0].Error, "failed to enable inventory tracking")
}

func TestBulkLoad_RazonPorDefecto(t *testing.T) {
	var recibido dtos.AdjustStockTxParams
	repo := repoBulkBase(true)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.NoError(t, err)
	assert.Equal(t, "bulk_load", recibido.Reason)
	assert.Equal(t, "bulk_load", recibido.ReferenceType)
	assert.Contains(t, recibido.Notes, "SKU-1")
}

func TestBulkLoad_RazonExplicita(t *testing.T) {
	var recibido dtos.AdjustStockTxParams
	repo := repoBulkBase(true)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	dto := bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10})
	dto.Reason = "inventario inicial"

	_, err := uc.BulkLoadInventory(context.Background(), dto)

	require.NoError(t, err)
	assert.Equal(t, "inventario inicial", recibido.Reason)
	assert.Equal(t, "bulk_load", recibido.ReferenceType, "el tipo de referencia siempre es bulk_load")
}

func TestBulkLoad_FallaAjuste_FallaElItem(t *testing.T) {
	repo := repoBulkBase(true)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		return nil, errors.New("deadlock")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.NoError(t, err)
	assert.Equal(t, 1, got.FailureCount)
	assert.Contains(t, got.Items[0].Error, "failed to adjust stock")
}

func TestBulkLoad_ReportaCantidadAnteriorYNueva(t *testing.T) {
	repo := repoBulkBase(true)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		return &dtos.AdjustStockTxResult{
			Movement: &entities.StockMovement{PreviousQty: 4, NewQty: 14},
		}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.NoError(t, err)
	assert.True(t, got.Items[0].Success)
	assert.Equal(t, 4, got.Items[0].PreviousQty)
	assert.Equal(t, 14, got.Items[0].NewQty)
}

func TestBulkLoad_ActualizaNivelesMinimoYMaximo(t *testing.T) {
	var nivelGuardado *entities.InventoryLevel
	repo := repoBulkBase(true)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		return &dtos.AdjustStockTxResult{Level: &entities.InventoryLevel{ID: 1}}, nil
	}
	repo.UpdateLevelFn = func(ctx context.Context, level *entities.InventoryLevel) error {
		nivelGuardado = level
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{
		SKU:          "SKU-1",
		Quantity:     10,
		MinStock:     intPtrInv(5),
		MaxStock:     intPtrInv(100),
		ReorderPoint: intPtrInv(20),
	}))

	require.NoError(t, err)
	require.NotNil(t, nivelGuardado)
	assert.Equal(t, 5, *nivelGuardado.MinStock)
	assert.Equal(t, 100, *nivelGuardado.MaxStock)
	assert.Equal(t, 20, *nivelGuardado.ReorderPoint)
}

func TestBulkLoad_SinNiveles_NoActualizaElNivel(t *testing.T) {
	actualizado := false
	repo := repoBulkBase(true)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		return &dtos.AdjustStockTxResult{Level: &entities.InventoryLevel{ID: 1}}, nil
	}
	repo.UpdateLevelFn = func(ctx context.Context, level *entities.InventoryLevel) error {
		actualizado = true
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.BulkLoadInventory(context.Background(), bulkDTO(request.BulkLoadItem{SKU: "SKU-1", Quantity: 10}))

	require.NoError(t, err)
	assert.False(t, actualizado)
}

func TestBulkLoad_VariosItems_CuentaExitosYFallos(t *testing.T) {
	repo := repoBulkBase(true)
	repo.GetProductBySKUFn = func(ctx context.Context, sku string, businessID uint) (string, string, bool, error) {
		if sku == "SKU-MALO" {
			return "", "", false, errors.New("not found")
		}
		return "prod-" + sku, "Producto", true, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.BulkLoadInventory(context.Background(), bulkDTO(
		request.BulkLoadItem{SKU: "SKU-1", Quantity: 10},
		request.BulkLoadItem{SKU: "SKU-MALO", Quantity: 5},
		request.BulkLoadItem{SKU: "SKU-2", Quantity: 3},
		request.BulkLoadItem{SKU: "SKU-3", Quantity: -1},
	))

	require.NoError(t, err)
	assert.Equal(t, 4, got.TotalItems)
	assert.Equal(t, 2, got.SuccessCount)
	assert.Equal(t, 2, got.FailureCount)
	assert.Len(t, got.Items, 4)
}

func TestExportKardex_SumaEntradasYSalidas(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetKardexFn: func(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
			return []entities.KardexEntry{
				{Quantity: 100, RunningBalance: 100},
				{Quantity: -30, RunningBalance: 70},
				{Quantity: 50, RunningBalance: 120},
				{Quantity: -20, RunningBalance: 100},
			}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ExportKardex(context.Background(), request.KardexExportDTO{
		BusinessID: 10,
		ProductID:  "prod-1",
	})

	require.NoError(t, err)
	assert.Equal(t, 150, got.TotalIn, "100 + 50 de entradas")
	assert.Equal(t, 50, got.TotalOut, "30 + 20 de salidas, en positivo")
	assert.Equal(t, 100, got.FinalBalance, "el saldo final es el del ultimo movimiento")
	assert.Len(t, got.Entries, 4)
}

func TestExportKardex_SinMovimientos_TodoEnCero(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetKardexFn: func(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ExportKardex(context.Background(), request.KardexExportDTO{BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, 0, got.TotalIn)
	assert.Equal(t, 0, got.TotalOut)
	assert.Equal(t, 0, got.FinalBalance)
}

func TestExportKardex_MovimientoCeroCuentaComoSalida(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetKardexFn: func(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
			return []entities.KardexEntry{{Quantity: 0, RunningBalance: 5}}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ExportKardex(context.Background(), request.KardexExportDTO{BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, 0, got.TotalIn)
	assert.Equal(t, 0, got.TotalOut)
	assert.Equal(t, 5, got.FinalBalance)
}

func TestExportKardex_PropagaLosFiltros(t *testing.T) {
	var recibido dtos.KardexQueryParams
	repo := &mocks.RepositoryMock{
		GetKardexFn: func(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
			recibido = params
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	desde := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	hasta := time.Date(2026, 1, 31, 0, 0, 0, 0, time.UTC)
	_, err := uc.ExportKardex(context.Background(), request.KardexExportDTO{
		BusinessID:  10,
		ProductID:   "prod-1",
		WarehouseID: 5,
		From:        &desde,
		To:          &hasta,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(10), recibido.BusinessID)
	assert.Equal(t, "prod-1", recibido.ProductID)
	require.NotNil(t, recibido.From)
	require.NotNil(t, recibido.To)
	assert.True(t, recibido.From.Equal(desde))
	assert.True(t, recibido.To.Equal(hasta))
}

func TestExportKardex_FallaConsulta_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetKardexFn: func(ctx context.Context, params dtos.KardexQueryParams) ([]entities.KardexEntry, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ExportKardex(context.Background(), request.KardexExportDTO{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestMergeLPN_MueveLasLineasYDisuelveElOrigen(t *testing.T) {
	var lineasAgregadas []*entities.LicensePlateLine
	var disuelta uint
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return &entities.LicensePlate{ID: id, Status: "active"}, nil
		},
		ListLPNLinesFn: func(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error) {
			return []entities.LicensePlateLine{
				{ID: 10, LpnID: lpnID, ProductID: "prod-1", Qty: 5},
				{ID: 11, LpnID: lpnID, ProductID: "prod-2", Qty: 3},
			}, nil
		},
		AddLPNLineFn: func(ctx context.Context, line *entities.LicensePlateLine) (*entities.LicensePlateLine, error) {
			lineasAgregadas = append(lineasAgregadas, line)
			return line, nil
		},
		DissolveLPNFn: func(ctx context.Context, businessID, id uint) error {
			disuelta = id
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.MergeLPN(context.Background(), request.MergeLPNDTO{
		BusinessID:  10,
		SourceLpnID: 1,
		TargetLpnID: 2,
	})

	require.NoError(t, err)
	require.Len(t, lineasAgregadas, 2)
	for _, l := range lineasAgregadas {
		assert.Equal(t, uint(0), l.ID, "las lineas se recrean, no se reusan los IDs")
		assert.Equal(t, uint(2), l.LpnID, "las lineas quedan en la LPN destino")
		assert.Equal(t, uint(10), l.BusinessID)
	}
	assert.Equal(t, uint(1), disuelta, "la LPN origen se disuelve tras la fusion")
}

func TestMergeLPN_OrigenInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.MergeLPN(context.Background(), request.MergeLPNDTO{BusinessID: 10, SourceLpnID: 1, TargetLpnID: 2})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestMergeLPN_FallaListarLineas_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return &entities.LicensePlate{ID: id}, nil
		},
		ListLPNLinesFn: func(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.MergeLPN(context.Background(), request.MergeLPNDTO{BusinessID: 10, SourceLpnID: 1, TargetLpnID: 2})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestMergeLPN_SinLineas_IgualDisuelveElOrigen(t *testing.T) {
	disuelta := false
	repo := &mocks.RepositoryMock{
		GetLPNByIDFn: func(ctx context.Context, businessID, id uint) (*entities.LicensePlate, error) {
			return &entities.LicensePlate{ID: id}, nil
		},
		ListLPNLinesFn: func(ctx context.Context, lpnID uint) ([]entities.LicensePlateLine, error) {
			return nil, nil
		},
		DissolveLPNFn: func(ctx context.Context, businessID, id uint) error {
			disuelta = true
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.MergeLPN(context.Background(), request.MergeLPNDTO{BusinessID: 10, SourceLpnID: 1, TargetLpnID: 2})

	require.NoError(t, err)
	assert.True(t, disuelta)
}

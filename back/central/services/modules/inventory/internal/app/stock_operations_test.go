package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func itemsDePrueba() []dtos.OrderInventoryItem {
	return []dtos.OrderInventoryItem{{ProductID: "prod-1", SKU: "SKU-1", Quantity: 5}}
}

func repoStockBase(rastrea bool) *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetDefaultWarehouseIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 5, nil
		},
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", rastrea, nil
		},
		ReserveStockTxFn: func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
			return &dtos.ReserveStockTxResult{Reserved: params.Quantity, Sufficient: true}, nil
		},
		ReleaseStockTxFn: func(ctx context.Context, params dtos.ReleaseTxParams) error { return nil },
		ConfirmSaleTxFn:  func(ctx context.Context, params dtos.ConfirmSaleTxParams) error { return nil },
		ReturnStockTxFn:  func(ctx context.Context, params dtos.ReturnStockTxParams) error { return nil },
	}
}

func TestResolveWarehouse_BodegaExplicita_NoConsultaLaPorDefecto(t *testing.T) {
	consultada := false
	repo := repoStockBase(true)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		consultada = true
		return 5, nil
	}
	uc := buildUseCase(repo, nil, nil)

	bodega := uint(9)
	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, &bodega, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, uint(9), got.WarehouseID)
	assert.False(t, consultada)
}

func TestResolveWarehouse_BodegaCero_CaeALaPorDefecto(t *testing.T) {
	repo := repoStockBase(true)
	uc := buildUseCase(repo, nil, nil)

	cero := uint(0)
	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, &cero, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, uint(5), got.WarehouseID)
}

func TestReserveStock_SinBodegaPorDefecto_RetornaError(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		return 0, errors.New("sin bodega")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, domainerrors.ErrNoDefaultWarehouse)
	assert.Nil(t, got)
}

func TestReserveStock_SinTipoDeMovimiento_RetornaError(t *testing.T) {
	dbErr := errors.New("catalogo vacio")
	repo := repoStockBase(true)
	repo.GetMovementTypeIDByCodeFn = func(ctx context.Context, code string) (uint, error) {
		return 0, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestReserveStock_ProductoInexistente_MarcaItemSinExito(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetProductByIDFn = func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
		return "", "", false, errors.New("not found")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	require.Len(t, got.ItemResults, 1)
	assert.Equal(t, "producto no encontrado", got.ItemResults[0].ErrorMessage)
	assert.False(t, got.ItemResults[0].Sufficient)
	assert.True(t, got.Success, "un producto inexistente no marca la operacion completa como fallida")
}

func TestReserveStock_ProductoSinRastreo_SeDaPorProcesado(t *testing.T) {
	reservado := false
	repo := repoStockBase(false)
	repo.ReserveStockTxFn = func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
		reservado = true
		return &dtos.ReserveStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, 5, got.ItemResults[0].Processed)
	assert.True(t, got.ItemResults[0].Sufficient)
	assert.False(t, reservado, "un producto sin rastreo de inventario no toca el stock")
}

func TestReserveStock_Exitoso(t *testing.T) {
	var recibido dtos.ReserveStockTxParams
	repo := repoStockBase(true)
	repo.ReserveStockTxFn = func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
		recibido = params
		return &dtos.ReserveStockTxResult{Reserved: 5, Sufficient: true}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, 5, got.ItemResults[0].Processed)
	assert.True(t, got.ItemResults[0].Sufficient)
	assert.Equal(t, "order-1", recibido.OrderID)
	assert.Equal(t, uint(5), recibido.WarehouseID)
	assert.Equal(t, 5, recibido.Quantity)
}

func TestReserveStock_StockInsuficiente_ReservaParcial(t *testing.T) {
	repo := repoStockBase(true)
	repo.ReserveStockTxFn = func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
		return &dtos.ReserveStockTxResult{Reserved: 2, Sufficient: false}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, 2, got.ItemResults[0].Processed)
	assert.False(t, got.ItemResults[0].Sufficient)
	assert.True(t, got.Success, "una reserva parcial no marca la operacion como fallida")
}

func TestReserveStock_FallaTransaccion_MarcaOperacionFallida(t *testing.T) {
	repo := repoStockBase(true)
	repo.ReserveStockTxFn = func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
		return nil, errors.New("deadlock")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.False(t, got.Success)
	assert.Contains(t, got.ItemResults[0].ErrorMessage, "deadlock")
}

func TestReserveStock_VariosItems_UnoFallaOtroNo(t *testing.T) {
	repo := repoStockBase(true)
	repo.ReserveStockTxFn = func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
		if params.ProductID == "prod-malo" {
			return nil, errors.New("sin stock")
		}
		return &dtos.ReserveStockTxResult{Reserved: params.Quantity, Sufficient: true}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, []dtos.OrderInventoryItem{
		{ProductID: "prod-bueno", Quantity: 3},
		{ProductID: "prod-malo", Quantity: 2},
	})

	require.NoError(t, err)
	require.Len(t, got.ItemResults, 2)
	assert.Equal(t, 3, got.ItemResults[0].Processed)
	assert.False(t, got.Success)
}

func TestReserveStock_PublicaEventoSegunSuficiencia(t *testing.T) {
	cases := []struct {
		nombre     string
		suficiente bool
		wantEvento string
	}{
		{nombre: "stock suficiente", suficiente: true, wantEvento: "inventory.reserved"},
		{nombre: "stock insuficiente", suficiente: false, wantEvento: "inventory.insufficient"},
	}

	for _, tc := range cases {
		t.Run(tc.nombre, func(t *testing.T) {
			pub := &publisherSeguro{}
			repo := repoStockBase(true)
			repo.ReserveStockTxFn = func(ctx context.Context, params dtos.ReserveStockTxParams) (*dtos.ReserveStockTxResult, error) {
				return &dtos.ReserveStockTxResult{Reserved: 5, Sufficient: tc.suficiente}, nil
			}
			uc := New(repo, nil, pub, &mocks.LoggerMock{})

			_, err := uc.ReserveStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

			require.NoError(t, err)
			assert.Eventually(t, func() bool {
				return len(pub.tipos()) >= 1
			}, 2e9, 5e6)
			assert.Contains(t, pub.tipos(), tc.wantEvento)
		})
	}
}

func TestReleaseStock_ProductoInexistente_ContinuaSinFallar(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetProductByIDFn = func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
		return "", "", false, errors.New("not found")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReleaseStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, "producto no encontrado", got.ItemResults[0].ErrorMessage)
	assert.True(t, got.Success)
}

func TestReleaseStock_ProductoSinRastreo_SeDaPorProcesado(t *testing.T) {
	liberado := false
	repo := repoStockBase(false)
	repo.ReleaseStockTxFn = func(ctx context.Context, params dtos.ReleaseTxParams) error {
		liberado = true
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReleaseStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, 5, got.ItemResults[0].Processed)
	assert.False(t, liberado)
}

func TestReleaseStock_Exitoso(t *testing.T) {
	var recibido dtos.ReleaseTxParams
	repo := repoStockBase(true)
	repo.ReleaseStockTxFn = func(ctx context.Context, params dtos.ReleaseTxParams) error {
		recibido = params
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReleaseStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, 5, got.ItemResults[0].Processed)
	assert.True(t, got.ItemResults[0].Sufficient)
	assert.Equal(t, "order-1", recibido.OrderID)
	assert.Equal(t, 5, recibido.Quantity)
}

func TestReleaseStock_FallaTransaccion_MarcaFallida(t *testing.T) {
	repo := repoStockBase(true)
	repo.ReleaseStockTxFn = func(ctx context.Context, params dtos.ReleaseTxParams) error {
		return errors.New("nada que liberar")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReleaseStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.False(t, got.Success)
	assert.Contains(t, got.ItemResults[0].ErrorMessage, "nada que liberar")
}

func TestReleaseStock_SinBodega_RetornaError(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		return 0, errors.New("sin bodega")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReleaseStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, domainerrors.ErrNoDefaultWarehouse)
	assert.Nil(t, got)
}

func TestReleaseStock_SinTipoDeMovimiento_RetornaError(t *testing.T) {
	dbErr := errors.New("catalogo vacio")
	repo := repoStockBase(true)
	repo.GetMovementTypeIDByCodeFn = func(ctx context.Context, code string) (uint, error) {
		return 0, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReleaseStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestConfirmSale_ProductoInexistente_ContinuaSinFallar(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetProductByIDFn = func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
		return "", "", false, errors.New("not found")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmSaleForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, "producto no encontrado", got.ItemResults[0].ErrorMessage)
}

func TestConfirmSale_ProductoSinRastreo_SeDaPorProcesado(t *testing.T) {
	confirmado := false
	repo := repoStockBase(false)
	repo.ConfirmSaleTxFn = func(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
		confirmado = true
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmSaleForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.Equal(t, 5, got.ItemResults[0].Processed)
	assert.False(t, confirmado)
}

func TestConfirmSale_Exitoso(t *testing.T) {
	var recibido dtos.ConfirmSaleTxParams
	repo := repoStockBase(true)
	repo.ConfirmSaleTxFn = func(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
		recibido = params
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmSaleForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, "order-1", recibido.OrderID)
	assert.Equal(t, 5, recibido.Quantity)
}

func TestConfirmSale_FallaTransaccion_MarcaFallida(t *testing.T) {
	repo := repoStockBase(true)
	repo.ConfirmSaleTxFn = func(ctx context.Context, params dtos.ConfirmSaleTxParams) error {
		return errors.New("reserva no encontrada")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmSaleForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.False(t, got.Success)
}

func TestConfirmSale_SinBodega_RetornaError(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		return 0, errors.New("sin bodega")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmSaleForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, domainerrors.ErrNoDefaultWarehouse)
	assert.Nil(t, got)
}

func TestConfirmSale_SinTipoDeMovimiento_RetornaError(t *testing.T) {
	dbErr := errors.New("catalogo vacio")
	repo := repoStockBase(true)
	repo.GetMovementTypeIDByCodeFn = func(ctx context.Context, code string) (uint, error) {
		return 0, dbErr
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConfirmSaleForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestReturnStock_Exitoso(t *testing.T) {
	var recibido dtos.ReturnStockTxParams
	repo := repoStockBase(true)
	repo.ReturnStockTxFn = func(ctx context.Context, params dtos.ReturnStockTxParams) error {
		recibido = params
		return nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReturnStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.True(t, got.Success)
	assert.Equal(t, "order-1", recibido.OrderID)
	assert.Equal(t, 5, recibido.Quantity)
}

func TestReturnStock_FallaTransaccion_MarcaFallida(t *testing.T) {
	repo := repoStockBase(true)
	repo.ReturnStockTxFn = func(ctx context.Context, params dtos.ReturnStockTxParams) error {
		return errors.New("tx failed")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReturnStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	require.NoError(t, err)
	assert.False(t, got.Success)
}

func TestReturnStock_SinBodega_RetornaError(t *testing.T) {
	repo := repoStockBase(true)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		return 0, errors.New("sin bodega")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ReturnStockForOrder(context.Background(), "order-1", 10, nil, itemsDePrueba())

	assert.ErrorIs(t, err, domainerrors.ErrNoDefaultWarehouse)
	assert.Nil(t, got)
}

func TestOperacionesDeStock_SinItems_ResultadoExitosoVacio(t *testing.T) {
	uc := buildUseCase(repoStockBase(true), nil, nil)
	ctx := context.Background()

	reserva, err := uc.ReserveStockForOrder(ctx, "order-1", 10, nil, nil)
	require.NoError(t, err)
	assert.True(t, reserva.Success)
	assert.Empty(t, reserva.ItemResults)

	liberacion, err := uc.ReleaseStockForOrder(ctx, "order-1", 10, nil, nil)
	require.NoError(t, err)
	assert.True(t, liberacion.Success)

	confirmacion, err := uc.ConfirmSaleForOrder(ctx, "order-1", 10, nil, nil)
	require.NoError(t, err)
	assert.True(t, confirmacion.Success)

	devolucion, err := uc.ReturnStockForOrder(ctx, "order-1", 10, nil, nil)
	require.NoError(t, err)
	assert.True(t, devolucion.Success)
}

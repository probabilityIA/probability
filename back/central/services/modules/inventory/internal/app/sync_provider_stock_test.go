package app

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type publisherSeguro struct {
	mu      sync.Mutex
	eventos []ports.InventoryEvent
}

func (p *publisherSeguro) PublishInventoryEvent(ctx context.Context, event ports.InventoryEvent) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.eventos = append(p.eventos, event)
	return nil
}

func (p *publisherSeguro) tipos() []string {
	p.mu.Lock()
	defer p.mu.Unlock()
	out := make([]string, len(p.eventos))
	for i, e := range p.eventos {
		out[i] = e.EventType
	}
	return out
}

func repoParaSync(productoExiste, rastrea bool, cantidadActual int, adjustErr error) *mocks.RepositoryMock {
	return &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			if code == "inbound" {
				return 1, nil
			}
			return 2, nil
		},
		GetProductBySKUFn: func(ctx context.Context, sku string, businessID uint) (string, string, bool, error) {
			if !productoExiste {
				return "", "", false, errors.New("not found")
			}
			return "prod-001", "Producto", rastrea, nil
		},
		GetDefaultWarehouseIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 5, nil
		},
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			return []entities.InventoryLevel{{WarehouseID: 5, Quantity: cantidadActual}}, nil
		},
		AdjustStockTxFn: func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
			if adjustErr != nil {
				return nil, adjustErr
			}
			return &dtos.AdjustStockTxResult{}, nil
		},
	}
}

func syncDTO(items ...request.ProviderStockSyncItem) request.ProviderStockSyncDTO {
	return request.ProviderStockSyncDTO{
		BusinessID:    10,
		IntegrationID: 197,
		Provider:      "siigo",
		CorrelationID: "corr-1",
		Items:         items,
	}
}

func TestSyncProviderStock_SinItems_ResultadoVacio(t *testing.T) {
	uc := buildUseCase(repoParaSync(true, true, 0, nil), nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO())

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, 0, got.Total)
	assert.Equal(t, 0, got.Updated)
}

func TestSyncProviderStock_ProductoInexistente_SeOmite(t *testing.T) {
	uc := buildUseCase(repoParaSync(false, true, 0, nil), nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-X", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Skipped)
	assert.Equal(t, 0, got.Updated)
	assert.Equal(t, 0, got.Failed)
}

func TestSyncProviderStock_ProductoSinRastreo_SeOmite(t *testing.T) {
	uc := buildUseCase(repoParaSync(true, false, 0, nil), nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Skipped, "un producto que no rastrea inventario no se sincroniza")
}

func TestSyncProviderStock_MismaCantidad_SinCambios(t *testing.T) {
	ajustes := 0
	repo := repoParaSync(true, true, 10, nil)
	adjustOriginal := repo.AdjustStockTxFn
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		ajustes++
		return adjustOriginal(ctx, params)
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Unchanged)
	assert.Equal(t, 0, ajustes, "sin diferencia no se toca el inventario")
}

func TestSyncProviderStock_CantidadMayor_AjustaConEntrada(t *testing.T) {
	var recibido dtos.AdjustStockTxParams
	repo := repoParaSync(true, true, 10, nil)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 25},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Updated)
	assert.Equal(t, 15, recibido.Quantity, "el ajuste es el delta, no la cantidad absoluta")
	assert.Equal(t, uint(1), recibido.MovementTypeID, "delta positivo usa el tipo inbound")
	assert.Equal(t, "provider_sync", recibido.ReferenceType)
}

func TestSyncProviderStock_CantidadMenor_AjustaConSalida(t *testing.T) {
	var recibido dtos.AdjustStockTxParams
	repo := repoParaSync(true, true, 30, nil)
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 12},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Updated)
	assert.Equal(t, -18, recibido.Quantity)
	assert.Equal(t, uint(2), recibido.MovementTypeID, "delta negativo usa el tipo outbound")
}

func TestSyncProviderStock_SinBodega_Falla(t *testing.T) {
	repo := repoParaSync(true, true, 0, nil)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		return 0, errors.New("sin bodega")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Failed)
}

func TestSyncProviderStock_BodegaExplicita_NoConsultaLaPorDefecto(t *testing.T) {
	consultoDefault := false
	var recibido dtos.AdjustStockTxParams
	repo := repoParaSync(true, true, 0, nil)
	repo.GetDefaultWarehouseIDFn = func(ctx context.Context, businessID uint) (uint, error) {
		consultoDefault = true
		return 5, nil
	}
	repo.GetProductInventoryFn = func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
		return []entities.InventoryLevel{{WarehouseID: 9, Quantity: 3}}, nil
	}
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10, WarehouseID: 9},
	))

	require.NoError(t, err)
	assert.False(t, consultoDefault)
	assert.Equal(t, uint(9), recibido.WarehouseID)
	assert.Equal(t, 7, recibido.Quantity)
}

func TestSyncProviderStock_FallaTipoDeMovimiento_MarcaFallido(t *testing.T) {
	repo := repoParaSync(true, true, 0, nil)
	repo.GetMovementTypeIDByCodeFn = func(ctx context.Context, code string) (uint, error) {
		return 0, errors.New("catalogo vacio")
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Failed)
}

func TestSyncProviderStock_FallaAjuste_MarcaFallido(t *testing.T) {
	uc := buildUseCase(repoParaSync(true, true, 0, errors.New("tx failed")), nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 1, got.Failed)
	assert.Equal(t, 0, got.Updated)
}

func TestSyncProviderStock_FallaLeerInventario_AsumeCero(t *testing.T) {
	var recibido dtos.AdjustStockTxParams
	repo := repoParaSync(true, true, 0, nil)
	repo.GetProductInventoryFn = func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
		return nil, errors.New("db down")
	}
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 10, recibido.Quantity,
		"si no se puede leer el inventario actual se asume cero y se ajusta la cantidad completa")
}

func TestSyncProviderStock_BodegaSinNivel_AsumeCero(t *testing.T) {
	var recibido dtos.AdjustStockTxParams
	repo := repoParaSync(true, true, 0, nil)
	repo.GetProductInventoryFn = func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
		return []entities.InventoryLevel{{WarehouseID: 99, Quantity: 50}}, nil
	}
	repo.AdjustStockTxFn = func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
		recibido = params
		return &dtos.AdjustStockTxResult{}, nil
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Equal(t, 10, recibido.Quantity, "el nivel de otra bodega no cuenta")
}

func TestSyncProviderStock_MezclaDeResultados(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 1, nil
		},
		GetProductBySKUFn: func(ctx context.Context, sku string, businessID uint) (string, string, bool, error) {
			switch sku {
			case "SKU-OK":
				return "prod-ok", "Producto", true, nil
			case "SKU-IGUAL":
				return "prod-igual", "Producto", true, nil
			default:
				return "", "", false, errors.New("not found")
			}
		},
		GetDefaultWarehouseIDFn: func(ctx context.Context, businessID uint) (uint, error) {
			return 5, nil
		},
		GetProductInventoryFn: func(ctx context.Context, params dtos.GetProductInventoryParams) ([]entities.InventoryLevel, error) {
			if params.ProductID == "prod-igual" {
				return []entities.InventoryLevel{{WarehouseID: 5, Quantity: 7}}, nil
			}
			return []entities.InventoryLevel{{WarehouseID: 5, Quantity: 0}}, nil
		},
		AdjustStockTxFn: func(ctx context.Context, params dtos.AdjustStockTxParams) (*dtos.AdjustStockTxResult, error) {
			return &dtos.AdjustStockTxResult{}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-OK", Quantity: 10},
		request.ProviderStockSyncItem{SKU: "SKU-IGUAL", Quantity: 7},
		request.ProviderStockSyncItem{SKU: "SKU-NO-EXISTE", Quantity: 3},
	))

	require.NoError(t, err)
	assert.Equal(t, 3, got.Total)
	assert.Equal(t, 1, got.Updated)
	assert.Equal(t, 1, got.Unchanged)
	assert.Equal(t, 1, got.Skipped)
}

func TestSyncProviderStock_PublicaEventosDeInicioYFin(t *testing.T) {
	pub := &publisherSeguro{}
	uc := New(repoParaSync(true, true, 0, nil), nil, pub, &mocks.LoggerMock{})

	_, err := uc.SyncProviderStock(context.Background(), syncDTO(
		request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
	))

	require.NoError(t, err)
	assert.Eventually(t, func() bool {
		return len(pub.tipos()) >= 3
	}, 2e9, 5e6)

	tipos := pub.tipos()
	assert.Contains(t, tipos, "inventory.sync.started")
	assert.Contains(t, tipos, "inventory.sync.progress")
	assert.Contains(t, tipos, "inventory.sync.completed")
}

func TestSyncProviderStock_SinPublisher_NoRompe(t *testing.T) {
	uc := buildUseCase(repoParaSync(true, true, 0, nil), nil, nil)

	assert.NotPanics(t, func() {
		_, _ = uc.SyncProviderStock(context.Background(), syncDTO(
			request.ProviderStockSyncItem{SKU: "SKU-1", Quantity: 10},
		))
	})
}

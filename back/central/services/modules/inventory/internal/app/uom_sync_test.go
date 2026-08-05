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

func uomsDeProducto() []entities.ProductUoM {
	return []entities.ProductUoM{
		{UomCode: "und", ConversionFactor: 1, IsBase: true},
		{UomCode: "caja", ConversionFactor: 12},
		{UomCode: "pallet", ConversionFactor: 144},
	}
}

func TestListUoMs_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListUoMsFn: func(ctx context.Context) ([]entities.UnitOfMeasure, error) {
			return []entities.UnitOfMeasure{{Code: "und"}, {Code: "caja"}}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ListUoMs(context.Background())

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestListProductUoMs_PropagaLosFiltros(t *testing.T) {
	var recibido dtos.ListProductUoMParams
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			recibido = params
			return uomsDeProducto(), nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ListProductUoMs(context.Background(), 10, "prod-1")

	require.NoError(t, err)
	assert.Len(t, got, 3)
	assert.Equal(t, uint(10), recibido.BusinessID)
	assert.Equal(t, "prod-1", recibido.ProductID)
}

func TestCreateProductUoM_ProductoInexistente_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "", "", false, errors.New("not found")
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateProductUoM(context.Background(), request.CreateProductUoMDTO{
		ProductID:  "prod-1",
		BusinessID: 10,
	})

	assert.ErrorIs(t, err, domainerrors.ErrProductNotFound)
	assert.Nil(t, got)
}

func TestCreateProductUoM_UnidadInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("uom no existe")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		GetUoMByCodeFn: func(ctx context.Context, code string) (*entities.UnitOfMeasure, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.CreateProductUoM(context.Background(), request.CreateProductUoMDTO{
		ProductID:  "prod-1",
		BusinessID: 10,
		UomCode:    "inventada",
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateProductUoM_FactorInvalidoUsaUno(t *testing.T) {
	for _, factor := range []float64{0, -3} {
		var guardada *entities.ProductUoM
		repo := &mocks.RepositoryMock{
			GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
				return "SKU-1", "Producto", true, nil
			},
			GetUoMByCodeFn: func(ctx context.Context, code string) (*entities.UnitOfMeasure, error) {
				return &entities.UnitOfMeasure{ID: 2, Code: code}, nil
			},
			CreateProductUoMFn: func(ctx context.Context, pu *entities.ProductUoM) (*entities.ProductUoM, error) {
				guardada = pu
				return pu, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, err := uc.CreateProductUoM(context.Background(), request.CreateProductUoMDTO{
			ProductID:        "prod-1",
			BusinessID:       10,
			UomCode:          "caja",
			ConversionFactor: factor,
		})

		require.NoError(t, err)
		assert.InDelta(t, 1.0, guardada.ConversionFactor, 0.001,
			"un factor de conversion invalido cae a 1 en vez de guardarse")
	}
}

func TestCreateProductUoM_Exitoso(t *testing.T) {
	var guardada *entities.ProductUoM
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, productID string, businessID uint) (string, string, bool, error) {
			return "SKU-1", "Producto", true, nil
		},
		GetUoMByCodeFn: func(ctx context.Context, code string) (*entities.UnitOfMeasure, error) {
			return &entities.UnitOfMeasure{ID: 2, Code: code}, nil
		},
		CreateProductUoMFn: func(ctx context.Context, pu *entities.ProductUoM) (*entities.ProductUoM, error) {
			guardada = pu
			return pu, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateProductUoM(context.Background(), request.CreateProductUoMDTO{
		ProductID:        "prod-1",
		BusinessID:       10,
		UomCode:          "caja",
		ConversionFactor: 12,
		IsBase:           false,
		Barcode:          "7701234567890",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(2), guardada.UomID)
	assert.InDelta(t, 12.0, guardada.ConversionFactor, 0.001)
	assert.Equal(t, "7701234567890", guardada.Barcode)
	assert.True(t, guardada.IsActive, "una unidad nueva nace activa")
}

func TestDeleteProductUoM_DelegaEnRepositorio(t *testing.T) {
	var borrada uint
	repo := &mocks.RepositoryMock{
		DeleteProductUoMFn: func(ctx context.Context, businessID, id uint) error {
			borrada = id
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	require.NoError(t, uc.DeleteProductUoM(context.Background(), 10, 3))
	assert.Equal(t, uint(3), borrada)
}

func TestConvertUoM_SinUnidadesConfiguradas_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{ProductID: "prod-1", BusinessID: 10})

	assert.ErrorIs(t, err, domainerrors.ErrProductUoMNotFound)
	assert.Nil(t, got)
}

func TestConvertUoM_CajaAUnidades(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			return uomsDeProducto(), nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{
		ProductID:   "prod-1",
		BusinessID:  10,
		FromUomCode: "caja",
		ToUomCode:   "und",
		Quantity:    3,
	})

	require.NoError(t, err)
	assert.InDelta(t, 36.0, got.ConvertedQty, 0.001, "3 cajas de 12 son 36 unidades")
	assert.InDelta(t, 36.0, got.BaseUnitQuantity, 0.001)
	assert.Equal(t, "und", got.BaseUomCode)
}

func TestConvertUoM_UnidadesACaja(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			return uomsDeProducto(), nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{
		ProductID:   "prod-1",
		BusinessID:  10,
		FromUomCode: "und",
		ToUomCode:   "caja",
		Quantity:    24,
	})

	require.NoError(t, err)
	assert.InDelta(t, 2.0, got.ConvertedQty, 0.001, "24 unidades son 2 cajas")
}

func TestConvertUoM_EntreUnidadesNoBase(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			return uomsDeProducto(), nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{
		ProductID:   "prod-1",
		BusinessID:  10,
		FromUomCode: "pallet",
		ToUomCode:   "caja",
		Quantity:    2,
	})

	require.NoError(t, err)
	assert.InDelta(t, 24.0, got.ConvertedQty, 0.001, "2 pallets de 144 son 24 cajas de 12")
	assert.InDelta(t, 288.0, got.BaseUnitQuantity, 0.001)
}

func TestConvertUoM_UnidadDesconocida_RetornaError(t *testing.T) {
	cases := []struct {
		nombre string
		desde  string
		hasta  string
	}{
		{nombre: "origen desconocido", desde: "inventada", hasta: "und"},
		{nombre: "destino desconocido", desde: "und", hasta: "inventada"},
	}

	for _, tc := range cases {
		t.Run(tc.nombre, func(t *testing.T) {
			repo := &mocks.RepositoryMock{
				ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
					return uomsDeProducto(), nil
				},
			}
			uc := buildUseCase(repo, nil, nil)

			got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{
				ProductID:   "prod-1",
				BusinessID:  10,
				FromUomCode: tc.desde,
				ToUomCode:   tc.hasta,
				Quantity:    1,
			})

			assert.ErrorIs(t, err, domainerrors.ErrUomConversion)
			assert.Nil(t, got)
		})
	}
}

func TestConvertUoM_SinUnidadBase_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			return []entities.ProductUoM{
				{UomCode: "caja", ConversionFactor: 12},
				{UomCode: "und", ConversionFactor: 1},
			}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{
		ProductID:   "prod-1",
		BusinessID:  10,
		FromUomCode: "caja",
		ToUomCode:   "und",
		Quantity:    1,
	})

	assert.ErrorIs(t, err, domainerrors.ErrUomConversion,
		"sin unidad base marcada no se puede convertir")
	assert.Nil(t, got)
}

func TestConvertUoM_FallaListar_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ListProductUoMsFn: func(ctx context.Context, params dtos.ListProductUoMParams) ([]entities.ProductUoM, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ConvertUoM(context.Background(), request.ConvertUoMDTO{ProductID: "prod-1", BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestComputePayloadHash_EsDeterministaYSensible(t *testing.T) {
	a := map[string]any{"sku": "SKU-1", "qty": 10}
	b := map[string]any{"qty": 10, "sku": "SKU-1"}
	c := map[string]any{"sku": "SKU-1", "qty": 11}

	assert.Equal(t, computePayloadHash(a), computePayloadHash(b),
		"el mismo contenido produce el mismo hash sin importar el orden de las claves")
	assert.NotEqual(t, computePayloadHash(a), computePayloadHash(c))
	assert.Len(t, computePayloadHash(a), 64, "sha256 en hexadecimal son 64 caracteres")
}

func TestComputePayloadHash_PayloadVacio(t *testing.T) {
	assert.NotEmpty(t, computePayloadHash(nil))
	assert.NotEmpty(t, computePayloadHash(map[string]any{}))
}

func TestInboundSync_PayloadRepetido_SeMarcaDuplicado(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		GetSyncLogByHashFn: func(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error) {
			return &entities.InventorySyncLog{ID: 1, PayloadHash: hash}, nil
		},
		CreateSyncLogFn: func(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error) {
			creado = true
			return log, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.InboundSync(context.Background(), request.InboundSyncDTO{
		BusinessID:    10,
		IntegrationID: 197,
		Payload:       map[string]any{"sku": "SKU-1"},
	})

	require.NoError(t, err)
	assert.True(t, got.Duplicate)
	assert.False(t, creado, "un payload ya recibido no se vuelve a registrar")
}

func TestInboundSync_PayloadNuevo_SeRegistra(t *testing.T) {
	var guardado *entities.InventorySyncLog
	var direccionConsultada string
	repo := &mocks.RepositoryMock{
		GetSyncLogByHashFn: func(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error) {
			direccionConsultada = direction
			return nil, nil
		},
		CreateSyncLogFn: func(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error) {
			guardado = log
			return log, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.InboundSync(context.Background(), request.InboundSyncDTO{
		BusinessID:    10,
		IntegrationID: 197,
		Payload:       map[string]any{"sku": "SKU-1"},
	})

	require.NoError(t, err)
	assert.False(t, got.Duplicate)
	require.NotNil(t, guardado)
	assert.Equal(t, "in", direccionConsultada)
	assert.Equal(t, "in", guardado.Direction)
	assert.Equal(t, "received", guardado.Status)
	assert.Len(t, guardado.PayloadHash, 64)
	require.NotNil(t, guardado.IntegrationID)
	assert.Equal(t, uint(197), *guardado.IntegrationID)
	assert.NotNil(t, guardado.SyncedAt)
}

func TestInboundSync_FallaConsultarHash_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetSyncLogByHashFn: func(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.InboundSync(context.Background(), request.InboundSyncDTO{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestInboundSync_FallaRegistrar_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		GetSyncLogByHashFn: func(ctx context.Context, businessID uint, direction, hash string) (*entities.InventorySyncLog, error) {
			return nil, nil
		},
		CreateSyncLogFn: func(ctx context.Context, log *entities.InventorySyncLog) (*entities.InventorySyncLog, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.InboundSync(context.Background(), request.InboundSyncDTO{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListSyncLogs_NormalizaPaginacion(t *testing.T) {
	var recibido dtos.ListSyncLogsParams
	repo := &mocks.RepositoryMock{
		ListSyncLogsFn: func(ctx context.Context, params dtos.ListSyncLogsParams) ([]entities.InventorySyncLog, int64, error) {
			recibido = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, _, err := uc.ListSyncLogs(context.Background(), dtos.ListSyncLogsParams{Page: 0, PageSize: 400})

	require.NoError(t, err)
	assert.Equal(t, 1, recibido.Page)
	assert.Equal(t, 10, recibido.PageSize)
}

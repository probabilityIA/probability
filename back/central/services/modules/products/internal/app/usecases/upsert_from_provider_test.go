package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func f(v float64) *float64 {
	return &v
}

func upsertDTO() *domain.ProductProviderUpsertDTO {
	return &domain.ProductProviderUpsertDTO{
		BusinessID:     10,
		IntegrationID:  197,
		SKU:            "SKU-1",
		Name:           "Camisa",
		ExternalID:     "EXT-1",
		Price:          25000,
		TrackInventory: true,
	}
}

func TestUpsertFromProvider_DTONulo_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	err := uc.UpsertFromProvider(context.Background(), nil)

	assert.ErrorIs(t, err, domain.ErrInvalidProductData)
}

func TestUpsertFromProvider_SinBusinessID_RetornaError(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			consultado = true
			return nil, nil
		},
	}
	uc := New(repo)

	dto := upsertDTO()
	dto.BusinessID = 0

	err := uc.UpsertFromProvider(context.Background(), dto)

	assert.ErrorIs(t, err, domain.ErrInvalidProductData)
	assert.False(t, consultado)
}

func TestUpsertFromProvider_SinSKU_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	dto := upsertDTO()
	dto.SKU = ""

	err := uc.UpsertFromProvider(context.Background(), dto)

	assert.ErrorIs(t, err, domain.ErrInvalidProductData)
}

func TestUpsertFromProvider_ProductoNuevo_LoCreaEnCOPYActivo(t *testing.T) {
	var creado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			product.ID = "prod-nuevo"
			creado = product
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return false, nil
		},
	}
	uc := New(repo)

	err := uc.UpsertFromProvider(context.Background(), upsertDTO())

	require.NoError(t, err)
	require.NotNil(t, creado)
	assert.Equal(t, "SKU-1", creado.SKU)
	assert.Equal(t, "Camisa", creado.Name)
	assert.Equal(t, "COP", creado.Currency, "los productos de proveedor nacen en pesos colombianos")
	assert.Equal(t, "active", creado.Status)
	assert.True(t, creado.IsActive)
	assert.True(t, creado.TrackInventory)
}

func TestUpsertFromProvider_ProductoNuevo_CreaLaRelacionConElCanal(t *testing.T) {
	var mapeoProducto string
	var mapeoIntegracion uint
	var mapeoExterno string
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			product.ID = "prod-nuevo"
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return false, nil
		},
		AddProductIntegrationFn: func(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
			mapeoProducto, mapeoIntegracion, mapeoExterno = productID, integrationID, externalProductID
			return &domain.ProductBusinessIntegration{}, nil
		},
	}
	uc := New(repo)

	err := uc.UpsertFromProvider(context.Background(), upsertDTO())

	require.NoError(t, err)
	assert.Equal(t, "prod-nuevo", mapeoProducto)
	assert.Equal(t, uint(197), mapeoIntegracion,
		"todo producto creado desde una integracion queda asociado a su canal de origen")
	assert.Equal(t, "EXT-1", mapeoExterno)
}

func TestUpsertFromProvider_FallaCrearProducto_NoMapeaIntegracion(t *testing.T) {
	dbErr := errors.New("insert failed")
	mapeado := false
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return nil, domain.ErrProductNotFound
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			return dbErr
		},
		AddProductIntegrationFn: func(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
			mapeado = true
			return nil, nil
		},
	}
	uc := New(repo)

	err := uc.UpsertFromProvider(context.Background(), upsertDTO())

	require.Error(t, err)
	assert.False(t, mapeado)
}

func TestUpsertFromProvider_ErrorDistintoDeNoEncontrado_SePropaga(t *testing.T) {
	dbErr := errors.New("db down")
	creado := false
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return nil, dbErr
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			creado = true
			return nil
		},
	}
	uc := New(repo)

	err := uc.UpsertFromProvider(context.Background(), upsertDTO())

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, creado, "un fallo de base de datos no debe crear un producto duplicado")
}

func TestUpsertFromProvider_ProductoExistente_LoActualiza(t *testing.T) {
	var actualizado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return &domain.Product{ID: "prod-1", BusinessID: businessID, SKU: sku, Name: "Nombre viejo"}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return &domain.Product{ID: id, BusinessID: businessID, SKU: "SKU-1", Name: "Nombre viejo"}, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			actualizado = product
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return true, nil
		},
	}
	uc := New(repo)

	err := uc.UpsertFromProvider(context.Background(), upsertDTO())

	require.NoError(t, err)
	require.NotNil(t, actualizado)
	assert.Equal(t, "Camisa", actualizado.Name)
}

func TestUpsertFromProvider_NoPisaLaImagenExistente(t *testing.T) {
	var actualizado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return &domain.Product{ID: "prod-1", SKU: sku, ImageURL: "https://s3/imagen-propia.jpg"}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return &domain.Product{ID: id, SKU: "SKU-1", ImageURL: "https://s3/imagen-propia.jpg"}, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			actualizado = product
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return true, nil
		},
	}
	uc := New(repo)

	dto := upsertDTO()
	dto.ImageURL = "https://proveedor/otra.jpg"

	err := uc.UpsertFromProvider(context.Background(), dto)

	require.NoError(t, err)
	require.NotNil(t, actualizado)
	assert.Equal(t, "https://s3/imagen-propia.jpg", actualizado.ImageURL,
		"si el producto ya tiene imagen propia, el proveedor no la reemplaza")
}

func TestUpsertFromProvider_AsignaImagenSiNoTenia(t *testing.T) {
	var actualizado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return &domain.Product{ID: "prod-1", SKU: sku, ImageURL: ""}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return &domain.Product{ID: id, SKU: "SKU-1", ImageURL: ""}, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			actualizado = product
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return true, nil
		},
	}
	uc := New(repo)

	dto := upsertDTO()
	dto.ImageURL = "https://proveedor/foto.jpg"

	err := uc.UpsertFromProvider(context.Background(), dto)

	require.NoError(t, err)
	assert.Equal(t, "https://proveedor/foto.jpg", actualizado.ImageURL)
}

func TestUpsertFromProvider_FallaActualizar_PropagaError(t *testing.T) {
	dbErr := errors.New("update failed")
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return &domain.Product{ID: "prod-1", SKU: sku}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return &domain.Product{ID: id, SKU: "SKU-1"}, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			return dbErr
		},
	}
	uc := New(repo)

	err := uc.UpsertFromProvider(context.Background(), upsertDTO())

	require.Error(t, err)
}

func TestEnsureIntegrationMapping_SinIntegracion_NoHaceNada(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			consultado = true
			return false, nil
		},
	}
	uc := New(repo)

	require.NoError(t, uc.ensureIntegrationMapping(context.Background(), "prod-1", 0, "EXT-1"))
	assert.False(t, consultado)
}

func TestEnsureIntegrationMapping_SinProducto_NoHaceNada(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			consultado = true
			return false, nil
		},
	}
	uc := New(repo)

	require.NoError(t, uc.ensureIntegrationMapping(context.Background(), "", 197, "EXT-1"))
	assert.False(t, consultado)
}

func TestEnsureIntegrationMapping_YaExiste_NoDuplica(t *testing.T) {
	agregado := false
	repo := &mocks.RepositoryMock{
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return true, nil
		},
		AddProductIntegrationFn: func(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
			agregado = true
			return nil, nil
		},
	}
	uc := New(repo)

	require.NoError(t, uc.ensureIntegrationMapping(context.Background(), "prod-1", 197, "EXT-1"))
	assert.False(t, agregado, "la relacion producto-canal es idempotente")
}

func TestEnsureIntegrationMapping_FallaVerificar_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	assert.ErrorIs(t, uc.ensureIntegrationMapping(context.Background(), "prod-1", 197, "EXT-1"), dbErr)
}

func TestEnsureIntegrationMapping_FallaAgregar_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return false, nil
		},
		AddProductIntegrationFn: func(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	assert.ErrorIs(t, uc.ensureIntegrationMapping(context.Background(), "prod-1", 197, "EXT-1"), dbErr)
}

func TestUpsertFromProvider_UnidadesDePesoYDimension(t *testing.T) {
	var actualizado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return &domain.Product{ID: "prod-1", SKU: sku}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return &domain.Product{ID: id, SKU: "SKU-1"}, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			actualizado = product
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return true, nil
		},
	}
	uc := New(repo)

	dto := upsertDTO()
	dto.Weight = f(1.5)
	dto.WeightUnit = "kg"
	dto.Length = f(10)
	dto.Width = f(20)
	dto.Height = f(30)
	dto.DimensionUnit = "cm"

	err := uc.UpsertFromProvider(context.Background(), dto)

	require.NoError(t, err)
	require.NotNil(t, actualizado)
	assert.Equal(t, "kg", actualizado.WeightUnit)
	assert.Equal(t, "cm", actualizado.DimensionUnit)
}

func TestUpsertFromProvider_SinPeso_NoAsignaUnidadDePeso(t *testing.T) {
	var actualizado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductBySKUFn: func(ctx context.Context, businessID uint, sku string) (*domain.Product, error) {
			return &domain.Product{ID: "prod-1", SKU: sku, WeightUnit: "lb"}, nil
		},
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return &domain.Product{ID: id, SKU: "SKU-1", WeightUnit: "lb"}, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			actualizado = product
			return nil
		},
		ProductIntegrationExistsFn: func(ctx context.Context, productID string, integrationID uint) (bool, error) {
			return true, nil
		},
	}
	uc := New(repo)

	dto := upsertDTO()
	dto.WeightUnit = "kg"

	err := uc.UpsertFromProvider(context.Background(), dto)

	require.NoError(t, err)
	assert.Equal(t, "lb", actualizado.WeightUnit,
		"sin peso nuevo la unidad existente no se toca")
}

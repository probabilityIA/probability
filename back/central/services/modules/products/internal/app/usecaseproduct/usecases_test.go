package usecaseproduct

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/datatypes"
)

func s(v string) *string {
	return &v
}

func u(v uint) *uint {
	return &v
}

func productoBase(id string) *domain.Product {
	return &domain.Product{
		ID:         id,
		BusinessID: 10,
		SKU:        "SKU-1",
		Name:       "Camisa",
		Price:      25000,
		Currency:   "COP",
		Status:     "active",
		IsActive:   true,
	}
}

func crearRequest() *domain.CreateProductRequest {
	return &domain.CreateProductRequest{
		BusinessID: 10,
		SKU:        "SKU-1",
		Name:       "Camisa",
		Price:      25000,
		Currency:   "COP",
		Status:     "active",
		IsActive:   true,
	}
}

func TestCreateProduct_SKUDuplicado_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		ProductExistsFn: func(ctx context.Context, businessID uint, sku string) (bool, error) {
			return true, nil
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			creado = true
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.CreateProduct(context.Background(), crearRequest())

	assert.ErrorIs(t, err, domain.ErrProductAlreadyExists)
	assert.Nil(t, got)
	assert.False(t, creado)
}

func TestCreateProduct_FallaVerificarSKU_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ProductExistsFn: func(ctx context.Context, businessID uint, sku string) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.CreateProduct(context.Background(), crearRequest())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateProduct_Exitoso(t *testing.T) {
	var creado *domain.Product
	repo := &mocks.RepositoryMock{
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			product.ID = "prod-1"
			creado = product
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.CreateProduct(context.Background(), crearRequest())

	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, creado)
	assert.Equal(t, "SKU-1", creado.SKU)
	assert.Equal(t, "Camisa", creado.Name)
	assert.Nil(t, creado.FamilyID, "sin familia el producto no queda asociado a ninguna")
}

func TestCreateProduct_ConFamiliaExistente(t *testing.T) {
	var creado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return &domain.ProductFamily{ID: familyID, BusinessID: businessID, Name: "Camisas"}, nil
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			creado = product
			return nil
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.FamilyID = u(7)

	_, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, creado.FamilyID)
	assert.Equal(t, uint(7), *creado.FamilyID)
}

func TestCreateProduct_FamiliaInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("family not found")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.FamilyID = u(7)

	got, err := uc.CreateProduct(context.Background(), req)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateProduct_CreaFamiliaEnLinea(t *testing.T) {
	var familiaCreada *domain.ProductFamily
	repo := &mocks.RepositoryMock{
		CreateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			family.ID = 9
			familiaCreada = family
			return nil
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			return nil
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.Family = &domain.CreateProductFamilyRequest{Name: "Camisas nuevas"}

	_, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, familiaCreada)
	assert.Equal(t, "Camisas nuevas", familiaCreada.Name)
	assert.True(t, familiaCreada.IsActive, "la familia creada en linea nace activa")
	assert.Equal(t, uint(10), familiaCreada.BusinessID)
}

func TestCreateProduct_FallaCrearFamiliaEnLinea_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	creado := false
	repo := &mocks.RepositoryMock{
		CreateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			return dbErr
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			creado = true
			return nil
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.Family = &domain.CreateProductFamilyRequest{Name: "Camisas"}

	got, err := uc.CreateProduct(context.Background(), req)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.False(t, creado)
}

func TestCreateProduct_AtributosDeVarianteInvalidos_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{}
	uc := New(repo)

	req := crearRequest()
	req.VariantAttributes = datatypes.JSON([]byte(`{roto`))

	got, err := uc.CreateProduct(context.Background(), req)

	assert.ErrorIs(t, err, domain.ErrInvalidProductData)
	assert.Nil(t, got)
}

func TestCreateProduct_VarianteDuplicadaEnLaFamilia_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return &domain.ProductFamily{ID: familyID}, nil
		},
		VariantExistsInFamilyFn: func(ctx context.Context, businessID, familyID uint, variantSignature string, excludeProductID *string) (bool, error) {
			return true, nil
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			creado = true
			return nil
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.FamilyID = u(7)
	req.VariantAttributes = datatypes.JSON([]byte(`{"talla":"M","color":"rojo"}`))

	got, err := uc.CreateProduct(context.Background(), req)

	assert.ErrorIs(t, err, domain.ErrVariantAlreadyExists)
	assert.Nil(t, got)
	assert.False(t, creado, "no se crea una segunda variante identica dentro de la misma familia")
}

func TestCreateProduct_VarianteSinFamilia_NoVerificaUnicidad(t *testing.T) {
	verificado := false
	repo := &mocks.RepositoryMock{
		VariantExistsInFamilyFn: func(ctx context.Context, businessID, familyID uint, variantSignature string, excludeProductID *string) (bool, error) {
			verificado = true
			return false, nil
		},
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			return nil
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.VariantAttributes = datatypes.JSON([]byte(`{"talla":"M"}`))

	_, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	assert.False(t, verificado, "sin familia no hay unicidad de variante que verificar")
}

func TestCreateProduct_GuardaLaFirmaCanonicaDeLaVariante(t *testing.T) {
	var creado *domain.Product
	repo := &mocks.RepositoryMock{
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			creado = product
			return nil
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.VariantAttributes = datatypes.JSON([]byte(`{"Talla":"  M  ","color":"rojo"}`))

	_, err := uc.CreateProduct(context.Background(), req)

	require.NoError(t, err)
	assert.Contains(t, creado.VariantSignature, `"talla":"M"`,
		"la firma guardada es la version canonica, no el payload crudo")
}

func TestCreateProduct_FallaVerificarVariante_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return &domain.ProductFamily{ID: familyID}, nil
		},
		VariantExistsInFamilyFn: func(ctx context.Context, businessID, familyID uint, variantSignature string, excludeProductID *string) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	req := crearRequest()
	req.FamilyID = u(7)
	req.VariantAttributes = datatypes.JSON([]byte(`{"talla":"M"}`))

	got, err := uc.CreateProduct(context.Background(), req)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreateProduct_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		CreateProductFn: func(ctx context.Context, product *domain.Product) error {
			return dbErr
		},
	}
	uc := New(repo)

	got, err := uc.CreateProduct(context.Background(), crearRequest())

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetProductByID_SinID_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	got, err := uc.GetProductByID(context.Background(), 10, "")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "product ID is required")
	assert.Nil(t, got)
}

func TestGetProductByID_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return nil, nil
		},
	}
	uc := New(repo)

	got, err := uc.GetProductByID(context.Background(), 10, "prod-1")

	assert.ErrorIs(t, err, domain.ErrProductNotFound)
	assert.Nil(t, got)
}

func TestGetProductByID_Existe_RetornaResponse(t *testing.T) {
	var pedidoBusiness uint
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			pedidoBusiness = businessID
			return productoBase(id), nil
		},
	}
	uc := New(repo)

	got, err := uc.GetProductByID(context.Background(), 10, "prod-1")

	require.NoError(t, err)
	assert.Equal(t, "prod-1", got.ID)
	assert.Equal(t, uint(10), pedidoBusiness,
		"la consulta siempre va filtrada por el negocio del token")
}

func TestGetProductByID_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.GetProductByID(context.Background(), 10, "prod-1")

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListProducts_NormalizaPaginacion(t *testing.T) {
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
			var gotPage, gotPageSize int
			repo := &mocks.RepositoryMock{
				ListProductsFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
					gotPage, gotPageSize = page, pageSize
					return nil, 0, nil
				},
			}
			uc := New(repo)

			got, err := uc.ListProducts(context.Background(), 10, tc.page, tc.pageSize, nil)

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, gotPage)
			assert.Equal(t, tc.wantPageSize, gotPageSize)
			assert.Equal(t, tc.wantPage, got.Page)
		})
	}
}

func TestListProducts_CalculaTotalDePaginas(t *testing.T) {
	cases := []struct {
		total     int64
		pageSize  int
		wantPages int
	}{
		{total: 0, pageSize: 10, wantPages: 0},
		{total: 1, pageSize: 10, wantPages: 1},
		{total: 10, pageSize: 10, wantPages: 1},
		{total: 11, pageSize: 10, wantPages: 2},
	}

	for _, tc := range cases {
		repo := &mocks.RepositoryMock{
			ListProductsFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
				return nil, tc.total, nil
			},
		}
		uc := New(repo)

		got, err := uc.ListProducts(context.Background(), 10, 1, tc.pageSize, nil)

		require.NoError(t, err)
		assert.Equal(t, tc.wantPages, got.TotalPages)
	}
}

func TestListProducts_FiltraPorNegocio(t *testing.T) {
	var pedidoBusiness uint
	repo := &mocks.RepositoryMock{
		ListProductsFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
			pedidoBusiness = businessID
			return []domain.Product{*productoBase("prod-1")}, 1, nil
		},
	}
	uc := New(repo)

	got, err := uc.ListProducts(context.Background(), 10, 1, 10, nil)

	require.NoError(t, err)
	assert.Len(t, got.Data, 1)
	assert.Equal(t, uint(10), pedidoBusiness)
}

func TestListProducts_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("query failed")
	repo := &mocks.RepositoryMock{
		ListProductsFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.Product, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.ListProducts(context.Background(), 10, 1, 10, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateProduct_SinID_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	got, err := uc.UpdateProduct(context.Background(), 10, "", &domain.UpdateProductRequest{})

	require.Error(t, err)
	assert.Nil(t, got)
}

func TestUpdateProduct_NoExiste_RetornaNotFound(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return nil, nil
		},
	}
	uc := New(repo)

	got, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{})

	assert.ErrorIs(t, err, domain.ErrProductNotFound)
	assert.Nil(t, got)
}

func TestUpdateProduct_SKUDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		ProductExistsFn: func(ctx context.Context, businessID uint, sku string) (bool, error) {
			return true, nil
		},
	}
	uc := New(repo)

	got, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{
		SKU: s("SKU-OTRO"),
	})

	assert.ErrorIs(t, err, domain.ErrProductAlreadyExists)
	assert.Nil(t, got)
}

func TestUpdateProduct_MismoSKU_NoVerificaDuplicado(t *testing.T) {
	verificado := false
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		ProductExistsFn: func(ctx context.Context, businessID uint, sku string) (bool, error) {
			verificado = true
			return true, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{
		SKU: s("SKU-1"),
	})

	require.NoError(t, err)
	assert.False(t, verificado, "guardar el mismo SKU no dispara la verificacion de duplicado")
}

func TestUpdateProduct_ActualizaCamposBasicos(t *testing.T) {
	var guardado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			guardado = product
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{
		Name:        s("Camiseta"),
		Description: s("Nueva descripcion"),
		Slug:        s("camiseta"),
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "Camiseta", guardado.Name)
	assert.Equal(t, "Nueva descripcion", guardado.Description)
	assert.Equal(t, "camiseta", guardado.Slug)
}

func TestUpdateProduct_CamposNulos_NoBorranLosExistentes(t *testing.T) {
	var guardado *domain.Product
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			guardado = product
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{})

	require.NoError(t, err)
	assert.Equal(t, "SKU-1", guardado.SKU)
	assert.Equal(t, "Camisa", guardado.Name)
	assert.InDelta(t, 25000.0, guardado.Price, 0.001)
}

func TestUpdateProduct_VarianteInvalida_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
	}
	uc := New(repo)

	got, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{
		VariantAttributes: datatypes.JSON([]byte(`{roto`)),
	})

	assert.ErrorIs(t, err, domain.ErrInvalidProductData)
	assert.Nil(t, got)
}

func TestUpdateProduct_VarianteDuplicada_SeExcluyeASiMismo(t *testing.T) {
	var excluido *string
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			p := productoBase(id)
			p.FamilyID = u(7)
			return p, nil
		},
		VariantExistsInFamilyFn: func(ctx context.Context, businessID, familyID uint, variantSignature string, excludeProductID *string) (bool, error) {
			excluido = excludeProductID
			return false, nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{
		VariantAttributes: datatypes.JSON([]byte(`{"talla":"M"}`)),
	})

	require.NoError(t, err)
	require.NotNil(t, excluido)
	assert.Equal(t, "prod-1", *excluido,
		"al actualizar se excluye el propio producto para no chocar consigo mismo")
}

func TestUpdateProduct_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("update failed")
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		UpdateProductFn: func(ctx context.Context, product *domain.Product) error {
			return dbErr
		},
	}
	uc := New(repo)

	got, err := uc.UpdateProduct(context.Background(), 10, "prod-1", &domain.UpdateProductRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestAddProductIntegration_DelegaEnRepositorio(t *testing.T) {
	var recibidoProducto string
	var recibidaIntegracion uint
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		AddProductIntegrationFn: func(ctx context.Context, productID string, integrationID uint, externalProductID string, externalVariantID, externalSKU, externalBarcode *string) (*domain.ProductBusinessIntegration, error) {
			recibidoProducto, recibidaIntegracion = productID, integrationID
			return &domain.ProductBusinessIntegration{}, nil
		},
	}
	uc := New(repo)

	got, err := uc.AddProductIntegration(context.Background(), 10, "prod-1", &domain.AddProductIntegrationRequest{
		IntegrationID:     197,
		ExternalProductID: "EXT-1",
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "prod-1", recibidoProducto)
	assert.Equal(t, uint(197), recibidaIntegracion)
}

func TestGetProductIntegrations_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		GetProductIntegrationsFn: func(ctx context.Context, productID string) ([]domain.ProductBusinessIntegration, error) {
			return []domain.ProductBusinessIntegration{{IntegrationID: 197}, {IntegrationID: 198}}, nil
		},
	}
	uc := New(repo)

	got, err := uc.GetProductIntegrations(context.Background(), 10, "prod-1")

	require.NoError(t, err)
	assert.Len(t, got, 2)
}

func TestRemoveProductIntegration_DelegaEnRepositorio(t *testing.T) {
	quitada := uint(0)
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		RemoveProductIntegrationFn: func(ctx context.Context, productID string, integrationID uint) error {
			quitada = integrationID
			return nil
		},
	}
	uc := New(repo)

	require.NoError(t, uc.RemoveProductIntegration(context.Background(), 10, "prod-1", 197))
	assert.Equal(t, uint(197), quitada)
}

func TestGetProductsByIntegration_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductsByIntegrationFn: func(ctx context.Context, integrationID uint) ([]domain.Product, error) {
			return []domain.Product{*productoBase("prod-1")}, nil
		},
	}
	uc := New(repo)

	got, err := uc.GetProductsByIntegration(context.Background(), 197)

	require.NoError(t, err)
	assert.Len(t, got, 1)
}

func TestDeleteProduct_DelegaEnRepositorio(t *testing.T) {
	var borrado string
	repo := &mocks.RepositoryMock{
		GetProductByIDFn: func(ctx context.Context, businessID uint, id string) (*domain.Product, error) {
			return productoBase(id), nil
		},
		DeleteProductFn: func(ctx context.Context, businessID uint, id string) error {
			borrado = id
			return nil
		},
	}
	uc := New(repo)

	err := uc.DeleteProduct(context.Background(), 10, "prod-1")

	require.NoError(t, err)
	assert.Equal(t, "prod-1", borrado)
}

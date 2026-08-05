package usecasefamily

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/products/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func s(v string) *string {
	return &v
}

func b(v bool) *bool {
	return &v
}

func familiaBase(id uint) *domain.ProductFamily {
	return &domain.ProductFamily{
		ID:          id,
		BusinessID:  10,
		Name:        "Camisas",
		Title:       "Camisas de algodon",
		Description: "Linea basica",
		Slug:        "camisas",
		Category:    "Ropa",
		Brand:       "Marca",
		Status:      "active",
		IsActive:    true,
	}
}

func TestCreateProductFamily_ActivaPorDefecto(t *testing.T) {
	var creada *domain.ProductFamily
	repo := &mocks.RepositoryMock{
		CreateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			family.ID = 1
			creada = family
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.CreateProductFamily(context.Background(), &domain.CreateProductFamilyStandaloneRequest{
		BusinessID:                 10,
		CreateProductFamilyRequest: domain.CreateProductFamilyRequest{Name: "Camisas"},
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	require.NotNil(t, creada)
	assert.True(t, creada.IsActive, "una familia nueva nace activa si no se indica lo contrario")
	assert.Equal(t, uint(10), creada.BusinessID)
	assert.Equal(t, "Camisas", creada.Name)
}

func TestCreateProductFamily_InactivaExplicita(t *testing.T) {
	var creada *domain.ProductFamily
	repo := &mocks.RepositoryMock{
		CreateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			creada = family
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.CreateProductFamily(context.Background(), &domain.CreateProductFamilyStandaloneRequest{
		BusinessID: 10,
		CreateProductFamilyRequest: domain.CreateProductFamilyRequest{
			Name:     "Camisas",
			IsActive: b(false),
		},
	})

	require.NoError(t, err)
	assert.False(t, creada.IsActive)
}

func TestCreateProductFamily_CopiaTodosLosCampos(t *testing.T) {
	var creada *domain.ProductFamily
	repo := &mocks.RepositoryMock{
		CreateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			creada = family
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.CreateProductFamily(context.Background(), &domain.CreateProductFamilyStandaloneRequest{
		BusinessID: 10,
		CreateProductFamilyRequest: domain.CreateProductFamilyRequest{
			Name:        "Camisas",
			Title:       "Camisas de algodon",
			Description: "Linea basica",
			Slug:        "camisas",
			Category:    "Ropa",
			Brand:       "Marca",
			ImageURL:    "https://s3/camisas.jpg",
			Status:      "active",
		},
	})

	require.NoError(t, err)
	assert.Equal(t, "Camisas de algodon", creada.Title)
	assert.Equal(t, "Linea basica", creada.Description)
	assert.Equal(t, "camisas", creada.Slug)
	assert.Equal(t, "Ropa", creada.Category)
	assert.Equal(t, "Marca", creada.Brand)
	assert.Equal(t, "https://s3/camisas.jpg", creada.ImageURL)
	assert.Equal(t, "active", creada.Status)
}

func TestCreateProductFamily_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		CreateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			return dbErr
		},
	}
	uc := New(repo)

	got, err := uc.CreateProductFamily(context.Background(), &domain.CreateProductFamilyStandaloneRequest{BusinessID: 10})

	require.Error(t, err)
	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetProductFamilyByID_IncluyeLasVariantes(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		ListProductsByFamilyIDFn: func(ctx context.Context, businessID, familyID uint) ([]domain.Product, error) {
			return []domain.Product{
				{ID: "prod-1", SKU: "SKU-1"},
				{ID: "prod-2", SKU: "SKU-2"},
			}, nil
		},
	}
	uc := New(repo)

	got, err := uc.GetProductFamilyByID(context.Background(), 10, 1)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.ID)
	assert.Len(t, got.Variants, 2)
	assert.Equal(t, "SKU-1", got.Variants[0].SKU)
}

func TestGetProductFamilyByID_FamiliaInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.GetProductFamilyByID(context.Background(), 10, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetProductFamilyByID_FallaListarVariantes_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		ListProductsByFamilyIDFn: func(ctx context.Context, businessID, familyID uint) ([]domain.Product, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.GetProductFamilyByID(context.Background(), 10, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListProductFamilies_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		name         string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{name: "pagina cero", page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "pagina negativa", page: -2, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{name: "page size cero", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{name: "page size sobre el maximo", page: 1, pageSize: 500, wantPage: 1, wantPageSize: 10},
		{name: "page size en el limite", page: 2, pageSize: 100, wantPage: 2, wantPageSize: 100},
		{name: "valores validos", page: 3, pageSize: 25, wantPage: 3, wantPageSize: 25},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var gotPage, gotPageSize int
			repo := &mocks.RepositoryMock{
				ListProductFamiliesFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error) {
					gotPage, gotPageSize = page, pageSize
					return nil, 0, nil
				},
			}
			uc := New(repo)

			got, err := uc.ListProductFamilies(context.Background(), 10, tc.page, tc.pageSize, nil)

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, gotPage)
			assert.Equal(t, tc.wantPageSize, gotPageSize)
			assert.Equal(t, tc.wantPage, got.Page)
			assert.Equal(t, tc.wantPageSize, got.PageSize)
		})
	}
}

func TestListProductFamilies_CalculaTotalDePaginas(t *testing.T) {
	cases := []struct {
		total     int64
		pageSize  int
		wantPages int
	}{
		{total: 0, pageSize: 10, wantPages: 0},
		{total: 1, pageSize: 10, wantPages: 1},
		{total: 10, pageSize: 10, wantPages: 1},
		{total: 11, pageSize: 10, wantPages: 2},
		{total: 95, pageSize: 10, wantPages: 10},
	}

	for _, tc := range cases {
		repo := &mocks.RepositoryMock{
			ListProductFamiliesFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error) {
				return nil, tc.total, nil
			},
		}
		uc := New(repo)

		got, err := uc.ListProductFamilies(context.Background(), 10, 1, tc.pageSize, nil)

		require.NoError(t, err)
		assert.Equal(t, tc.wantPages, got.TotalPages)
		assert.Equal(t, tc.total, got.Total)
	}
}

func TestListProductFamilies_MapeaCadaFamilia(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListProductFamiliesFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error) {
			return []domain.ProductFamily{*familiaBase(1), *familiaBase(2)}, 2, nil
		},
	}
	uc := New(repo)

	got, err := uc.ListProductFamilies(context.Background(), 10, 1, 10, nil)

	require.NoError(t, err)
	require.Len(t, got.Data, 2)
	assert.Equal(t, uint(1), got.Data[0].ID)
	assert.Equal(t, uint(2), got.Data[1].ID)
}

func TestListProductFamilies_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("query failed")
	repo := &mocks.RepositoryMock{
		ListProductFamiliesFn: func(ctx context.Context, businessID uint, page, pageSize int, filters map[string]interface{}) ([]domain.ProductFamily, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.ListProductFamilies(context.Background(), 10, 1, 10, nil)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateProductFamily_SinID_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	got, err := uc.UpdateProductFamily(context.Background(), 10, 0, &domain.UpdateProductFamilyRequest{})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "product family ID is required")
	assert.Nil(t, got)
}

func TestUpdateProductFamily_ActualizaTodosLosCampos(t *testing.T) {
	var guardada *domain.ProductFamily
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		UpdateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			guardada = family
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.UpdateProductFamily(context.Background(), 10, 1, &domain.UpdateProductFamilyRequest{
		Name:        s("Camisetas"),
		Title:       s("Camisetas nuevas"),
		Description: s("Otra linea"),
		Slug:        s("camisetas"),
		Category:    s("Deportivo"),
		Brand:       s("Otra marca"),
		ImageURL:    s("https://s3/otra.jpg"),
		Status:      s("draft"),
		IsActive:    b(false),
	})

	require.NoError(t, err)
	require.NotNil(t, guardada)
	assert.Equal(t, "Camisetas", guardada.Name)
	assert.Equal(t, "Camisetas nuevas", guardada.Title)
	assert.Equal(t, "Otra linea", guardada.Description)
	assert.Equal(t, "camisetas", guardada.Slug)
	assert.Equal(t, "Deportivo", guardada.Category)
	assert.Equal(t, "Otra marca", guardada.Brand)
	assert.Equal(t, "https://s3/otra.jpg", guardada.ImageURL)
	assert.Equal(t, "draft", guardada.Status)
	assert.False(t, guardada.IsActive)
}

func TestUpdateProductFamily_CamposNulos_NoBorranLosExistentes(t *testing.T) {
	var guardada *domain.ProductFamily
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		UpdateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			guardada = family
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.UpdateProductFamily(context.Background(), 10, 1, &domain.UpdateProductFamilyRequest{})

	require.NoError(t, err)
	assert.Equal(t, "Camisas", guardada.Name)
	assert.Equal(t, "Ropa", guardada.Category)
	assert.True(t, guardada.IsActive)
}

func TestUpdateProductFamily_FamiliaInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.UpdateProductFamily(context.Background(), 10, 1, &domain.UpdateProductFamilyRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateProductFamily_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("update failed")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		UpdateProductFamilyFn: func(ctx context.Context, family *domain.ProductFamily) error {
			return dbErr
		},
	}
	uc := New(repo)

	got, err := uc.UpdateProductFamily(context.Background(), 10, 1, &domain.UpdateProductFamilyRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDeleteProductFamily_SinID_RetornaError(t *testing.T) {
	uc := New(&mocks.RepositoryMock{})

	err := uc.DeleteProductFamily(context.Background(), 10, 0)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "product family ID is required")
}

func TestDeleteProductFamily_ConVariantesActivas_Rechaza(t *testing.T) {
	borrada := false
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		HasFamilyActiveVariantsFn: func(ctx context.Context, businessID, familyID uint) (bool, error) {
			return true, nil
		},
		DeleteProductFamilyFn: func(ctx context.Context, businessID, familyID uint) error {
			borrada = true
			return nil
		},
	}
	uc := New(repo)

	err := uc.DeleteProductFamily(context.Background(), 10, 1)

	assert.ErrorIs(t, err, domain.ErrFamilyHasActiveVariants)
	assert.False(t, borrada, "no se borra una familia que todavia tiene variantes activas")
}

func TestDeleteProductFamily_SinVariantes_SeBorra(t *testing.T) {
	var borradaID uint
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		HasFamilyActiveVariantsFn: func(ctx context.Context, businessID, familyID uint) (bool, error) {
			return false, nil
		},
		DeleteProductFamilyFn: func(ctx context.Context, businessID, familyID uint) error {
			borradaID = familyID
			return nil
		},
	}
	uc := New(repo)

	require.NoError(t, uc.DeleteProductFamily(context.Background(), 10, 1))
	assert.Equal(t, uint(1), borradaID)
}

func TestDeleteProductFamily_FamiliaInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	assert.ErrorIs(t, uc.DeleteProductFamily(context.Background(), 10, 1), dbErr)
}

func TestDeleteProductFamily_FallaVerificarVariantes_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		HasFamilyActiveVariantsFn: func(ctx context.Context, businessID, familyID uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := New(repo)

	assert.ErrorIs(t, uc.DeleteProductFamily(context.Background(), 10, 1), dbErr)
}

func TestDeleteProductFamily_FallaBorrar_PropagaError(t *testing.T) {
	dbErr := errors.New("delete failed")
	repo := &mocks.RepositoryMock{
		GetProductFamilyByIDFn: func(ctx context.Context, businessID, familyID uint) (*domain.ProductFamily, error) {
			return familiaBase(familyID), nil
		},
		HasFamilyActiveVariantsFn: func(ctx context.Context, businessID, familyID uint) (bool, error) {
			return false, nil
		},
		DeleteProductFamilyFn: func(ctx context.Context, businessID, familyID uint) error {
			return dbErr
		},
	}
	uc := New(repo)

	assert.ErrorIs(t, uc.DeleteProductFamily(context.Background(), 10, 1), dbErr)
}

func TestMapProductFamilyToResponse_Nil(t *testing.T) {
	assert.Nil(t, mapProductFamilyToResponse(nil))
}

func TestMapProductFamilyToResponse_CopiaLosCampos(t *testing.T) {
	got := mapProductFamilyToResponse(familiaBase(7))

	require.NotNil(t, got)
	assert.Equal(t, uint(7), got.ID)
	assert.Equal(t, "Camisas", got.Name)
	assert.Equal(t, "Ropa", got.Category)
	assert.True(t, got.IsActive)
}

func TestMapProductsToResponses_ListaVacia(t *testing.T) {
	assert.Empty(t, mapProductsToResponses(nil))
	assert.Empty(t, mapProductsToResponses([]domain.Product{}))
}

func TestMapProductsToResponses_MapeaCadaProducto(t *testing.T) {
	got := mapProductsToResponses([]domain.Product{
		{ID: "prod-1", SKU: "SKU-1", Name: "Camisa M"},
		{ID: "prod-2", SKU: "SKU-2", Name: "Camisa L"},
	})

	require.Len(t, got, 2)
	assert.Equal(t, "SKU-1", got[0].SKU)
	assert.Equal(t, "Camisa L", got[1].Name)
}

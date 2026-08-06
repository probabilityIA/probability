package app

import (
	"context"
	stderrors "errors"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/auth/resources/internal/domain"
	"github.com/secamc93/probability/back/central/services/auth/resources/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newResourcesUseCase(repo *mocks.RepositoryMock) IUseCaseResource {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo, mocks.NewSilentLogger())
}

func uintPtr(v uint) *uint { return &v }

func TestCreateResource_NombreObligatorio(t *testing.T) {
	for _, nombre := range []string{"", "   ", "\t\n"} {
		repo := &mocks.RepositoryMock{}

		got, err := newResourcesUseCase(repo).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: nombre})

		require.Error(t, err, "nombre=%q", nombre)
		assert.Contains(t, err.Error(), "nombre del recurso es obligatorio")
		assert.Nil(t, got)
		assert.Empty(t, repo.CreatedResources)
	}
}

func TestCreateResource_LimitesDeLongitud(t *testing.T) {
	t.Run("nombre en el limite se acepta", func(t *testing.T) {
		_, err := newResourcesUseCase(nil).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: strings.Repeat("a", 100)})
		require.NoError(t, err)
	})

	t.Run("nombre sobre el limite se rechaza", func(t *testing.T) {
		repo := &mocks.RepositoryMock{}
		_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: strings.Repeat("a", 101)})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no puede exceder 100")
		assert.Empty(t, repo.CreatedResources)
	})

	t.Run("descripcion en el limite se acepta", func(t *testing.T) {
		_, err := newResourcesUseCase(nil).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: "orders", Description: strings.Repeat("d", 500)})
		require.NoError(t, err)
	})

	t.Run("descripcion sobre el limite se rechaza", func(t *testing.T) {
		repo := &mocks.RepositoryMock{}
		_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: "orders", Description: strings.Repeat("d", 501)})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no puede exceder 500")
		assert.Empty(t, repo.CreatedResources)
	})
}

func TestCreateResource_NombreDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByNameFn: func(ctx context.Context, name string) (*domain.Resource, error) {
			return &domain.Resource{ID: 1, Name: name}, nil
		},
	}

	got, err := newResourcesUseCase(repo).CreateResource(context.Background(),
		domain.CreateResourceDTO{Name: "orders"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ya existe un recurso con el nombre 'orders'")
	assert.Nil(t, got)
	assert.Empty(t, repo.CreatedResources,
		"dos recursos con el mismo nombre harian ambiguo el permiso")
}

func TestCreateResource_ErrorAlBuscarDuplicado_NoBloquea(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByNameFn: func(ctx context.Context, name string) (*domain.Resource, error) {
			return nil, stderrors.New("no encontrado")
		},
	}

	got, err := newResourcesUseCase(repo).CreateResource(context.Background(),
		domain.CreateResourceDTO{Name: "orders"})

	require.NoError(t, err, "que la busqueda falle se interpreta como que no existe")
	assert.NotNil(t, got)
	assert.Len(t, repo.CreatedResources, 1)
}

func TestCreateResource_RecortaEspacios(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
		domain.CreateResourceDTO{Name: "  orders  ", Description: "  gestion de ordenes  "})

	require.NoError(t, err)
	require.Len(t, repo.CreatedResources, 1)
	assert.Equal(t, "orders", repo.CreatedResources[0].Name)
	assert.Equal(t, "gestion de ordenes", repo.CreatedResources[0].Description)
}

func TestCreateResource_SinBusinessTypeQuedaGenerico(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
		domain.CreateResourceDTO{Name: "orders", BusinessTypeID: nil})

	require.NoError(t, err)
	assert.Zero(t, repo.CreatedResources[0].BusinessTypeID,
		"sin tipo de negocio el recurso es generico y aplica a todos")
}

func TestCreateResource_ConBusinessTypeLoAsigna(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
		domain.CreateResourceDTO{Name: "orders", BusinessTypeID: uintPtr(4)})

	require.NoError(t, err)
	assert.Equal(t, uint(4), repo.CreatedResources[0].BusinessTypeID)
}

func TestCreateResource_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("constraint violation")

	t.Run("al crear", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			CreateResourceFn: func(ctx context.Context, r domain.Resource) (uint, error) { return 0, dbErr },
		}
		_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: "orders"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al recargar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
				return nil, dbErr
			},
		}
		_, err := newResourcesUseCase(repo).CreateResource(context.Background(),
			domain.CreateResourceDTO{Name: "orders"})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestCreateResource_DevuelveElRecursoRecargado(t *testing.T) {
	repo := &mocks.RepositoryMock{
		CreateResourceFn: func(ctx context.Context, r domain.Resource) (uint, error) { return 9, nil },
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{
				ID: id, Name: "orders", Description: "d",
				BusinessTypeID: 4, BusinessTypeName: "Retail",
			}, nil
		},
	}

	got, err := newResourcesUseCase(repo).CreateResource(context.Background(),
		domain.CreateResourceDTO{Name: "orders"})

	require.NoError(t, err)
	assert.Equal(t, uint(9), got.ID)
	assert.Equal(t, "Retail", got.BusinessTypeName,
		"la respuesta trae el nombre del tipo, que solo sale de la relectura")
}

func TestGetResources_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 10, 1, 10},
		{-3, 10, 1, 10},
		{1, 0, 1, 10},
		{1, -1, 1, 10},
		{2, 50, 2, 50},
	}

	for _, tc := range casos {
		var visto domain.ResourceFilters
		repo := &mocks.RepositoryMock{
			GetResourcesFn: func(ctx context.Context, f domain.ResourceFilters) ([]domain.Resource, int64, error) {
				visto = f
				return nil, 0, nil
			},
		}

		got, err := newResourcesUseCase(repo).GetResources(context.Background(),
			domain.ResourceFilters{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
		assert.Equal(t, tc.wantPage, got.Page)
	}
}

func TestGetResources_CalculaTotalPagesConRedondeoHaciaArriba(t *testing.T) {
	casos := []struct {
		total    int64
		pageSize int
		want     int
	}{
		{0, 10, 0},
		{1, 10, 1},
		{10, 10, 1},
		{11, 10, 2},
		{95, 10, 10},
	}

	for _, tc := range casos {
		repo := &mocks.RepositoryMock{
			GetResourcesFn: func(ctx context.Context, f domain.ResourceFilters) ([]domain.Resource, int64, error) {
				return nil, tc.total, nil
			},
		}

		got, err := newResourcesUseCase(repo).GetResources(context.Background(),
			domain.ResourceFilters{PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.want, got.TotalPages, "total=%d pageSize=%d", tc.total, tc.pageSize)
	}
}

func TestGetResources_PasaLosFiltrosTalCual(t *testing.T) {
	var visto domain.ResourceFilters
	repo := &mocks.RepositoryMock{
		GetResourcesFn: func(ctx context.Context, f domain.ResourceFilters) ([]domain.Resource, int64, error) {
			visto = f
			return nil, 0, nil
		},
	}

	_, err := newResourcesUseCase(repo).GetResources(context.Background(), domain.ResourceFilters{
		Page: 1, PageSize: 10, Name: "ord", Description: "gestion",
		BusinessTypeID: uintPtr(4), SortBy: "name", SortOrder: "asc",
	})

	require.NoError(t, err)
	assert.Equal(t, "ord", visto.Name)
	assert.Equal(t, "gestion", visto.Description)
	require.NotNil(t, visto.BusinessTypeID)
	assert.Equal(t, uint(4), *visto.BusinessTypeID)
	assert.Equal(t, "name", visto.SortBy)
	assert.Equal(t, "asc", visto.SortOrder)
}

func TestGetResources_MapeaTodosLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourcesFn: func(ctx context.Context, f domain.ResourceFilters) ([]domain.Resource, int64, error) {
			return []domain.Resource{
				{ID: 1, Name: "orders", Description: "d", BusinessTypeID: 4, BusinessTypeName: "Retail"},
			}, 1, nil
		},
	}

	got, err := newResourcesUseCase(repo).GetResources(context.Background(), domain.ResourceFilters{})

	require.NoError(t, err)
	require.Len(t, got.Resources, 1)
	assert.Equal(t, "orders", got.Resources[0].Name)
	assert.Equal(t, uint(4), got.Resources[0].BusinessTypeID)
	assert.Equal(t, "Retail", got.Resources[0].BusinessTypeName)
	assert.Equal(t, int64(1), got.Total)
}

func TestGetResources_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetResourcesFn: func(ctx context.Context, f domain.ResourceFilters) ([]domain.Resource, int64, error) {
			return nil, 0, dbErr
		},
	}

	got, err := newResourcesUseCase(repo).GetResources(context.Background(), domain.ResourceFilters{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetResourceByID_Inexistente_RetornaError(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) { return nil, nil },
	}

	got, err := newResourcesUseCase(repo).GetResourceByID(context.Background(), 404)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "404 no encontrado")
}

func TestGetResourceByID_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) { return nil, dbErr },
	}

	_, err := newResourcesUseCase(repo).GetResourceByID(context.Background(), 1)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetResourceByID_MapeaTodosLosCampos(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{
				ID: id, Name: "orders", Description: "d", BusinessTypeID: 4, BusinessTypeName: "Retail",
			}, nil
		},
	}

	got, err := newResourcesUseCase(repo).GetResourceByID(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), got.ID)
	assert.Equal(t, "orders", got.Name)
	assert.Equal(t, "Retail", got.BusinessTypeName)
}

func TestUpdateResource_ValidaAntesDeLeer(t *testing.T) {
	leido := false
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			leido = true
			return &domain.Resource{ID: id}, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "  "})

	require.Error(t, err)
	assert.False(t, leido)
}

func TestUpdateResource_Inexistente_NoActualiza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) { return nil, nil },
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 404,
		domain.UpdateResourceDTO{Name: "orders"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "404 no encontrado")
	assert.Empty(t, repo.UpdatedResources)
}

func TestUpdateResource_NombreDeOtroRecurso_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "viejo"}, nil
		},
		GetResourceByNameFn: func(ctx context.Context, name string) (*domain.Resource, error) {
			return &domain.Resource{ID: 99, Name: name}, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "ocupado"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "ya existe otro recurso con el nombre 'ocupado'")
	assert.Empty(t, repo.UpdatedResources)
}

func TestUpdateResource_SuPropioNombreNoCuentaComoDuplicado(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "viejo"}, nil
		},
		GetResourceByNameFn: func(ctx context.Context, name string) (*domain.Resource, error) {
			return &domain.Resource{ID: 1, Name: name}, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "nuevo"})

	require.NoError(t, err, "si el duplicado es el mismo recurso no hay conflicto")
	assert.Len(t, repo.UpdatedResources, 1)
}

func TestUpdateResource_MismoNombreNoDisparaChequeo(t *testing.T) {
	chequeado := false
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "orders"}, nil
		},
		GetResourceByNameFn: func(ctx context.Context, name string) (*domain.Resource, error) {
			chequeado = true
			return nil, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "  orders  "})

	require.NoError(t, err)
	assert.False(t, chequeado, "el nombre recortado coincide con el actual, no hay que buscar duplicados")
}

func TestUpdateResource_SinBusinessTypeConservaElActual(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "orders", BusinessTypeID: 4}, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "orders", BusinessTypeID: nil})

	require.NoError(t, err)
	require.Len(t, repo.UpdatedResources, 1)
	assert.Equal(t, uint(4), repo.UpdatedResources[0].BusinessTypeID,
		"omitir el tipo no debe convertir el recurso en generico")
}

func TestUpdateResource_ConBusinessTypeLoCambia(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "orders", BusinessTypeID: 4}, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "orders", BusinessTypeID: uintPtr(9)})

	require.NoError(t, err)
	assert.Equal(t, uint(9), repo.UpdatedResources[0].BusinessTypeID)
}

func TestUpdateResource_RecortaEspacios(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "viejo"}, nil
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "  nuevo  ", Description: "  desc  "})

	require.NoError(t, err)
	assert.Equal(t, "nuevo", repo.UpdatedResources[0].Name)
	assert.Equal(t, "desc", repo.UpdatedResources[0].Description)
}

func TestUpdateResource_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("deadlock")

	t.Run("al actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
				return &domain.Resource{ID: id, Name: "orders"}, nil
			},
			UpdateResourceFn: func(ctx context.Context, id uint, r domain.Resource) (string, error) {
				return "", dbErr
			},
		}
		_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
			domain.UpdateResourceDTO{Name: "orders"})
		assert.ErrorIs(t, err, dbErr)
	})

	t.Run("al leer para actualizar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
				return nil, dbErr
			},
		}
		_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
			domain.UpdateResourceDTO{Name: "orders"})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no encontrado")
	})
}

func TestUpdateResource_ErrorAlRecargar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	llamadas := 0
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			llamadas++
			if llamadas == 1 {
				return &domain.Resource{ID: id, Name: "orders"}, nil
			}
			return nil, dbErr
		},
	}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "orders"})

	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateResource_LimitesDeLongitud(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: strings.Repeat("a", 101)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no puede exceder 100")

	_, err = newResourcesUseCase(repo).UpdateResource(context.Background(), 1,
		domain.UpdateResourceDTO{Name: "orders", Description: strings.Repeat("d", 501)})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "no puede exceder 500")

	assert.Empty(t, repo.UpdatedResources)
}

func TestDeleteResource_Inexistente_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) { return nil, nil },
		DeleteResourceFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	got, err := newResourcesUseCase(repo).DeleteResource(context.Background(), 404)

	require.Error(t, err)
	assert.Empty(t, got)
	assert.Contains(t, err.Error(), "404 no encontrado")
	assert.False(t, borrado)
}

func TestDeleteResource_ErrorAlVerificar_NoBorra(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return nil, stderrors.New("timeout")
		},
		DeleteResourceFn: func(ctx context.Context, id uint) (string, error) {
			borrado = true
			return "", nil
		},
	}

	_, err := newResourcesUseCase(repo).DeleteResource(context.Background(), 1)

	require.Error(t, err)
	assert.False(t, borrado)
}

func TestDeleteResource_Exitoso_DevuelveElMensajeDelRepo(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "orders"}, nil
		},
		DeleteResourceFn: func(ctx context.Context, id uint) (string, error) {
			borradoID = id
			return "recurso eliminado", nil
		},
	}

	got, err := newResourcesUseCase(repo).DeleteResource(context.Background(), 5)

	require.NoError(t, err)
	assert.Equal(t, uint(5), borradoID)
	assert.Equal(t, "recurso eliminado", got)
}

func TestDeleteResource_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation: hay permisos usando el recurso")
	repo := &mocks.RepositoryMock{
		GetResourceByIDFn: func(ctx context.Context, id uint) (*domain.Resource, error) {
			return &domain.Resource{ID: id, Name: "orders"}, nil
		},
		DeleteResourceFn: func(ctx context.Context, id uint) (string, error) { return "", dbErr },
	}

	got, err := newResourcesUseCase(repo).DeleteResource(context.Background(), 5)

	assert.ErrorIs(t, err, dbErr)
	assert.Empty(t, got)
}

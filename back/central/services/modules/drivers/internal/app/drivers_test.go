package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/drivers/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/drivers/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDriversUseCase(repo *mocks.RepositoryMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo)
}

func uintPtr(v uint) *uint    { return &v }
func strPtr(v string) *string { return &v }

func TestCreateDriver_SinEstado_NaceActivo(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newDriversUseCase(repo).CreateDriver(context.Background(),
		dtos.CreateDriverDTO{BusinessID: 26, FirstName: "Juan"})

	require.NoError(t, err)
	require.NotNil(t, repo.Created)
	assert.Equal(t, "active", repo.Created.Status)
}

func TestCreateDriver_ConEstado_LoRespeta(t *testing.T) {
	for _, estado := range []string{"active", "inactive", "on_route", "suspended"} {
		repo := &mocks.RepositoryMock{}

		_, err := newDriversUseCase(repo).CreateDriver(context.Background(),
			dtos.CreateDriverDTO{BusinessID: 26, Status: estado})

		require.NoError(t, err)
		assert.Equal(t, estado, repo.Created.Status)
	}
}

func TestCreateDriver_SinIdentificacion_NoBuscaDuplicados(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newDriversUseCase(repo).CreateDriver(context.Background(),
		dtos.CreateDriverDTO{BusinessID: 26, FirstName: "Juan", Identification: ""})

	require.NoError(t, err)
	assert.Empty(t, repo.ExistsCalls,
		"sin cedula no hay nada por lo que verificar unicidad")
	assert.NotNil(t, repo.Created)
}

func TestCreateDriver_IdentificacionDuplicada_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ExistsByIdentificationFn: func(ctx context.Context, businessID uint, ident string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	got, err := newDriversUseCase(repo).CreateDriver(context.Background(),
		dtos.CreateDriverDTO{BusinessID: 26, Identification: "1020304050"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateIdentification)
	assert.Nil(t, got)
	assert.Nil(t, repo.Created)
}

func TestCreateDriver_LaUnicidadEsPorNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newDriversUseCase(repo).CreateDriver(context.Background(),
		dtos.CreateDriverDTO{BusinessID: 99, Identification: "1020304050"})

	require.NoError(t, err)
	require.Len(t, repo.ExistsCalls, 1)
	assert.Equal(t, uint(99), repo.ExistsCalls[0].BusinessID,
		"la misma cedula puede existir en dos negocios distintos")
	assert.Equal(t, "1020304050", repo.ExistsCalls[0].Identification)
	assert.Nil(t, repo.ExistsCalls[0].ExcludeID, "al crear no hay ID propio que excluir")
}

func TestCreateDriver_ErrorAlVerificarDuplicado_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ExistsByIdentificationFn: func(ctx context.Context, businessID uint, ident string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}

	_, err := newDriversUseCase(repo).CreateDriver(context.Background(),
		dtos.CreateDriverDTO{BusinessID: 26, Identification: "1020304050"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, repo.Created)
}

func TestCreateDriver_CopiaTodosLosCampos(t *testing.T) {
	vence := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{}

	_, err := newDriversUseCase(repo).CreateDriver(context.Background(), dtos.CreateDriverDTO{
		BusinessID: 26, FirstName: "Juan", LastName: "Perez",
		Email: "juan@x.com", Phone: "3001234567", Identification: "1020304050",
		PhotoURL: "s3://foto.png", LicenseType: "C2", LicenseExpiry: &vence,
		WarehouseID: uintPtr(3), Notes: strPtr("turno noche"),
	})

	require.NoError(t, err)
	d := repo.Created
	assert.Equal(t, uint(26), d.BusinessID)
	assert.Equal(t, "Juan", d.FirstName)
	assert.Equal(t, "Perez", d.LastName)
	assert.Equal(t, "juan@x.com", d.Email)
	assert.Equal(t, "3001234567", d.Phone)
	assert.Equal(t, "C2", d.LicenseType)
	require.NotNil(t, d.LicenseExpiry)
	assert.Equal(t, vence, *d.LicenseExpiry)
	require.NotNil(t, d.WarehouseID)
	assert.Equal(t, uint(3), *d.WarehouseID)
	assert.Equal(t, "turno noche", *d.Notes)
}

func TestCreateDriver_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, d *entities.Driver) (*entities.Driver, error) {
			return nil, dbErr
		},
	}

	got, err := newDriversUseCase(repo).CreateDriver(context.Background(),
		dtos.CreateDriverDTO{BusinessID: 26})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateDriver_ConservaNegocioYFechaDeCreacion(t *testing.T) {
	creado := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return &entities.Driver{
				ID: driverID, BusinessID: businessID, CreatedAt: creado, FirstName: "Viejo",
			}, nil
		},
	}

	_, err := newDriversUseCase(repo).UpdateDriver(context.Background(), dtos.UpdateDriverDTO{
		ID: 5, BusinessID: 26, FirstName: "Nuevo",
	})

	require.NoError(t, err)
	require.NotNil(t, repo.Updated)
	assert.Equal(t, uint(5), repo.Updated.ID)
	assert.Equal(t, uint(26), repo.Updated.BusinessID, "el negocio dueno no cambia por update")
	assert.Equal(t, creado, repo.Updated.CreatedAt)
	assert.Equal(t, "Nuevo", repo.Updated.FirstName)
}

func TestUpdateDriver_LeePorNegocioYPorID(t *testing.T) {
	var vistoNegocio, vistoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			vistoNegocio, vistoID = businessID, driverID
			return &entities.Driver{ID: driverID, BusinessID: businessID}, nil
		},
	}

	_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
		dtos.UpdateDriverDTO{ID: 5, BusinessID: 26})

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"un negocio no puede editar el conductor de otro pasando su ID")
	assert.Equal(t, uint(5), vistoID)
}

func TestUpdateDriver_MismaIdentificacion_NoDisparaChequeo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return &entities.Driver{ID: driverID, Identification: "1020304050"}, nil
		},
	}

	_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
		dtos.UpdateDriverDTO{ID: 5, BusinessID: 26, Identification: "1020304050"})

	require.NoError(t, err)
	assert.Empty(t, repo.ExistsCalls)
}

func TestUpdateDriver_IdentificacionVacia_NoDisparaChequeo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return &entities.Driver{ID: driverID, Identification: "1020304050"}, nil
		},
	}

	_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
		dtos.UpdateDriverDTO{ID: 5, BusinessID: 26, Identification: ""})

	require.NoError(t, err)
	assert.Empty(t, repo.ExistsCalls)
	assert.Empty(t, repo.Updated.Identification, "se guarda vacia, es lo que pidieron")
}

func TestUpdateDriver_ExcluyeElPropioIDAlBuscarDuplicados(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return &entities.Driver{ID: driverID, Identification: "vieja"}, nil
		},
	}

	_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
		dtos.UpdateDriverDTO{ID: 5, BusinessID: 26, Identification: "nueva"})

	require.NoError(t, err)
	require.Len(t, repo.ExistsCalls, 1)
	require.NotNil(t, repo.ExistsCalls[0].ExcludeID)
	assert.Equal(t, uint(5), *repo.ExistsCalls[0].ExcludeID)
}

func TestUpdateDriver_IdentificacionDeOtro_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return &entities.Driver{ID: driverID, Identification: "vieja"}, nil
		},
		ExistsByIdentificationFn: func(ctx context.Context, businessID uint, ident string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
		dtos.UpdateDriverDTO{ID: 5, BusinessID: 26, Identification: "ocupada"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateIdentification)
	assert.Nil(t, repo.Updated)
}

func TestUpdateDriver_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("conductor inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
				return nil, domainerrors.ErrDriverNotFound
			},
		}
		_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
			dtos.UpdateDriverDTO{ID: 404, BusinessID: 26})
		assert.ErrorIs(t, err, domainerrors.ErrDriverNotFound)
		assert.Nil(t, repo.Updated)
	})

	t.Run("al verificar duplicado", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
				return &entities.Driver{ID: driverID, Identification: "vieja"}, nil
			},
			ExistsByIdentificationFn: func(ctx context.Context, businessID uint, ident string, excludeID *uint) (bool, error) {
				return false, dbErr
			},
		}
		_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
			dtos.UpdateDriverDTO{ID: 5, BusinessID: 26, Identification: "nueva"})
		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, repo.Updated)
	})

	t.Run("al guardar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateFn: func(ctx context.Context, d *entities.Driver) (*entities.Driver, error) {
				return nil, dbErr
			},
		}
		_, err := newDriversUseCase(repo).UpdateDriver(context.Background(),
			dtos.UpdateDriverDTO{ID: 5, BusinessID: 26})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestDeleteDriver_VerificaAntesDeBorrar(t *testing.T) {
	var vistoNegocio, vistoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			vistoNegocio, vistoID = businessID, driverID
			return &entities.Driver{ID: driverID}, nil
		},
	}

	err := newDriversUseCase(repo).DeleteDriver(context.Background(), 26, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"no se puede borrar el conductor de otro negocio adivinando el ID")
	assert.Equal(t, uint(5), vistoID)
	assert.Equal(t, []uint{5}, repo.Deleted)
}

func TestDeleteDriver_Inexistente_NoBorra(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return nil, domainerrors.ErrDriverNotFound
		},
	}

	err := newDriversUseCase(repo).DeleteDriver(context.Background(), 26, 404)

	assert.ErrorIs(t, err, domainerrors.ErrDriverNotFound)
	assert.Empty(t, repo.Deleted)
}

func TestDeleteDriver_ErrorAlBorrar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation: tiene rutas asignadas")
	repo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, businessID, driverID uint) error { return dbErr },
	}

	err := newDriversUseCase(repo).DeleteDriver(context.Background(), 26, 5)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetDriver_FiltraPorNegocio(t *testing.T) {
	var vistoNegocio uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			vistoNegocio = businessID
			return &entities.Driver{ID: driverID, FirstName: "Juan"}, nil
		},
	}

	got, err := newDriversUseCase(repo).GetDriver(context.Background(), 26, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, "Juan", got.FirstName)
}

func TestGetDriver_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, driverID uint) (*entities.Driver, error) {
			return nil, domainerrors.ErrDriverNotFound
		},
	}

	_, err := newDriversUseCase(repo).GetDriver(context.Background(), 26, 404)

	assert.ErrorIs(t, err, domainerrors.ErrDriverNotFound)
}

func TestListDrivers_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		page, pageSize, wantPage, wantSize int
	}{
		{0, 20, 1, 20},
		{-5, 20, 1, 20},
		{1, 0, 1, 20},
		{1, -1, 1, 20},
		{1, 101, 1, 20},
		{1, 100, 1, 100},
		{3, 50, 3, 50},
	}

	for _, tc := range casos {
		var visto dtos.ListDriversParams
		repo := &mocks.RepositoryMock{
			ListFn: func(ctx context.Context, p dtos.ListDriversParams) ([]entities.Driver, int64, error) {
				visto = p
				return nil, 0, nil
			},
		}

		_, _, err := newDriversUseCase(repo).ListDrivers(context.Background(),
			dtos.ListDriversParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestListDriversParams_Offset(t *testing.T) {
	assert.Equal(t, 0, dtos.ListDriversParams{Page: 1, PageSize: 20}.Offset())
	assert.Equal(t, 20, dtos.ListDriversParams{Page: 2, PageSize: 20}.Offset())
	assert.Equal(t, 100, dtos.ListDriversParams{Page: 3, PageSize: 50}.Offset())
}

func TestListDrivers_PasaLosFiltros(t *testing.T) {
	var visto dtos.ListDriversParams
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, p dtos.ListDriversParams) ([]entities.Driver, int64, error) {
			visto = p
			return []entities.Driver{{ID: 1}}, 1, nil
		},
	}

	got, total, err := newDriversUseCase(repo).ListDrivers(context.Background(),
		dtos.ListDriversParams{BusinessID: 26, Search: "juan", Status: "active", Page: 1, PageSize: 20})

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto.BusinessID)
	assert.Equal(t, "juan", visto.Search)
	assert.Equal(t, "active", visto.Status)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestListDrivers_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, p dtos.ListDriversParams) ([]entities.Driver, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newDriversUseCase(repo).ListDrivers(context.Background(), dtos.ListDriversParams{})

	assert.ErrorIs(t, err, dbErr)
}

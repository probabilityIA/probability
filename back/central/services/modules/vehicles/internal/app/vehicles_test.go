package app

import (
	"context"
	stderrors "errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/vehicles/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/vehicles/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newVehiclesUseCase(repo *mocks.RepositoryMock) IUseCase {
	if repo == nil {
		repo = &mocks.RepositoryMock{}
	}
	return New(repo)
}

func TestCreateVehicle_SinEstado_NaceActivo(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
		dtos.CreateVehicleDTO{BusinessID: 26, Brand: "Chevrolet"})

	require.NoError(t, err)
	require.NotNil(t, repo.Created)
	assert.Equal(t, "active", repo.Created.Status)
}

func TestCreateVehicle_ConEstado_LoRespeta(t *testing.T) {
	for _, estado := range []string{"active", "inactive", "maintenance", "out_of_service"} {
		repo := &mocks.RepositoryMock{}

		_, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
			dtos.CreateVehicleDTO{BusinessID: 26, Status: estado})

		require.NoError(t, err)
		assert.Equal(t, estado, repo.Created.Status)
	}
}

func TestCreateVehicle_SinPlaca_NoBuscaDuplicados(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
		dtos.CreateVehicleDTO{BusinessID: 26, Brand: "Chevrolet", LicensePlate: ""})

	require.NoError(t, err)
	assert.Empty(t, repo.ExistsCalls,
		"sin placa no hay nada por lo que verificar unicidad")
	assert.NotNil(t, repo.Created)
}

func TestCreateVehicle_PlacaDuplicada_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ExistsByLicensePlateFn: func(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	got, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
		dtos.CreateVehicleDTO{BusinessID: 26, LicensePlate: "1020304050"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateLicensePlate)
	assert.Nil(t, got)
	assert.Nil(t, repo.Created)
}

func TestCreateVehicle_LaUnicidadEsPorNegocio(t *testing.T) {
	repo := &mocks.RepositoryMock{}

	_, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
		dtos.CreateVehicleDTO{BusinessID: 99, LicensePlate: "1020304050"})

	require.NoError(t, err)
	require.Len(t, repo.ExistsCalls, 1)
	assert.Equal(t, uint(99), repo.ExistsCalls[0].BusinessID,
		"la misma placa puede existir en dos negocios distintos")
	assert.Equal(t, "1020304050", repo.ExistsCalls[0].Plate)
	assert.Nil(t, repo.ExistsCalls[0].ExcludeID, "al crear no hay ID propio que excluir")
}

func TestCreateVehicle_ErrorAlVerificarDuplicado_NoCrea(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ExistsByLicensePlateFn: func(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}

	_, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
		dtos.CreateVehicleDTO{BusinessID: 26, LicensePlate: "1020304050"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, repo.Created)
}

func TestCreateVehicle_CopiaTodosLosCampos(t *testing.T) {
	vence := time.Date(2027, 6, 30, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{}

	_, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(), dtos.CreateVehicleDTO{
		BusinessID: 26, Type: "van", LicensePlate: "XYZ789",
		Brand: "Chevrolet", VehicleModel: "N300", Color: "blanco",
		PhotoURL: "s3://foto.png", InsuranceExpiry: &vence, RegistrationExpiry: &vence,
	})

	require.NoError(t, err)
	v := repo.Created
	assert.Equal(t, uint(26), v.BusinessID)
	assert.Equal(t, "van", v.Type)
	assert.Equal(t, "XYZ789", v.LicensePlate)
	assert.Equal(t, "Chevrolet", v.Brand)
	assert.Equal(t, "N300", v.VehicleModel)
	assert.Equal(t, "blanco", v.Color)
	require.NotNil(t, v.InsuranceExpiry)
	assert.Equal(t, vence, *v.InsuranceExpiry)
	require.NotNil(t, v.RegistrationExpiry)
	assert.Equal(t, vence, *v.RegistrationExpiry)
}

func TestCreateVehicle_ErrorAlPersistir_SePropaga(t *testing.T) {
	dbErr := stderrors.New("constraint violation")
	repo := &mocks.RepositoryMock{
		CreateFn: func(ctx context.Context, d *entities.Vehicle) (*entities.Vehicle, error) {
			return nil, dbErr
		},
	}

	got, err := newVehiclesUseCase(repo).CreateVehicle(context.Background(),
		dtos.CreateVehicleDTO{BusinessID: 26})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateVehicle_ConservaNegocioYFechaDeCreacion(t *testing.T) {
	creado := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return &entities.Vehicle{
				ID: vehicleID, BusinessID: businessID, CreatedAt: creado, Brand: "Viejo",
			}, nil
		},
	}

	_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(), dtos.UpdateVehicleDTO{
		ID: 5, BusinessID: 26, Brand: "Nuevo",
	})

	require.NoError(t, err)
	require.NotNil(t, repo.Updated)
	assert.Equal(t, uint(5), repo.Updated.ID)
	assert.Equal(t, uint(26), repo.Updated.BusinessID, "el negocio dueno no cambia por update")
	assert.Equal(t, creado, repo.Updated.CreatedAt)
	assert.Equal(t, "Nuevo", repo.Updated.Brand)
}

func TestUpdateVehicle_LeePorNegocioYPorID(t *testing.T) {
	var vistoNegocio, vistoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			vistoNegocio, vistoID = businessID, vehicleID
			return &entities.Vehicle{ID: vehicleID, BusinessID: businessID}, nil
		},
	}

	_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
		dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26})

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"un negocio no puede editar el vehiculo de otro pasando su ID")
	assert.Equal(t, uint(5), vistoID)
}

func TestUpdateVehicle_MismaPlaca_NoDisparaChequeo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return &entities.Vehicle{ID: vehicleID, LicensePlate: "1020304050"}, nil
		},
	}

	_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
		dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26, LicensePlate: "1020304050"})

	require.NoError(t, err)
	assert.Empty(t, repo.ExistsCalls)
}

func TestUpdateVehicle_PlacaVacia_NoDisparaChequeo(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return &entities.Vehicle{ID: vehicleID, LicensePlate: "1020304050"}, nil
		},
	}

	_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
		dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26, LicensePlate: ""})

	require.NoError(t, err)
	assert.Empty(t, repo.ExistsCalls)
	assert.Empty(t, repo.Updated.LicensePlate, "se guarda vacia, es lo que pidieron")
}

func TestUpdateVehicle_ExcluyeElPropioIDAlBuscarDuplicados(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return &entities.Vehicle{ID: vehicleID, LicensePlate: "vieja"}, nil
		},
	}

	_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
		dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26, LicensePlate: "nueva"})

	require.NoError(t, err)
	require.Len(t, repo.ExistsCalls, 1)
	require.NotNil(t, repo.ExistsCalls[0].ExcludeID)
	assert.Equal(t, uint(5), *repo.ExistsCalls[0].ExcludeID)
}

func TestUpdateVehicle_PlacaDeOtro_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return &entities.Vehicle{ID: vehicleID, LicensePlate: "vieja"}, nil
		},
		ExistsByLicensePlateFn: func(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}

	_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
		dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26, LicensePlate: "ocupada"})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateLicensePlate)
	assert.Nil(t, repo.Updated)
}

func TestUpdateVehicle_ErroresDelRepo_SePropagan(t *testing.T) {
	dbErr := stderrors.New("db caida")

	t.Run("vehiculo inexistente", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
				return nil, domainerrors.ErrVehicleNotFound
			},
		}
		_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
			dtos.UpdateVehicleDTO{ID: 404, BusinessID: 26})
		assert.ErrorIs(t, err, domainerrors.ErrVehicleNotFound)
		assert.Nil(t, repo.Updated)
	})

	t.Run("al verificar duplicado", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
				return &entities.Vehicle{ID: vehicleID, LicensePlate: "vieja"}, nil
			},
			ExistsByLicensePlateFn: func(ctx context.Context, businessID uint, plate string, excludeID *uint) (bool, error) {
				return false, dbErr
			},
		}
		_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
			dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26, LicensePlate: "nueva"})
		assert.ErrorIs(t, err, dbErr)
		assert.Nil(t, repo.Updated)
	})

	t.Run("al guardar", func(t *testing.T) {
		repo := &mocks.RepositoryMock{
			UpdateFn: func(ctx context.Context, d *entities.Vehicle) (*entities.Vehicle, error) {
				return nil, dbErr
			},
		}
		_, err := newVehiclesUseCase(repo).UpdateVehicle(context.Background(),
			dtos.UpdateVehicleDTO{ID: 5, BusinessID: 26})
		assert.ErrorIs(t, err, dbErr)
	})
}

func TestDeleteVehicle_VerificaAntesDeBorrar(t *testing.T) {
	var vistoNegocio, vistoID uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			vistoNegocio, vistoID = businessID, vehicleID
			return &entities.Vehicle{ID: vehicleID}, nil
		},
	}

	err := newVehiclesUseCase(repo).DeleteVehicle(context.Background(), 26, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio,
		"no se puede borrar el vehiculo de otro negocio adivinando el ID")
	assert.Equal(t, uint(5), vistoID)
	assert.Equal(t, []uint{5}, repo.Deleted)
}

func TestDeleteVehicle_Inexistente_NoBorra(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return nil, domainerrors.ErrVehicleNotFound
		},
	}

	err := newVehiclesUseCase(repo).DeleteVehicle(context.Background(), 26, 404)

	assert.ErrorIs(t, err, domainerrors.ErrVehicleNotFound)
	assert.Empty(t, repo.Deleted)
}

func TestDeleteVehicle_ErrorAlBorrar_SePropaga(t *testing.T) {
	dbErr := stderrors.New("fk violation: tiene rutas asignadas")
	repo := &mocks.RepositoryMock{
		DeleteFn: func(ctx context.Context, businessID, vehicleID uint) error { return dbErr },
	}

	err := newVehiclesUseCase(repo).DeleteVehicle(context.Background(), 26, 5)

	assert.ErrorIs(t, err, dbErr)
}

func TestGetVehicle_FiltraPorNegocio(t *testing.T) {
	var vistoNegocio uint
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			vistoNegocio = businessID
			return &entities.Vehicle{ID: vehicleID, Brand: "Chevrolet"}, nil
		},
	}

	got, err := newVehiclesUseCase(repo).GetVehicle(context.Background(), 26, 5)

	require.NoError(t, err)
	assert.Equal(t, uint(26), vistoNegocio)
	assert.Equal(t, "Chevrolet", got.Brand)
}

func TestGetVehicle_ErrorDelRepo_SePropaga(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetByIDFn: func(ctx context.Context, businessID, vehicleID uint) (*entities.Vehicle, error) {
			return nil, domainerrors.ErrVehicleNotFound
		},
	}

	_, err := newVehiclesUseCase(repo).GetVehicle(context.Background(), 26, 404)

	assert.ErrorIs(t, err, domainerrors.ErrVehicleNotFound)
}

func TestListVehicles_NormalizaPaginacion(t *testing.T) {
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
		var visto dtos.ListVehiclesParams
		repo := &mocks.RepositoryMock{
			ListFn: func(ctx context.Context, p dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error) {
				visto = p
				return nil, 0, nil
			},
		}

		_, _, err := newVehiclesUseCase(repo).ListVehicles(context.Background(),
			dtos.ListVehiclesParams{Page: tc.page, PageSize: tc.pageSize})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, visto.Page)
		assert.Equal(t, tc.wantSize, visto.PageSize)
	}
}

func TestListVehiclesParams_Offset(t *testing.T) {
	assert.Equal(t, 0, dtos.ListVehiclesParams{Page: 1, PageSize: 20}.Offset())
	assert.Equal(t, 20, dtos.ListVehiclesParams{Page: 2, PageSize: 20}.Offset())
	assert.Equal(t, 100, dtos.ListVehiclesParams{Page: 3, PageSize: 50}.Offset())
}

func TestListVehicles_PasaLosFiltros(t *testing.T) {
	var visto dtos.ListVehiclesParams
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, p dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error) {
			visto = p
			return []entities.Vehicle{{ID: 1}}, 1, nil
		},
	}

	got, total, err := newVehiclesUseCase(repo).ListVehicles(context.Background(),
		dtos.ListVehiclesParams{BusinessID: 26, Search: "XYZ", Status: "active", Page: 1, PageSize: 20})

	require.NoError(t, err)
	assert.Equal(t, uint(26), visto.BusinessID)
	assert.Equal(t, "XYZ", visto.Search)
	assert.Equal(t, "active", visto.Status)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestListVehicles_ErrorDelRepo_SePropaga(t *testing.T) {
	dbErr := stderrors.New("timeout")
	repo := &mocks.RepositoryMock{
		ListFn: func(ctx context.Context, p dtos.ListVehiclesParams) ([]entities.Vehicle, int64, error) {
			return nil, 0, dbErr
		},
	}

	_, _, err := newVehiclesUseCase(repo).ListVehicles(context.Background(), dtos.ListVehiclesParams{})

	assert.ErrorIs(t, err, dbErr)
}

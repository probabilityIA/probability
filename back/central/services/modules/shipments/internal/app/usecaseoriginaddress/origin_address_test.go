package usecaseoriginaddress

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func s(v string) *string {
	return &v
}

func b(v bool) *bool {
	return &v
}

func addressOf(id, businessID uint) *domain.OriginAddress {
	return &domain.OriginAddress{
		ID:         id,
		BusinessID: businessID,
		Alias:      "Bodega",
		City:       "Bogota",
	}
}

func TestCreate_PrimeraDireccion_QuedaPorDefecto(t *testing.T) {
	var saved *domain.OriginAddress
	setDefaultLlamado := false
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			return []domain.OriginAddress{}, nil
		},
		CreateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			address.ID = 1
			saved = address
			return nil
		},
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error {
			setDefaultLlamado = true
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.Create(context.Background(), 36, domain.CreateOriginAddressRequest{
		Alias:     "Bodega Norte",
		IsDefault: false,
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, saved.IsDefault)
	assert.False(t, setDefaultLlamado)
	assert.Equal(t, uint(36), saved.BusinessID)
}

func TestCreate_SegundaDireccionMarcadaPorDefecto_ActualizaLaPorDefecto(t *testing.T) {
	var bizID, addrID uint
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			return []domain.OriginAddress{*addressOf(1, 36)}, nil
		},
		CreateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			address.ID = 2
			return nil
		},
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error {
			bizID, addrID = businessID, addressID
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.Create(context.Background(), 36, domain.CreateOriginAddressRequest{IsDefault: true})

	require.NoError(t, err)
	assert.Equal(t, uint(2), got.ID)
	assert.Equal(t, uint(36), bizID)
	assert.Equal(t, uint(2), addrID)
}

func TestCreate_SegundaDireccionSinDefault_NoTocaLaPorDefecto(t *testing.T) {
	setDefaultLlamado := false
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			return []domain.OriginAddress{*addressOf(1, 36)}, nil
		},
		CreateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			address.ID = 2
			return nil
		},
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error {
			setDefaultLlamado = true
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.Create(context.Background(), 36, domain.CreateOriginAddressRequest{IsDefault: false})

	require.NoError(t, err)
	assert.False(t, setDefaultLlamado)
}

func TestCreate_CopiaTodosLosCampos(t *testing.T) {
	var saved *domain.OriginAddress
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			return []domain.OriginAddress{}, nil
		},
		CreateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			saved = address
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.Create(context.Background(), 36, domain.CreateOriginAddressRequest{
		Alias:        "Bodega Norte",
		Company:      "Probability SAS",
		FirstName:    "Ana",
		LastName:     "Gomez",
		Email:        "ana@test.com",
		Phone:        "3001234567",
		Street:       "Calle 1 # 2-3",
		Suburb:       "Chapinero",
		CityDaneCode: "11001",
		City:         "Bogota",
		State:        "Cundinamarca",
		PostalCode:   "110111",
	})

	require.NoError(t, err)
	assert.Equal(t, "Bodega Norte", saved.Alias)
	assert.Equal(t, "Probability SAS", saved.Company)
	assert.Equal(t, "Ana", saved.FirstName)
	assert.Equal(t, "Gomez", saved.LastName)
	assert.Equal(t, "ana@test.com", saved.Email)
	assert.Equal(t, "3001234567", saved.Phone)
	assert.Equal(t, "Calle 1 # 2-3", saved.Street)
	assert.Equal(t, "Chapinero", saved.Suburb)
	assert.Equal(t, "11001", saved.CityDaneCode)
	assert.Equal(t, "Bogota", saved.City)
	assert.Equal(t, "Cundinamarca", saved.State)
	assert.Equal(t, "110111", saved.PostalCode)
}

func TestCreate_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("insert failed")
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			return []domain.OriginAddress{}, nil
		},
		CreateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error { return dbErr },
	}
	uc := New(repo)

	got, err := uc.Create(context.Background(), 36, domain.CreateOriginAddressRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreate_FallaMarcarPorDefecto_PropagaError(t *testing.T) {
	dbErr := errors.New("set default failed")
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			return []domain.OriginAddress{*addressOf(1, 36)}, nil
		},
		CreateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			address.ID = 2
			return nil
		},
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error { return dbErr },
	}
	uc := New(repo)

	got, err := uc.Create(context.Background(), 36, domain.CreateOriginAddressRequest{IsDefault: true})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestList_DelegaEnRepositorio(t *testing.T) {
	var pedido uint
	repo := &mocks.RepositoryMock{
		ListOriginAddressesByBusinessFn: func(ctx context.Context, businessID uint) ([]domain.OriginAddress, error) {
			pedido = businessID
			return []domain.OriginAddress{*addressOf(1, 36), *addressOf(2, 36)}, nil
		},
	}
	uc := New(repo)

	got, err := uc.List(context.Background(), 36)

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, uint(36), pedido)
}

func TestGetByID_DireccionDeOtroNegocio_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
	}
	uc := New(repo)

	got, err := uc.GetByID(context.Background(), 1, 99)

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "no tienes permiso")
}

func TestGetByID_DireccionPropia_Retorna(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
	}
	uc := New(repo)

	got, err := uc.GetByID(context.Background(), 1, 36)

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(1), got.ID)
}

func TestGetByID_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.GetByID(context.Background(), 1, 36)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdate_DireccionDeOtroNegocio_Rechaza(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			actualizado = true
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.Update(context.Background(), 1, 99, domain.UpdateOriginAddressRequest{Alias: s("hackeado")})

	require.Error(t, err)
	assert.Nil(t, got)
	assert.Contains(t, err.Error(), "no tienes permiso")
	assert.False(t, actualizado)
}

func TestUpdate_ActualizaTodosLosCamposOpcionales(t *testing.T) {
	var saved *domain.OriginAddress
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			saved = address
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{
		Alias:        s("Bodega Sur"),
		Company:      s("Otra SAS"),
		FirstName:    s("Luis"),
		LastName:     s("Diaz"),
		Email:        s("luis@test.com"),
		Phone:        s("3009999999"),
		Street:       s("Carrera 9 # 8-7"),
		Suburb:       s("Kennedy"),
		CityDaneCode: s("05001"),
		City:         s("Medellin"),
		State:        s("Antioquia"),
		PostalCode:   s("050001"),
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Bodega Sur", saved.Alias)
	assert.Equal(t, "Otra SAS", saved.Company)
	assert.Equal(t, "Luis", saved.FirstName)
	assert.Equal(t, "Diaz", saved.LastName)
	assert.Equal(t, "luis@test.com", saved.Email)
	assert.Equal(t, "3009999999", saved.Phone)
	assert.Equal(t, "Carrera 9 # 8-7", saved.Street)
	assert.Equal(t, "Kennedy", saved.Suburb)
	assert.Equal(t, "05001", saved.CityDaneCode)
	assert.Equal(t, "Medellin", saved.City)
	assert.Equal(t, "Antioquia", saved.State)
	assert.Equal(t, "050001", saved.PostalCode)
}

func TestUpdate_SinCampos_NoCambiaNada(t *testing.T) {
	var saved *domain.OriginAddress
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error {
			saved = address
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{})

	require.NoError(t, err)
	assert.Equal(t, "Bodega", saved.Alias)
	assert.Equal(t, "Bogota", saved.City)
}

func TestUpdate_MarcarPorDefecto_LlamaSetDefault(t *testing.T) {
	var bizID, addrID uint
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error { return nil },
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error {
			bizID, addrID = businessID, addressID
			return nil
		},
	}
	uc := New(repo)

	got, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{IsDefault: b(true)})

	require.NoError(t, err)
	assert.True(t, got.IsDefault)
	assert.Equal(t, uint(36), bizID)
	assert.Equal(t, uint(1), addrID)
}

func TestUpdate_IsDefaultFalse_NoLlamaSetDefault(t *testing.T) {
	llamado := false
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error { return nil },
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error {
			llamado = true
			return nil
		},
	}
	uc := New(repo)

	_, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{IsDefault: b(false)})

	require.NoError(t, err)
	assert.False(t, llamado)
}

func TestUpdate_FallaGuardar_PropagaError(t *testing.T) {
	dbErr := errors.New("update failed")
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn: func(ctx context.Context, address *domain.OriginAddress) error { return dbErr },
	}
	uc := New(repo)

	got, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdate_FallaSetDefault_PropagaError(t *testing.T) {
	dbErr := errors.New("set default failed")
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		UpdateOriginAddressFn:     func(ctx context.Context, address *domain.OriginAddress) error { return nil },
		SetDefaultOriginAddressFn: func(ctx context.Context, businessID, addressID uint) error { return dbErr },
	}
	uc := New(repo)

	got, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{IsDefault: b(true)})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdate_FallaLeer_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	got, err := uc.Update(context.Background(), 1, 36, domain.UpdateOriginAddressRequest{})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDelete_DireccionDeOtroNegocio_Rechaza(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		DeleteOriginAddressFn: func(ctx context.Context, id uint) error {
			borrado = true
			return nil
		},
	}
	uc := New(repo)

	err := uc.Delete(context.Background(), 1, 99)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no tienes permiso")
	assert.False(t, borrado)
}

func TestDelete_DireccionPorDefecto_Rechaza(t *testing.T) {
	borrado := false
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			a := addressOf(id, 36)
			a.IsDefault = true
			return a, nil
		},
		DeleteOriginAddressFn: func(ctx context.Context, id uint) error {
			borrado = true
			return nil
		},
	}
	uc := New(repo)

	err := uc.Delete(context.Background(), 1, 36)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "no puedes eliminar la direcci")
	assert.False(t, borrado)
}

func TestDelete_Exitoso(t *testing.T) {
	var borradoID uint
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		DeleteOriginAddressFn: func(ctx context.Context, id uint) error {
			borradoID = id
			return nil
		},
	}
	uc := New(repo)

	err := uc.Delete(context.Background(), 7, 36)

	require.NoError(t, err)
	assert.Equal(t, uint(7), borradoID)
}

func TestDelete_FallaLeer_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return nil, dbErr
		},
	}
	uc := New(repo)

	err := uc.Delete(context.Background(), 1, 36)

	assert.ErrorIs(t, err, dbErr)
}

func TestDelete_FallaBorrar_PropagaError(t *testing.T) {
	dbErr := errors.New("delete failed")
	repo := &mocks.RepositoryMock{
		GetOriginAddressByIDFn: func(ctx context.Context, id uint) (*domain.OriginAddress, error) {
			return addressOf(id, 36), nil
		},
		DeleteOriginAddressFn: func(ctx context.Context, id uint) error { return dbErr },
	}
	uc := New(repo)

	err := uc.Delete(context.Background(), 1, 36)

	assert.ErrorIs(t, err, dbErr)
}

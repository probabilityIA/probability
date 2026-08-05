package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func s(v string) *string {
	return &v
}

func clienteBase(id uint) *entities.Client {
	return &entities.Client{
		ID:         id,
		BusinessID: 10,
		Name:       "Ana Gomez",
		Email:      s("ana@test.com"),
		Phone:      "3001234567",
		Dni:        s("1020304050"),
	}
}

func TestCreateClient_EmailDuplicado_Rechaza(t *testing.T) {
	creado := false
	repo := &mockRepo{
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			return true, nil
		},
		createFn: func(ctx context.Context, c *entities.Client) (*entities.Client, error) {
			creado = true
			return c, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{
		BusinessID: 10,
		Name:       "Ana",
		Email:      s("ana@test.com"),
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateEmail)
	assert.Nil(t, got)
	assert.False(t, creado)
}

func TestCreateClient_DniDuplicado_Rechaza(t *testing.T) {
	repo := &mockRepo{
		existsByDniFn: func(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{
		BusinessID: 10,
		Name:       "Ana",
		Dni:        s("1020304050"),
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateDni)
	assert.Nil(t, got)
}

func TestCreateClient_VerificaDuplicadosDentroDelNegocio(t *testing.T) {
	var negocioEmail, negocioDni uint
	repo := &mockRepo{
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			negocioEmail = businessID
			assert.Nil(t, excludeID, "al crear no se excluye ningun cliente")
			return false, nil
		},
		existsByDniFn: func(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
			negocioDni = businessID
			return false, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{
		BusinessID: 36,
		Name:       "Ana",
		Email:      s("ana@test.com"),
		Dni:        s("1020304050"),
	})

	require.NoError(t, err)
	assert.Equal(t, uint(36), negocioEmail,
		"la unicidad es por negocio, no global")
	assert.Equal(t, uint(36), negocioDni)
}

func TestCreateClient_EmailVacio_NoVerificaNiGuarda(t *testing.T) {
	verificado := false
	var guardado *entities.Client
	repo := &mockRepo{
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			verificado = true
			return false, nil
		},
		createFn: func(ctx context.Context, c *entities.Client) (*entities.Client, error) {
			guardado = c
			return c, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{
		BusinessID: 10,
		Name:       "Ana",
		Email:      s(""),
	})

	require.NoError(t, err)
	assert.False(t, verificado)
	assert.Nil(t, guardado.Email,
		"un email vacio se normaliza a nulo para no chocar con el indice unico")
}

func TestCreateClient_SinEmail_NoVerifica(t *testing.T) {
	verificado := false
	repo := &mockRepo{
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			verificado = true
			return false, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{BusinessID: 10, Name: "Ana"})

	require.NoError(t, err)
	assert.False(t, verificado)
}

func TestCreateClient_Exitoso_CopiaLosCampos(t *testing.T) {
	var guardado *entities.Client
	repo := &mockRepo{
		createFn: func(ctx context.Context, c *entities.Client) (*entities.Client, error) {
			guardado = c
			return c, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{
		BusinessID: 10,
		Name:       "Ana Gomez",
		Email:      s("ana@test.com"),
		Phone:      "3001234567",
		Dni:        s("1020304050"),
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, uint(10), guardado.BusinessID)
	assert.Equal(t, "Ana Gomez", guardado.Name)
	assert.Equal(t, "ana@test.com", *guardado.Email)
	assert.Equal(t, "3001234567", guardado.Phone)
	assert.Equal(t, "1020304050", *guardado.Dni)
}

func TestCreateClient_FallaVerificarEmail_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mockRepo{
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.CreateClient(context.Background(), dtos.CreateClientDTO{
		BusinessID: 10,
		Email:      s("ana@test.com"),
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateClient_ClienteInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return nil, dbErr
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{ID: 1, BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateClient_MismoEmail_NoVerificaDuplicado(t *testing.T) {
	verificado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			verificado = true
			return true, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Email:      s("ana@test.com"),
	})

	require.NoError(t, err)
	assert.False(t, verificado, "guardar el mismo email no dispara la verificacion")
}

func TestUpdateClient_EmailNuevoDuplicado_Rechaza(t *testing.T) {
	var excluido *uint
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			excluido = excludeID
			return true, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Email:      s("otro@test.com"),
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateEmail)
	assert.Nil(t, got)
	require.NotNil(t, excluido)
	assert.Equal(t, uint(1), *excluido,
		"al actualizar se excluye el propio cliente de la verificacion")
}

func TestUpdateClient_MismoDni_NoVerifica(t *testing.T) {
	verificado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		existsByDniFn: func(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
			verificado = true
			return true, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Dni:        s("1020304050"),
	})

	require.NoError(t, err)
	assert.False(t, verificado)
}

func TestUpdateClient_DniNuevoDuplicado_Rechaza(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		existsByDniFn: func(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
			return true, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Dni:        s("9999999999"),
	})

	assert.ErrorIs(t, err, domainerrors.ErrDuplicateDni)
	assert.Nil(t, got)
}

func TestUpdateClient_ClienteSinDniPrevio_VerificaElNuevo(t *testing.T) {
	verificado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			c := clienteBase(clientID)
			c.Dni = nil
			return c, nil
		},
		existsByDniFn: func(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
			verificado = true
			return false, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Dni:        s("1020304050"),
	})

	require.NoError(t, err)
	assert.True(t, verificado, "asignar un DNI a un cliente que no tenia si verifica duplicados")
}

func TestUpdateClient_ClienteSinEmailPrevio_VerificaElNuevo(t *testing.T) {
	verificado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			c := clienteBase(clientID)
			c.Email = nil
			return c, nil
		},
		existsByEmailFn: func(ctx context.Context, businessID uint, email string, excludeID *uint) (bool, error) {
			verificado = true
			return false, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Email:      s("nuevo@test.com"),
	})

	require.NoError(t, err)
	assert.True(t, verificado)
}

func TestUpdateClient_EmailVacio_SeGuardaComoNulo(t *testing.T) {
	var guardado *entities.Client
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		updateFn: func(ctx context.Context, c *entities.Client) (*entities.Client, error) {
			guardado = c
			return c, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Name:       "Ana",
		Email:      s(""),
	})

	require.NoError(t, err)
	assert.Nil(t, guardado.Email)
}

func TestUpdateClient_Exitoso_ActualizaLosCampos(t *testing.T) {
	var guardado *entities.Client
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		updateFn: func(ctx context.Context, c *entities.Client) (*entities.Client, error) {
			guardado = c
			return c, nil
		},
	}
	uc := newTestUseCase(repo)

	_, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Name:       "Ana Maria",
		Phone:      "3009999999",
	})

	require.NoError(t, err)
	assert.Equal(t, "Ana Maria", guardado.Name)
	assert.Equal(t, "3009999999", guardado.Phone)
}

func TestUpdateClient_FallaVerificarDni_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		existsByDniFn: func(ctx context.Context, businessID uint, dni string, excludeID *uint) (bool, error) {
			return false, dbErr
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.UpdateClient(context.Background(), dtos.UpdateClientDTO{
		ID:         1,
		BusinessID: 10,
		Dni:        s("9999999999"),
	})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetClient_EnriqueceConElResumen(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return &entities.CustomerSummary{TotalOrders: 12, TotalSpent: 450000}, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.GetClient(context.Background(), 10, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(12), got.OrderCount)
	assert.InDelta(t, 450000.0, got.TotalSpent, 0.001)
}

func TestGetClient_FallaElResumen_DevuelveElClienteIgual(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return nil, errors.New("db down")
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.GetClient(context.Background(), 10, 1)

	require.NoError(t, err, "un fallo en el resumen no impide ver la ficha del cliente")
	require.NotNil(t, got)
	assert.Equal(t, int64(0), got.OrderCount)
}

func TestGetClient_SinResumen_NoRompe(t *testing.T) {
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return nil, nil
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.GetClient(context.Background(), 10, 1)

	require.NoError(t, err)
	assert.Equal(t, int64(0), got.OrderCount)
}

func TestGetClient_ClienteInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return nil, dbErr
		},
	}
	uc := newTestUseCase(repo)

	got, err := uc.GetClient(context.Background(), 10, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestListClients_NormalizaPaginacion(t *testing.T) {
	casos := []struct {
		nombre       string
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{nombre: "pagina cero", page: 0, pageSize: 20, wantPage: 1, wantPageSize: 20},
		{nombre: "pagina negativa", page: -3, pageSize: 20, wantPage: 1, wantPageSize: 20},
		{nombre: "page size cero usa 20", page: 1, pageSize: 0, wantPage: 1, wantPageSize: 20},
		{nombre: "page size sobre el maximo usa 20", page: 1, pageSize: 500, wantPage: 1, wantPageSize: 20},
		{nombre: "page size en el limite", page: 2, pageSize: 100, wantPage: 2, wantPageSize: 100},
		{nombre: "valores validos", page: 3, pageSize: 50, wantPage: 3, wantPageSize: 50},
	}

	for _, tc := range casos {
		t.Run(tc.nombre, func(t *testing.T) {
			var recibido dtos.ListClientsParams
			repo := &mockRepo{
				listFn: func(ctx context.Context, p dtos.ListClientsParams) ([]entities.Client, int64, error) {
					recibido = p
					return nil, 0, nil
				},
			}
			uc := newTestUseCase(repo)

			_, _, err := uc.ListClients(context.Background(), dtos.ListClientsParams{
				Page:     tc.page,
				PageSize: tc.pageSize,
			})

			require.NoError(t, err)
			assert.Equal(t, tc.wantPage, recibido.Page)
			assert.Equal(t, tc.wantPageSize, recibido.PageSize)
		})
	}
}

func TestListClients_PropagaLosResultados(t *testing.T) {
	repo := &mockRepo{
		listFn: func(ctx context.Context, p dtos.ListClientsParams) ([]entities.Client, int64, error) {
			return []entities.Client{*clienteBase(1), *clienteBase(2)}, 2, nil
		},
	}
	uc := newTestUseCase(repo)

	got, total, err := uc.ListClients(context.Background(), dtos.ListClientsParams{BusinessID: 10})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
}

func TestListClients_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("query failed")
	repo := &mockRepo{
		listFn: func(ctx context.Context, p dtos.ListClientsParams) ([]entities.Client, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := newTestUseCase(repo)

	_, _, err := uc.ListClients(context.Background(), dtos.ListClientsParams{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
}

func TestDeleteClient_ConOrdenes_Rechaza(t *testing.T) {
	borrado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return &entities.CustomerSummary{TotalOrders: 3}, nil
		},
		deleteFn: func(ctx context.Context, businessID, clientID uint) error {
			borrado = true
			return nil
		},
	}
	uc := newTestUseCase(repo)

	err := uc.DeleteClient(context.Background(), 10, 1)

	assert.ErrorIs(t, err, domainerrors.ErrClientHasOrders)
	assert.Contains(t, err.Error(), "3 orden")
	assert.False(t, borrado, "no se borra un cliente con historial de ordenes")
}

func TestDeleteClient_SinOrdenes_SeBorra(t *testing.T) {
	var borradoID uint
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return &entities.CustomerSummary{TotalOrders: 0}, nil
		},
		deleteFn: func(ctx context.Context, businessID, clientID uint) error {
			borradoID = clientID
			return nil
		},
	}
	uc := newTestUseCase(repo)

	require.NoError(t, uc.DeleteClient(context.Background(), 10, 1))
	assert.Equal(t, uint(1), borradoID)
}

func TestDeleteClient_SinResumen_SeBorra(t *testing.T) {
	borrado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return nil, nil
		},
		deleteFn: func(ctx context.Context, businessID, clientID uint) error {
			borrado = true
			return nil
		},
	}
	uc := newTestUseCase(repo)

	require.NoError(t, uc.DeleteClient(context.Background(), 10, 1))
	assert.True(t, borrado)
}

func TestDeleteClient_ClienteInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return nil, dbErr
		},
	}
	uc := newTestUseCase(repo)

	assert.ErrorIs(t, uc.DeleteClient(context.Background(), 10, 1), dbErr)
}

func TestDeleteClient_FallaConsultarResumen_NoBorra(t *testing.T) {
	dbErr := errors.New("db down")
	borrado := false
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		getCustomerSummaryFn: func(ctx context.Context, businessID, customerID uint) (*entities.CustomerSummary, error) {
			return nil, dbErr
		},
		deleteFn: func(ctx context.Context, businessID, clientID uint) error {
			borrado = true
			return nil
		},
	}
	uc := newTestUseCase(repo)

	err := uc.DeleteClient(context.Background(), 10, 1)

	assert.ErrorIs(t, err, dbErr)
	assert.False(t, borrado,
		"sin poder verificar el historial no se borra al cliente")
}

func TestDeleteClient_FallaBorrar_PropagaError(t *testing.T) {
	dbErr := errors.New("delete failed")
	repo := &mockRepo{
		getByIDFn: func(ctx context.Context, businessID, clientID uint) (*entities.Client, error) {
			return clienteBase(clientID), nil
		},
		deleteFn: func(ctx context.Context, businessID, clientID uint) error {
			return dbErr
		},
	}
	uc := newTestUseCase(repo)

	assert.ErrorIs(t, uc.DeleteClient(context.Background(), 10, 1), dbErr)
}

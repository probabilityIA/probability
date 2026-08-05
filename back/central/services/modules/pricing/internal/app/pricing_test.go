package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func u(v uint) *uint {
	return &v
}

func f(v float64) *float64 {
	return &v
}

func buildUseCase(repo *mocks.RepositoryMock) IUseCase {
	return New(repo, mocks.NewSilentLogger())
}

func targetGrupo() dtos.CatalogPriceTarget {
	return dtos.CatalogPriceTarget{BusinessID: 10, ClientGroupID: u(3)}
}

func targetCliente() dtos.CatalogPriceTarget {
	return dtos.CatalogPriceTarget{BusinessID: 10, ClientID: u(77)}
}

func TestValidateTarget_SinDestino_RetornaError(t *testing.T) {
	err := validateTarget(dtos.CatalogPriceTarget{BusinessID: 10})

	assert.ErrorIs(t, err, domainerrors.ErrTargetRequired)
}

func TestValidateTarget_AmbosDestinos_RetornaError(t *testing.T) {
	err := validateTarget(dtos.CatalogPriceTarget{BusinessID: 10, ClientGroupID: u(3), ClientID: u(77)})

	assert.ErrorIs(t, err, domainerrors.ErrTargetAmbiguous)
}

func TestValidateTarget_SoloGrupo_EsValido(t *testing.T) {
	assert.NoError(t, validateTarget(targetGrupo()))
}

func TestValidateTarget_SoloCliente_EsValido(t *testing.T) {
	assert.NoError(t, validateTarget(targetCliente()))
}

func TestListCatalogPrices_DestinoInvalido_NoConsultaElRepositorio(t *testing.T) {
	consultado := false
	repo := &mocks.RepositoryMock{
		ListCatalogPricesFn: func(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error) {
			consultado = true
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo)

	got, total, err := uc.ListCatalogPrices(context.Background(), dtos.ListCatalogPricesParams{
		Target: dtos.CatalogPriceTarget{BusinessID: 10},
	})

	assert.ErrorIs(t, err, domainerrors.ErrTargetRequired)
	assert.Nil(t, got)
	assert.Equal(t, int64(0), total)
	assert.False(t, consultado)
}

func TestListCatalogPrices_DestinoValido_Consulta(t *testing.T) {
	var recibido dtos.ListCatalogPricesParams
	repo := &mocks.RepositoryMock{
		ListCatalogPricesFn: func(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error) {
			recibido = params
			return []entities.CatalogPriceRow{{ProductID: "prod-1"}}, 1, nil
		},
	}
	uc := buildUseCase(repo)

	got, total, err := uc.ListCatalogPrices(context.Background(), dtos.ListCatalogPricesParams{
		Target: targetGrupo(),
		Search: "camisa",
	})

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
	assert.Equal(t, "camisa", recibido.Search)
	require.NotNil(t, recibido.Target.ClientGroupID)
	assert.Equal(t, uint(3), *recibido.Target.ClientGroupID)
}

func TestListCatalogPrices_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		ListCatalogPricesFn: func(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error) {
			return nil, 0, dbErr
		},
	}
	uc := buildUseCase(repo)

	_, _, err := uc.ListCatalogPrices(context.Background(), dtos.ListCatalogPricesParams{Target: targetGrupo()})

	assert.ErrorIs(t, err, dbErr)
}

func TestSaveCatalogPrices_DestinoInvalido_NoGuarda(t *testing.T) {
	guardado := false
	repo := &mocks.RepositoryMock{
		SaveCatalogPricesFn: func(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
			guardado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	err := uc.SaveCatalogPrices(context.Background(), dtos.SaveCatalogPricesDTO{
		Target: dtos.CatalogPriceTarget{BusinessID: 10, ClientGroupID: u(3), ClientID: u(77)},
	})

	assert.ErrorIs(t, err, domainerrors.ErrTargetAmbiguous)
	assert.False(t, guardado)
}

func TestSaveCatalogPrices_PrecioNegativo_Rechaza(t *testing.T) {
	guardado := false
	repo := &mocks.RepositoryMock{
		SaveCatalogPricesFn: func(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
			guardado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	err := uc.SaveCatalogPrices(context.Background(), dtos.SaveCatalogPricesDTO{
		Target: targetGrupo(),
		Items: []dtos.SaveCatalogPriceItem{
			{ProductID: "prod-1", Price: f(1000)},
			{ProductID: "prod-2", Price: f(-1)},
		},
	})

	assert.ErrorIs(t, err, domainerrors.ErrInvalidPrice)
	assert.False(t, guardado, "un solo precio negativo aborta el lote completo")
}

func TestSaveCatalogPrices_PrecioCero_EsValido(t *testing.T) {
	guardado := false
	repo := &mocks.RepositoryMock{
		SaveCatalogPricesFn: func(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
			guardado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	err := uc.SaveCatalogPrices(context.Background(), dtos.SaveCatalogPricesDTO{
		Target: targetGrupo(),
		Items:  []dtos.SaveCatalogPriceItem{{ProductID: "prod-1", Price: f(0)}},
	})

	require.NoError(t, err)
	assert.True(t, guardado)
}

func TestSaveCatalogPrices_PrecioNulo_SeAcepta(t *testing.T) {
	guardado := false
	repo := &mocks.RepositoryMock{
		SaveCatalogPricesFn: func(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
			guardado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	err := uc.SaveCatalogPrices(context.Background(), dtos.SaveCatalogPricesDTO{
		Target: targetCliente(),
		Items:  []dtos.SaveCatalogPriceItem{{ProductID: "prod-1", Price: nil}},
	})

	require.NoError(t, err, "un precio nulo limpia el precio especial del producto")
	assert.True(t, guardado)
}

func TestSaveCatalogPrices_SinItems_SeAcepta(t *testing.T) {
	uc := buildUseCase(&mocks.RepositoryMock{})

	err := uc.SaveCatalogPrices(context.Background(), dtos.SaveCatalogPricesDTO{Target: targetGrupo()})

	require.NoError(t, err)
}

func TestSaveCatalogPrices_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		SaveCatalogPricesFn: func(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
			return dbErr
		},
	}
	uc := buildUseCase(repo)

	err := uc.SaveCatalogPrices(context.Background(), dtos.SaveCatalogPricesDTO{Target: targetGrupo()})

	assert.ErrorIs(t, err, dbErr)
}

func TestCreateClientGroup_NombreVacio_Rechaza(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		CreateClientGroupFn: func(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
			creado = true
			return group, nil
		},
	}
	uc := buildUseCase(repo)

	for _, nombre := range []string{"", "   ", "\t\n"} {
		got, err := uc.CreateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
			BusinessID: 10,
			Name:       nombre,
		})

		assert.ErrorIs(t, err, domainerrors.ErrGroupNameRequired)
		assert.Nil(t, got)
	}
	assert.False(t, creado)
}

func TestCreateClientGroup_NombreDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GroupNameExistsFn: func(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
			return true, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.CreateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		BusinessID: 10,
		Name:       "Mayoristas",
	})

	assert.ErrorIs(t, err, domainerrors.ErrGroupNameDuplicate)
	assert.Nil(t, got)
}

func TestCreateClientGroup_VerificaDuplicadoConGrupoCero(t *testing.T) {
	var grupoConsultado uint
	repo := &mocks.RepositoryMock{
		GroupNameExistsFn: func(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
			grupoConsultado = groupID
			return false, nil
		},
	}
	uc := buildUseCase(repo)

	_, err := uc.CreateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		BusinessID: 10,
		Name:       "Mayoristas",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(0), grupoConsultado,
		"al crear se compara contra todos los grupos, sin excluir ninguno")
}

func TestCreateClientGroup_RecortaEspacios(t *testing.T) {
	var guardado *entities.ClientGroup
	repo := &mocks.RepositoryMock{
		CreateClientGroupFn: func(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
			guardado = group
			return group, nil
		},
	}
	uc := buildUseCase(repo)

	_, err := uc.CreateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		BusinessID:  10,
		Name:        "  Mayoristas  ",
		Description: "  Clientes al por mayor  ",
		Color:       "  #FF0000  ",
		IsActive:    true,
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "Mayoristas", guardado.Name)
	assert.Equal(t, "Clientes al por mayor", guardado.Description)
	assert.Equal(t, "#FF0000", guardado.Color)
	assert.Equal(t, uint(10), guardado.BusinessID)
	assert.True(t, guardado.IsActive)
}

func TestCreateClientGroup_FallaVerificar_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GroupNameExistsFn: func(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
			return false, dbErr
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.CreateClientGroup(context.Background(), dtos.SaveClientGroupDTO{BusinessID: 10, Name: "X"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestUpdateClientGroup_NombreVacio_Rechaza(t *testing.T) {
	uc := buildUseCase(&mocks.RepositoryMock{})

	got, err := uc.UpdateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		ID:         3,
		BusinessID: 10,
		Name:       "   ",
	})

	assert.ErrorIs(t, err, domainerrors.ErrGroupNameRequired)
	assert.Nil(t, got)
}

func TestUpdateClientGroup_SeExcluyeASiMismoAlVerificarDuplicado(t *testing.T) {
	var grupoConsultado uint
	repo := &mocks.RepositoryMock{
		GroupNameExistsFn: func(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
			grupoConsultado = groupID
			return false, nil
		},
	}
	uc := buildUseCase(repo)

	_, err := uc.UpdateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		ID:         3,
		BusinessID: 10,
		Name:       "Mayoristas",
	})

	require.NoError(t, err)
	assert.Equal(t, uint(3), grupoConsultado,
		"al actualizar se excluye el propio grupo para no chocar consigo mismo")
}

func TestUpdateClientGroup_NombreDuplicado_Rechaza(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GroupNameExistsFn: func(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
			return true, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.UpdateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		ID:         3,
		BusinessID: 10,
		Name:       "Mayoristas",
	})

	assert.ErrorIs(t, err, domainerrors.ErrGroupNameDuplicate)
	assert.Nil(t, got)
}

func TestUpdateClientGroup_ConservaElID(t *testing.T) {
	var guardado *entities.ClientGroup
	repo := &mocks.RepositoryMock{
		UpdateClientGroupFn: func(ctx context.Context, group *entities.ClientGroup) (*entities.ClientGroup, error) {
			guardado = group
			return group, nil
		},
	}
	uc := buildUseCase(repo)

	_, err := uc.UpdateClientGroup(context.Background(), dtos.SaveClientGroupDTO{
		ID:         3,
		BusinessID: 10,
		Name:       "  Mayoristas  ",
		IsActive:   false,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(3), guardado.ID)
	assert.Equal(t, "Mayoristas", guardado.Name)
	assert.False(t, guardado.IsActive)
}

func TestUpdateClientGroup_FallaVerificar_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GroupNameExistsFn: func(ctx context.Context, businessID, groupID uint, name string) (bool, error) {
			return false, dbErr
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.UpdateClientGroup(context.Background(), dtos.SaveClientGroupDTO{ID: 3, BusinessID: 10, Name: "X"})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestGetClientGroup_DelegaEnRepositorio(t *testing.T) {
	var pedidoBusiness, pedidoGrupo uint
	repo := &mocks.RepositoryMock{
		GetClientGroupFn: func(ctx context.Context, businessID, groupID uint) (*entities.ClientGroup, error) {
			pedidoBusiness, pedidoGrupo = businessID, groupID
			return &entities.ClientGroup{ID: groupID, BusinessID: businessID}, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GetClientGroup(context.Background(), 10, 3)

	require.NoError(t, err)
	assert.Equal(t, uint(3), got.ID)
	assert.Equal(t, uint(10), pedidoBusiness)
	assert.Equal(t, uint(3), pedidoGrupo)
}

func TestListClientGroups_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListClientGroupsFn: func(ctx context.Context, params dtos.ListClientGroupsParams) ([]entities.ClientGroup, int64, error) {
			return []entities.ClientGroup{{ID: 1}, {ID: 2}}, 2, nil
		},
	}
	uc := buildUseCase(repo)

	got, total, err := uc.ListClientGroups(context.Background(), dtos.ListClientGroupsParams{BusinessID: 10})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
}

func TestDeleteClientGroup_DelegaEnRepositorio(t *testing.T) {
	var borrado uint
	repo := &mocks.RepositoryMock{
		DeleteClientGroupFn: func(ctx context.Context, businessID, groupID uint) error {
			borrado = groupID
			return nil
		},
	}
	uc := buildUseCase(repo)

	require.NoError(t, uc.DeleteClientGroup(context.Background(), 10, 3))
	assert.Equal(t, uint(3), borrado)
}

func TestListGroupMembers_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListGroupMembersFn: func(ctx context.Context, params dtos.ListGroupMembersParams) ([]entities.ClientSummary, int64, error) {
			return []entities.ClientSummary{{ID: 1}}, 1, nil
		},
	}
	uc := buildUseCase(repo)

	got, total, err := uc.ListGroupMembers(context.Background(), dtos.ListGroupMembersParams{BusinessID: 10, ClientGroupID: 3})

	require.NoError(t, err)
	assert.Len(t, got, 1)
	assert.Equal(t, int64(1), total)
}

func TestListAvailableClients_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		ListAvailableClientsFn: func(ctx context.Context, params dtos.ListAvailableClientsParams) ([]entities.ClientSummary, int64, error) {
			return []entities.ClientSummary{{ID: 5}, {ID: 6}}, 2, nil
		},
	}
	uc := buildUseCase(repo)

	got, total, err := uc.ListAvailableClients(context.Background(), dtos.ListAvailableClientsParams{BusinessID: 10})

	require.NoError(t, err)
	assert.Len(t, got, 2)
	assert.Equal(t, int64(2), total)
}

func TestAddGroupMembers_SinClientes_Rechaza(t *testing.T) {
	agregado := false
	repo := &mocks.RepositoryMock{
		AddGroupMembersFn: func(ctx context.Context, dto dtos.AddGroupMembersDTO) error {
			agregado = true
			return nil
		},
	}
	uc := buildUseCase(repo)

	err := uc.AddGroupMembers(context.Background(), dtos.AddGroupMembersDTO{
		BusinessID:    10,
		ClientGroupID: 3,
	})

	assert.ErrorIs(t, err, domainerrors.ErrNoClients)
	assert.False(t, agregado)
}

func TestAddGroupMembers_ConClientes_Agrega(t *testing.T) {
	var recibido dtos.AddGroupMembersDTO
	repo := &mocks.RepositoryMock{
		AddGroupMembersFn: func(ctx context.Context, dto dtos.AddGroupMembersDTO) error {
			recibido = dto
			return nil
		},
	}
	uc := buildUseCase(repo)

	err := uc.AddGroupMembers(context.Background(), dtos.AddGroupMembersDTO{
		BusinessID:    10,
		ClientGroupID: 3,
		ClientIDs:     []uint{1, 2, 3},
	})

	require.NoError(t, err)
	assert.Equal(t, []uint{1, 2, 3}, recibido.ClientIDs)
	assert.Equal(t, uint(3), recibido.ClientGroupID)
}

func TestAddGroupMembers_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		AddGroupMembersFn: func(ctx context.Context, dto dtos.AddGroupMembersDTO) error {
			return dbErr
		},
	}
	uc := buildUseCase(repo)

	err := uc.AddGroupMembers(context.Background(), dtos.AddGroupMembersDTO{
		BusinessID: 10,
		ClientIDs:  []uint{1},
	})

	assert.ErrorIs(t, err, dbErr)
}

func TestRemoveGroupMember_DelegaEnRepositorio(t *testing.T) {
	var quitadoGrupo, quitadoCliente uint
	repo := &mocks.RepositoryMock{
		RemoveGroupMemberFn: func(ctx context.Context, businessID, groupID, clientID uint) error {
			quitadoGrupo, quitadoCliente = groupID, clientID
			return nil
		},
	}
	uc := buildUseCase(repo)

	require.NoError(t, uc.RemoveGroupMember(context.Background(), 10, 3, 77))
	assert.Equal(t, uint(3), quitadoGrupo)
	assert.Equal(t, uint(77), quitadoCliente)
}

func TestGetEffectivePrice_DelegaEnRepositorio(t *testing.T) {
	var recibido dtos.EffectivePriceParams
	repo := &mocks.RepositoryMock{
		GetEffectivePriceFn: func(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error) {
			recibido = params
			return &entities.EffectivePrice{ProductID: params.ProductID, FinalPrice: 25000}, nil
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GetEffectivePrice(context.Background(), dtos.EffectivePriceParams{
		BusinessID: 10,
		ProductID:  "prod-1",
		ClientID:   77,
	})

	require.NoError(t, err)
	require.NotNil(t, got)
	assert.InDelta(t, 25000.0, got.FinalPrice, 0.001)
	assert.Equal(t, "prod-1", recibido.ProductID)
	assert.Equal(t, uint(77), recibido.ClientID)
	assert.Equal(t, uint(10), recibido.BusinessID)
}

func TestGetEffectivePrice_FallaRepositorio_PropagaError(t *testing.T) {
	dbErr := errors.New("db down")
	repo := &mocks.RepositoryMock{
		GetEffectivePriceFn: func(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo)

	got, err := uc.GetEffectivePrice(context.Background(), dtos.EffectivePriceParams{BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

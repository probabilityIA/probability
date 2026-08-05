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

func uintPtrInv(v uint) *uint {
	return &v
}

func TestCreateCrossDockLink_CantidadInvalida_RetornaError(t *testing.T) {
	creado := false
	repo := &mocks.RepositoryMock{
		CreateCrossDockLinkFn: func(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
			creado = true
			return l, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	for _, qty := range []int{0, -3} {
		got, err := uc.CreateCrossDockLink(context.Background(), request.CreateCrossDockLinkDTO{
			BusinessID: 10,
			Quantity:   qty,
		})

		assert.ErrorIs(t, err, domainerrors.ErrInvalidQuantity)
		assert.Nil(t, got)
	}
	assert.False(t, creado)
}

func TestCreateCrossDockLink_Exitoso_NacePendiente(t *testing.T) {
	var guardado *entities.CrossDockLink
	repo := &mocks.RepositoryMock{
		CreateCrossDockLinkFn: func(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
			guardado = l
			return l, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateCrossDockLink(context.Background(), request.CreateCrossDockLinkDTO{
		BusinessID:        10,
		InboundShipmentID: uintPtrInv(100),
		OutboundOrderID:   "order-uuid",
		ProductID:         "prod-1",
		Quantity:          25,
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "pending", guardado.Status)
	require.NotNil(t, guardado.InboundShipmentID)
	assert.Equal(t, uint(100), *guardado.InboundShipmentID)
	assert.Equal(t, "order-uuid", guardado.OutboundOrderID)
	assert.Equal(t, "prod-1", guardado.ProductID)
	assert.Equal(t, 25, guardado.Quantity)
}

func TestListCrossDockLinks_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 250, wantPage: 1, wantPageSize: 10},
		{page: 3, pageSize: 20, wantPage: 3, wantPageSize: 20},
	}

	for _, tc := range cases {
		var recibido dtos.ListCrossDockLinksParams
		repo := &mocks.RepositoryMock{
			ListCrossDockLinksFn: func(ctx context.Context, params dtos.ListCrossDockLinksParams) ([]entities.CrossDockLink, int64, error) {
				recibido = params
				return nil, 0, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, _, err := uc.ListCrossDockLinks(context.Background(), dtos.ListCrossDockLinksParams{
			Page:     tc.page,
			PageSize: tc.pageSize,
		})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, recibido.Page)
		assert.Equal(t, tc.wantPageSize, recibido.PageSize)
	}
}

func TestExecuteCrossDock_YaEjecutado_Rechaza(t *testing.T) {
	actualizado := false
	repo := &mocks.RepositoryMock{
		GetCrossDockLinkByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error) {
			return &entities.CrossDockLink{ID: id, Status: "executed"}, nil
		},
		UpdateCrossDockLinkFn: func(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
			actualizado = true
			return l, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ExecuteCrossDock(context.Background(), request.ExecuteCrossDockDTO{
		BusinessID: 10,
		LinkID:     1,
	})

	assert.ErrorIs(t, err, domainerrors.ErrCrossDockExecuted)
	assert.Nil(t, got)
	assert.False(t, actualizado, "un cross-dock no se ejecuta dos veces")
}

func TestExecuteCrossDock_Pendiente_SeEjecuta(t *testing.T) {
	var guardado *entities.CrossDockLink
	repo := &mocks.RepositoryMock{
		GetCrossDockLinkByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error) {
			return &entities.CrossDockLink{ID: id, Status: "pending"}, nil
		},
		UpdateCrossDockLinkFn: func(ctx context.Context, l *entities.CrossDockLink) (*entities.CrossDockLink, error) {
			guardado = l
			return l, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.ExecuteCrossDock(context.Background(), request.ExecuteCrossDockDTO{
		BusinessID: 10,
		LinkID:     1,
	})

	require.NoError(t, err)
	assert.Equal(t, "executed", guardado.Status)
	assert.NotNil(t, guardado.ExecutedAt)
}

func TestExecuteCrossDock_Inexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetCrossDockLinkByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CrossDockLink, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ExecuteCrossDock(context.Background(), request.ExecuteCrossDockDTO{BusinessID: 10, LinkID: 1})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestCreatePutawayRule_EstrategiaPorDefecto(t *testing.T) {
	var guardada *entities.PutawayRule
	repo := &mocks.RepositoryMock{
		CreatePutawayRuleFn: func(ctx context.Context, r *entities.PutawayRule) (*entities.PutawayRule, error) {
			guardada = r
			return r, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreatePutawayRule(context.Background(), request.CreatePutawayRuleDTO{
		BusinessID:   10,
		TargetZoneID: 3,
	})

	require.NoError(t, err)
	require.NotNil(t, guardada)
	assert.Equal(t, "nearest_empty", guardada.Strategy)
	assert.Equal(t, uint(3), guardada.TargetZoneID)
}

func TestCreatePutawayRule_EstrategiaExplicita(t *testing.T) {
	var guardada *entities.PutawayRule
	repo := &mocks.RepositoryMock{
		CreatePutawayRuleFn: func(ctx context.Context, r *entities.PutawayRule) (*entities.PutawayRule, error) {
			guardada = r
			return r, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreatePutawayRule(context.Background(), request.CreatePutawayRuleDTO{
		BusinessID: 10,
		Strategy:   "fixed_zone",
		Priority:   5,
		IsActive:   true,
	})

	require.NoError(t, err)
	assert.Equal(t, "fixed_zone", guardada.Strategy)
	assert.Equal(t, 5, guardada.Priority)
	assert.True(t, guardada.IsActive)
}

func TestListPutawayRules_NormalizaPaginacion(t *testing.T) {
	var recibido dtos.ListPutawayRulesParams
	repo := &mocks.RepositoryMock{
		ListPutawayRulesFn: func(ctx context.Context, params dtos.ListPutawayRulesParams) ([]entities.PutawayRule, int64, error) {
			recibido = params
			return nil, 0, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, _, err := uc.ListPutawayRules(context.Background(), dtos.ListPutawayRulesParams{Page: 0, PageSize: 800})

	require.NoError(t, err)
	assert.Equal(t, 1, recibido.Page)
	assert.Equal(t, 10, recibido.PageSize)
}

func TestUpdatePutawayRule_ActualizaCampos(t *testing.T) {
	zona := uint(9)
	prioridad := 2
	activa := false
	producto := "prod-1"

	var guardada *entities.PutawayRule
	repo := &mocks.RepositoryMock{
		GetPutawayRuleByIDFn: func(ctx context.Context, businessID, id uint) (*entities.PutawayRule, error) {
			return &entities.PutawayRule{ID: id, TargetZoneID: 3, Priority: 1, Strategy: "nearest_empty", IsActive: true}, nil
		},
		UpdatePutawayRuleFn: func(ctx context.Context, r *entities.PutawayRule) (*entities.PutawayRule, error) {
			guardada = r
			return r, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdatePutawayRule(context.Background(), request.UpdatePutawayRuleDTO{
		ID:           1,
		BusinessID:   10,
		ProductID:    &producto,
		TargetZoneID: &zona,
		Priority:     &prioridad,
		Strategy:     "fixed_zone",
		IsActive:     &activa,
	})

	require.NoError(t, err)
	assert.Equal(t, "prod-1", *guardada.ProductID)
	assert.Equal(t, uint(9), guardada.TargetZoneID)
	assert.Equal(t, 2, guardada.Priority)
	assert.Equal(t, "fixed_zone", guardada.Strategy)
	assert.False(t, guardada.IsActive)
}

func TestUpdatePutawayRule_SinCambios_ConservaLoExistente(t *testing.T) {
	var guardada *entities.PutawayRule
	repo := &mocks.RepositoryMock{
		GetPutawayRuleByIDFn: func(ctx context.Context, businessID, id uint) (*entities.PutawayRule, error) {
			return &entities.PutawayRule{ID: id, TargetZoneID: 3, Priority: 1, Strategy: "nearest_empty", IsActive: true}, nil
		},
		UpdatePutawayRuleFn: func(ctx context.Context, r *entities.PutawayRule) (*entities.PutawayRule, error) {
			guardada = r
			return r, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdatePutawayRule(context.Background(), request.UpdatePutawayRuleDTO{ID: 1, BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, uint(3), guardada.TargetZoneID)
	assert.Equal(t, 1, guardada.Priority)
	assert.Equal(t, "nearest_empty", guardada.Strategy)
	assert.True(t, guardada.IsActive)
}

func TestUpdatePutawayRule_Inexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetPutawayRuleByIDFn: func(ctx context.Context, businessID, id uint) (*entities.PutawayRule, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdatePutawayRule(context.Background(), request.UpdatePutawayRuleDTO{ID: 1, BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDeletePutawayRule_DelegaEnRepositorio(t *testing.T) {
	var borrada uint
	repo := &mocks.RepositoryMock{
		DeletePutawayRuleFn: func(ctx context.Context, businessID, ruleID uint) error {
			borrada = ruleID
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	require.NoError(t, uc.DeletePutawayRule(context.Background(), 10, 4))
	assert.Equal(t, uint(4), borrada)
}

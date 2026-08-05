package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/app/request"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/inventory/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCountPlan_ValoresPorDefecto(t *testing.T) {
	var guardado *entities.CycleCountPlan
	repo := &mocks.RepositoryMock{
		CreateCountPlanFn: func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
			guardado = p
			return p, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateCountPlan(context.Background(), request.CreateCountPlanDTO{
		BusinessID:  10,
		WarehouseID: 5,
		Name:        "Conteo mensual",
	})

	require.NoError(t, err)
	require.NotNil(t, guardado)
	assert.Equal(t, "abc", guardado.Strategy, "la estrategia por defecto es abc")
	assert.Equal(t, 30, guardado.FrequencyDays, "la frecuencia por defecto es 30 dias")
	assert.Equal(t, uint(10), guardado.BusinessID)
	assert.Equal(t, uint(5), guardado.WarehouseID)
}

func TestCreateCountPlan_FrecuenciaInvalidaUsaElDefault(t *testing.T) {
	for _, frecuencia := range []int{0, -5} {
		var guardado *entities.CycleCountPlan
		repo := &mocks.RepositoryMock{
			CreateCountPlanFn: func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
				guardado = p
				return p, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, err := uc.CreateCountPlan(context.Background(), request.CreateCountPlanDTO{
			BusinessID:    10,
			FrequencyDays: frecuencia,
		})

		require.NoError(t, err)
		assert.Equal(t, 30, guardado.FrequencyDays)
	}
}

func TestCreateCountPlan_ValoresExplicitosSeRespetan(t *testing.T) {
	proxima := time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	var guardado *entities.CycleCountPlan
	repo := &mocks.RepositoryMock{
		CreateCountPlanFn: func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
			guardado = p
			return p, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.CreateCountPlan(context.Background(), request.CreateCountPlanDTO{
		BusinessID:    10,
		Name:          "Conteo semanal",
		Strategy:      "location",
		FrequencyDays: 7,
		NextRunAt:     &proxima,
		IsActive:      true,
	})

	require.NoError(t, err)
	assert.Equal(t, "location", guardado.Strategy)
	assert.Equal(t, 7, guardado.FrequencyDays)
	assert.True(t, guardado.IsActive)
	assert.True(t, guardado.NextRunAt.Equal(proxima))
}

func TestListCountPlans_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 400, wantPage: 1, wantPageSize: 10},
		{page: 2, pageSize: 40, wantPage: 2, wantPageSize: 40},
	}

	for _, tc := range cases {
		var recibido dtos.ListCycleCountPlansParams
		repo := &mocks.RepositoryMock{
			ListCountPlansFn: func(ctx context.Context, params dtos.ListCycleCountPlansParams) ([]entities.CycleCountPlan, int64, error) {
				recibido = params
				return nil, 0, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, _, err := uc.ListCountPlans(context.Background(), dtos.ListCycleCountPlansParams{
			Page:     tc.page,
			PageSize: tc.pageSize,
		})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, recibido.Page)
		assert.Equal(t, tc.wantPageSize, recibido.PageSize)
	}
}

func TestGetCountPlan_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return &entities.CycleCountPlan{ID: id, BusinessID: businessID}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GetCountPlan(context.Background(), 10, 2)

	require.NoError(t, err)
	assert.Equal(t, uint(2), got.ID)
}

func TestUpdateCountPlan_ActualizaCamposOpcionales(t *testing.T) {
	bodega := uint(9)
	frecuencia := 15
	activo := false
	proxima := time.Date(2026, 10, 1, 0, 0, 0, 0, time.UTC)

	var guardado *entities.CycleCountPlan
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return &entities.CycleCountPlan{
				ID:            id,
				WarehouseID:   5,
				Name:          "Original",
				Strategy:      "abc",
				FrequencyDays: 30,
				IsActive:      true,
			}, nil
		},
		UpdateCountPlanFn: func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
			guardado = p
			return p, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateCountPlan(context.Background(), request.UpdateCountPlanDTO{
		ID:            1,
		BusinessID:    10,
		WarehouseID:   &bodega,
		Name:          "Actualizado",
		Strategy:      "location",
		FrequencyDays: &frecuencia,
		NextRunAt:     &proxima,
		IsActive:      &activo,
	})

	require.NoError(t, err)
	assert.Equal(t, uint(9), guardado.WarehouseID)
	assert.Equal(t, "Actualizado", guardado.Name)
	assert.Equal(t, "location", guardado.Strategy)
	assert.Equal(t, 15, guardado.FrequencyDays)
	assert.False(t, guardado.IsActive)
	assert.True(t, guardado.NextRunAt.Equal(proxima))
}

func TestUpdateCountPlan_CamposVacios_NoBorranLosExistentes(t *testing.T) {
	var guardado *entities.CycleCountPlan
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return &entities.CycleCountPlan{
				ID:            id,
				WarehouseID:   5,
				Name:          "Original",
				Strategy:      "abc",
				FrequencyDays: 30,
				IsActive:      true,
			}, nil
		},
		UpdateCountPlanFn: func(ctx context.Context, p *entities.CycleCountPlan) (*entities.CycleCountPlan, error) {
			guardado = p
			return p, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.UpdateCountPlan(context.Background(), request.UpdateCountPlanDTO{ID: 1, BusinessID: 10})

	require.NoError(t, err)
	assert.Equal(t, uint(5), guardado.WarehouseID)
	assert.Equal(t, "Original", guardado.Name)
	assert.Equal(t, "abc", guardado.Strategy)
	assert.Equal(t, 30, guardado.FrequencyDays)
	assert.True(t, guardado.IsActive)
}

func TestUpdateCountPlan_PlanInexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetCountPlanByIDFn: func(ctx context.Context, businessID, id uint) (*entities.CycleCountPlan, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.UpdateCountPlan(context.Background(), request.UpdateCountPlanDTO{ID: 1, BusinessID: 10})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestDeleteCountPlan_DelegaEnRepositorio(t *testing.T) {
	var borrado uint
	repo := &mocks.RepositoryMock{
		DeleteCountPlanFn: func(ctx context.Context, businessID, id uint) error {
			borrado = id
			return nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	require.NoError(t, uc.DeleteCountPlan(context.Background(), 10, 2))
	assert.Equal(t, uint(2), borrado)
}

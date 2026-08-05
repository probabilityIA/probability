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

func TestListDiscrepancias_NormalizaPaginacion(t *testing.T) {
	cases := []struct {
		page         int
		pageSize     int
		wantPage     int
		wantPageSize int
	}{
		{page: 0, pageSize: 10, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 0, wantPage: 1, wantPageSize: 10},
		{page: 1, pageSize: 999, wantPage: 1, wantPageSize: 10},
		{page: 4, pageSize: 30, wantPage: 4, wantPageSize: 30},
	}

	for _, tc := range cases {
		var recibido dtos.ListDiscrepanciesParams
		repo := &mocks.RepositoryMock{
			ListDiscrepanciesFn: func(ctx context.Context, params dtos.ListDiscrepanciesParams) ([]entities.InventoryDiscrepancy, int64, error) {
				recibido = params
				return nil, 0, nil
			},
		}
		uc := buildUseCase(repo, nil, nil)

		_, _, err := uc.ListDiscrepancies(context.Background(), dtos.ListDiscrepanciesParams{
			Page:     tc.page,
			PageSize: tc.pageSize,
		})

		require.NoError(t, err)
		assert.Equal(t, tc.wantPage, recibido.Page)
		assert.Equal(t, tc.wantPageSize, recibido.PageSize)
	}
}

func TestGetDiscrepancy_DelegaEnRepositorio(t *testing.T) {
	repo := &mocks.RepositoryMock{
		GetDiscrepancyByIDFn: func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
			return &entities.InventoryDiscrepancy{ID: id, BusinessID: businessID}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.GetDiscrepancy(context.Background(), 10, 3)

	require.NoError(t, err)
	assert.Equal(t, uint(3), got.ID)
	assert.Equal(t, uint(10), got.BusinessID)
}

func TestApproveDiscrepancy_UsaElTipoCountAdjustment(t *testing.T) {
	var codigoPedido string
	var params dtos.ApproveDiscrepancyTxParams
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			codigoPedido = code
			return 8, nil
		},
		ApproveDiscrepancyTxFn: func(ctx context.Context, p dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error) {
			params = p
			return &entities.InventoryDiscrepancy{ID: p.DiscrepancyID, Status: "approved"}, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ApproveDiscrepancy(context.Background(), request.ApproveDiscrepancyDTO{
		BusinessID:    10,
		DiscrepancyID: 3,
		ReviewerID:    7,
		Notes:         "conteo verificado",
	})

	require.NoError(t, err)
	assert.Equal(t, "count_adjustment", codigoPedido)
	assert.Equal(t, uint(8), params.MovementTypeID)
	assert.Equal(t, uint(10), params.BusinessID)
	assert.Equal(t, uint(3), params.DiscrepancyID)
	assert.Equal(t, uint(7), params.ReviewerID)
	assert.Equal(t, "conteo verificado", params.Notes)
	assert.Equal(t, "approved", got.Status)
}

func TestApproveDiscrepancy_SinTipoDeMovimiento_NoAjusta(t *testing.T) {
	dbErr := errors.New("catalogo vacio")
	ajustado := false
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 0, dbErr
		},
		ApproveDiscrepancyTxFn: func(ctx context.Context, p dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error) {
			ajustado = true
			return nil, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ApproveDiscrepancy(context.Background(), request.ApproveDiscrepancyDTO{BusinessID: 10, DiscrepancyID: 3})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
	assert.False(t, ajustado, "sin tipo de movimiento no se toca el inventario")
}

func TestApproveDiscrepancy_FallaTransaccion_PropagaError(t *testing.T) {
	dbErr := errors.New("tx failed")
	repo := &mocks.RepositoryMock{
		GetMovementTypeIDByCodeFn: func(ctx context.Context, code string) (uint, error) {
			return 8, nil
		},
		ApproveDiscrepancyTxFn: func(ctx context.Context, p dtos.ApproveDiscrepancyTxParams) (*entities.InventoryDiscrepancy, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.ApproveDiscrepancy(context.Background(), request.ApproveDiscrepancyDTO{BusinessID: 10, DiscrepancyID: 3})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

func TestRejectDiscrepancy_YaResuelta_Rechaza(t *testing.T) {
	for _, estado := range []string{"approved", "rejected"} {
		t.Run(estado, func(t *testing.T) {
			actualizada := false
			repo := &mocks.RepositoryMock{
				GetDiscrepancyByIDFn: func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
					return &entities.InventoryDiscrepancy{ID: id, Status: estado}, nil
				},
				UpdateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
					actualizada = true
					return d, nil
				},
			}
			uc := buildUseCase(repo, nil, nil)

			got, err := uc.RejectDiscrepancy(context.Background(), request.RejectDiscrepancyDTO{
				BusinessID:    10,
				DiscrepancyID: 3,
				ReviewerID:    7,
			})

			assert.ErrorIs(t, err, domainerrors.ErrDiscrepancyResolved)
			assert.Nil(t, got)
			assert.False(t, actualizada)
		})
	}
}

func TestRejectDiscrepancy_Abierta_SeRechaza(t *testing.T) {
	var guardada *entities.InventoryDiscrepancy
	repo := &mocks.RepositoryMock{
		GetDiscrepancyByIDFn: func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
			return &entities.InventoryDiscrepancy{ID: id, Status: "open"}, nil
		},
		UpdateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
			guardada = d
			return d, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.RejectDiscrepancy(context.Background(), request.RejectDiscrepancyDTO{
		BusinessID:    10,
		DiscrepancyID: 3,
		ReviewerID:    7,
		Reason:        "error de digitacion",
	})

	require.NoError(t, err)
	assert.Equal(t, "rejected", guardada.Status)
	require.NotNil(t, guardada.ReviewedByID)
	assert.Equal(t, uint(7), *guardada.ReviewedByID)
	assert.NotNil(t, guardada.ReviewedAt)
	assert.Equal(t, "error de digitacion", guardada.Notes)
}

func TestRejectDiscrepancy_SinMotivo_ConservaLasNotas(t *testing.T) {
	var guardada *entities.InventoryDiscrepancy
	repo := &mocks.RepositoryMock{
		GetDiscrepancyByIDFn: func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
			return &entities.InventoryDiscrepancy{ID: id, Status: "open", Notes: "nota original"}, nil
		},
		UpdateDiscrepancyFn: func(ctx context.Context, d *entities.InventoryDiscrepancy) (*entities.InventoryDiscrepancy, error) {
			guardada = d
			return d, nil
		},
	}
	uc := buildUseCase(repo, nil, nil)

	_, err := uc.RejectDiscrepancy(context.Background(), request.RejectDiscrepancyDTO{
		BusinessID:    10,
		DiscrepancyID: 3,
		ReviewerID:    7,
	})

	require.NoError(t, err)
	assert.Equal(t, "nota original", guardada.Notes)
}

func TestRejectDiscrepancy_Inexistente_PropagaError(t *testing.T) {
	dbErr := errors.New("not found")
	repo := &mocks.RepositoryMock{
		GetDiscrepancyByIDFn: func(ctx context.Context, businessID, id uint) (*entities.InventoryDiscrepancy, error) {
			return nil, dbErr
		},
	}
	uc := buildUseCase(repo, nil, nil)

	got, err := uc.RejectDiscrepancy(context.Background(), request.RejectDiscrepancyDTO{BusinessID: 10, DiscrepancyID: 3})

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, got)
}

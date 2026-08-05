package app

import (
	"context"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/pay/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetWalletKPISelection_Exitoso(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)

	repo.On("GetWalletKPISelection", ctx).Return(&entities.WalletKPISelection{
		ID:                  1,
		SelectedBusinessIDs: []uint{4, 8},
	}, nil)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetWalletKPISelection(ctx)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)
	assert.Equal(t, []uint{4, 8}, resp.SelectedBusinessIDs)
}

func TestGetWalletKPISelection_FallaRepositorio_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dbErr := errors.New("db down")

	repo.On("GetWalletKPISelection", ctx).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetWalletKPISelection(ctx)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

func TestGetFinancialStats_FechaInicioInvalida_RetornaError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetFinancialStats(context.Background(), &dtos.FinancialStatsDTO{
		StartDate: "01-01-2026",
		EndDate:   "2026-01-31",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid start_date")
	repo.AssertNotCalled(t, "GetFinancialStats", mock.Anything, mock.Anything)
}

func TestGetFinancialStats_FechaFinInvalida_RetornaError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetFinancialStats(context.Background(), &dtos.FinancialStatsDTO{
		StartDate: "2026-01-01",
		EndDate:   "no-es-fecha",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "invalid end_date")
	repo.AssertNotCalled(t, "GetFinancialStats", mock.Anything, mock.Anything)
}

func TestGetFinancialStats_RangoInvertido_RetornaError(t *testing.T) {
	repo := new(mocks.RepositoryMock)
	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetFinancialStats(context.Background(), &dtos.FinancialStatsDTO{
		StartDate: "2026-02-01",
		EndDate:   "2026-01-01",
	})

	require.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "end_date cannot be before start_date")
	repo.AssertNotCalled(t, "GetFinancialStats", mock.Anything, mock.Anything)
}

func TestGetFinancialStats_MismoDia_EsRangoValido(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dto := &dtos.FinancialStatsDTO{StartDate: "2026-01-15", EndDate: "2026-01-15"}

	repo.On("GetFinancialStats", ctx, dto).Return(&dtos.FinancialStatsResponse{TotalIncome: 0}, nil)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetFinancialStats(ctx, dto)

	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestGetFinancialStats_Exitoso_DelegaEnRepositorio(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dto := &dtos.FinancialStatsDTO{StartDate: "2026-01-01", EndDate: "2026-01-31"}

	repo.On("GetFinancialStats", ctx, dto).Return(&dtos.FinancialStatsResponse{
		TotalIncome: 1500000,
	}, nil)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetFinancialStats(ctx, dto)

	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.InDelta(t, 1500000.0, resp.TotalIncome, 0.001)
	repo.AssertExpectations(t)
}

func TestGetFinancialStats_FallaRepositorio_PropagaError(t *testing.T) {
	ctx := context.Background()
	repo := new(mocks.RepositoryMock)
	dto := &dtos.FinancialStatsDTO{StartDate: "2026-01-01", EndDate: "2026-01-31"}
	dbErr := errors.New("query timeout")

	repo.On("GetFinancialStats", ctx, dto).Return(nil, dbErr)

	uc := newWalletUseCaseForTest(repo)

	resp, err := uc.GetFinancialStats(ctx, dto)

	assert.ErrorIs(t, err, dbErr)
	assert.Nil(t, resp)
}

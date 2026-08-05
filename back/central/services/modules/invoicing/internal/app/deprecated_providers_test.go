package app

import (
	"context"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/mocks"
	"github.com/stretchr/testify/assert"
)

func useCaseParaDeprecados() *useCase {
	return &useCase{log: mocks.NewSilentLogger()}
}

func TestCreateProvider_EstaDeprecado(t *testing.T) {
	uc := useCaseParaDeprecados()

	got, err := uc.CreateProvider(context.Background(), &dtos.CreateProviderDTO{})

	assert.ErrorIs(t, err, ErrProviderManagementDeprecated)
	assert.Nil(t, got)
}

func TestUpdateProvider_EstaDeprecado(t *testing.T) {
	uc := useCaseParaDeprecados()

	err := uc.UpdateProvider(context.Background(), 1, &dtos.UpdateProviderDTO{})

	assert.ErrorIs(t, err, ErrProviderManagementDeprecated)
}

func TestGetProvider_EstaDeprecado(t *testing.T) {
	uc := useCaseParaDeprecados()

	got, err := uc.GetProvider(context.Background(), 1)

	assert.ErrorIs(t, err, ErrProviderManagementDeprecated)
	assert.Nil(t, got)
}

func TestListProviders_EstaDeprecado(t *testing.T) {
	uc := useCaseParaDeprecados()

	got, err := uc.ListProviders(context.Background(), 36)

	assert.ErrorIs(t, err, ErrProviderManagementDeprecated)
	assert.Nil(t, got)
}

func TestTestProviderConnection_EstaDeprecado(t *testing.T) {
	uc := useCaseParaDeprecados()

	err := uc.TestProviderConnection(context.Background(), 1)

	assert.ErrorIs(t, err, ErrProviderManagementDeprecated)
}

func TestListProviderTypes_EstaDeprecado(t *testing.T) {
	uc := useCaseParaDeprecados()

	got, err := uc.ListProviderTypes(context.Background())

	assert.ErrorIs(t, err, ErrProviderManagementDeprecated)
	assert.Nil(t, got)
}

func TestProveedoresDeprecados_MensajeApuntaAIntegrationsCore(t *testing.T) {
	assert.Contains(t, ErrProviderManagementDeprecated.Error(), "integrations/core")
}

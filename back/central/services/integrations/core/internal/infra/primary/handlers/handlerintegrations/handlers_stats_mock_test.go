package handlerintegrations

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core/internal/domain"
)

func (m *mockIntegrationUseCase) GetIntegrationStats(ctx context.Context, businessID uint) ([]domain.IntegrationStats, error) {
	return nil, nil
}

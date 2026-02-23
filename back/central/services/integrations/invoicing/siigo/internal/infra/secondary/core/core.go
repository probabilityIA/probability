package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
)

// SiigoCore adapta el use case de Siigo al contrato core.IIntegrationContract.
type SiigoCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *SiigoCore {
	return &SiigoCore{useCase: useCase}
}

func (s *SiigoCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return s.useCase.TestConnection(ctx, config, credentials)
}

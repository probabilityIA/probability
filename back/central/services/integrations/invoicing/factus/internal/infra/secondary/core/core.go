package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/ports"
)

// FactusCore adapta el use case de Factus al contrato core.IIntegrationContract.
type FactusCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *FactusCore {
	return &FactusCore{useCase: useCase}
}

func (f *FactusCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return f.useCase.TestConnection(ctx, config, credentials)
}

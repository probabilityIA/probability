package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/alegra/internal/domain/ports"
)

// AlegraCore adapta el use case de Alegra al contrato core.IIntegrationContract.
type AlegraCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *AlegraCore {
	return &AlegraCore{useCase: useCase}
}

func (a *AlegraCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return a.useCase.TestConnection(ctx, config, credentials)
}

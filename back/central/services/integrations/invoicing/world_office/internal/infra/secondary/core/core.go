package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/world_office/internal/domain/ports"
)

// WorldOfficeCore adapta el use case de World Office al contrato core.IIntegrationContract.
type WorldOfficeCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *WorldOfficeCore {
	return &WorldOfficeCore{useCase: useCase}
}

func (w *WorldOfficeCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return w.useCase.TestConnection(ctx, config, credentials)
}

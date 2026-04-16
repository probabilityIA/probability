package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/helisa/internal/domain/ports"
)

// HelisaCore adapta el use case de Helisa al contrato core.IIntegrationContract.
type HelisaCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *HelisaCore {
	return &HelisaCore{useCase: useCase}
}

func (h *HelisaCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return h.useCase.TestConnection(ctx, config, credentials)
}

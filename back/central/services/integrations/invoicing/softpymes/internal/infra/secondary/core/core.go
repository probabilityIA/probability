package core

import (
	"context"

	integrationcore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
)

// SoftpymesCore adapta el use case de Softpymes al contrato core.IIntegrationContract.
type SoftpymesCore struct {
	integrationcore.BaseIntegration
	useCase ports.IInvoiceUseCase
}

func New(useCase ports.IInvoiceUseCase) *SoftpymesCore {
	return &SoftpymesCore{useCase: useCase}
}

func (s *SoftpymesCore) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	return s.useCase.TestConnection(ctx, config, credentials)
}

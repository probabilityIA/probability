package dian

import (
	"context"
	"fmt"

	integrationsCore "github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/ports"
)

const factusTypeCode = "factus"

type IntegrationAdapter struct {
	core integrationsCore.IIntegrationCore
}

func NewIntegrationAdapter(core integrationsCore.IIntegrationCore) ports.IDianIntegration {
	return &IntegrationAdapter{core: core}
}

func (a *IntegrationAdapter) GetGlobalIntegrationID(ctx context.Context) (uint, error) {
	return a.core.GetGlobalIntegrationIDByTypeCode(ctx, factusTypeCode)
}

func (a *IntegrationAdapter) GetConfig(ctx context.Context, integrationID uint) (map[string]interface{}, error) {
	return a.core.GetIntegrationConfig(ctx, fmt.Sprintf("%d", integrationID))
}

func (a *IntegrationAdapter) Ensure(ctx context.Context, name string, config, credentials map[string]interface{}) (uint, error) {
	return a.core.EnsureGlobalIntegration(ctx, factusTypeCode, name, config, credentials)
}

func (a *IntegrationAdapter) TestConnection(ctx context.Context, config, credentials map[string]interface{}) error {
	contract, ok := a.core.GetRegisteredIntegration(integrationsCore.IntegrationTypeFactus)
	if !ok {
		return fmt.Errorf("proveedor factus no registrado")
	}
	return contract.TestConnection(ctx, config, credentials)
}

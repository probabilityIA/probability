package storefront

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/core"
)

type provider struct {
	core.BaseIntegration
}

// New creates a minimal integration provider for the Tienda (storefront with login).
// No external API, no consumers, no queues — just on/off control.
func New() core.IIntegrationContract {
	return &provider{}
}

func (p *provider) TestConnection(_ context.Context, _ map[string]interface{}, _ map[string]interface{}) error {
	return nil // No external API to test
}

package mocks

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

// OrderPublisherMock mock de domain.OrderPublisher para tests unitarios.
// Permite capturar las órdenes publicadas y simular errores de publicación.
type OrderPublisherMock struct {
	PublishFn func(ctx context.Context, order *canonical.ProbabilityOrderDTO) error
	// Published almacena las órdenes publicadas durante el test.
	Published []*canonical.ProbabilityOrderDTO
}

// Verificar en tiempo de compilación que implementa la interfaz.
var _ domain.OrderPublisher = (*OrderPublisherMock)(nil)

func (m *OrderPublisherMock) Publish(ctx context.Context, order *canonical.ProbabilityOrderDTO) error {
	if m.PublishFn != nil {
		return m.PublishFn(ctx, order)
	}
	m.Published = append(m.Published, order)
	return nil
}

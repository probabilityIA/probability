package usecases

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/test/internal/domain"
)

// GenerateAndPublishOrders genera múltiples órdenes aleatorias y las publica a RabbitMQ
func (uc *UseCases) GenerateAndPublishOrders(ctx context.Context, req *domain.GenerateOrderRequest) (*domain.GenerateOrderResponse, error) {
	response := &domain.GenerateOrderResponse{
		OrderIDs: make([]string, 0, req.Count),
	}

	for i := 0; i < req.Count; i++ {
		// Generar orden aleatoria
		order := uc.generator.GenerateRandomOrder(req)
		response.Generated++

		// Publicar a RabbitMQ
		if err := uc.publisher.PublishCanonicalOrder(ctx, order); err != nil {
			response.Failed++
			continue
		}

		response.Published++
		response.OrderIDs = append(response.OrderIDs, order.ExternalID)
	}

	if response.Failed > 0 {
		return response, fmt.Errorf("generated %d orders, published %d, failed %d", response.Generated, response.Published, response.Failed)
	}

	return response, nil
}

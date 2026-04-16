package handlers

import (
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/redis"
)

// Handlers contiene todos los handlers del módulo shipments
type Handlers struct {
	uc              *usecases.UseCases
	transportPub    domain.ITransportRequestPublisher // Async: quote, generate, track, cancel
	carrierResolver domain.ICarrierResolver           // Resolves active shipping carrier per business
	redisClient     redis.IRedis                      // Used for synchronous quote polling
}

// New crea una nueva instancia de Handlers
func New(uc *usecases.UseCases, transportPub domain.ITransportRequestPublisher, carrierResolver domain.ICarrierResolver, redisClient redis.IRedis) *Handlers {
	return &Handlers{
		uc:              uc,
		transportPub:    transportPub,
		carrierResolver: carrierResolver,
		redisClient:     redisClient,
	}
}

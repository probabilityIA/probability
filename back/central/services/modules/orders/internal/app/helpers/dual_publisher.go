package helpers

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

// PublishEventDual publica un evento tanto en Redis como RabbitMQ con snapshot completo
// Redis: Pub/Sub para tiempo real (no bloqueante, ignora errores)
// RabbitMQ: Garantía de entrega (goroutine para no bloquear respuesta HTTP)
// El parámetro order permite construir OrderSnapshot completo sin consultas adicionales
func PublishEventDual(
	ctx context.Context,
	event *entities.OrderEvent,
	order *entities.ProbabilityOrder,
	redisPublisher ports.IOrderEventPublisher,
	rabbitPublisher ports.IOrderRabbitPublisher,
	logger log.ILogger,
) {
	// Contexto background para ejecución asíncrona
	bgCtx := context.Background()

	// 1. Publicar a Redis (no bloqueante, ignora errores)
	if redisPublisher != nil {
		go func() {
			if err := redisPublisher.PublishOrderEvent(bgCtx, event, order); err != nil {
				logger.Warn(bgCtx).
					Err(err).
					Str("event_type", string(event.Type)).
					Str("order_id", event.OrderID).
					Msg("⚠️ Error al publicar evento a Redis (continuando)")
			}
		}()
	}

	// 2. Publicar a RabbitMQ (en goroutine para no bloquear respuesta HTTP)
	if rabbitPublisher != nil {
		go func() {
			if err := rabbitPublisher.PublishOrderEvent(bgCtx, event, order); err != nil {
				logger.Error(bgCtx).
					Err(err).
					Str("event_type", string(event.Type)).
					Str("order_id", event.OrderID).
					Msg("❌ Error al publicar evento a RabbitMQ")
			} else {
				logger.Info(bgCtx).
					Str("event_type", string(event.Type)).
					Str("order_id", event.OrderID).
					Msg("✅ Evento publicado en ambos canales (Redis + RabbitMQ)")
			}
		}()
	}
}

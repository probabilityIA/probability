package consumer

import (
	"context"
	"encoding/json"
	"fmt"

	domain "github.com/secamc93/probability/back/central/services/modules/ai_sales/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type OrderResultConsumer struct {
	queue                rabbitmq.IQueue
	responsePublisher    domain.IAIResponsePublisher
	sessionCache         domain.ISessionCache
	persistencePublisher domain.IAIPersistencePublisher
	log                  log.ILogger
}

func NewOrderResultConsumer(
	queue rabbitmq.IQueue,
	responsePublisher domain.IAIResponsePublisher,
	sessionCache domain.ISessionCache,
	persistencePublisher domain.IAIPersistencePublisher,
	logger log.ILogger,
) *OrderResultConsumer {
	return &OrderResultConsumer{
		queue:                queue,
		responsePublisher:    responsePublisher,
		sessionCache:         sessionCache,
		persistencePublisher: persistencePublisher,
		log:                  logger,
	}
}

func (c *OrderResultConsumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueAIOrderResult

	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
	}

	if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
		return fmt.Errorf("failed to start consumer %s: %w", queueName, err)
	}

	c.log.Info().Str("queue", queueName).Msg("AI Order Result consumer iniciado")
	return nil
}

func (c *OrderResultConsumer) handleMessage(messageBody []byte) error {
	ctx := context.Background()

	var dto domain.OrderResultDTO
	if err := json.Unmarshal(messageBody, &dto); err != nil {
		c.log.Error().Err(err).Msg("Error deserializando order result")
		return nil
	}

	if dto.PhoneNumber == "" || dto.BusinessID == 0 {
		c.log.Warn().Msg("Order result sin phone_number o business_id, descartando")
		return nil
	}

	var text string
	if dto.Success {
		text = fmt.Sprintf("Tu pedido %s fue creado exitosamente! En breve recibiras mas detalles.", dto.OrderNumber)
	} else {
		text = fmt.Sprintf("Hubo un problema procesando tu pedido: %s. Por favor intenta de nuevo.", dto.ErrorMessage)
	}

	if err := c.responsePublisher.PublishResponse(ctx, dto.PhoneNumber, dto.BusinessID, text); err != nil {
		c.log.Error().Err(err).
			Str("phone", dto.PhoneNumber).
			Msg("Error publicando resultado de orden a WhatsApp")
	}

	c.persistOrderResult(ctx, dto.PhoneNumber, text)

	return nil
}

func (c *OrderResultConsumer) persistOrderResult(ctx context.Context, phoneNumber, text string) {
	if c.persistencePublisher == nil || c.sessionCache == nil {
		return
	}

	session, err := c.sessionCache.Get(ctx, phoneNumber)
	if err != nil || session == nil {
		return
	}

	if err := c.persistencePublisher.PublishMessageLog(ctx, session.ID, phoneNumber, "outgoing", text); err != nil {
		c.log.Error().Err(err).Msg("Error persistiendo order result en historial")
	}
}

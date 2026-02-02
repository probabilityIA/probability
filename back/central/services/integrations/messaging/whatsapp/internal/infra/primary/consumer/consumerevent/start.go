package consumerevent

import "context"

// Start inicia el consumidor de eventos
func (c *consumer) Start(ctx context.Context) error {
	pubsub := c.redisClient.Client(ctx).Subscribe(ctx, c.channel)
	defer pubsub.Close()

	ch := pubsub.Channel()

	for {
		select {
		case <-ctx.Done():
			c.logger.Info().Msg("Stopping WhatsApp order event consumer")
			return ctx.Err()
		case msg := <-ch:
			if msg == nil {
				continue
			}
			c.handleOrderEvent(ctx, msg.Payload)
		}
	}
}

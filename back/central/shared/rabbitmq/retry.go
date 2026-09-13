package rabbitmq

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	MaxDeliveryAttempts = 5
	DeadLetterQueue     = "probability.dlq"

	headerRetryCount  = "x-retry-count"
	headerOriginQueue = "x-origin-queue"
	headerLastError   = "x-last-error"
	headerFirstFailed = "x-first-failed-at"
)

var retryTiers = []time.Duration{
	10 * time.Second,
	60 * time.Second,
	300 * time.Second,
}

func retryQueueName(d time.Duration) string {
	return fmt.Sprintf("probability.retry.%ds", int(d.Seconds()))
}

func nextRetryDelay(attempts int) time.Duration {
	if attempts < 1 {
		attempts = 1
	}
	idx := attempts - 1
	if idx >= len(retryTiers) {
		idx = len(retryTiers) - 1
	}
	return retryTiers[idx]
}

func retryCountFromHeaders(headers amqp.Table) int {
	if headers == nil {
		return 0
	}
	switch v := headers[headerRetryCount].(type) {
	case int32:
		return int(v)
	case int64:
		return int(v)
	case int:
		return v
	case float64:
		return int(v)
	}
	return 0
}

func (r *rabbitMQ) declareRetryInfrastructure(ch *amqp.Channel) error {
	if _, err := ch.QueueDeclare(DeadLetterQueue, true, false, false, false, nil); err != nil {
		return fmt.Errorf("declarando la DLQ %s: %w", DeadLetterQueue, err)
	}

	for _, tier := range retryTiers {
		name := retryQueueName(tier)

		if err := ch.ExchangeDeclare(name, "fanout", true, false, false, false, nil); err != nil {
			return fmt.Errorf("declarando el exchange de reintento %s: %w", name, err)
		}

		if _, err := ch.QueueDeclare(name, true, false, false, false, amqp.Table{
			"x-message-ttl":          int32(tier.Milliseconds()),
			"x-dead-letter-exchange": "",
		}); err != nil {
			return fmt.Errorf("declarando la cola de reintento %s: %w", name, err)
		}

		if err := ch.QueueBind(name, "", name, false, nil); err != nil {
			return fmt.Errorf("enlazando la cola de reintento %s: %w", name, err)
		}
	}

	return nil
}

func (r *rabbitMQ) republish(exchange, routingKey string, msg amqp.Delivery, headers amqp.Table) error {
	r.retryMu.Lock()
	defer r.retryMu.Unlock()

	if r.retryChannel == nil {
		return fmt.Errorf("canal de reintentos no inicializado")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return r.retryChannel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  contentTypeOr(msg.ContentType),
		Body:         msg.Body,
		Headers:      headers,
		DeliveryMode: amqp.Persistent,
	})
}

func contentTypeOr(actual string) string {
	if actual == "" {
		return "application/json"
	}
	return actual
}

func headersForRetry(queueName string, msg amqp.Delivery, handlerErr error, attempts int) amqp.Table {
	headers := amqp.Table{}
	for k, v := range msg.Headers {
		headers[k] = v
	}

	headers[headerRetryCount] = int32(attempts)
	headers[headerOriginQueue] = queueName
	headers[headerLastError] = truncar(handlerErr.Error(), 500)
	if _, yaEstaba := headers[headerFirstFailed]; !yaEstaba {
		headers[headerFirstFailed] = time.Now().UTC().Format(time.RFC3339)
	}

	return headers
}

func truncar(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max]
}

func (r *rabbitMQ) handleFailedMessage(queueName string, msg amqp.Delivery, handlerErr error) {
	attempts := retryCountFromHeaders(msg.Headers) + 1
	headers := headersForRetry(queueName, msg, handlerErr, attempts)

	if attempts >= MaxDeliveryAttempts {
		if err := r.republish("", DeadLetterQueue, msg, headers); err != nil {
			r.logger.Error().
				Err(err).
				Str("queue", queueName).
				Int("attempts", attempts).
				Msg("No se pudo mover el mensaje a la DLQ - se reencola para no perderlo")
			msg.Nack(false, true)
			return
		}

		r.logger.Error().
			Err(handlerErr).
			Str("queue", queueName).
			Str("dlq", DeadLetterQueue).
			Int("attempts", attempts).
			Msg("Mensaje agotado: movido a la DLQ y sacado de circulacion")
		msg.Ack(false)
		return
	}

	delay := nextRetryDelay(attempts)
	exchange := retryQueueName(delay)

	if err := r.republish(exchange, queueName, msg, headers); err != nil {
		r.logger.Error().
			Err(err).
			Str("queue", queueName).
			Int("attempts", attempts).
			Msg("No se pudo programar el reintento - se reencola")
		msg.Nack(false, true)
		return
	}

	r.logger.Warn().
		Err(handlerErr).
		Str("queue", queueName).
		Int("attempts", attempts).
		Int("max_attempts", MaxDeliveryAttempts).
		Dur("retry_in", delay).
		Msg("Fallo el procesamiento: reintento programado")
	msg.Ack(false)
}

package consumerwalletalert

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	whaErrors "github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerwalletalert/request"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const templateName = "reporte_saldo_billetera_v2"

func (c *consumer) Start(ctx context.Context) error {
	queueName := rabbitmq.QueueWalletBalanceAlertRequested
	if err := c.queue.DeclareQueue(queueName, true); err != nil {
		c.log.Error().Err(err).Str("queue", queueName).Msg("Error declaring queue")
		return err
	}

	go func() {
		if err := c.queue.Consume(ctx, queueName, c.handleMessage); err != nil {
			c.log.Error().Err(err).Msg("Error consuming wallet balance alert queue")
		}
	}()

	return nil
}

func (c *consumer) handleMessage(messageBody []byte) error {
	var event request.WalletBalanceAlertEvent
	if err := json.Unmarshal(messageBody, &event); err != nil {
		c.log.Warn().Err(err).Msg("Malformed wallet balance alert message - discarding (ACK)")
		return nil
	}

	if strings.TrimSpace(event.PhoneNumber) == "" {
		c.log.Warn().
			Uint("business_id", event.BusinessID).
			Msg("Wallet balance alert sin telefono de bodega principal - skipping")
		return nil
	}

	variables := map[string]string{
		"1": orDefault(event.BusinessName, "tu negocio"),
		"2": formatCOP(event.Balance),
	}

	messageID, err := c.useCase.SendPlatformTemplate(
		context.Background(),
		templateName,
		event.PhoneNumber,
		variables,
		"",
		event.BusinessID,
	)
	if err != nil {
		if whaErrors.IsNonRetryable(err) {
			c.log.Warn().Err(err).
				Uint("business_id", event.BusinessID).
				Str("phone_number", event.PhoneNumber).
				Msg("Wallet balance alert skipped - non-retryable error (ACK)")
			return nil
		}
		c.log.Error().Err(err).
			Uint("business_id", event.BusinessID).
			Str("phone_number", event.PhoneNumber).
			Msg("Error sending wallet balance alert - will be retried")
		return err
	}

	c.log.Info().
		Uint("business_id", event.BusinessID).
		Str("message_id", messageID).
		Msg("Wallet balance alert sent successfully")

	return nil
}

func orDefault(value, defaultValue string) string {
	if strings.TrimSpace(value) == "" {
		return defaultValue
	}
	return value
}

func formatCOP(amount float64) string {
	intVal := int64(amount)
	s := fmt.Sprintf("%d", intVal)
	formatted := ""
	for i, ch := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			formatted += "."
		}
		formatted += string(ch)
	}
	return "$" + formatted
}

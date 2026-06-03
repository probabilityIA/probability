package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/app/usecases"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
	"github.com/secamc93/probability/back/central/shared/redis"
)

const (
	QueueOrderCreatedForShipments = rabbitmq.QueueOrdersToShipments
	autoGenerateGuideConfigKey    = "auto_generate_guide_enabled"
)

type orderCreatedMessage struct {
	EventType     string `json:"event_type"`
	OrderID       string `json:"order_id"`
	BusinessID    *uint  `json:"business_id"`
	IntegrationID *uint  `json:"integration_id"`
}

type OrderCreatedConsumer struct {
	queue        rabbitmq.IQueue
	uc           *usecases.UseCases
	transportPub domain.ITransportRequestPublisher
	redisClient  redis.IRedis
	log          log.ILogger
}

func NewOrderCreatedConsumer(
	queue rabbitmq.IQueue,
	uc *usecases.UseCases,
	transportPub domain.ITransportRequestPublisher,
	redisClient redis.IRedis,
	logger log.ILogger,
) *OrderCreatedConsumer {
	return &OrderCreatedConsumer{
		queue:        queue,
		uc:           uc,
		transportPub: transportPub,
		redisClient:  redisClient,
		log:          logger.WithModule("shipments.order_created_consumer"),
	}
}

func (c *OrderCreatedConsumer) Start(ctx context.Context) error {
	if err := c.queue.DeclareExchange(rabbitmq.ExchangeOrderEvents, "fanout", true); err != nil {
		return fmt.Errorf("failed to declare order events exchange: %w", err)
	}
	if err := c.queue.DeclareQueue(QueueOrderCreatedForShipments, true); err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}
	if err := c.queue.BindQueue(QueueOrderCreatedForShipments, rabbitmq.ExchangeOrderEvents, ""); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	c.log.Info(ctx).Str("queue", QueueOrderCreatedForShipments).Msg("Starting order created consumer for shipments")

	if err := c.queue.Consume(ctx, QueueOrderCreatedForShipments, c.handle); err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}
	return nil
}

func (c *OrderCreatedConsumer) handle(message []byte) error {
	ctx := context.Background()

	var msg orderCreatedMessage
	if err := json.Unmarshal(message, &msg); err != nil {
		c.log.Error(ctx).Err(err).Msg("Failed to unmarshal order event")
		return nil
	}
	if msg.EventType != "order.created" || msg.OrderID == "" {
		return nil
	}

	c.processAssociation(ctx, &msg)
	return nil
}

func (c *OrderCreatedConsumer) processAssociation(ctx context.Context, msg *orderCreatedMessage) {
	repo := c.uc.Repo()

	existing, _ := repo.GetShipmentsByOrderID(ctx, msg.OrderID)
	for i := range existing {
		if existing[i].GuideURL != nil && *existing[i].GuideURL != "" {
			return
		}
	}

	sel, err := repo.GetOrderSelectedShipping(ctx, msg.OrderID)
	if err != nil || sel == nil {
		return
	}
	if !strings.HasPrefix(strings.ToLower(strings.TrimSpace(sel.Source)), "probability") {
		return
	}

	quoteID, rateIdx, ok := parseQuoteServiceCode(sel.Code)
	if !ok {
		c.log.Warn(ctx).Str("order_id", msg.OrderID).Str("code", sel.Code).Msg("Probability shipping line without parseable quote code")
		return
	}

	quote, err := repo.GetSavedQuoteByID(ctx, quoteID)
	if err != nil || quote == nil {
		c.log.Warn(ctx).Str("order_id", msg.OrderID).Uint("quote_id", quoteID).Msg("Saved quote not found for order")
		return
	}
	if quote.Status == domain.QuoteStatusGuideGenerated {
		return
	}

	carrier, idRate := rateFromQuote(quote, rateIdx)

	orderUUID := msg.OrderID
	quote.OrderUUID = &orderUUID
	quote.SelectedCarrier = carrier
	quote.SelectedServiceCode = sel.Code
	quote.SelectedIDRate = idRate
	quote.Status = domain.QuoteStatusAssociated

	if err := repo.UpdateSavedQuote(ctx, quote); err != nil {
		c.log.Error(ctx).Err(err).Str("order_id", msg.OrderID).Uint("quote_id", quoteID).Msg("Failed to associate saved quote to order")
		return
	}

	c.log.Info(ctx).
		Str("order_id", msg.OrderID).
		Uint("quote_id", quoteID).
		Str("carrier", carrier).
		Msg("Saved quote associated to order")

	c.maybeAutoGenerate(ctx, msg, quote, rateIdx)
}

func parseQuoteServiceCode(code string) (uint, int, bool) {
	parts := strings.Split(strings.TrimSpace(code), "-")
	if len(parts) != 3 || parts[0] != "pq" {
		return 0, 0, false
	}
	id, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return 0, 0, false
	}
	idx, err := strconv.Atoi(parts[2])
	if err != nil || idx < 0 {
		return 0, 0, false
	}
	return uint(id), idx, true
}

func rateFromQuote(quote *domain.SavedQuote, rateIdx int) (string, *int64) {
	rate := rateAt(quote, rateIdx)
	if rate == nil {
		return "", nil
	}
	carrier, _ := rate["carrier"].(string)
	return carrier, idRateFromMap(rate)
}

func rateAt(quote *domain.SavedQuote, rateIdx int) map[string]interface{} {
	if rateIdx < 0 || rateIdx >= len(quote.Rates) {
		return nil
	}
	return quote.Rates[rateIdx]
}

func idRateFromMap(rate map[string]interface{}) *int64 {
	switch v := rate["idRate"].(type) {
	case float64:
		n := int64(v)
		return &n
	case int64:
		return &v
	case int:
		n := int64(v)
		return &n
	}
	return nil
}

package router

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/core"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

const (
	QueueTransportRequests = rabbitmq.QueueTransportRequests

	QueueEnvioclickRequests = rabbitmq.QueueTransportEnvioclickRequests
	QueueEnviameRequests    = rabbitmq.QueueTransportEnviameRequests
	QueueTuRequests         = rabbitmq.QueueTransportTuRequests
	QueueMiPaqueteRequests  = rabbitmq.QueueTransportMiPaqueteRequests
	QueueShipitRequests     = rabbitmq.QueueTransportShipitRequests
)

type transportRequestHeader struct {
	IntegrationTypeID uint      `json:"integration_type_id"`
	Provider          string    `json:"provider"`
	Operation         string    `json:"operation"`
	CorrelationID     string    `json:"correlation_id"`
	Timestamp         time.Time `json:"timestamp"`
}

type Bundle struct {
	rabbit rabbitmq.IQueue
	log    log.ILogger
}

func New(
	logger log.ILogger,
	rabbit rabbitmq.IQueue,
) *Bundle {
	logger = logger.WithModule("transport.router")

	b := &Bundle{
		rabbit: rabbit,
		log:    logger,
	}

	if rabbit == nil {
		logger.Warn(context.Background()).
			Msg("RabbitMQ no disponible, router de transporte deshabilitado")
		return b
	}

	go func() {
		ctx := context.Background()
		logger.Info(ctx).Msg("Starting transport router in background...")
		if err := b.startRouter(ctx); err != nil {
			logger.Error(ctx).Err(err).Msg("Transport router failed to start or stopped with error")
		}
	}()

	logger.Info(context.Background()).Msg("Transport router initialized")

	return b
}

func (b *Bundle) startRouter(ctx context.Context) error {
	if b.rabbit == nil {
		return fmt.Errorf("rabbitmq client is nil")
	}

	if err := b.rabbit.DeclareQueue(QueueTransportRequests, true); err != nil {
		b.log.Error(ctx).Err(err).Str("queue", QueueTransportRequests).Msg("Failed to declare transport.requests queue")
		return err
	}

	carrierQueues := []string{
		QueueEnvioclickRequests,
		QueueEnviameRequests,
		QueueTuRequests,
		QueueMiPaqueteRequests,
		QueueShipitRequests,
	}
	for _, q := range carrierQueues {
		if err := b.rabbit.DeclareQueue(q, true); err != nil {
			b.log.Warn(ctx).Err(err).Str("queue", q).Msg("Failed to declare carrier queue")
		}
	}

	b.log.Info(ctx).
		Str("queue", QueueTransportRequests).
		Msg("Transport router listening")

	return b.rabbit.Consume(ctx, QueueTransportRequests, b.handleTransportRequest)
}

func (b *Bundle) handleTransportRequest(message []byte) error {
	ctx := context.Background()

	var header transportRequestHeader
	if err := json.Unmarshal(message, &header); err != nil {
		b.log.Error(ctx).
			Err(err).
			Str("body", string(message)).
			Msg("Failed to unmarshal transport request header")
		return err
	}

	b.log.Info(ctx).
		Uint("integration_type_id", header.IntegrationTypeID).
		Str("provider", header.Provider).
		Str("operation", header.Operation).
		Str("correlation_id", header.CorrelationID).
		Msg("Routing transport request")

	targetQueue := b.getCarrierQueue(header.IntegrationTypeID)
	if targetQueue == "" {
		b.log.Error(ctx).
			Uint("integration_type_id", header.IntegrationTypeID).
			Str("provider", header.Provider).
			Msg("Unknown carrier - cannot route transport request (message discarded)")
		return nil
	}

	if err := b.rabbit.Publish(ctx, targetQueue, message); err != nil {
		b.log.Error(ctx).
			Err(err).
			Str("target_queue", targetQueue).
			Msg("Failed to forward transport request to carrier queue")
		return err
	}

	b.log.Info(ctx).
		Str("provider", header.Provider).
		Str("target_queue", targetQueue).
		Msg("Transport request forwarded")

	return nil
}

func (b *Bundle) getCarrierQueue(integrationTypeID uint) string {
	switch int(integrationTypeID) {
	case core.IntegrationTypeEnvioClick:
		return QueueEnvioclickRequests
	case core.IntegrationTypeEnviame:
		return QueueEnviameRequests
	case core.IntegrationTypeTu:
		return QueueTuRequests
	case core.IntegrationTypeMiPaquete:
		return QueueMiPaqueteRequests
	case core.IntegrationTypeShipit:
		return QueueShipitRequests
	default:
		return ""
	}
}

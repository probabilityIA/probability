package consumer

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/mocks"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

func TestPublicaGuiaGeneradaParaFacturacion(t *testing.T) {
	var colaUsada string
	var payload []byte

	queue := &mocks.RabbitMQMock{
		PublishFn: func(_ context.Context, queueName string, message []byte) error {
			colaUsada = queueName
			payload = message
			return nil
		},
	}
	repo := &mocks.RepositoryMock{
		GetOrderIntegrationIDFn: func(context.Context, string) (uint, error) { return 221, nil },
	}

	c := &ResponseConsumer{queue: queue, repo: repo, log: &mocks.LoggerMock{}}
	c.publishGuideGeneratedForInvoicing(context.Background(), "orden-uuid", 46)

	if colaUsada != rabbitmq.QueueOrdersToInvoicing {
		t.Fatalf("se esperaba publicar en %q, se publico en %q", rabbitmq.QueueOrdersToInvoicing, colaUsada)
	}

	var evento map[string]interface{}
	if err := json.Unmarshal(payload, &evento); err != nil {
		t.Fatalf("el evento no es json valido: %v", err)
	}
	if evento["event_type"] != "order.guide_generated" {
		t.Errorf("event_type inesperado: %v", evento["event_type"])
	}
	if evento["order_id"] != "orden-uuid" {
		t.Errorf("order_id inesperado: %v", evento["order_id"])
	}
	if evento["business_id"] != float64(46) || evento["integration_id"] != float64(221) {
		t.Errorf("negocio o integracion mal: %v / %v", evento["business_id"], evento["integration_id"])
	}
	if evento["timestamp"] == nil {
		t.Error("falta el timestamp")
	}
}

func TestSinIntegracionNoPublicaEventoIncompleto(t *testing.T) {
	publico := false
	queue := &mocks.RabbitMQMock{
		PublishFn: func(context.Context, string, []byte) error { publico = true; return nil },
	}
	repo := &mocks.RepositoryMock{
		GetOrderIntegrationIDFn: func(context.Context, string) (uint, error) {
			return 0, errors.New("orden sin integracion")
		},
	}

	c := &ResponseConsumer{queue: queue, repo: repo, log: &mocks.LoggerMock{}}
	c.publishGuideGeneratedForInvoicing(context.Background(), "orden-uuid", 46)

	if publico {
		t.Fatal("sin integration_id no se debe publicar: facturacion no sabria que config aplicar")
	}
}

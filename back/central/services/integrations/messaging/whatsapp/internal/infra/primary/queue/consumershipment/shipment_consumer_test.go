package consumershipment

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumershipment/request"
)

func guideEvent() request.ShipmentGuideEvent {
	return request.ShipmentGuideEvent{
		CustomerName:   "Cristian camilo Herrera diaz",
		BusinessName:   "Viga ropa deportiva",
		OrderNumber:    "15789",
		TrackingNumber: "240061196214",
		Carrier:        "INTERRAPIDISIMO",
		CodTotal:       73367,
		CodCarrierFee:  5365,
	}
}

func TestGuiaNoDuplicaLaComisionDelCheckout(t *testing.T) {
	event := guideEvent()
	event.CodIncludesShipping = true

	template, vars := buildGuideVariables(event)
	if template != "guia_envio_generada_cod" {
		t.Fatalf("plantilla = %q, se esperaba guia_envio_generada_cod", template)
	}
	if vars["6"] != "$73.367" {
		t.Errorf("valor a recaudar = %q, se esperaba $73.367 (lo que declara la guia)", vars["6"])
	}
}

func TestGuiaDeOrdenManualSumaLaComision(t *testing.T) {
	event := guideEvent()
	event.CodTotal = 135929
	event.CodCarrierFee = 9451

	_, vars := buildGuideVariables(event)
	if vars["6"] != "$145.380" {
		t.Errorf("valor a recaudar = %q, se esperaba $145.380", vars["6"])
	}
}

func TestGuiaSinContraEntregaUsaLaPlantillaSinValor(t *testing.T) {
	event := guideEvent()
	event.CodTotal = 0

	template, vars := buildGuideVariables(event)
	if template != "guia_envio_generada" {
		t.Fatalf("plantilla = %q, se esperaba guia_envio_generada", template)
	}
	if err := entities.ValidateTemplateVariables(template, vars); err != nil {
		t.Errorf("validacion fallida: %v", err)
	}
}

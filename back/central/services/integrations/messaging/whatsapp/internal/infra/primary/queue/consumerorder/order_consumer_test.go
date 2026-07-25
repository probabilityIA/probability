package consumerorder

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/infra/primary/queue/consumerorder/request"
)

func vigaEvent() request.OrderConfirmationEvent {
	return request.OrderConfirmationEvent{
		OrderNumber:    "VIG-0046",
		CustomerName:   "Luisa Orozco",
		BusinessName:   "Viga ropa deportiva",
		TotalAmount:    54900,
		CodTotal:       66552.40,
		IsCod:          true,
		Currency:       "COP",
		TrackingNumber: "84151691723",
		Carrier:        "COORDINADORA",
		ShippingStreet: "Calle 56b #21a-33",
		ShippingCity:   "MEDELLIN",
		ShippingState:  "ANTIOQUIA",
		ItemsSummary:   "- Camiseta\n  Cant: 1",
	}
}

func TestValorARecaudarUsaCodTotalNoTotalAmount(t *testing.T) {
	tests := []struct {
		template string
		varKey   string
	}{
		{"confirmacion_pedido_contraentrega", "9"},
		{"pedido_en_reparto_cod", "6"},
		{"pedido_entregado_cod", "11"},
	}

	for _, tt := range tests {
		t.Run(tt.template, func(t *testing.T) {
			vars := buildVariables(tt.template, vigaEvent())

			got := vars[tt.varKey]
			if got != "$66.552" {
				t.Errorf("valor a recaudar = %q, se esperaba %q (cod_total, no total_amount)", got, "$66.552")
			}
			if got == "$54.900" {
				t.Error("se esta enviando total_amount en vez de cod_total")
			}
		})
	}
}

func TestFallbackATotalAmountSinCodTotal(t *testing.T) {
	event := vigaEvent()
	event.CodTotal = 0
	event.IsCod = false

	vars := buildVariables("pedido_en_reparto_cod", event)
	if vars["6"] != "$54.900" {
		t.Errorf("fallback = %q, se esperaba %q", vars["6"], "$54.900")
	}
}

func TestVariablesCoincidenConPlantillasAprobadas(t *testing.T) {
	templates := []string{
		"confirmacion_pedido_contraentrega",
		"confirmacion_pedido",
		"guia_envio_generada_cod",
		"pedido_en_reparto_cod",
		"pedido_entregado_cod",
	}

	for _, name := range templates {
		t.Run(name, func(t *testing.T) {
			def, ok := entities.GetTemplateDefinition(name)
			if !ok {
				t.Fatalf("plantilla %q no esta en el catalogo", name)
			}
			if name == "guia_envio_generada_cod" {
				return
			}
			vars := buildVariables(name, vigaEvent())
			if len(vars) != len(def.Variables) {
				t.Errorf("se construyeron %d variables, la plantilla declara %d",
					len(vars), len(def.Variables))
			}
			if err := entities.ValidateTemplateVariables(name, vars); err != nil {
				t.Errorf("validacion fallida: %v", err)
			}
		})
	}
}

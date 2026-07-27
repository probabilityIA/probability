package channel

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

func TestIsCODEventUsaFlagIsCod(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
		want bool
	}{
		{
			name: "orden manual de viga: efectivo + is_cod",
			data: map[string]any{"payment_method_id": uint(5), "is_cod": true},
			want: true,
		},
		{
			name: "payment_method_id 6 sin flag sigue siendo cod",
			data: map[string]any{"payment_method_id": uint(6)},
			want: true,
		},
		{
			name: "efectivo sin is_cod no es cod",
			data: map[string]any{"payment_method_id": uint(5), "is_cod": false},
			want: false,
		},
		{
			name: "float64 desde json",
			data: map[string]any{"payment_method_id": float64(5), "is_cod": true},
			want: true,
		},
		{
			name: "sin datos",
			data: map[string]any{},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entities.Event{Data: tt.data}.IsCOD()
			if got != tt.want {
				t.Errorf("IsCOD() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}

func TestEventCodeToTemplateNameCOD(t *testing.T) {
	tests := []struct {
		eventCode string
		isCOD     bool
		codFinal  bool
		want      string
	}{
		{"order.created", true, true, "confirmacion_pedido_contraentrega"},
		{"order.created", true, false, "confirmacion_pedido_contraentrega_sin_valor"},
		{"order.created", false, false, "confirmacion_pedido"},
		{"order.shipped", true, false, "pedido_en_reparto_cod"},
		{"order.shipped", false, false, "pedido_en_reparto"},
		{"order.delivered", true, false, "pedido_entregado_cod"},
		{"order.delivered", false, false, "pedido_entregado"},
	}

	for _, tt := range tests {
		got := eventCodeToTemplateName(tt.eventCode, tt.isCOD, tt.codFinal)
		if got != tt.want {
			t.Errorf("eventCodeToTemplateName(%q, %v, %v) = %q, se esperaba %q",
				tt.eventCode, tt.isCOD, tt.codFinal, got, tt.want)
		}
	}
}

func TestCodAmountIsFinal(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
		want bool
	}{
		{"total ya incluye el envio", map[string]any{"cod_includes_shipping": true}, true},
		{"el envio se suma despues", map[string]any{"cod_includes_shipping": false}, false},
		{"sin el campo", map[string]any{}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := (entities.Event{Data: tt.data}).CodAmountIsFinal(); got != tt.want {
				t.Errorf("CodAmountIsFinal() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}

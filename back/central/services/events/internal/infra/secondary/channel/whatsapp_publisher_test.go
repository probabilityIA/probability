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
		want      string
	}{
		{"order.created", true, "confirmacion_pedido_contraentrega"},
		{"order.created", false, "confirmacion_pedido"},
		{"order.shipped", true, "pedido_en_reparto_cod"},
		{"order.shipped", false, "pedido_en_reparto"},
		{"order.delivered", true, "pedido_entregado_cod"},
		{"order.delivered", false, "pedido_entregado"},
	}

	for _, tt := range tests {
		got := eventCodeToTemplateName(tt.eventCode, tt.isCOD)
		if got != tt.want {
			t.Errorf("eventCodeToTemplateName(%q, %v) = %q, se esperaba %q",
				tt.eventCode, tt.isCOD, got, tt.want)
		}
	}
}

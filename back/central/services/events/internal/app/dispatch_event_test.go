package app

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
)

func TestValidateConditionsOrderConfirmationRequiresCOD(t *testing.T) {
	confirmation := entities.CachedNotificationConfig{
		NotificationTypeID: dtos.NotificationTypeWhatsApp,
		EventCode:          dtos.EventCodeOrderCreated,
	}

	tests := []struct {
		name   string
		config entities.CachedNotificationConfig
		data   map[string]interface{}
		want   bool
	}{
		{
			name:   "confirmacion con is_cod true se envia",
			config: confirmation,
			data:   map[string]interface{}{"is_cod": true},
			want:   true,
		},
		{
			name:   "confirmacion con metodo de pago contra entrega se envia",
			config: confirmation,
			data:   map[string]interface{}{"payment_method_id": float64(entities.PaymentMethodIDCOD)},
			want:   true,
		},
		{
			name:   "confirmacion de orden pagada se descarta",
			config: confirmation,
			data:   map[string]interface{}{"is_cod": false, "payment_method_id": float64(1)},
			want:   false,
		},
		{
			name:   "confirmacion sin datos de pago se descarta",
			config: confirmation,
			data:   map[string]interface{}{},
			want:   false,
		},
		{
			name: "otros eventos de whatsapp no exigen contra entrega",
			config: entities.CachedNotificationConfig{
				NotificationTypeID: dtos.NotificationTypeWhatsApp,
				EventCode:          "order.shipped",
			},
			data: map[string]interface{}{"is_cod": false},
			want: true,
		},
		{
			name: "sse de orden creada no exige contra entrega",
			config: entities.CachedNotificationConfig{
				NotificationTypeID: dtos.NotificationTypeSSE,
				EventCode:          dtos.EventCodeOrderCreated,
			},
			data: map[string]interface{}{"is_cod": false},
			want: true,
		},
	}

	dispatcher := &EventDispatcher{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := entities.Event{Data: tt.data}
			if got := dispatcher.validateConditions(event, tt.config); got != tt.want {
				t.Errorf("validateConditions() = %v, se esperaba %v", got, tt.want)
			}
		})
	}
}

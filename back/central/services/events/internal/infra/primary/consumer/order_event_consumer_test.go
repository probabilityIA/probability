package consumer

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

type dispatcherEspia struct {
	recibido entities.Event
	llamado  bool
}

func (d *dispatcherEspia) HandleEvent(ctx context.Context, event entities.Event) error {
	d.recibido = event
	d.llamado = true
	return nil
}

func consumerDePrueba() (*OrderEventConsumer, *dispatcherEspia) {
	espia := &dispatcherEspia{}
	return &OrderEventConsumer{dispatcher: espia, logger: log.New()}, espia
}

func mensajeDeOrden(t *testing.T, extra map[string]any) []byte {
	t.Helper()
	orden := map[string]any{
		"id":              "abc-123",
		"order_number":    "DEMO-1",
		"customer_phone":  "573001234567",
		"shipping_street": "av calle 145 # 128 - 40",
		"shipping_city":   "Bogota",
		"is_cod":          true,
	}
	for k, v := range extra {
		orden[k] = v
	}
	businessID := uint(26)
	integrationID := uint(35)
	cuerpo, err := json.Marshal(map[string]any{
		"event_id":       "evt-1",
		"event_type":     "order.created",
		"business_id":    businessID,
		"integration_id": integrationID,
		"order":          orden,
	})
	require.NoError(t, err)
	return cuerpo
}

func TestLasCoordenadasDeEntregaLleganAlEvento(t *testing.T) {
	c, espia := consumerDePrueba()

	err := c.handleMessage(context.Background(), mensajeDeOrden(t, map[string]any{
		"shipping_lat": 4.7533418,
		"shipping_lng": -74.1091901,
	}))

	require.NoError(t, err)
	require.True(t, espia.llamado)
	assert.Equal(t, 4.7533418, espia.recibido.Data["shipping_lat"],
		"sin las coordenadas el mapa de la confirmacion no se puede generar")
	assert.Equal(t, -74.1091901, espia.recibido.Data["shipping_lng"])
}

func TestSinCoordenadasElEventoNoLasInventa(t *testing.T) {
	c, espia := consumerDePrueba()

	err := c.handleMessage(context.Background(), mensajeDeOrden(t, nil))

	require.NoError(t, err)
	require.True(t, espia.llamado)
	_, hayLat := espia.recibido.Data["shipping_lat"]
	_, hayLng := espia.recibido.Data["shipping_lng"]
	assert.False(t, hayLat)
	assert.False(t, hayLng)
}

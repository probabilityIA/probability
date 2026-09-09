package app

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/secamc93/probability/back/central/services/events/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/events/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

const canalWhatsApp = uint(dtos.NotificationTypeWhatsApp)

func eventoOrdenCreada() entities.Event {
	return entities.Event{
		ID:            "evt-1",
		Type:          dtos.EventCodeOrderCreated,
		Category:      "order",
		BusinessID:    26,
		IntegrationID: 7,
		Data:          map[string]interface{}{"is_cod": true},
	}
}

func configDe(id uint, canal uint, eventCode string) entities.CachedNotificationConfig {
	return entities.CachedNotificationConfig{
		ID:                 id,
		IntegrationID:      7,
		NotificationTypeID: canal,
		Enabled:            true,
		EventCode:          eventCode,
	}
}

func TestApplyVariantsSinVarianteDejaLasConfigsIgual(t *testing.T) {
	d, cache := dispatcherDePrueba()
	base := []entities.CachedNotificationConfig{configDe(1, canalWhatsApp, dtos.EventCodeOrderCreated)}
	cache.porTrigger["order.created_with_map"] = nil

	resultado := d.applyVariants(context.Background(), eventoOrdenCreada(), base)

	assert.Equal(t, base, resultado)
}

func TestApplyVariantsReemplazaLaConfigDelMismoCanal(t *testing.T) {
	d, cache := dispatcherDePrueba()
	base := []entities.CachedNotificationConfig{configDe(1, canalWhatsApp, dtos.EventCodeOrderCreated)}
	cache.porTrigger["order.created_with_map"] = []entities.CachedNotificationConfig{
		configDe(99, canalWhatsApp, "order.created_with_map"),
	}

	resultado := d.applyVariants(context.Background(), eventoOrdenCreada(), base)

	assert.Len(t, resultado, 1, "un solo mensaje por canal, nunca los dos")
	assert.Equal(t, uint(99), resultado[0].ID)
	assert.Equal(t, "order.created_with_map", resultado[0].EventCode)
}

func TestApplyVariantsNoTocaOtrosCanales(t *testing.T) {
	d, cache := dispatcherDePrueba()
	base := []entities.CachedNotificationConfig{
		configDe(1, canalWhatsApp, dtos.EventCodeOrderCreated),
		configDe(2, uint(dtos.NotificationTypeSSE), dtos.EventCodeOrderCreated),
	}
	cache.porTrigger["order.created_with_map"] = []entities.CachedNotificationConfig{
		configDe(99, canalWhatsApp, "order.created_with_map"),
	}

	resultado := d.applyVariants(context.Background(), eventoOrdenCreada(), base)

	ids := []uint{}
	for _, c := range resultado {
		ids = append(ids, c.ID)
	}
	assert.ElementsMatch(t, []uint{2, 99}, ids, "SSE sobrevive, WhatsApp se reemplaza")
}

func TestApplyVariantsIgnoraEventosSinVariante(t *testing.T) {
	d, _ := dispatcherDePrueba()
	evento := eventoOrdenCreada()
	evento.Type = "order.delivered"
	base := []entities.CachedNotificationConfig{configDe(1, canalWhatsApp, "order.delivered")}

	resultado := d.applyVariants(context.Background(), evento, base)

	assert.Equal(t, base, resultado)
}

type cacheDePrueba struct {
	porTrigger map[string][]entities.CachedNotificationConfig
}

func (c *cacheDePrueba) GetActiveConfigsByIntegrationAndTrigger(ctx context.Context, integrationID uint, trigger string) ([]entities.CachedNotificationConfig, error) {
	return c.porTrigger[trigger], nil
}

func dispatcherDePrueba() (*EventDispatcher, *cacheDePrueba) {
	cache := &cacheDePrueba{porTrigger: map[string][]entities.CachedNotificationConfig{}}
	return &EventDispatcher{
		configCache: cache,
		logger:      log.New(),
	}, cache
}

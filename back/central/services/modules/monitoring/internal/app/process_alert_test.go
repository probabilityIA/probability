package app

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/mocks"
)

func nuevoCaso(pub *mocks.AlertPublisherMock) ports.IUseCase {
	return New(pub, mocks.NewSilentLogger(), nil)
}

func TestProcessGrafanaAlert_TraduceElAlertnameAUnTipoLegible(t *testing.T) {
	casos := []struct {
		alertname string
		esperado  string
	}{
		{"RAM_ALTA", "RAM"},
		{"CPU_ALTO", "CPU"},
		{"DISCO_ALTO", "Disco"},
	}

	for _, caso := range casos {
		t.Run(caso.alertname, func(t *testing.T) {
			pub := &mocks.AlertPublisherMock{}
			uc := nuevoCaso(pub)

			err := uc.ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
				Status: "firing",
				Alerts: []dtos.GrafanaAlertDTO{{
					Status:      "firing",
					Labels:      map[string]string{"alertname": caso.alertname},
					Annotations: map[string]string{"summary": "uso al 95%"},
				}},
			})

			require.NoError(t, err)
			require.Len(t, pub.Published, 1)
			assert.Equal(t, caso.esperado, pub.Published[0].AlertType)
		})
	}
}

func TestProcessGrafanaAlert_UnAlertnameDesconocidoViajaTalCual(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{{
			Status: "firing",
			Labels: map[string]string{"alertname": "LATENCIA_RDS"},
		}},
	})

	require.NoError(t, err)
	require.Len(t, pub.Published, 1)
	assert.Equal(t, "LATENCIA_RDS", pub.Published[0].AlertType,
		"una alerta nueva de Grafana no se pierde por no estar en el mapa: se publica con su nombre crudo")
}

func TestProcessGrafanaAlert_SinAlertnameNiSummaryElEventoQuedaVacioPeroSePublica(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{{Status: "firing"}},
	})

	require.NoError(t, err)
	require.Len(t, pub.Published, 1)
	assert.Empty(t, pub.Published[0].AlertType,
		"sin alertname el tipo queda vacio; el consumidor de la cola debe tolerarlo")
	assert.Empty(t, pub.Published[0].Summary)
}

func TestProcessGrafanaAlert_SiNoHaySummaryUsaElValueString(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{{
			Status:      "firing",
			Labels:      map[string]string{"alertname": "RAM_ALTA"},
			ValueString: "value=93.4",
		}},
	})

	require.NoError(t, err)
	require.Equal(t, "value=93.4", pub.Published[0].Summary,
		"el valueString es el respaldo cuando la regla de Grafana no define annotation summary")
}

func TestProcessGrafanaAlert_ElSummaryTienePrioridadSobreElValueString(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{{
			Status:      "firing",
			Labels:      map[string]string{"alertname": "RAM_ALTA"},
			Annotations: map[string]string{"summary": "RAM al 95% en back-central"},
			ValueString: "value=95",
		}},
	})

	require.NoError(t, err)
	assert.Equal(t, "RAM al 95% en back-central", pub.Published[0].Summary)
}

func TestProcessGrafanaAlert_SoloPublicaLasAlertasEnFiring(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Status: "resolved",
		Alerts: []dtos.GrafanaAlertDTO{
			{Status: "resolved", Labels: map[string]string{"alertname": "RAM_ALTA"}},
			{Status: "firing", Labels: map[string]string{"alertname": "CPU_ALTO"}},
			{Status: "", Labels: map[string]string{"alertname": "DISCO_ALTO"}},
		},
	})

	require.NoError(t, err)
	require.Len(t, pub.Published, 1,
		"resolved y estados vacios se descartan: solo firing genera notificacion")
	assert.Equal(t, "CPU", pub.Published[0].AlertType)
}

func TestProcessGrafanaAlert_UnWebhookConVariasAlertasPublicaUnaPorCadaUna(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{
			{Status: "firing", Labels: map[string]string{"alertname": "RAM_ALTA"}},
			{Status: "firing", Labels: map[string]string{"alertname": "CPU_ALTO"}},
			{Status: "firing", Labels: map[string]string{"alertname": "DISCO_ALTO"}},
		},
	})

	require.NoError(t, err)
	assert.Len(t, pub.Published, 3)
}

func TestProcessGrafanaAlert_UnWebhookSinAlertasNoPublicaNada(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{Status: "firing"})

	require.NoError(t, err)
	assert.Empty(t, pub.Published)
}

func TestProcessGrafanaAlert_ConservaElStartsAtComoMomentoDelDisparo(t *testing.T) {
	disparo := time.Date(2026, 8, 3, 14, 30, 0, 0, time.UTC)
	pub := &mocks.AlertPublisherMock{}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{{
			Status:   "firing",
			Labels:   map[string]string{"alertname": "RAM_ALTA"},
			StartsAt: disparo,
		}},
	})

	require.NoError(t, err)
	assert.Equal(t, disparo, pub.Published[0].FiredAt)
	assert.Equal(t, "firing", pub.Published[0].Status)
}

func TestProcessGrafanaAlert_SiFallaLaPublicacionNoDevuelveErrorParaQueGrafanaNoReintente(t *testing.T) {
	pub := &mocks.AlertPublisherMock{
		PublishFn: func(ctx context.Context, event entities.AlertEvent) error {
			return errors.New("rabbitmq caido")
		},
	}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{{Status: "firing", Labels: map[string]string{"alertname": "RAM_ALTA"}}},
	})

	assert.NoError(t, err,
		"se traga el error a proposito: devolverlo haria que Grafana reintente el webhook en loop. La alerta se pierde sin DLQ")
}

func TestProcessGrafanaAlert_UnFalloDePublicacionNoAbortaLasAlertasSiguientes(t *testing.T) {
	pub := &mocks.AlertPublisherMock{}
	pub.PublishFn = func(ctx context.Context, event entities.AlertEvent) error {
		if event.AlertType == "RAM" {
			return errors.New("rabbitmq caido")
		}
		return nil
	}

	err := nuevoCaso(pub).ProcessGrafanaAlert(context.Background(), dtos.GrafanaWebhookDTO{
		Alerts: []dtos.GrafanaAlertDTO{
			{Status: "firing", Labels: map[string]string{"alertname": "RAM_ALTA"}},
			{Status: "firing", Labels: map[string]string{"alertname": "CPU_ALTO"}},
		},
	})

	require.NoError(t, err)
	require.Len(t, pub.Published, 2,
		"el fallo de la primera no corta el bucle: la segunda alerta si llega a la cola")
}

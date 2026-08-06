package handlers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/secamc93/probability/back/central/services/modules/monitoring/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/testkit"
)

type useCaseMock struct {
	ProcessFn func(ctx context.Context, dto dtos.GrafanaWebhookDTO) error

	Recibidos []dtos.GrafanaWebhookDTO
}

func (m *useCaseMock) ProcessGrafanaAlert(ctx context.Context, dto dtos.GrafanaWebhookDTO) error {
	m.Recibidos = append(m.Recibidos, dto)
	if m.ProcessFn != nil {
		return m.ProcessFn(ctx, dto)
	}
	return nil
}

const payloadValido = `{"status":"firing","title":"RAM alta","alerts":[
	{"status":"firing","labels":{"alertname":"RAM_ALTA"},"annotations":{"summary":"95%"},
	 "startsAt":"2026-08-05T12:00:00Z","valueString":"value=95"}]}`

func llamar(t *testing.T, uc *useCaseMock, secreto, auth, cuerpo string) (int, map[string]any) {
	t.Helper()
	h := New(uc, testkit.NewSilentLogger(), testkit.NewConfig("GRAFANA_WEBHOOK_SECRET", secreto))

	c, rec := testkit.Peticion(t, http.MethodPost, "/webhooks/grafana", []byte(cuerpo))
	if auth != "" {
		c.Request.Header.Set("Authorization", auth)
	}

	h.WebhookGrafana(c)
	return rec.Code, testkit.CuerpoJSON(t, rec)
}

func TestWebhookGrafana_ConElTokenCorrectoProcesaLaAlerta(t *testing.T) {
	uc := &useCaseMock{}

	code, body := llamar(t, uc, "s3cr3t", "Bearer s3cr3t", payloadValido)

	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "received", body["status"])
	require.Len(t, uc.Recibidos, 1)
	assert.Equal(t, "firing", uc.Recibidos[0].Status)
	assert.Equal(t, "RAM alta", uc.Recibidos[0].Title)
	require.Len(t, uc.Recibidos[0].Alerts, 1)
	assert.Equal(t, "RAM_ALTA", uc.Recibidos[0].Alerts[0].Labels["alertname"])
	assert.Equal(t, "value=95", uc.Recibidos[0].Alerts[0].ValueString)
}

func TestWebhookGrafana_RechazaUnTokenQueNoCoincide(t *testing.T) {
	casos := map[string]string{
		"token equivocado":   "Bearer otro",
		"sin encabezado":     "",
		"sin prefijo":        "s3cr3t",
		"prefijo incorrecto": "Token s3cr3t",
	}

	for nombre, auth := range casos {
		t.Run(nombre, func(t *testing.T) {
			uc := &useCaseMock{}

			code, _ := llamar(t, uc, "s3cr3t", auth, payloadValido)

			assert.Equal(t, http.StatusUnauthorized, code)
			assert.Empty(t, uc.Recibidos,
				"sin token valido la alerta ni siquiera llega al caso de uso")
		})
	}
}

func TestWebhookGrafana_ToleraDiferenciasDeMayusculasYEspaciosEnElToken(t *testing.T) {
	uc := &useCaseMock{}

	code, _ := llamar(t, uc, "s3cr3t", "  bearer s3cr3t  ", payloadValido)

	assert.Equal(t, http.StatusOK, code,
		"la comparacion es case-insensitive y recorta espacios: evita falsos negativos por como Grafana arma el header")
	assert.Len(t, uc.Recibidos, 1)
}

func TestWebhookGrafana_SinSecretoConfiguradoAceptaCualquierPeticion(t *testing.T) {
	uc := &useCaseMock{}

	code, _ := llamar(t, uc, "", "", payloadValido)

	assert.Equal(t, http.StatusOK, code)
	assert.Len(t, uc.Recibidos, 1,
		"falla ABIERTO: con GRAFANA_WEBHOOK_SECRET vacio no valida nada. Se documenta como modo dev, pero en produccion la variable tiene que estar puesta")
}

func TestWebhookGrafana_RechazaUnCuerpoQueNoEsJSON(t *testing.T) {
	uc := &useCaseMock{}

	code, body := llamar(t, uc, "s3cr3t", "Bearer s3cr3t", `{no es json`)

	assert.Equal(t, http.StatusBadRequest, code)
	assert.Equal(t, "invalid payload", body["error"])
	assert.Empty(t, uc.Recibidos)
}

func TestWebhookGrafana_RespondeOKAunqueElCasoDeUsoFalle(t *testing.T) {
	uc := &useCaseMock{
		ProcessFn: func(ctx context.Context, dto dtos.GrafanaWebhookDTO) error {
			return errors.New("rabbitmq caido")
		},
	}

	code, body := llamar(t, uc, "s3cr3t", "Bearer s3cr3t", payloadValido)

	assert.Equal(t, http.StatusOK, code,
		"siempre 200 para que Grafana no reintente en loop; el fallo solo queda en el log y la alerta se pierde")
	assert.Equal(t, "received", body["status"])
}

func TestWebhookGrafana_UnPayloadSinAlertasNoRevienta(t *testing.T) {
	uc := &useCaseMock{}

	code, _ := llamar(t, uc, "", "", `{"status":"ok","title":"nada"}`)

	assert.Equal(t, http.StatusOK, code)
	require.Len(t, uc.Recibidos, 1)
	assert.Empty(t, uc.Recibidos[0].Alerts)
}

func TestWebhookGrafana_TrasladaTodasLasAlertasDelPayload(t *testing.T) {
	uc := &useCaseMock{}

	code, _ := llamar(t, uc, "", "", `{"status":"firing","alerts":[
		{"status":"firing","labels":{"alertname":"RAM_ALTA"}},
		{"status":"resolved","labels":{"alertname":"CPU_ALTO"}}]}`)

	assert.Equal(t, http.StatusOK, code)
	require.Len(t, uc.Recibidos[0].Alerts, 2,
		"el handler no filtra por estado: ese criterio vive en el caso de uso")
}

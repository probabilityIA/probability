package app

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	bcErrors "github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/errors"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/domain/ports"
	"github.com/secamc93/probability/back/central/services/integrations/pay/bancolombia/internal/mocks"
)

const secretoDePrueba = "secreto"

func hmacHex(body []byte, secreto string) string {
	mac := hmac.New(sha256.New, []byte(secreto))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func hmacB64(body []byte, secreto string) string {
	mac := hmac.New(sha256.New, []byte(secreto))
	mac.Write(body)
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func hmacSobreBodyB64Hex(body []byte, secreto string) string {
	mac := hmac.New(sha256.New, []byte(secreto))
	mac.Write([]byte(base64.StdEncoding.EncodeToString(body)))
	return hex.EncodeToString(mac.Sum(nil))
}

func hmacSobreBodyB64B64(body []byte, secreto string) string {
	mac := hmac.New(sha256.New, []byte(secreto))
	mac.Write([]byte(base64.StdEncoding.EncodeToString(body)))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func casoWebhook(repo *mocks.IntegrationRepositoryMock, pub *mocks.WebhookPublisherMock) ports.IWebhookUseCase {
	if repo == nil {
		repo = &mocks.IntegrationRepositoryMock{}
	}
	if pub == nil {
		pub = &mocks.WebhookPublisherMock{}
	}
	return NewWebhookUseCase(repo, pub, mocks.NewSilentLogger())
}

func sinSecreto() *mocks.IntegrationRepositoryMock {
	return &mocks.IntegrationRepositoryMock{
		GetConfigForModeFn: func(ctx context.Context, testMode bool) (*ports.BancolombiaConfig, error) {
			return &ports.BancolombiaConfig{Environment: "sandbox", WebhookSecret: ""}, nil
		},
	}
}

func TestHandleIncomingWebhook_AceptaLasCuatroCodificacionesDeFirma(t *testing.T) {
	body := []byte(`{"eventId":"evt-1","transferState":"APPROVED"}`)
	casos := map[string]string{
		"hmac hex sobre el body":     hmacHex(body, secretoDePrueba),
		"hmac base64 sobre el body":  hmacB64(body, secretoDePrueba),
		"hmac hex sobre el body b64": hmacSobreBodyB64Hex(body, secretoDePrueba),
		"hmac b64 sobre el body b64": hmacSobreBodyB64B64(body, secretoDePrueba),
		"con prefijo sha256=":        "sha256=" + hmacHex(body, secretoDePrueba),
		"con espacios":               "  " + hmacHex(body, secretoDePrueba) + " ",
	}

	for nombre, firma := range casos {
		t.Run(nombre, func(t *testing.T) {
			pub := &mocks.WebhookPublisherMock{}

			err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), firma, body, false)

			require.NoError(t, err,
				"Bancolombia no dejo cerrado el formato de la firma, por eso se aceptan las cuatro variantes en vez de perder notificaciones")
			require.Len(t, pub.Publicados, 1)
			assert.True(t, pub.Publicados[0].SignatureValid)
		})
	}
}

func TestHandleIncomingWebhook_RechazaUnaFirmaConSecretoEquivocado(t *testing.T) {
	body := []byte(`{"eventId":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(nil, pub).HandleIncomingWebhook(context.Background(), hmacHex(body, "otro"), body, false)

	assert.ErrorIs(t, err, bcErrors.ErrInvalidSignature)
	assert.Empty(t, pub.Publicados)
}

func TestHandleIncomingWebhook_ConSecretoConfiguradoUnaFirmaVaciaSeRechaza(t *testing.T) {
	body := []byte(`{"eventId":"evt-1"}`)

	err := casoWebhook(nil, nil).HandleIncomingWebhook(context.Background(), "   ", body, false)

	assert.ErrorIs(t, err, bcErrors.ErrInvalidSignature)
}

func TestHandleIncomingWebhook_SinSecretoConfiguradoAceptaCualquierPayloadSinFirma(t *testing.T) {
	body := []byte(`{"eventId":"evt-1","transferState":"APPROVED","transferAmount":999999}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(sinSecreto(), pub).HandleIncomingWebhook(context.Background(), "", body, false)

	require.NoError(t, err)
	require.Len(t, pub.Publicados, 1)
	assert.False(t, pub.Publicados[0].SignatureValid,
		"falla ABIERTO: si no hay webhook_secret cargado se publica igual, con signature_valid=false. El consumidor DEBE mirar ese flag antes de dar un pago por recibido")
}

func TestHandleIncomingWebhook_SinSecretoIgnoraLaFirmaQueLlegue(t *testing.T) {
	body := []byte(`{"eventId":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(sinSecreto(), pub).HandleIncomingWebhook(context.Background(), "firma-inventada", body, false)

	require.NoError(t, err)
	assert.False(t, pub.Publicados[0].SignatureValid)
}

func TestHandleIncomingWebhook_RechazaUnCuerpoVacioAntesDeConsultarLaConfig(t *testing.T) {
	repo := &mocks.IntegrationRepositoryMock{}

	err := casoWebhook(repo, nil).HandleIncomingWebhook(context.Background(), "firma", nil, false)

	require.Error(t, err)
	assert.Empty(t, repo.ModosConsultados)
}

func TestHandleIncomingWebhook_UsaLasCredencialesDelModoQueLlega(t *testing.T) {
	body := []byte(`{"eventId":"evt-1"}`)

	for _, esPrueba := range []bool{true, false} {
		repo := &mocks.IntegrationRepositoryMock{}
		pub := &mocks.WebhookPublisherMock{}

		err := casoWebhook(repo, pub).HandleIncomingWebhook(context.Background(), hmacHex(body, secretoDePrueba), body, esPrueba)

		require.NoError(t, err)
		assert.Equal(t, []bool{esPrueba}, repo.ModosConsultados)
		assert.Equal(t, esPrueba, pub.Publicados[0].IsTest)
	}
}

func TestHandleIncomingWebhook_SiNoSePuedeLeerLaConfigNoPublicaNada(t *testing.T) {
	fallo := errors.New("timeout")
	repo := &mocks.IntegrationRepositoryMock{
		GetConfigForModeFn: func(ctx context.Context, testMode bool) (*ports.BancolombiaConfig, error) {
			return nil, fallo
		},
	}
	pub := &mocks.WebhookPublisherMock{}

	err := casoWebhook(repo, pub).HandleIncomingWebhook(context.Background(), "f", []byte(`{}`), false)

	assert.ErrorIs(t, err, fallo)
	assert.Empty(t, pub.Publicados)
}

func TestHandleIncomingWebhook_RechazaUnJSONInvalido(t *testing.T) {
	body := []byte(`{no es json`)

	err := casoWebhook(sinSecreto(), nil).HandleIncomingWebhook(context.Background(), "", body, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "decode bancolombia webhook envelope")
}

func TestHandleIncomingWebhook_SiFallaLaPublicacionDevuelveError(t *testing.T) {
	pub := &mocks.WebhookPublisherMock{
		PublishWebhookEventFn: func(ctx context.Context, msg *ports.BancolombiaWebhookMessage) error {
			return errors.New("rabbitmq caido")
		},
	}

	err := casoWebhook(sinSecreto(), pub).HandleIncomingWebhook(context.Background(), "", []byte(`{"eventId":"e"}`), false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "publish bancolombia webhook event")
}

func TestHandleIncomingWebhook_ConservaElPayloadCrudo(t *testing.T) {
	body := []byte(`{"eventId":"evt-1"}`)
	pub := &mocks.WebhookPublisherMock{}

	require.NoError(t, casoWebhook(sinSecreto(), pub).HandleIncomingWebhook(context.Background(), "", body, false))

	assert.Equal(t, body, pub.Publicados[0].RawPayload,
		"el crudo permite reconstruir la conciliacion si el parseo aplano un campo mal")
}

func publicar(t *testing.T, crudo string) *ports.BancolombiaWebhookMessage {
	t.Helper()
	pub := &mocks.WebhookPublisherMock{}
	require.NoError(t, casoWebhook(sinSecreto(), pub).HandleIncomingWebhook(context.Background(), "", []byte(crudo), false))
	require.Len(t, pub.Publicados, 1)
	return &pub.Publicados[0]
}

func TestParseo_AplanaLosObjetosAnidadosConocidos(t *testing.T) {
	casos := map[string]string{
		"en la raiz":            `{"eventId":"evt-1"}`,
		"dentro de header":      `{"header":{"eventId":"evt-1"}}`,
		"dentro de data":        `{"data":{"eventId":"evt-1"}}`,
		"dentro de transaction": `{"transaction":{"eventId":"evt-1"}}`,
		"dentro de transfer":    `{"transfer":{"eventId":"evt-1"}}`,
		"dentro de payload":     `{"payload":{"eventId":"evt-1"}}`,
	}

	for nombre, crudo := range casos {
		t.Run(nombre, func(t *testing.T) {
			assert.Equal(t, "evt-1", publicar(t, crudo).EventID,
				"Bancolombia manda el mismo dato en niveles distintos segun el ambiente")
		})
	}
}

func TestParseo_TambienAplanaElPrimerElementoDeUnArreglo(t *testing.T) {
	msg := publicar(t, `{"data":[{"eventId":"evt-1","transaction":{"transferCode":"CUS-9"},"merchant":{"currency":"COP"}}]}`)

	assert.Equal(t, "evt-1", msg.EventID)
	assert.Equal(t, "CUS-9", msg.TransferCode)
	assert.Equal(t, "COP", msg.Currency)
}

func TestParseo_UnArregloVacioNoRevienta(t *testing.T) {
	msg := publicar(t, `{"data":[]}`)

	assert.NotEmpty(t, msg.EventID, "cae al id derivado del hash")
}

func TestParseo_LaRaizGanaSobreLoAnidado(t *testing.T) {
	msg := publicar(t, `{"eventId":"raiz","data":{"eventId":"anidado"}}`)

	assert.Equal(t, "raiz", msg.EventID,
		"mergeMap no sobreescribe lo que ya existe, asi que el nivel mas externo manda")
}

func TestParseo_SinEventIdDerivaUnoDelHashDelCuerpo(t *testing.T) {
	body := `{"transferState":"APPROVED"}`
	msg := publicar(t, body)

	sum := sha256.Sum256([]byte(body))
	assert.Equal(t, "sha256:"+hex.EncodeToString(sum[:]), msg.EventID,
		"el hash del cuerpo sirve de clave de idempotencia cuando el banco no manda id, pero dos pagos identicos colisionarian")
}

func TestParseo_SinTipoUsaElTipoPorDefectoDeQR(t *testing.T) {
	assert.Equal(t, "qr.transfer.notification", publicar(t, `{"eventId":"e"}`).Type)
}

func TestParseo_AceptaLosAliasDeCadaCampo(t *testing.T) {
	casos := []struct {
		nombre   string
		crudo    string
		revisar  func(*ports.BancolombiaWebhookMessage) string
		esperado string
	}{
		{"event_id", `{"event_id":"E1"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.EventID }, "E1"},
		{"messageId", `{"messageId":"E2"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.EventID }, "E2"},
		{"notificationId", `{"notificationId":"E3"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.EventID }, "E3"},
		{"eventType", `{"eventId":"e","eventType":"T1"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.Type }, "T1"},
		{"notificationType", `{"eventId":"e","notificationType":"T2"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.Type }, "T2"},
		{"status", `{"eventId":"e","status":"APPROVED"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.TransferState }, "APPROVED"},
		{"transactionState", `{"eventId":"e","transactionState":"REJECTED"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.TransferState }, "REJECTED"},
		{"paymentReference", `{"eventId":"e","paymentReference":"ORD-1"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.Reference }, "ORD-1"},
		{"invoiceNumber", `{"eventId":"e","invoiceNumber":"ORD-2"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.Reference }, "ORD-2"},
		{"cus", `{"eventId":"e","cus":"CUS-1"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.TransferCode }, "CUS-1"},
		{"approvalNumber", `{"eventId":"e","approvalNumber":"AP-1"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.TransferCode }, "AP-1"},
		{"currencyCode", `{"eventId":"e","currencyCode":"COP"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.Currency }, "COP"},
		{"originAccount", `{"eventId":"e","originAccount":"1234"}`, func(m *ports.BancolombiaWebhookMessage) string { return m.PayerAccount }, "1234"},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			assert.Equal(t, caso.esperado, caso.revisar(publicar(t, caso.crudo)))
		})
	}
}

func TestParseo_ElOrdenDeLosAliasDefineCualGana(t *testing.T) {
	msg := publicar(t, `{"eventId":"E1","event_id":"E2","messageId":"E3","id":"E4"}`)

	assert.Equal(t, "E1", msg.EventID,
		"eventId es el primero de la lista, por eso gana aunque lleguen todos")
}

func TestParseo_UnCampoDeTextoVacioNoCuentaYSePasaAlSiguienteAlias(t *testing.T) {
	msg := publicar(t, `{"eventId":"  ","event_id":"E2"}`)

	assert.Equal(t, "E2", msg.EventID,
		"un campo con solo espacios se trata como ausente")
}

func TestParseo_UnIdNumericoSeConvierteATexto(t *testing.T) {
	assert.Equal(t, "12345", publicar(t, `{"eventId":12345}`).EventID)
	assert.Equal(t, "12345.5", publicar(t, `{"eventId":12345.5}`).EventID)
}

func TestParseo_LeeElMontoComoNumeroOComoTexto(t *testing.T) {
	casos := []struct {
		crudo    string
		esperado float64
	}{
		{`{"eventId":"e","transferAmount":50000}`, 50000},
		{`{"eventId":"e","amount":"50000"}`, 50000},
		{`{"eventId":"e","amount":"  50000.75  "}`, 50000.75},
		{`{"eventId":"e","totalAmount":123}`, 123},
		{`{"eventId":"e","amount":"no es numero"}`, 0},
		{`{"eventId":"e"}`, 0},
		{`{"eventId":"e","amount":null}`, 0},
	}

	for _, caso := range casos {
		assert.InDelta(t, caso.esperado, publicar(t, caso.crudo).Amount, 0.001,
			"un monto ilegible se publica como cero en vez de abortar: el consumidor debe validarlo contra la orden")
	}
}

func TestParseo_LeeLaFechaEnVariosFormatos(t *testing.T) {
	casos := []struct {
		nombre   string
		crudo    string
		esperado time.Time
	}{
		{"RFC3339", `{"eventId":"e","transferDate":"2026-08-05T12:00:00Z"}`, time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)},
		{"fecha y hora sin zona", `{"eventId":"e","date":"2026-08-05 12:00:00"}`, time.Date(2026, 8, 5, 12, 0, 0, 0, time.UTC)},
		{"solo fecha", `{"eventId":"e","date":"2026-08-05"}`, time.Date(2026, 8, 5, 0, 0, 0, 0, time.UTC)},
		{"epoch en segundos", `{"eventId":"e","timestamp":"1785931200"}`, time.Unix(1785931200, 0)},
		{"epoch en milisegundos", `{"eventId":"e","timestamp":"1785931200000"}`, time.UnixMilli(1785931200000)},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			msg := publicar(t, caso.crudo)
			require.NotNil(t, msg.OccurredAt)
			assert.True(t, msg.OccurredAt.Equal(caso.esperado),
				"esperaba %v, recibi %v", caso.esperado, *msg.OccurredAt)
		})
	}
}

func TestParseo_UnaFechaIlegibleQuedaNula(t *testing.T) {
	casos := []string{
		`{"eventId":"e"}`,
		`{"eventId":"e","date":"ayer"}`,
		`{"eventId":"e","timestamp":"0"}`,
		`{"eventId":"e","timestamp":"-5"}`,
	}

	for _, crudo := range casos {
		assert.Nil(t, publicar(t, crudo).OccurredAt,
			"nunca se sustituye por time.Now(): la hora del banco no es la de recepcion")
	}
}

func TestWebhookSecretKey_SoloEsValidoConSecreto(t *testing.T) {
	var nulo *ports.BancolombiaConfig
	_, ok := nulo.WebhookSecretKey()
	assert.False(t, ok)

	_, ok = (&ports.BancolombiaConfig{WebhookSecret: ""}).WebhookSecretKey()
	assert.False(t, ok)

	secreto, ok := (&ports.BancolombiaConfig{WebhookSecret: "abc"}).WebhookSecretKey()
	assert.True(t, ok)
	assert.Equal(t, "abc", secreto)
}

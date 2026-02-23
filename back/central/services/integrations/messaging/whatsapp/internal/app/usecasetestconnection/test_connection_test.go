package usecasetestconnection

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/env"
	"github.com/secamc93/probability/back/central/shared/log"
)

// ---------------------------------------------------------------------------
// Mocks locales
// ---------------------------------------------------------------------------

// testDiscardLog es el logger base de zerolog que descarta todos los mensajes
var testDiscardLog = zerolog.New(io.Discard)

// testLogger descarta todos los logs para que los tests no sean ruidosos
type testLogger struct{}

var _ log.ILogger = (*testLogger)(nil)

func (l *testLogger) Info(ctx ...context.Context) *zerolog.Event {
	return testDiscardLog.Info()
}
func (l *testLogger) Error(ctx ...context.Context) *zerolog.Event {
	return testDiscardLog.Error()
}
func (l *testLogger) Warn(ctx ...context.Context) *zerolog.Event {
	return testDiscardLog.Warn()
}
func (l *testLogger) Debug(ctx ...context.Context) *zerolog.Event {
	return testDiscardLog.Debug()
}
func (l *testLogger) Fatal(ctx ...context.Context) *zerolog.Event {
	return testDiscardLog.WithLevel(zerolog.NoLevel)
}
func (l *testLogger) Panic(ctx ...context.Context) *zerolog.Event {
	return testDiscardLog.WithLevel(zerolog.NoLevel)
}
func (l *testLogger) With() zerolog.Context {
	return testDiscardLog.With()
}
func (l *testLogger) WithService(s string) log.ILogger   { return l }
func (l *testLogger) WithModule(m string) log.ILogger    { return l }
func (l *testLogger) WithBusinessID(id uint) log.ILogger { return l }

// testConfig implementa env.IConfig para tests
type testConfig struct {
	values map[string]string
}

func (c *testConfig) Get(key string) string {
	if c.values == nil {
		return ""
	}
	return c.values[key]
}

// testWhatsApp implementa ports.IWhatsApp para tests
type testWhatsApp struct {
	sendMessageFn func(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error)
}

func (m *testWhatsApp) SendMessage(ctx context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error) {
	if m.sendMessageFn != nil {
		return m.sendMessageFn(ctx, phoneNumberID, msg, accessToken)
	}
	return "wamid.test.ok", nil
}

// waFactory construye una factory que siempre retorna el mismo mock de IWhatsApp
func waFactory(wa ports.IWhatsApp) func(env.IConfig, log.ILogger) ports.IWhatsApp {
	return func(_ env.IConfig, _ log.ILogger) ports.IWhatsApp {
		return wa
	}
}

// newUseCase es un helper que crea el use case con config y logger de test
func newUseCase() *TestConnectionUseCase {
	cfg := &testConfig{values: map[string]string{"WHATSAPP_URL": "https://graph.facebook.com"}}
	return New(cfg, &testLogger{})
}

// ---------------------------------------------------------------------------
// TestConnection — sin test_phone_number (solo valida credenciales)
// ---------------------------------------------------------------------------

func TestTestConnection_SoloValidacionCredenciales_Exito(t *testing.T) {
	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id": "123456789",
		// sin test_phone_number
	}
	credentials := map[string]interface{}{
		"access_token": "EAAGxyz",
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(&testWhatsApp{}))

	if err != nil {
		t.Errorf("TestConnection() error inesperado cuando solo valida credenciales: %v", err)
	}
}

func TestTestConnection_SoloValidacionCredenciales_PhoneNumberIDComoFloat(t *testing.T) {
	// Meta API puede enviar phone_number_id como float64 en JSON
	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id": float64(987654321),
	}
	credentials := map[string]interface{}{
		"access_token": "EAAGxyz",
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(&testWhatsApp{}))

	if err != nil {
		t.Errorf("TestConnection() no debería fallar con phone_number_id como float64: %v", err)
	}
}

// ---------------------------------------------------------------------------
// TestConnection — con test_phone_number (envía mensaje hello_world)
// ---------------------------------------------------------------------------

func TestTestConnection_ConTestPhone_EnviaMensaje_Exito(t *testing.T) {
	var capturedPhoneNumberID uint
	var capturedToken string
	var capturedTemplateName string

	waClient := &testWhatsApp{
		sendMessageFn: func(_ context.Context, phoneNumberID uint, msg entities.TemplateMessage, accessToken string) (string, error) {
			capturedPhoneNumberID = phoneNumberID
			capturedToken = accessToken
			capturedTemplateName = msg.Template.Name
			return "wamid.hello.001", nil
		},
	}

	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id":   "111222333",
		"test_phone_number": "+573001234567",
	}
	credentials := map[string]interface{}{
		"access_token": "token-de-prueba",
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(waClient))

	if err != nil {
		t.Fatalf("TestConnection() error inesperado: %v", err)
	}
	if capturedPhoneNumberID != 111222333 {
		t.Errorf("phoneNumberID = %d, quería %d", capturedPhoneNumberID, 111222333)
	}
	if capturedToken != "token-de-prueba" {
		t.Errorf("accessToken = %q, quería %q", capturedToken, "token-de-prueba")
	}
	if capturedTemplateName != "hello_world" {
		t.Errorf("template.Name = %q, quería %q", capturedTemplateName, "hello_world")
	}
}

func TestTestConnection_ConTestPhone_ErrorEnvioMensaje(t *testing.T) {
	expectedErr := errors.New("credenciales de WhatsApp inválidas")

	waClient := &testWhatsApp{
		sendMessageFn: func(_ context.Context, _ uint, _ entities.TemplateMessage, _ string) (string, error) {
			return "", expectedErr
		},
	}

	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id":   "111222333",
		"test_phone_number": "+573001234567",
	}
	credentials := map[string]interface{}{
		"access_token": "token-invalido",
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(waClient))

	if err == nil {
		t.Fatal("TestConnection() esperaba error cuando el cliente WhatsApp falla")
	}
	if !errors.Is(err, expectedErr) {
		t.Errorf("TestConnection() error = %v, quería wrapping de %v", err, expectedErr)
	}
}

// ---------------------------------------------------------------------------
// TestConnection — errores de validación de parámetros
// ---------------------------------------------------------------------------

func TestTestConnection_ErrorAccessTokenVacio(t *testing.T) {
	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id": "123456",
	}
	credentials := map[string]interface{}{
		"access_token": "", // vacío
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(&testWhatsApp{}))

	if err == nil {
		t.Fatal("TestConnection() esperaba error cuando access_token está vacío")
	}
}

func TestTestConnection_ErrorAccessTokenAusente(t *testing.T) {
	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id": "123456",
	}
	credentials := map[string]interface{}{
		// sin access_token
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(&testWhatsApp{}))

	if err == nil {
		t.Fatal("TestConnection() esperaba error cuando access_token no está presente")
	}
}

func TestTestConnection_ErrorPhoneNumberIDAusente(t *testing.T) {
	uc := newUseCase()

	config := map[string]interface{}{
		// sin phone_number_id
	}
	credentials := map[string]interface{}{
		"access_token": "token-ok",
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(&testWhatsApp{}))

	if err == nil {
		t.Fatal("TestConnection() esperaba error cuando phone_number_id no está presente")
	}
}

func TestTestConnection_ErrorPhoneNumberIDInvalidoConTestPhone(t *testing.T) {
	// Cuando hay test_phone_number y phone_number_id no puede convertirse a uint
	uc := newUseCase()

	config := map[string]interface{}{
		"phone_number_id":   "no_es_numero",
		"test_phone_number": "+573001234567",
	}
	credentials := map[string]interface{}{
		"access_token": "token-ok",
	}

	err := uc.TestConnection(context.Background(), config, credentials, waFactory(&testWhatsApp{}))

	if err == nil {
		t.Fatal("TestConnection() esperaba error cuando phone_number_id no es numérico")
	}
}

// ---------------------------------------------------------------------------
// tempEnvConfig (implementación interna)
// ---------------------------------------------------------------------------

func TestTempEnvConfig_Get(t *testing.T) {
	cfg := &tempEnvConfig{
		values: map[string]string{
			"WHATSAPP_URL":   "https://graph.facebook.com",
			"WHATSAPP_TOKEN": "mi-token",
		},
	}

	if got := cfg.Get("WHATSAPP_URL"); got != "https://graph.facebook.com" {
		t.Errorf("Get(WHATSAPP_URL) = %q, quería %q", got, "https://graph.facebook.com")
	}
	if got := cfg.Get("CLAVE_INEXISTENTE"); got != "" {
		t.Errorf("Get(CLAVE_INEXISTENTE) = %q, quería string vacío", got)
	}
}

func TestTempEnvConfig_GetInt_RetornaCero(t *testing.T) {
	cfg := &tempEnvConfig{values: map[string]string{"KEY": "123"}}
	if got := cfg.GetInt("KEY"); got != 0 {
		t.Errorf("GetInt() = %d, quería 0 (no implementado)", got)
	}
}

func TestTempEnvConfig_GetBool_RetornaFalse(t *testing.T) {
	cfg := &tempEnvConfig{values: map[string]string{"KEY": "true"}}
	if got := cfg.GetBool("KEY"); got != false {
		t.Errorf("GetBool() = %v, quería false (no implementado)", got)
	}
}

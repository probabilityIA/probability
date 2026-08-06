package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNonRetryable_UnNegocioSinIntegracionDeWhatsAppNoSeReintenta(t *testing.T) {
	err := fmt.Errorf("error obteniendo configuracion de WhatsApp: %w",
		fmt.Errorf("no se encontró integración WhatsApp para business 26: registro no encontrado"))

	assert.True(t, IsNonRetryable(err),
		"es el caso del incidente del 2026-08-06: sin integracion configurada ningun reintento va a funcionar, y reencolar genera un bucle de 1500 intentos por segundo")
}

func TestIsNonRetryable_LosFallosPermanentesDeMetaNoSeReintentan(t *testing.T) {
	casos := map[string]string{
		"credenciales ausentes":  "credenciales no encontradas para business 26",
		"sin phone number id":    "phone_number_id no encontrado",
		"sin access token":       "access_token no encontrado",
		"clave ausente en cache": "key not found",
		"telefono invalido":      "número de teléfono inválido: faltan digitos",
		"parametro faltante":     "Required parameter is missing",
		"plantilla inexistente":  "template does not exist in es_CO",
	}

	for nombre, mensaje := range casos {
		t.Run(nombre, func(t *testing.T) {
			assert.True(t, IsNonRetryable(fmt.Errorf("%s", mensaje)),
				"reintentar no cambia el resultado; el mensaje se descarta con ACK y queda en el log")
		})
	}
}

func TestIsNonRetryable_LosCodigosDeErrorDeMetaQueSonDefinitivos(t *testing.T) {
	codigos := []string{"131008", "132001", "131009", "132018"}

	for _, codigo := range codigos {
		err := fmt.Errorf(`Meta respondio {"error":{"code":%s}}`, codigo)
		assert.True(t, IsNonRetryable(err),
			"el codigo %s de la Graph API es permanente: plantilla o parametro mal, no una caida", codigo)
	}
}

func TestIsNonRetryable_ReconoceLosErroresTipadosDelDominio(t *testing.T) {
	assert.True(t, IsNonRetryable(&ErrTemplateNotFound{TemplateName: "reporte_saldo_billetera"}),
		"la plantilla no existe: reintentar no la va a crear")
	assert.True(t, IsNonRetryable(&ErrMissingVariable{TemplateName: "t", VariableName: "saldo", VariableKey: "2"}),
		"al mensaje le falta un dato: el reintento mandaria el mismo mensaje incompleto")
}

func TestIsNonRetryable_ReconoceLosErroresTipadosAunqueVenganEnvueltos(t *testing.T) {
	envuelto := fmt.Errorf("error enviando plantilla: %w", &ErrTemplateNotFound{TemplateName: "x"})

	assert.True(t, IsNonRetryable(envuelto),
		"los casos de uso envuelven el error con contexto; errors.As tiene que atravesar esa capa")
}

func TestIsNonRetryable_LosFallosTransitoriosSiSeReintentan(t *testing.T) {
	casos := map[string]string{
		"meta caida":          "Meta respondio 503 Service Unavailable",
		"timeout de red":      "context deadline exceeded",
		"conexion rechazada":  "dial tcp: connection refused",
		"limite de tasa":      "rate limit exceeded, try again later",
		"error 500 generico":  "internal server error",
		"base de datos caida": "driver: bad connection",
	}

	for nombre, mensaje := range casos {
		t.Run(nombre, func(t *testing.T) {
			assert.False(t, IsNonRetryable(fmt.Errorf("%s", mensaje)),
				"el negocio SI tiene WhatsApp configurado y el fallo es pasajero: hay que reencolar hasta que vuelva")
		})
	}
}

func TestIsNonRetryable_UnErrorNuloNoEsPermanente(t *testing.T) {
	assert.False(t, IsNonRetryable(nil))
}

func TestIsNonRetryable_UnErrorDesconocidoSeReintentaPorDefecto(t *testing.T) {
	assert.False(t, IsNonRetryable(fmt.Errorf("algo raro que nadie previo")),
		"falla del lado del reintento: es preferible reencolar de mas que descartar un mensaje que si se podia entregar")
}

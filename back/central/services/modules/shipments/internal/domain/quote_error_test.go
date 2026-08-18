package domain

import (
	"strings"
	"testing"
)

func TestSanitizeGuideError_EmailInvalido(t *testing.T) {
	raw := `{"status":"NOT OK","status_messages":[{"error":["Unprocessed Entity.",{"destination":{"email":["Dato invalido, el campo debe ser de tipo email"]}}]}]}`

	got := SanitizeGuideError(raw)

	if !strings.Contains(strings.ToLower(got), "correo") {
		t.Fatalf("esperaba un mensaje sobre el correo, obtuve: %q", got)
	}
	if MentionsProvider(got) {
		t.Fatalf("el mensaje expone el proveedor: %q", got)
	}
}

func TestSanitizeGuideError_NuncaExponeElProveedor(t *testing.T) {
	entradas := []string{
		"Error de Transporte: EnvioClick devolvio 500",
		"dial tcp api.envioclickpro.com.co:443 connect: connection refused",
		"mipaquete respondio con un error inesperado",
		"shipit unavailable",
		"",
		"algo raro que nadie mapeo",
	}

	for _, in := range entradas {
		got := SanitizeGuideError(in)
		if MentionsProvider(got) {
			t.Fatalf("entrada %q produjo un mensaje que expone el proveedor: %q", in, got)
		}
		if strings.TrimSpace(got) == "" {
			t.Fatalf("entrada %q produjo un mensaje vacio", in)
		}
	}
}

func TestSanitizeGuideError_SinCoincidenciaUsaElGenerico(t *testing.T) {
	got := SanitizeGuideError("codigo interno 9182 sin significado conocido")

	if got != GenericGuideErrorMessage {
		t.Fatalf("esperaba el mensaje generico, obtuve: %q", got)
	}
}

func TestSanitizeGuideError_SaldoInsuficiente(t *testing.T) {
	got := SanitizeGuideError("Insufficient balance to create the shipment")

	if !strings.Contains(strings.ToLower(got), "saldo") {
		t.Fatalf("esperaba un mensaje sobre el saldo, obtuve: %q", got)
	}
}

func TestSafeQuoteMessage_DescartaTextoConProveedor(t *testing.T) {
	got := SafeQuoteMessage("la cuenta de EnvioClick no tiene cobertura")

	if got != GenericGuideErrorMessage {
		t.Fatalf("esperaba el mensaje generico, obtuve: %q", got)
	}
}

func TestSafeQuoteMessage_ConservaMensajeLimpio(t *testing.T) {
	in := "la transportadora elegida en el checkout ya no esta disponible, genera la guia manualmente"

	if got := SafeQuoteMessage(in); got != in {
		t.Fatalf("esperaba %q, obtuve %q", in, got)
	}
}

package client

import (
	"strings"
	"testing"
)

func TestMensajesRealesDeProduccion(t *testing.T) {
	casos := []struct {
		nombre   string
		body     string
		contiene []string
	}{
		{
			nombre:   "direccion de destino fuera de rango",
			body:     `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.",{"destination":{"address":["El campo debe tener entre 2 y 50 caracteres."]}}]}]}`,
			contiene: []string{"direcci", "destino", "entre 2 y 50 caracteres"},
		},
		{
			nombre:   "barrio de destino fuera de rango",
			body:     `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.",{"destination":{"suburb":["El campo debe tener entre 2 y 30 caracteres."]}}]}]}`,
			contiene: []string{"barrio", "destino", "entre 2 y 30 caracteres"},
		},
		{
			nombre:   "descripcion fuera de rango",
			body:     `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.",{"description":["El campo debe tener entre 2 y 25 caracteres."]}]}]}`,
			contiene: []string{"descripci", "entre 2 y 25 caracteres"},
		},
		{
			nombre:   "correo de destino invalido",
			body:     `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.",{"destination":{"email":["Dato inválido, el campo debe ser de tipo email"]}}]}]}`,
			contiene: []string{"correo", "destinatario"},
		},
		{
			nombre:   "saldo insuficiente",
			body:     `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.","No tiene suficiente crédito disponible en su cuenta. Necesita hacer un depósito en www.envioclickpro.com.co"]}]}`,
			contiene: []string{"Saldo insuficiente", "billetera"},
		},
		{
			nombre:   "credenciales invalidas",
			body:     `{"status":"NOT OK","status_codes":[401],"status_messages":[{"error":["Unauthorized.","API Key inválido. Revisa que tu API key sea válido. Puedes encontrar tu  API Key en tu cuenta de www.envioclickpro.com.co"]}]}`,
			contiene: []string{"credenciales", "integraci"},
		},
		{
			nombre:   "guia no generada",
			body:     `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.","No se pudo generar la guía"]}]}`,
			contiene: []string{"cobertura"},
		},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			got := parseEnvioClickError([]byte(caso.body))
			for _, fragmento := range caso.contiene {
				if !strings.Contains(strings.ToLower(got), strings.ToLower(fragmento)) {
					t.Fatalf("el mensaje %q no menciona %q", got, fragmento)
				}
			}
		})
	}
}

func TestNuncaSeFiltraElNombreDelProveedor(t *testing.T) {
	bodies := []string{
		`{"status":"NOT OK","status_codes":[401],"status_messages":[{"error":["Unauthorized.","API Key inválido. Puedes encontrar tu API Key en tu cuenta de www.envioclickpro.com.co"]}]}`,
		`{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.","Necesita hacer un depósito en www.envioclickpro.com.co"]}]}`,
		`{"status":"NOT OK","status_codes":[500],"status_messages":[{"error":["Algo raro paso en EnvioClick Pro, contacta a envioclickpro.com.co"]}]}`,
		`una respuesta cruda de EnvioClick que no es json`,
		`{"status":"NOT OK","status_messages":[{"error":[{"destination":{"address":["Revisa en www.envioclickpro.com.co"]}}]}]}`,
	}

	prohibidos := []string{"envioclick", "envioclickpro", "envioclickpro.com.co"}

	for _, body := range bodies {
		got := strings.ToLower(parseEnvioClickError([]byte(body)))
		for _, prohibido := range prohibidos {
			if strings.Contains(got, prohibido) {
				t.Fatalf("el mensaje %q filtra el nombre del proveedor (%q)", got, prohibido)
			}
		}
	}
}

func TestVariosCamposEnUnSoloError(t *testing.T) {
	body := `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.",{"origin":{"email":["Dato inválido, el campo debe ser de tipo email"],"phone":["El campo debe tener entre 7 y 10 caracteres."],"suburb":["El campo debe tener entre 2 y 30 caracteres."]},"destination":{"address":["El campo debe tener entre 2 y 50 caracteres."]}}]}]}`

	got := parseEnvioClickError([]byte(body))

	if !strings.Contains(got, ";") {
		t.Fatalf("con varios campos invalidos el mensaje deberia enumerarlos: %q", got)
	}
	if !strings.Contains(strings.ToLower(got), "origen") {
		t.Fatalf("el mensaje deberia mencionar los campos de origen: %q", got)
	}
	if strings.Count(got, ";") > 2 {
		t.Fatalf("el mensaje no deberia enumerar mas de tres campos: %q", got)
	}
}

func TestNoDevuelveJsonCrudoAlUsuario(t *testing.T) {
	body := `{"status":"NOT OK","status_codes":[422],"status_messages":[{"error":["Unprocessed Entity.",{"destination":{"address":["El campo debe tener entre 2 y 50 caracteres."]}}]}]}`

	got := parseEnvioClickError([]byte(body))

	for _, basura := range []string{"{", "}", "[", "]", "status_messages", "Unprocessed"} {
		if strings.Contains(got, basura) {
			t.Fatalf("el mensaje %q todavia expone %q", got, basura)
		}
	}
}

package handlers

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func i64ptr(v int64) *int64 { return &v }

func TestBuildRetryPayload_ReemplazaElDestinoPlaceholderConElDeLaOrden(t *testing.T) {
	quote := &domain.SavedQuote{
		SelectedCarrier: "COORDINADORA",
		SelectedIDRate:  i64ptr(25300507),
		RequestPayload: map[string]interface{}{
			"destination": map[string]interface{}{
				"firstName": "Destinatario",
				"lastName":  "Default",
				"email":     "",
				"phone":     "",
				"address":   "carrera 46 # 50 -21c",
				"city":      "GUADALAJARA DE BUGA",
				"daneCode":  "76111000",
			},
			"packages": []interface{}{map[string]interface{}{"weight": 1.0}},
		},
	}
	recipient := &domain.OrderRecipient{
		FullName:  "Ferney Gonzalez",
		FirstName: "Ferney",
		LastName:  "Gonzalez",
		Email:     "ferney@example.com",
		Phone:     "+57 315 570 7677",
		Address:   "carrera 11#17-27 barrio Juan mayor",
		City:      "GUADALAJARA DE BUGA",
	}

	payload := buildRetryPayload(quote, recipient, "order-uuid")

	dest, ok := payload["destination"].(map[string]interface{})
	if !ok {
		t.Fatal("se esperaba un bloque destination en el payload")
	}
	if dest["email"] != "ferney@example.com" {
		t.Errorf("email esperado del destinatario real, obtenido: %v", dest["email"])
	}
	if dest["phone"] != "3155707677" {
		t.Errorf("telefono esperado normalizado a 10 digitos, obtenido: %v", dest["phone"])
	}
	if dest["firstName"] != "Ferney" || dest["lastName"] != "Gonzalez" {
		t.Errorf("se esperaba el nombre real, obtenido: %v %v", dest["firstName"], dest["lastName"])
	}
	if dest["address"] != "carrera 11#17-27 barrio Juan mayor" {
		t.Errorf("se esperaba la direccion de la orden, obtenida: %v", dest["address"])
	}
	if dest["daneCode"] != "76111000" {
		t.Errorf("el daneCode define la tarifa y debe conservarse, obtenido: %v", dest["daneCode"])
	}
	if payload["idRate"] != int64(25300507) {
		t.Errorf("se esperaba reusar la tarifa cotizada, obtenida: %v", payload["idRate"])
	}
	if payload["carrier"] != "COORDINADORA" {
		t.Errorf("se esperaba reusar la transportadora elegida, obtenida: %v", payload["carrier"])
	}
	if payload["packages"] == nil {
		t.Error("se esperaba conservar los paquetes de la cotizacion")
	}
}

func TestBuildRetryPayload_NoMutaLaCotizacionOriginal(t *testing.T) {
	original := map[string]interface{}{
		"firstName": "Destinatario",
		"email":     "",
	}
	quote := &domain.SavedQuote{
		SelectedIDRate: i64ptr(1),
		RequestPayload: map[string]interface{}{"destination": original},
	}
	recipient := &domain.OrderRecipient{
		FirstName: "Ferney",
		Email:     "ferney@example.com",
		Phone:     "3155707677",
		Address:   "calle 1",
		City:      "BUGA",
	}

	buildRetryPayload(quote, recipient, "order-uuid")

	if original["email"] != "" || original["firstName"] != "Destinatario" {
		t.Error("el payload guardado de la cotizacion no debe mutarse al construir el reintento")
	}
}

func TestMissingRecipientData(t *testing.T) {
	casos := []struct {
		nombre    string
		recipient domain.OrderRecipient
		falla     bool
	}{
		{"completo", domain.OrderRecipient{Email: "a@b.com", Phone: "3155707677", Address: "calle 1"}, false},
		{"sin email", domain.OrderRecipient{Email: "", Phone: "3155707677", Address: "calle 1"}, true},
		{"email invalido", domain.OrderRecipient{Email: "noesunemail", Phone: "3155707677", Address: "calle 1"}, true},
		{"email sin dominio", domain.OrderRecipient{Email: "a@b", Phone: "3155707677", Address: "calle 1"}, true},
		{"email corto valido", domain.OrderRecipient{Email: "a@b.co", Phone: "3155707677", Address: "calle 1"}, false},
		{"telefono corto", domain.OrderRecipient{Email: "a@b.com", Phone: "123", Address: "calle 1"}, true},
		{"telefono con indicativo", domain.OrderRecipient{Email: "a@b.com", Phone: "573155707677", Address: "calle 1"}, false},
		{"sin direccion", domain.OrderRecipient{Email: "a@b.com", Phone: "3155707677", Address: "  "}, true},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			msg := missingRecipientData(&c.recipient)
			if c.falla && msg == "" {
				t.Error("se esperaba un mensaje de dato faltante")
			}
			if !c.falla && msg != "" {
				t.Errorf("no se esperaba mensaje, obtenido: %q", msg)
			}
		})
	}
}

package templates

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func TestValidateMetaRules(t *testing.T) {
	cases := []struct {
		name      string
		body      string
		header    string
		footer    string
		buttons   []entities.TemplateButton
		variables int
		wantErr   bool
	}{
		{name: "valida", body: "Hola {{1}}, tu pedido {{2}} ya va en camino.", header: "Tu pedido", variables: 2},
		{name: "sin variables", body: "Gracias por tu compra"},
		{name: "empieza con variable", body: "{{1}}, gracias por tu compra hoy", variables: 1, wantErr: true},
		{name: "termina con variable", body: "Gracias por tu compra {{1}}", variables: 1, wantErr: true},
		{name: "termina con punto", body: "Gracias por tu compra {{1}}.", variables: 1},
		{name: "variables seguidas", body: "Hola querido cliente {{1}} {{2}} gracias por todo.", variables: 2, wantErr: true},
		{name: "variable mal escrita", body: "Hola {{nombre}}, gracias por tu compra", wantErr: true},
		{name: "llave suelta", body: "Hola {{1}, gracias por tu compra", wantErr: true},
		{name: "poco texto", body: "Hola {{1}} y {{2}}.", variables: 2, wantErr: true},
		{name: "encabezado con emoji", body: "Gracias por tu compra", header: "Hola \U0001F44B", wantErr: true},
		{name: "encabezado con variable", body: "Gracias por tu compra", header: "Hola {{1}}", wantErr: true},
		{name: "encabezado con formato", body: "Gracias por tu compra", header: "*Oferta*", wantErr: true},
		{name: "encabezado con salto", body: "Gracias por tu compra", header: "Hola\namigo", wantErr: true},
		{name: "pie con variable", body: "Gracias por tu compra", footer: "Atentamente {{1}}", wantErr: true},
		{name: "boton con variable", body: "Gracias por tu compra", buttons: []entities.TemplateButton{{Text: "Ver {{1}}"}}, wantErr: true},
		{name: "cuerpo con emoji", body: "Hola {{1}} \U0001F44B tu pedido ya va en camino.", variables: 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateMetaRules(tc.body, tc.header, tc.footer, tc.buttons, tc.variables)
			if tc.wantErr && err == nil {
				t.Fatal("se esperaba error")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("no se esperaba error: %v", err)
			}
		})
	}
}

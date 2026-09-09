package usecasetemplates

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/ports"
)

func TestIsPlatformOnlyTemplate(t *testing.T) {
	platformOnly := []string{"alerta_servidor", "recuperacion_codigo", "reporte_saldo_billetera_v2", "hello_world"}
	for _, name := range platformOnly {
		if !isPlatformOnlyTemplate(name) {
			t.Fatalf("la plantilla %s debe quedar en el WABA de la plataforma", name)
		}
	}

	delNegocio := []string{"confirmacion_pedido", "guia_envio_generada", "prueba_conexion", "pedido_entregado_cod"}
	for _, name := range delNegocio {
		if isPlatformOnlyTemplate(name) {
			t.Fatalf("la plantilla %s debe replicarse al WABA del negocio", name)
		}
	}
}

func TestComponentsForCreateDejaIgualLasNoAutenticacion(t *testing.T) {
	template := ports.TemplateDefinitionRemote{
		Category: "UTILITY",
		Components: []map[string]any{
			{"type": "BODY", "text": "Hola {{1}}"},
			{"type": "BUTTONS", "buttons": []any{map[string]any{"type": "QUICK_REPLY", "text": "Confirmar"}}},
		},
	}

	got := componentsForCreate(template)
	if len(got) != 2 || got[0]["text"] != "Hola {{1}}" {
		t.Fatalf("los componentes de una plantilla normal no deben cambiar: %v", got)
	}
}

func TestComponentsForCreateReconstruyeAutenticacion(t *testing.T) {
	template := ports.TemplateDefinitionRemote{
		Category: "AUTHENTICATION",
		Components: []map[string]any{
			{
				"type":                        "BODY",
				"text":                        "Tu código de verificación es *{{1}}*.",
				"add_security_recommendation": true,
				"example":                     map[string]any{"body_text": []any{[]any{"123456"}}},
			},
			{"type": "FOOTER", "text": "Vence en 10 minutos.", "code_expiration_minutes": float64(10)},
			{"type": "BUTTONS", "buttons": []any{map[string]any{
				"type": "URL",
				"text": "Copiar código",
				"url":  "https://www.whatsapp.com/otp/code/?otp_type=COPY_CODE&code=otp{{1}}",
			}}},
		},
	}

	got := componentsForCreate(template)
	if len(got) != 3 {
		t.Fatalf("se esperaban 3 componentes, llegaron %d: %v", len(got), got)
	}

	body := got[0]
	if _, tieneTexto := body["text"]; tieneTexto {
		t.Fatal("Meta rechaza el BODY con texto en plantillas de autenticación")
	}
	if _, tieneEjemplo := body["example"]; tieneEjemplo {
		t.Fatal("Meta rechaza el example en plantillas de autenticación")
	}
	if body["add_security_recommendation"] != true {
		t.Fatalf("se perdió add_security_recommendation: %v", body)
	}

	footer := got[1]
	if _, tieneTexto := footer["text"]; tieneTexto {
		t.Fatal("el FOOTER de autenticación solo lleva code_expiration_minutes")
	}
	if footer["code_expiration_minutes"] != 10 {
		t.Fatalf("code_expiration_minutes incorrecto: %v", footer)
	}

	buttons, ok := got[2]["buttons"].([]map[string]any)
	if !ok || len(buttons) != 1 {
		t.Fatalf("se esperaba un botón OTP: %v", got[2])
	}
	if buttons[0]["type"] != "OTP" || buttons[0]["otp_type"] != "COPY_CODE" {
		t.Fatalf("el botón debe ser OTP COPY_CODE: %v", buttons[0])
	}
	if _, tieneURL := buttons[0]["url"]; tieneURL {
		t.Fatal("el botón OTP no lleva url")
	}
}

func TestCandidatesForEventCubreTodosLosEventosDelCatalogo(t *testing.T) {
	eventos := []string{
		"order.created",
		"order.created_with_map",
		"order.shipped",
		"order.delivered",
		"order.canceled",
		"shipment.guide_generated",
		"wallet.low_balance",
	}

	for _, evento := range eventos {
		if len(candidatesForEvent(evento)) == 0 {
			t.Errorf("el evento %q no tiene plantillas candidatas: la vista previa saldria vacia", evento)
		}
	}
}

func TestCandidatesForEventConMapaApuntaALaPlantillaCorrecta(t *testing.T) {
	candidatos := candidatesForEvent("order.created_with_map")
	if len(candidatos) != 1 {
		t.Fatalf("se esperaba 1 candidato, hay %d", len(candidatos))
	}
	if candidatos[0].name != "confirmacion_pedido_contraentrega_mapa" {
		t.Errorf("apunta a %q", candidatos[0].name)
	}
	if candidatos[0].condition == "" {
		t.Error("el candidato debe explicar cuando aplica, es lo que lee el usuario en la UI")
	}
}

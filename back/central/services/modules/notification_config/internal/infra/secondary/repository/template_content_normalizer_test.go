package repository

import "testing"

func TestParseTemplateContentFormats(t *testing.T) {
	cases := []struct {
		name, templateName, content, wantName string
		wantParams                            map[string]string
		wantOK                                bool
	}{
		{"campana con parametros", "26_dice_que_no", "26_dice_que_no: Andres", "26_dice_que_no", map[string]string{"1": "Andres"}, true},
		{"varios parametros", "promo", "promo: Ana | Tienda", "promo", map[string]string{"1": "Ana", "2": "Tienda"}, true},
		{"sin parametros", "26_hola", "26_hola", "26_hola", map[string]string{}, true},
		{"plantilla sin variables", "pedido_confirmado_v2", "Plantilla: pedido_confirmado_v2", "pedido_confirmado_v2", map[string]string{}, true},
		{"plantilla con variables", "confirmacion_pedido", "Plantilla: confirmacion_pedido (variables: map[1:Sebastian Camacho 2:Probability])", "confirmacion_pedido", map[string]string{"1": "Sebastian Camacho", "2": "Probability"}, true},
		{"formato legacy", "", "Template: guia_envio_generada, Variables: map[1:Ana 2:Demo]", "guia_envio_generada", map[string]string{"1": "Ana", "2": "Demo"}, true},
		{"texto ya renderizado", "confirmacion_pedido_contraentrega", "Hola Ana, tu pedido fue recibido", "", nil, false},
		{"texto libre", "", "Hola, como vas?", "", nil, false},
	}

	for _, tc := range cases {
		name, params, ok := parseTemplateContent(tc.templateName, tc.content)
		if ok != tc.wantOK || name != tc.wantName {
			t.Fatalf("%s: got (%q, %v) want (%q, %v)", tc.name, name, ok, tc.wantName, tc.wantOK)
		}
		for k, v := range tc.wantParams {
			if params[k] != v {
				t.Fatalf("%s: param %s = %q, want %q", tc.name, k, params[k], v)
			}
		}
	}
}

func TestResolveMessageContentRendersHeaderBodyFooter(t *testing.T) {
	texts := map[string]templateText{
		"26_dice_que_no": {Header: "hola", Body: "dices que no ? {{1}}", Footer: "que mal"},
	}

	got := resolveMessageContent(texts, "26_dice_que_no", "26_dice_que_no: Andres")
	want := "hola\n\ndices que no ? Andres\n\nque mal"
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}

func TestResolveMessageContentMissingParamUsesDash(t *testing.T) {
	texts := map[string]templateText{"saludo": {Body: "Hola {{1}}, somos {{2}}"}}

	if got := resolveMessageContent(texts, "saludo", "saludo: Ana"); got != "Hola Ana, somos -" {
		t.Fatalf("got %q", got)
	}
}

func TestResolveMessageContentKeepsUnknownTemplate(t *testing.T) {
	content := "Plantilla: pedido_confirmado_v2"
	if got := resolveMessageContent(map[string]templateText{}, "pedido_confirmado_v2", content); got != content {
		t.Fatalf("unknown template must keep stored content, got %q", got)
	}
}

func TestResolveMessageContentKeepsRenderedText(t *testing.T) {
	texts := map[string]templateText{"confirmacion": {Body: "otra cosa"}}
	content := "Hola Ana, tu pedido fue recibido"
	if got := resolveMessageContent(texts, "confirmacion", content); got != content {
		t.Fatalf("already rendered text must stay, got %q", got)
	}
}

func TestResolveMessageContentUsesFallbackBodies(t *testing.T) {
	got := resolveMessageContent(map[string]templateText{}, "", "Template: guia_envio_generada, Variables: map[1:Ana 2:Demo 3:DEM-1 4:G1 5:Envia]")
	if got == "" || got[:8] != "Hola Ana" {
		t.Fatalf("fallback body not rendered: %q", got)
	}
}

func TestResolveMessageReturnsTemplateButtons(t *testing.T) {
	texts := map[string]templateText{
		"confirmacion": {Body: "Hola {{1}}", Buttons: []templateButton{{Text: "Confirmar pedido", Type: "QUICK_REPLY"}}},
	}

	content, buttons := resolveMessage(texts, "confirmacion", "confirmacion: Ana")
	if content != "Hola Ana" || len(buttons) != 1 || buttons[0].Text != "Confirmar pedido" {
		t.Fatalf("got %q %+v", content, buttons)
	}

	rendered, buttons := resolveMessage(texts, "confirmacion", "Hola Ana, texto ya guardado")
	if rendered != "Hola Ana, texto ya guardado" || len(buttons) != 1 {
		t.Fatalf("already rendered text must still carry buttons, got %q %+v", rendered, buttons)
	}
}

func TestParseTemplateButtonsAddsOptOutToMarketing(t *testing.T) {
	buttons := parseTemplateButtons(`[{"URL": "", "Text": "Si", "Type": "QUICK_REPLY"}]`, "MARKETING")
	if len(buttons) != 2 || buttons[1].Text != "Dejar de recibir" {
		t.Fatalf("marketing must end with the opt-out button, got %+v", buttons)
	}

	again := parseTemplateButtons(`[{"Text": "Dejar de recibir", "Type": "QUICK_REPLY"}]`, "MARKETING")
	if len(again) != 1 {
		t.Fatalf("opt-out must not be duplicated, got %+v", again)
	}

	if utility := parseTemplateButtons("", "UTILITY"); utility != nil {
		t.Fatalf("utility without buttons must be nil, got %+v", utility)
	}
}

func TestTemplateNamesToResolveDeduplicates(t *testing.T) {
	names := templateNamesToResolve([][2]string{
		{"a", "a: 1"},
		{"a", "a: 2"},
		{"", "texto libre"},
		{"b", "Plantilla: b"},
		{"c", "Hola, texto ya renderizado"},
	})
	if len(names) != 3 || names[0] != "a" || names[1] != "b" || names[2] != "c" {
		t.Fatalf("unexpected names: %v", names)
	}
}

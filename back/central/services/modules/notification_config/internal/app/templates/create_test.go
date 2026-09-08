package templates

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

func baseDTO() dtos.CreateTemplateDTO {
	return dtos.CreateTemplateDTO{
		BusinessID: 26,
		Name:       "vuelve_pronto",
		Language:   "es",
		Category:   entities.TemplateCategoryMarketing,
		BodyText:   "Hola {{1}}, hace {{2}} dias que no te vemos.",
		Variables: []dtos.TemplateVariableDTO{
			{Position: 1, Source: "customer.first_name"},
			{Position: 2, Source: "customer.days_inactive"},
		},
	}
}

func TestBuildTemplateOK(t *testing.T) {
	template, err := buildTemplate(baseDTO())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}
	if template.Status != entities.TemplateStatusDraft {
		t.Fatalf("status esperado draft, llego %s", template.Status)
	}
	if len(template.Variables) != 2 {
		t.Fatalf("se esperaban 2 variables, llegaron %d", len(template.Variables))
	}
	if template.Variables[0].Position != 1 || template.Variables[1].Position != 2 {
		t.Fatal("las variables no quedaron ordenadas por posicion")
	}
}

func TestBuildTemplateRejectsNameWithUppercase(t *testing.T) {
	dto := baseDTO()
	dto.Name = "Vuelve Pronto"

	if _, err := buildTemplate(dto); err == nil {
		t.Fatal("se esperaba error por nombre invalido")
	}
}

func TestBuildTemplateRejectsVariableMismatch(t *testing.T) {
	dto := baseDTO()
	dto.Variables = dto.Variables[:1]

	if _, err := buildTemplate(dto); err == nil {
		t.Fatal("se esperaba error porque el cuerpo usa 2 variables y se declaro 1")
	}
}

func TestBuildTemplateRejectsUnknownVariableSource(t *testing.T) {
	dto := baseDTO()
	dto.Variables[1].Source = "customer.tarjeta_de_credito"

	if _, err := buildTemplate(dto); err == nil {
		t.Fatal("se esperaba error por variable no permitida")
	}
}

func TestBuildTemplateRejectsDuplicatePosition(t *testing.T) {
	dto := baseDTO()
	dto.Variables[1].Position = 1

	if _, err := buildTemplate(dto); err == nil {
		t.Fatal("se esperaba error por posicion duplicada")
	}
}

func TestMarketingTemplateGetsOptOutButton(t *testing.T) {
	template, err := buildTemplate(baseDTO())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if !componentsHaveOptOut(template.Components) {
		t.Fatal("una plantilla MARKETING sin botones debe recibir el boton de baja")
	}
}

func TestUtilityTemplateDoesNotGetOptOutButton(t *testing.T) {
	dto := baseDTO()
	dto.Category = entities.TemplateCategoryUtility

	template, err := buildTemplate(dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	if componentsHaveOptOut(template.Components) {
		t.Fatal("una plantilla UTILITY no debe recibir el boton de baja")
	}
}

func TestOptOutButtonIsNotDuplicated(t *testing.T) {
	dto := baseDTO()
	dto.Buttons = []dtos.TemplateButtonDTO{{Type: "QUICK_REPLY", Text: OptOutButtonText}}

	template, err := buildTemplate(dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	count := 0
	for _, component := range template.Components {
		if component["type"] != "BUTTONS" {
			continue
		}
		buttons, _ := component["buttons"].([]map[string]any)
		for _, button := range buttons {
			if text, _ := button["text"].(string); text == OptOutButtonText {
				count++
			}
		}
	}

	if count != 1 {
		t.Fatalf("se esperaba un solo boton de baja, hay %d", count)
	}
}

func TestBodyExampleUsesFallbackOrLabel(t *testing.T) {
	dto := baseDTO()
	dto.Variables[0].Fallback = "Ana"

	template, err := buildTemplate(dto)
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	for _, component := range template.Components {
		if component["type"] != "BODY" {
			continue
		}
		example, ok := component["example"].(map[string]any)
		if !ok {
			t.Fatal("el cuerpo con variables debe llevar example")
		}
		rows, ok := example["body_text"].([][]string)
		if !ok || len(rows) != 1 || len(rows[0]) != 2 {
			t.Fatalf("example mal formado: %#v", example)
		}
		if rows[0][0] != "Ana" {
			t.Fatalf("se esperaba el fallback Ana, llego %s", rows[0][0])
		}
		return
	}

	t.Fatal("no se encontro el componente BODY")
}

func componentsHaveOptOut(components []map[string]any) bool {
	for _, component := range components {
		if component["type"] != "BUTTONS" {
			continue
		}
		buttons, _ := component["buttons"].([]map[string]any)
		for _, button := range buttons {
			if text, _ := button["text"].(string); text == OptOutButtonText {
				return true
			}
		}
	}
	return false
}

func TestBodyExampleUsesRealisticSampleNotLabel(t *testing.T) {
	template, err := buildTemplate(baseDTO())
	if err != nil {
		t.Fatalf("no se esperaba error: %v", err)
	}

	for _, component := range template.Components {
		if component["type"] != "BODY" {
			continue
		}
		example := component["example"].(map[string]any)
		rows := example["body_text"].([][]string)

		if rows[0][0] == "Nombre del cliente" || rows[0][1] == "Dias sin comprar" {
			t.Fatalf("el example no puede ser la etiqueta de la variable: %v", rows[0])
		}
		if rows[0][0] != "Ana" || rows[0][1] != "45" {
			t.Fatalf("se esperaban valores de muestra reales, llegaron %v", rows[0])
		}
		return
	}

	t.Fatal("no se encontro el componente BODY")
}

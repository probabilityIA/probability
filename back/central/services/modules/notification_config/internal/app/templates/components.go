package templates

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

const OptOutButtonText = "Dejar de recibir"

var sampleValues = map[string]string{
	"customer.first_name":    "Ana",
	"customer.full_name":     "Ana Ramirez",
	"customer.days_inactive": "45",
	"customer.last_product":  "Camiseta blanca",
	"customer.total_orders":  "3",
	"business.name":          "Mi Tienda",
	"sender.name":            "Isabel Rojas",
	"campaign.name":          "Ruta 30",
}

func BuildMetaComponents(template *entities.WhatsappTemplate) []map[string]any {
	components := make([]map[string]any, 0, 4)

	if header := strings.TrimSpace(template.HeaderText); header != "" {
		components = append(components, map[string]any{
			"type":   "HEADER",
			"format": "TEXT",
			"text":   header,
		})
	}

	body := map[string]any{
		"type": "BODY",
		"text": template.BodyText,
	}
	if example := bodyExample(template); len(example) > 0 {
		body["example"] = map[string]any{"body_text": [][]string{example}}
	}
	components = append(components, body)

	if footer := strings.TrimSpace(template.FooterText); footer != "" {
		components = append(components, map[string]any{
			"type": "FOOTER",
			"text": footer,
		})
	}

	if buttons := buildButtons(template); len(buttons) > 0 {
		components = append(components, map[string]any{
			"type":    "BUTTONS",
			"buttons": buttons,
		})
	}

	return components
}

func bodyExample(template *entities.WhatsappTemplate) []string {
	if len(template.Variables) == 0 {
		return nil
	}

	example := make([]string, 0, len(template.Variables))
	for _, variable := range template.Variables {
		value := strings.TrimSpace(variable.Fallback)
		if value == "" {
			value = sampleValues[variable.Source]
		}
		if value == "" {
			value = "ejemplo"
		}
		example = append(example, value)
	}

	return example
}

func buildButtons(template *entities.WhatsappTemplate) []map[string]any {
	buttons := make([]map[string]any, 0, len(template.Buttons)+1)

	for _, button := range template.Buttons {
		switch strings.ToUpper(strings.TrimSpace(button.Type)) {
		case "URL":
			buttons = append(buttons, map[string]any{
				"type": "URL",
				"text": button.Text,
				"url":  button.URL,
			})
		default:
			buttons = append(buttons, map[string]any{
				"type": "QUICK_REPLY",
				"text": button.Text,
			})
		}
	}

	if template.Category == entities.TemplateCategoryMarketing && !hasOptOut(buttons) {
		buttons = append(buttons, map[string]any{
			"type": "QUICK_REPLY",
			"text": OptOutButtonText,
		})
	}

	return buttons
}

func hasOptOut(buttons []map[string]any) bool {
	for _, button := range buttons {
		text, _ := button["text"].(string)
		if strings.EqualFold(strings.TrimSpace(text), OptOutButtonText) {
			return true
		}
	}
	return false
}

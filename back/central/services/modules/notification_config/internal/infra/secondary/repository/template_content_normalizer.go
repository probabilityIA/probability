package repository

import (
	"context"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

type templateButton struct {
	Text string
	Type string
}

type templateText struct {
	Header  string
	Body    string
	Footer  string
	Buttons []templateButton
}

const (
	optOutButtonLabel      = "Dejar de recibir"
	templateCategoryMarket = "MARKETING"
)

var fallbackTemplateBodies = map[string]string{
	"guia_envio_generada": "Hola {{1}} \U0001F44B\n" +
		"Somos {{2}}. Tu pedido {{3}} ya fue despachado \U0001F4E6\n\n" +
		"\U0001F4D1 Gu\u00eda: {{4}}\n" +
		"\U0001F69A Transportadora: {{5}}\n\n" +
		"Gracias por tu compra.",
	"confirmacion_pedido_contraentrega": "Hola {{1}} \U0001F44B\n" +
		"Recibimos tu pedido en {{2}}.\n\n" +
		"\U0001F9FE Pedido: {{3}}\n" +
		"\U0001F4CD Env\u00edo a: {{4}}\n" +
		"\U0001F6D2 Productos: {{5}}\n\n" +
		"\u00bfConfirmas tu pedido?",
}

var (
	legacyContentPattern    = regexp.MustCompile(`(?s)^Template: ([^,]+), Variables: map\[(.*)\]$`)
	plantillaVarsPattern    = regexp.MustCompile(`(?s)^Plantilla: (\S+) \(variables: map\[(.*)\]\)$`)
	plantillaNoVarsPattern  = regexp.MustCompile(`^Plantilla: (\S+)$`)
	legacyKeyRegex          = regexp.MustCompile(`(?:^|\s)(\d+):`)
	templatePlaceholderExpr = regexp.MustCompile(`\{\{\s*(\w+)\s*\}\}`)
)

func parseTemplateContent(templateName, content string) (string, map[string]string, bool) {
	trimmed := strings.TrimSpace(content)
	name := strings.TrimSpace(templateName)

	if match := legacyContentPattern.FindStringSubmatch(trimmed); match != nil {
		return firstNonEmpty(strings.TrimSpace(match[1]), name), parseLegacyVariables(match[2]), true
	}
	if match := plantillaVarsPattern.FindStringSubmatch(trimmed); match != nil {
		return firstNonEmpty(strings.TrimSpace(match[1]), name), parseLegacyVariables(match[2]), true
	}
	if match := plantillaNoVarsPattern.FindStringSubmatch(trimmed); match != nil {
		return firstNonEmpty(strings.TrimSpace(match[1]), name), map[string]string{}, true
	}

	if name == "" {
		return "", nil, false
	}
	if trimmed == name {
		return name, map[string]string{}, true
	}
	if rest, ok := strings.CutPrefix(trimmed, name+":"); ok {
		params := map[string]string{}
		for i, value := range strings.Split(rest, " | ") {
			params[strconv.Itoa(i+1)] = strings.TrimSpace(value)
		}
		return name, params, true
	}

	return "", nil, false
}

func renderTemplateText(text templateText, params map[string]string) string {
	fill := func(raw string) string {
		return templatePlaceholderExpr.ReplaceAllStringFunc(raw, func(token string) string {
			key := templatePlaceholderExpr.FindStringSubmatch(token)[1]
			if value, ok := params[key]; ok && value != "" {
				return value
			}
			return "-"
		})
	}

	parts := make([]string, 0, 3)
	if header := strings.TrimSpace(fill(text.Header)); header != "" {
		parts = append(parts, header)
	}
	if body := strings.TrimSpace(fill(text.Body)); body != "" {
		parts = append(parts, body)
	}
	if footer := strings.TrimSpace(fill(text.Footer)); footer != "" {
		parts = append(parts, footer)
	}
	return strings.Join(parts, "\n\n")
}

func resolveMessageContent(texts map[string]templateText, templateName, content string) string {
	rendered, _ := resolveMessage(texts, templateName, content)
	return rendered
}

func resolveMessage(texts map[string]templateText, templateName, content string) (string, []templateButton) {
	name, params, ok := parseTemplateContent(templateName, content)
	if !ok {
		if text, found := texts[strings.TrimSpace(templateName)]; found {
			return content, text.Buttons
		}
		return content, nil
	}

	text, found := texts[name]
	if !found {
		body, legacy := fallbackTemplateBodies[name]
		if !legacy {
			return content, nil
		}
		text = templateText{Body: body}
	}

	if rendered := renderTemplateText(text, params); rendered != "" {
		return rendered, text.Buttons
	}
	return content, text.Buttons
}

func templateNamesToResolve(pairs [][2]string) []string {
	seen := map[string]struct{}{}
	names := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		name, _, ok := parseTemplateContent(pair[0], pair[1])
		if !ok {
			name = strings.TrimSpace(pair[0])
		}
		if name == "" {
			continue
		}
		if _, dup := seen[name]; dup {
			continue
		}
		seen[name] = struct{}{}
		names = append(names, name)
	}
	return names
}

type templateTextRow struct {
	Name       string
	BusinessID *uint
	Category   string
	HeaderType string
	HeaderText string
	BodyText   string
	FooterText string
	Buttons    string
}

type storedTemplateButton struct {
	Text string `json:"Text"`
	Type string `json:"Type"`
}

func parseTemplateButtons(raw, category string) []templateButton {
	buttons := []templateButton{}
	var stored []storedTemplateButton
	if strings.TrimSpace(raw) != "" && json.Unmarshal([]byte(raw), &stored) == nil {
		for _, button := range stored {
			if text := strings.TrimSpace(button.Text); text != "" {
				buttons = append(buttons, templateButton{Text: text, Type: strings.ToUpper(strings.TrimSpace(button.Type))})
			}
		}
	}

	if strings.EqualFold(category, templateCategoryMarket) {
		hasOptOut := false
		for _, button := range buttons {
			if strings.EqualFold(button.Text, optOutButtonLabel) {
				hasOptOut = true
				break
			}
		}
		if !hasOptOut {
			buttons = append(buttons, templateButton{Text: optOutButtonLabel, Type: "QUICK_REPLY"})
		}
	}

	if len(buttons) == 0 {
		return nil
	}
	return buttons
}

func (q *messageAuditQuerier) loadTemplateTexts(ctx context.Context, businessID uint, names []string) map[string]templateText {
	texts := map[string]templateText{}
	if len(names) == 0 {
		return texts
	}

	var rows []templateTextRow
	err := q.db.Conn(ctx).
		Table("whatsapp_templates").
		Select("name, business_id, category, header_type, header_text, body_text, footer_text, COALESCE(buttons::text, '') AS buttons").
		Where("deleted_at IS NULL AND name IN ? AND (business_id = ? OR business_id IS NULL)", names, businessID).
		Order("business_id NULLS LAST, updated_at DESC").
		Scan(&rows).Error
	if err != nil {
		q.logger.Warn().Err(err).Msg("No se pudo cargar el texto de las plantillas: se muestra el contenido guardado")
		return texts
	}

	for _, row := range rows {
		if _, exists := texts[row.Name]; exists {
			continue
		}
		header := ""
		if strings.EqualFold(row.HeaderType, "TEXT") || row.HeaderType == "" {
			header = row.HeaderText
		}
		texts[row.Name] = templateText{
			Header:  header,
			Body:    row.BodyText,
			Footer:  row.FooterText,
			Buttons: parseTemplateButtons(row.Buttons, row.Category),
		}
	}

	return texts
}

func parseLegacyVariables(raw string) map[string]string {
	out := make(map[string]string)
	if raw == "" {
		return out
	}

	positions := legacyKeyRegex.FindAllStringIndex(raw, -1)
	if len(positions) == 0 {
		return out
	}

	for i, pos := range positions {
		keyStart := pos[0]
		valueStart := pos[1]
		valueEnd := len(raw)
		if i+1 < len(positions) {
			valueEnd = positions[i+1][0]
		}
		key := strings.TrimSpace(raw[keyStart : valueStart-1])
		value := strings.TrimSpace(raw[valueStart:valueEnd])
		if key != "" {
			out[key] = value
		}
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}
	return ""
}

package repository

import (
	"regexp"
	"strings"
)

var templateBodies = map[string]string{
	"guia_envio_generada": "Hola {{1}} 👋\n" +
		"Somos {{2}}. Tu pedido {{3}} ya fue despachado 📦\n\n" +
		"📑 Guía: {{4}}\n" +
		"🚚 Transportadora: {{5}}\n\n" +
		"Gracias por tu compra.",
	"confirmacion_pedido_contraentrega": "Hola {{1}} 👋\n" +
		"Recibimos tu pedido en {{2}}.\n\n" +
		"🧾 Pedido: {{3}}\n" +
		"📍 Envío a: {{4}}\n" +
		"🛒 Productos: {{5}}\n\n" +
		"¿Confirmas tu pedido?",
}

var legacyContentPattern = regexp.MustCompile(`^Template: (?P<name>[^,]+), Variables: map\[(?P<vars>.*)\]$`)
var legacyKeyRegex = regexp.MustCompile(`(?:^|\s)(\d+):`)

func normalizeMessageContent(templateName, content string) string {
	trimmed := strings.TrimSpace(content)
	if trimmed == "" {
		return content
	}

	match := legacyContentPattern.FindStringSubmatch(trimmed)
	if match == nil {
		return content
	}

	name := strings.TrimSpace(match[1])
	if name == "" {
		name = templateName
	}

	body, ok := templateBodies[name]
	if !ok {
		return content
	}

	variables := parseLegacyVariables(match[2])
	rendered := body
	for key, value := range variables {
		rendered = strings.ReplaceAll(rendered, "{{"+key+"}}", value)
	}
	return rendered
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

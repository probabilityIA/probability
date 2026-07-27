package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

type meliErrorCause struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type meliErrorBody struct {
	Message string           `json:"message"`
	Error   string           `json:"error"`
	Cause   []meliErrorCause `json:"cause"`
}

var meliErrorHints = map[string]string{
	"missing.fashion_grid.grid_id.values": "falta la guia de talles que Mercado Libre exige en categorias de moda (atributo GRID_ID)",
	"item.attributes.missing":             "faltan atributos obligatorios de la categoria",
	"item.price.invalid":                  "el precio no es valido para Mercado Libre",
	"item.title.invalid":                  "el titulo no es valido para Mercado Libre",
	"item.pictures.empty":                 "la publicacion necesita al menos una imagen",
}

func describeMeliError(status int, body []byte) string {
	var parsed meliErrorBody
	if err := json.Unmarshal(body, &parsed); err != nil {
		return fmt.Sprintf("Mercado Libre respondio %d: %s", status, truncateForMessage(string(body), 180))
	}

	seen := make(map[string]bool, len(parsed.Cause))
	reasons := make([]string, 0, len(parsed.Cause))
	for _, cause := range parsed.Cause {
		text := hintForCause(cause)
		if text == "" || seen[text] {
			continue
		}
		seen[text] = true
		reasons = append(reasons, text)
		if len(reasons) == 2 {
			break
		}
	}

	if len(reasons) == 0 {
		fallback := strings.TrimSpace(firstNonEmpty(parsed.Message, parsed.Error))
		if fallback == "" {
			return fmt.Sprintf("Mercado Libre respondio %d sin detalle", status)
		}
		return truncateForMessage(fallback, 180)
	}

	return truncateForMessage(strings.Join(reasons, " | "), 220)
}

func hintForCause(cause meliErrorCause) string {
	if hint, ok := meliErrorHints[cause.Code]; ok {
		return hint
	}
	message := strings.TrimSpace(cause.Message)
	if message != "" {
		return message
	}
	return strings.TrimSpace(cause.Code)
}

func truncateForMessage(value string, max int) string {
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	return string(runes[:max]) + "..."
}

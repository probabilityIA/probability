package templates

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

var (
	leadingPlaceholderPattern  = regexp.MustCompile(`^\{\{\d+\}\}`)
	trailingPlaceholderPattern = regexp.MustCompile(`\{\{\d+\}\}$`)
	adjacentPlaceholderPattern = regexp.MustCompile(`\{\{\d+\}\}\s*\{\{\d+\}\}`)
	emojiPattern               = regexp.MustCompile(`[\x{1F000}-\x{1FAFF}\x{2600}-\x{27BF}\x{2B00}-\x{2BFF}\x{FE0F}]`)
	wordPattern                = regexp.MustCompile(`[\p{L}\p{N}]+`)
)

func validateMetaRules(body, header, footer string, buttons []entities.TemplateButton, variableCount int) error {
	stripped := placeholderPattern.ReplaceAllString(body, " ")
	if strings.Contains(stripped, "{{") || strings.Contains(stripped, "}}") {
		return fmt.Errorf("las variables se escriben con numero entre llaves dobles: {{1}}, {{2}}")
	}
	if leadingPlaceholderPattern.MatchString(body) {
		return fmt.Errorf("el cuerpo no puede empezar con una variable: Meta lo rechaza, pon texto antes")
	}
	if trailingPlaceholderPattern.MatchString(body) {
		return fmt.Errorf("el cuerpo no puede terminar con una variable: Meta lo rechaza, pon texto o un punto despues")
	}
	if adjacentPlaceholderPattern.MatchString(body) {
		return fmt.Errorf("hay dos variables seguidas: Meta exige texto entre una variable y otra")
	}
	if variableCount > 0 {
		fixedWords := len(wordPattern.FindAllString(stripped, -1))
		if fixedWords < 2*variableCount+1 {
			return fmt.Errorf("el cuerpo tiene demasiadas variables para tan poco texto: con %d variable(s) escribe al menos %d palabras fijas", variableCount, 2*variableCount+1)
		}
	}

	if header != "" {
		if strings.Contains(header, "{{") {
			return fmt.Errorf("el encabezado no admite variables")
		}
		if strings.ContainsAny(header, "\n\r") {
			return fmt.Errorf("el encabezado no admite saltos de linea")
		}
		if emojiPattern.MatchString(header) {
			return fmt.Errorf("el encabezado no admite emojis: Meta lo rechaza")
		}
		if strings.ContainsAny(header, "*_~") {
			return fmt.Errorf("el encabezado no admite formato (*, _, ~)")
		}
	}

	if footer != "" {
		if strings.Contains(footer, "{{") {
			return fmt.Errorf("el pie no admite variables")
		}
		if strings.ContainsAny(footer, "\n\r") {
			return fmt.Errorf("el pie no admite saltos de linea")
		}
	}

	for _, button := range buttons {
		if strings.Contains(button.Text, "{{") {
			return fmt.Errorf("el boton %q no admite variables", button.Text)
		}
		if strings.ContainsAny(button.Text, "\n\r") {
			return fmt.Errorf("el boton %q no admite saltos de linea", button.Text)
		}
	}

	return nil
}

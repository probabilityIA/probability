package usecasemessaging

import (
	"regexp"
	"strings"
)

var templateParamWhitespace = regexp.MustCompile(`\s+`)

func SanitizeTemplateParam(value string) string {
	if value == "" {
		return ""
	}
	return strings.TrimSpace(templateParamWhitespace.ReplaceAllString(value, " "))
}

func SanitizeTemplateVariables(variables map[string]string) map[string]string {
	if variables == nil {
		return nil
	}
	sanitized := make(map[string]string, len(variables))
	for k, v := range variables {
		sanitized[k] = SanitizeTemplateParam(v)
	}
	return sanitized
}

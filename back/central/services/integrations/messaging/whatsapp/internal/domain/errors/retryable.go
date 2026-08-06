package errors

import (
	stderrors "errors"
	"strings"
)

var nonRetryablePhrases = []string{
	"key not found",
	"no se encontró integración",
	"credenciales no encontradas",
	"número de teléfono inválido",
	"phone_number_id no encontrado",
	"access_token no encontrado",
	"Required parameter is missing",
	"does not exist in",
	"131008",
	"132001",
	"131009",
	"132018",
}

func IsNonRetryable(err error) bool {
	if err == nil {
		return false
	}

	var templateNotFound *ErrTemplateNotFound
	if stderrors.As(err, &templateNotFound) {
		return true
	}

	var missingVar *ErrMissingVariable
	if stderrors.As(err, &missingVar) {
		return true
	}

	errMsg := err.Error()
	for _, phrase := range nonRetryablePhrases {
		if strings.Contains(errMsg, phrase) {
			return true
		}
	}

	return false
}

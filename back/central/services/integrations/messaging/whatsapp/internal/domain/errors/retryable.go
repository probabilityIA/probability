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
	"con número propio mal configurada",
	"waba_id no configurado",
	"Required parameter is missing",
	"does not exist in",
	"131008",
	"132001",
	"131009",
	"132018",
}

var nonRetryableMetaCodes = map[int]bool{
	190:    true,
	131026: true,
	131030: true,
	131031: true,
	132000: true,
	132001: true,
	132005: true,
	132012: true,
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

	var metaErr *MetaGraphError
	if stderrors.As(err, &metaErr) {
		if nonRetryableMetaCodes[metaErr.Code] {
			return true
		}
		if metaErr.StatusCode == 401 || metaErr.StatusCode == 403 {
			return true
		}
	}

	errMsg := err.Error()
	for _, phrase := range nonRetryablePhrases {
		if strings.Contains(errMsg, phrase) {
			return true
		}
	}

	return false
}

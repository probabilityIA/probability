package client

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	siigoerrors "github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/errors"
)

type siigoAPIError struct {
	Code    string   `json:"Code"`
	Message string   `json:"Message"`
	Params  []string `json:"Params"`
	Detail  string   `json:"Detail"`
}

type siigoAPIErrorEnvelope struct {
	Status    int             `json:"Status"`
	Errors    []siigoAPIError `json:"Errors"`
	ErrorsAlt []siigoAPIError `json:"errors"`
}

var indiceItemRegex = regexp.MustCompile(`items\[(\d+)\]`)

func parseSiigoErrors(body []byte) []siigoAPIError {
	var env siigoAPIErrorEnvelope
	if err := json.Unmarshal(body, &env); err != nil {
		return nil
	}
	if len(env.Errors) > 0 {
		return env.Errors
	}
	return env.ErrorsAlt
}

func errorSiigo(body []byte, statusCode int, operacion string) error {
	return siigoerrors.NewProviderError(
		siigoErrorCode(body, statusCode),
		siigoErrorMessage(body, statusCode, operacion),
	)
}

func siigoErrorCode(body []byte, statusCode int) string {
	errores := parseSiigoErrors(body)
	if len(errores) == 0 {
		if statusCode >= 500 {
			return "provider_unavailable"
		}
		return "provider_rejected"
	}

	e := errores[0]
	campo := ""
	if len(e.Params) > 0 {
		campo = e.Params[0]
	}

	switch e.Code {
	case "parameter_inactive":
		switch {
		case strings.HasPrefix(campo, "document"):
			return "document_inactive"
		case strings.HasPrefix(campo, "items"):
			return "product_inactive"
		case strings.HasPrefix(campo, "customer"):
			return "customer_inactive"
		case strings.HasPrefix(campo, "payments"):
			return "payment_method_inactive"
		case strings.HasPrefix(campo, "seller"):
			return "seller_inactive"
		}
		return "parameter_inactive"
	case "document_settings":
		if strings.HasPrefix(campo, "stamp") {
			return "document_not_electronic"
		}
		return "document_settings"
	case "invalid_reference":
		return "invalid_reference"
	case "invalid_code":
		return "invalid_product_code"
	case "documents_service":
		return "provider_unavailable"
	case "unauthorized", "invalid_token":
		return "invalid_credentials"
	case "quota_exceeded", "too_many_requests":
		return "rate_limited"
	}

	if e.Code != "" {
		return e.Code
	}
	return "provider_rejected"
}

func siigoErrorMessage(body []byte, statusCode int, operacion string) string {
	errores := parseSiigoErrors(body)
	if len(errores) == 0 {
		cuerpo := strings.TrimSpace(string(body))
		if cuerpo == "" {
			return fmt.Sprintf("Siigo rechazo %s (codigo %d)", operacion, statusCode)
		}
		return fmt.Sprintf("Siigo rechazo %s (codigo %d): %s", operacion, statusCode, cuerpo)
	}

	partes := make([]string, 0, len(errores))
	for _, e := range errores {
		partes = append(partes, explicarSiigoError(e))
	}
	return fmt.Sprintf("Siigo rechazo %s: %s", operacion, strings.Join(partes, "; "))
}

func explicarSiigoError(e siigoAPIError) string {
	campo := ""
	if len(e.Params) > 0 {
		campo = e.Params[0]
	}
	valor := valorDelMensaje(e.Message)

	switch e.Code {
	case "parameter_inactive":
		switch {
		case strings.HasPrefix(campo, "document"):
			return fmt.Sprintf("el tipo de documento configurado (id %s) esta inactivo en Siigo: activarlo en Siigo o seleccionar otro en la configuracion de facturacion", valorODefecto(valor, "sin id"))
		case strings.HasPrefix(campo, "items"):
			return fmt.Sprintf("el producto con codigo %s%s esta inactivo en Siigo: activarlo o quitarlo de la orden", valorODefecto(valor, "desconocido"), sufijoItem(campo))
		case strings.HasPrefix(campo, "customer"):
			return fmt.Sprintf("el cliente %s esta inactivo en Siigo: activarlo para poder facturarle", valorODefecto(valor, "de la factura"))
		case strings.HasPrefix(campo, "payments"):
			return fmt.Sprintf("el medio de pago configurado (id %s) esta inactivo en Siigo", valorODefecto(valor, "sin id"))
		case strings.HasPrefix(campo, "seller"):
			return fmt.Sprintf("el vendedor configurado (id %s) esta inactivo en Siigo", valorODefecto(valor, "sin id"))
		default:
			return fmt.Sprintf("un dato esta inactivo en Siigo (%s): %s", campoODefecto(campo), e.Message)
		}
	case "document_settings":
		if strings.HasPrefix(campo, "stamp") {
			return "el tipo de documento configurado no esta habilitado para facturacion electronica en Siigo: revisar en Siigo la configuracion del documento (resolucion DIAN y envio electronico)"
		}
		return fmt.Sprintf("la configuracion del documento en Siigo no permite esta operacion (%s)", campoODefecto(campo))
	case "invalid_reference":
		return fmt.Sprintf("un dato configurado no existe en Siigo (%s): %s", campoODefecto(campo), e.Message)
	case "invalid_code":
		return fmt.Sprintf("codigo invalido en Siigo (%s): %s", campoODefecto(campo), e.Message)
	case "documents_service":
		return "el servicio de documentos de Siigo no esta disponible en este momento: reintentar en unos minutos"
	case "unauthorized", "invalid_token":
		return "Siigo rechazo las credenciales de la integracion: revisar usuario y access key"
	case "quota_exceeded", "too_many_requests":
		return "se supero el limite de peticiones de la API de Siigo: reintentar en unos minutos"
	}

	if campo != "" {
		return fmt.Sprintf("%s (%s)", e.Message, campo)
	}
	return e.Message
}

func valorDelMensaje(mensaje string) string {
	idx := strings.LastIndex(mensaje, ": ")
	if idx == -1 {
		return ""
	}
	return strings.TrimSpace(mensaje[idx+2:])
}

func valorODefecto(valor, defecto string) string {
	if valor == "" {
		return defecto
	}
	return valor
}

func campoODefecto(campo string) string {
	if campo == "" {
		return "sin campo"
	}
	return campo
}

func sufijoItem(campo string) string {
	m := indiceItemRegex.FindStringSubmatch(campo)
	if len(m) != 2 {
		return ""
	}
	return fmt.Sprintf(" (item %s de la factura)", m[1])
}

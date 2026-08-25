package domain

import "strings"

const GenericGuideErrorMessage = "No fue posible generar la guía con la transportadora. Revisa los datos del destinatario e inténtalo de nuevo."

const GenericQuoteErrorMessage = "No fue posible cotizar el envío con la transportadora. Revisa los datos de la orden e inténtalo de nuevo."

const GenericCancelErrorMessage = "No fue posible cancelar el envío con la transportadora. Inténtalo de nuevo en unos minutos."

const GenericTrackingErrorMessage = "No fue posible consultar el estado del envío con la transportadora. Inténtalo de nuevo en unos minutos."

var providerBlocklist = []string{
	"envioclick",
	"envio click",
	"envioclik",
	"envio clik",
	"enviame",
	"envíame",
	"mipaquete",
	"mi paquete",
	"shipit",
	"api.envioclickpro",
	"envioclickpro",
}

type guideErrorRule struct {
	keywords []string
	message  string
}

var guideErrorRules = []guideErrorRule{
	{
		keywords: []string{"email", "correo", "e-mail"},
		message:  "Falta un correo electrónico o no tiene un formato válido. Revisa el correo del destinatario en la orden y el de tu dirección de origen, y vuelve a generar la guía.",
	},
	{
		keywords: []string{"phone", "teléfono", "telefono", "celular", "móvil", "movil"},
		message:  "Un teléfono no tiene el formato correcto (entre 7 y 10 dígitos). Revisa el teléfono del destinatario en la orden y el de tu dirección de origen, y vuelve a generar la guía.",
	},
	{
		keywords: []string{"suburb", "barrio"},
		message:  "El barrio está vacío o excede el largo permitido (entre 2 y 30 caracteres). Revísalo en la orden y en tu dirección de origen, y vuelve a generar la guía.",
	},
	{
		keywords: []string{"dane", "cobertura", "coverage", "postal"},
		message:  "La transportadora no tiene cobertura para la ciudad de destino o la ciudad no está bien identificada. Revisa la dirección de la orden.",
	},
	{
		keywords: []string{"saldo", "balance", "insufficient", "fondos", "credito", "crédito"},
		message:  "No hay saldo suficiente para generar la guía. Recarga tu billetera e inténtalo de nuevo.",
	},
	{
		keywords: []string{"contentvalue", "declared value", "valor declarado", "seguro", "insurance"},
		message:  "El valor declarado de la mercancía no es válido o es insuficiente para el seguro.",
	},
	{
		keywords: []string{"weight", "peso"},
		message:  "El peso del paquete no es válido o excede el límite de la transportadora.",
	},
	{
		keywords: []string{"dimension", "height", "width", "length", "alto", "ancho", "largo"},
		message:  "Las dimensiones del paquete no son válidas o exceden el límite de la transportadora.",
	},
	{
		keywords: []string{"address", "dirección", "direccion", "street"},
		message:  "Una dirección está incompleta o no es válida. Revisa la dirección de destino en la orden y tu dirección de origen, y vuelve a generar la guía.",
	},
	{
		keywords: []string{"unauthorized", "forbidden", "credential", "token", "401", "403"},
		message:  "No fue posible autenticar con la transportadora. Revisa la configuración de la integración de envíos.",
	},
	{
		keywords: []string{"timeout", "deadline", "temporarily", "503", "502", "504"},
		message:  "La transportadora no respondió a tiempo. Vuelve a intentar la generación de la guía en unos minutos.",
	},
	{
		keywords: []string{"faltan datos", "missing", "requerido", "required"},
		message:  "Faltan datos obligatorios del destinatario o del paquete para generar la guía.",
	},
}

func SanitizeGuideError(raw string) string {
	lower := strings.ToLower(strings.TrimSpace(raw))
	if lower == "" {
		return GenericGuideErrorMessage
	}

	for _, rule := range guideErrorRules {
		for _, kw := range rule.keywords {
			if strings.Contains(lower, kw) {
				return rule.message
			}
		}
	}

	return GenericGuideErrorMessage
}

func MentionsProvider(text string) bool {
	lower := strings.ToLower(text)
	for _, term := range providerBlocklist {
		if strings.Contains(lower, term) {
			return true
		}
	}
	return false
}

func SafeMessage(text, fallback string) string {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" || MentionsProvider(trimmed) {
		return fallback
	}
	if len(trimmed) > 240 {
		return fallback
	}
	return trimmed
}

func SafeQuoteMessage(text string) string {
	return SafeMessage(text, GenericGuideErrorMessage)
}

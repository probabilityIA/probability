package client

import (
	"context"
	"encoding/json"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/shared/httpclient"
	"github.com/secamc93/probability/back/central/shared/log"
)

const (
	DefaultBaseURL = "https://api.envioclickpro.com.co/api/v2"

	generateTimeout = 120 * time.Second
	readTimeout     = 30 * time.Second
)

type Client struct {
	httpClient *httpclient.Client
	log        log.ILogger
}

func New(logger log.ILogger) domain.IEnvioClickClient {
	logger.Info(context.Background()).Msg("\u1f50d Creating EnvioClick HTTP client")

	httpConfig := httpclient.HTTPClientConfig{
		BaseURL:    DefaultBaseURL,
		Timeout:    generateTimeout,
		RetryCount: 2,
		RetryWait:  3 * time.Second,
		Debug:      true,
	}

	httpClient := httpclient.New(httpConfig, logger)
	httpClient.SetHeader("Accept", "application/json")
	httpClient.SetHeader("Content-Type", "application/json")

	return &Client{
		httpClient: httpClient,
		log:        logger.WithModule("envioclick.client"),
	}
}

type envioClickErrorResponse struct {
	StatusMessages []struct {
		Error []json.RawMessage `json:"error"`
	} `json:"status_messages"`
}

type fieldIssue struct {
	Path    string
	Message string
}

var carrierFieldLabels = map[string]string{
	"destination.address":   "La direcci\u00f3n de destino",
	"destination.suburb":    "El barrio de destino",
	"destination.city":      "La ciudad de destino",
	"destination.firstName": "El nombre del destinatario",
	"destination.lastName":  "El apellido del destinatario",
	"destination.company":   "La empresa de destino",
	"destination.email":     "El correo del destinatario",
	"destination.phone":     "El tel\u00e9fono del destinatario",
	"destination.zipCode":   "El c\u00f3digo postal de destino",
	"origin.address":        "La direcci\u00f3n de origen",
	"origin.suburb":         "El barrio de origen",
	"origin.city":           "La ciudad de origen",
	"origin.firstName":      "El nombre del remitente",
	"origin.lastName":       "El apellido del remitente",
	"origin.company":        "La empresa de origen",
	"origin.email":          "El correo de la bodega de origen",
	"origin.phone":          "El tel\u00e9fono de la bodega de origen",
	"origin.zipCode":        "El c\u00f3digo postal de origen",
	"description":           "La descripci\u00f3n del contenido",
	"contentValue":          "El valor declarado de la mercanc\u00eda",
	"myShipmentReference":   "La referencia del env\u00edo",
	"externalOrderId":       "El identificador de la orden",
}

var providerMentions = []string{
	"www.envioclickpro.com.co",
	"envioclickpro.com.co",
	"envioclickpro",
	"envioclick pro",
	"envioclick",
}

func stripProviderMentions(text string) string {
	cleaned := text
	for _, mention := range providerMentions {
		for {
			idx := strings.Index(strings.ToLower(cleaned), mention)
			if idx < 0 {
				break
			}
			cleaned = cleaned[:idx] + cleaned[idx+len(mention):]
		}
	}
	cleaned = strings.ReplaceAll(cleaned, "  ", " ")
	cleaned = strings.Trim(cleaned, " .,:;-")
	return cleaned
}

func collectFieldIssues(raw json.RawMessage, prefix string, issues *[]fieldIssue) {
	var texto string
	if err := json.Unmarshal(raw, &texto); err == nil {
		if prefix != "" {
			*issues = append(*issues, fieldIssue{Path: prefix, Message: texto})
		}
		return
	}

	var lista []json.RawMessage
	if err := json.Unmarshal(raw, &lista); err == nil {
		for _, item := range lista {
			collectFieldIssues(item, prefix, issues)
		}
		return
	}

	var objeto map[string]json.RawMessage
	if err := json.Unmarshal(raw, &objeto); err == nil {
		claves := make([]string, 0, len(objeto))
		for clave := range objeto {
			claves = append(claves, clave)
		}
		sort.Strings(claves)
		for _, clave := range claves {
			hijo := clave
			if prefix != "" {
				hijo = prefix + "." + clave
			}
			collectFieldIssues(objeto[clave], hijo, issues)
		}
	}
}

func labelForField(path string) string {
	if label, ok := carrierFieldLabels[path]; ok {
		return label
	}
	if idx := strings.LastIndex(path, "."); idx >= 0 {
		if label, ok := carrierFieldLabels[path[idx+1:]]; ok {
			return label
		}
	}
	return ""
}

func lowerFirst(text string) string {
	runes := []rune(text)
	if len(runes) == 0 {
		return text
	}
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func describeFieldIssue(issue fieldIssue) string {
	mensaje := stripProviderMentions(issue.Message)
	if mensaje == "" {
		return ""
	}

	label := labelForField(issue.Path)
	if label == "" {
		return mensaje
	}

	if idx := strings.Index(mensaje, "El campo"); idx == 0 {
		return strings.TrimSpace(label + mensaje[len("El campo"):])
	}
	for _, marcador := range []string{"el campo", "El campo"} {
		if idx := strings.Index(mensaje, marcador); idx >= 0 {
			return strings.TrimSpace(mensaje[:idx] + lowerFirst(label) + mensaje[idx+len(marcador):])
		}
	}

	return label + ": " + mensaje
}

func parseEnvioClickError(body []byte) string {
	var errorResp envioClickErrorResponse
	if err := json.Unmarshal(body, &errorResp); err == nil {
		var issues []fieldIssue
		var sueltos []string

		for _, msg := range errorResp.StatusMessages {
			for _, raw := range msg.Error {
				var texto string
				if err := json.Unmarshal(raw, &texto); err == nil {
					sueltos = append(sueltos, texto)
					continue
				}
				collectFieldIssues(raw, "", &issues)
			}
		}

		if detalle := describeIssues(issues); detalle != "" {
			return detalle
		}
		if len(sueltos) > 0 {
			return mapEnvioClickError(strings.Join(sueltos, " "))
		}
	}

	return mapEnvioClickError(string(body))
}

func describeIssues(issues []fieldIssue) string {
	descripciones := make([]string, 0, len(issues))
	for _, issue := range issues {
		if texto := describeFieldIssue(issue); texto != "" {
			descripciones = append(descripciones, texto)
		}
	}

	if len(descripciones) == 0 {
		return ""
	}
	if len(descripciones) == 1 {
		return "Datos del env\u00edo incompletos: " + descripciones[0]
	}

	restantes := ""
	if len(descripciones) > 3 {
		sobran := len(descripciones) - 3
		if sobran == 1 {
			restantes = " (y 1 campo m\u00e1s)"
		} else {
			restantes = " (y " + strconv.Itoa(sobran) + " campos m\u00e1s)"
		}
		descripciones = descripciones[:3]
	}
	return "Datos del env\u00edo incompletos: " + strings.Join(descripciones, "; ") + restantes
}

func mapEnvioClickError(originalErr string) string {
	lowerErr := strings.ToLower(originalErr)

	if strings.Contains(lowerErr, "cr\u00e9dito") || strings.Contains(lowerErr, "credito") ||
		strings.Contains(lowerErr, "dep\u00f3sito") || strings.Contains(lowerErr, "deposito") ||
		strings.Contains(lowerErr, "saldo") {
		return "Saldo insuficiente para generar la gu\u00eda: recarga la billetera e intenta de nuevo"
	}
	if strings.Contains(lowerErr, "api key") || strings.Contains(lowerErr, "unauthorized") ||
		strings.Contains(lowerErr, "no autorizado") {
		return "Las credenciales de la transportadora no son v\u00e1lidas: revisa la configuraci\u00f3n de la integraci\u00f3n"
	}
	if (strings.Contains(lowerErr, "destination") || strings.Contains(lowerErr, "destino")) && strings.Contains(lowerErr, "dane") {
		return "El c\u00f3digo DANE de destino no es v\u00e1lido para esta transportadora"
	}
	if (strings.Contains(lowerErr, "origin") || strings.Contains(lowerErr, "origen")) && strings.Contains(lowerErr, "dane") {
		return "El c\u00f3digo DANE de origen no es v\u00e1lido para esta transportadora"
	}
	if strings.Contains(lowerErr, "contentvalue") || strings.Contains(lowerErr, "declared value") || strings.Contains(lowerErr, "valor") {
		return "El valor declarado de la mercanc\u00eda no es v\u00e1lido o es insuficiente para el seguro"
	}
	if strings.Contains(lowerErr, "weight") || strings.Contains(lowerErr, "peso") {
		return "El peso indicado no es v\u00e1lido o excede los l\u00edmites"
	}
	if strings.Contains(lowerErr, "dimensions") || strings.Contains(lowerErr, "height") || strings.Contains(lowerErr, "width") || strings.Contains(lowerErr, "length") {
		return "Dimensiones del paquete inv\u00e1lidas"
	}
	if strings.Contains(lowerErr, "phone") || strings.Contains(lowerErr, "tel\u00e9fono") || strings.Contains(lowerErr, "celular") {
		return "Formato de tel\u00e9fono incorrecto (debe tener entre 7 y 10 d\u00edgitos)"
	}
	if strings.Contains(lowerErr, "no se pudo generar la gu\u00eda") || strings.Contains(lowerErr, "no se pudo generar la guia") {
		return "La transportadora no pudo generar la gu\u00eda: revisa cobertura entre origen y destino"
	}
	if strings.Contains(lowerErr, "missing") || strings.Contains(lowerErr, "requerido") || strings.Contains(lowerErr, "falta") {
		return "Faltan datos obligatorios en la solicitud"
	}
	if strings.Contains(lowerErr, "unprocessed entity") || strings.Contains(lowerErr, "unprocessable") {
		return "La transportadora rechaz\u00f3 los datos del env\u00edo: revisa direcci\u00f3n, barrio y contacto de origen y destino"
	}

	limpio := stripProviderMentions(originalErr)
	if limpio == "" {
		return "La transportadora rechaz\u00f3 la solicitud"
	}
	return "La transportadora rechaz\u00f3 la solicitud: " + limpio
}

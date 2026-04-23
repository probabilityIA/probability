package entities

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
)

type TemplateDefinition struct {
	Name         string
	Language     string
	Variables    []string
	HasButtons   bool
	ButtonLabels []string
	Description  string
	Body         string
}

// Templates contiene el catálogo completo de las 11 plantillas aprobadas
var Templates = map[string]TemplateDefinition{
	"confirmacion_pedido_contraentrega": {
		Name:     "confirmacion_pedido_contraentrega",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_orden",
			"direccion",
			"productos",
		},
		HasButtons: true,
		ButtonLabels: []string{
			"Confirmar pedido",
			"No confirmar",
		},
		Description: "Plantilla inicial de confirmación de pedido contra entrega",
		Body: "Hola {{1}} 👋\n" +
			"Recibimos tu pedido en {{2}}.\n\n" +
			"🧾 Pedido: {{3}}\n" +
			"📍 Envío a: {{4}}\n" +
			"🛒 Productos: {{5}}\n\n" +
			"¿Confirmas tu pedido?",
	},
	"pedido_confirmado_v2": {
		Name:     "pedido_confirmado_v2",
		Language: "es",
		Variables: []string{
			"nombre",        // {{1}}
			"numero_pedido", // {{2}}
			"tienda",        // {{3}}
			"direccion",     // {{4}}
			"productos",     // {{5}}
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmación de que el pedido ha sido confirmado exitosamente (con resumen e iconos)",
	},
	"menu_no_confirmacion": {
		Name:     "menu_no_confirmacion",
		Language: "es",
		Variables: []string{
			"numero_pedido", // {{1}}
		},
		HasButtons: true,
		ButtonLabels: []string{
			"Presentar novedad",
			"Cancelar pedido",
			"Asesor",
		},
		Description: "Menú de opciones cuando el usuario no confirma el pedido",
	},
	"tipo_novedad_pedido": {
		Name:         "tipo_novedad_pedido",
		Language:     "es",
		Variables:    []string{}, // No tiene variables
		HasButtons:   true,
		ButtonLabels: []string{
			"Cambio de dirección",
			"Cambio de productos",
			"Cambio medio de pago",
		},
		Description: "Menú de selección del tipo de novedad",
	},
	"confirmar_cancelacion_pedido": {
		Name:     "confirmar_cancelacion_pedido",
		Language: "es",
		Variables: []string{
			"numero_pedido", // {{1}}
		},
		HasButtons: true,
		ButtonLabels: []string{
			"Sí, cancelar",
			"No, volver",
		},
		Description: "Confirmación antes de cancelar el pedido",
	},
	"motivo_cancelacion_pedido": {
		Name:         "motivo_cancelacion_pedido",
		Language:     "es",
		Variables:    []string{}, // No tiene variables
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Solicita al usuario que escriba el motivo de cancelación (texto libre)",
	},
	"pedido_cancelado": {
		Name:     "pedido_cancelado",
		Language: "es",
		Variables: []string{
			"numero_pedido", // {{1}}
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmación de que el pedido ha sido cancelado",
	},
	"novedad_cambio_direccion": {
		Name:         "novedad_cambio_direccion",
		Language:     "es",
		Variables:    []string{}, // No tiene variables
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmación de recepción de solicitud de cambio de dirección",
	},
	"novedad_cambio_productos": {
		Name:         "novedad_cambio_productos",
		Language:     "es",
		Variables:    []string{}, // No tiene variables
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmación de recepción de solicitud de cambio de productos",
	},
	"novedad_cambio_medio_pago": {
		Name:         "novedad_cambio_medio_pago",
		Language:     "es",
		Variables:    []string{}, // No tiene variables
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmación de recepción de solicitud de cambio de medio de pago",
	},
	"handoff_asesor": {
		Name:         "handoff_asesor",
		Language:     "es",
		Variables:    []string{}, // No tiene variables
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Mensaje de espera mientras se conecta con un asesor humano",
	},
	"prueba_conexion": {
		Name:         "prueba_conexion",
		Language:     "es",
		Variables:    []string{},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Mensaje de prueba para verificar conexión de WhatsApp",
	},
	"alerta_servidor": {
		Name:     "alerta_servidor",
		Language: "es",
		Variables: []string{
			"tipo_alerta", // {{1}} ej: "RAM"
			"descripcion", // {{2}} ej: "87.3% - supera umbral de 85%"
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Alerta de monitoreo del servidor para el administrador",
	},
	"guia_envio_generada": {
		Name:     "guia_envio_generada",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_pedido",
			"numero_guia",
			"transportadora",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Notificación de guía de envío generada con datos de tracking",
		Body: "Hola {{1}} 👋\n" +
			"Somos {{2}}. Tu pedido {{3}} ya fue despachado 📦\n\n" +
			"📑 Guía: {{4}}\n" +
			"🚚 Transportadora: {{5}}\n\n" +
			"Gracias por tu compra.",
	},
}

func RenderTemplateBody(templateName string, variables map[string]string) string {
	tpl, ok := Templates[templateName]
	if !ok || tpl.Body == "" {
		return ""
	}
	body := tpl.Body
	for key, value := range variables {
		body = strings.ReplaceAll(body, "{{"+key+"}}", value)
	}
	return body
}

// GetTemplateDefinition retorna la definición de una plantilla por su nombre
func GetTemplateDefinition(templateName string) (TemplateDefinition, bool) {
	template, exists := Templates[templateName]
	return template, exists
}

// ValidateTemplateVariables verifica que se provean todas las variables requeridas
func ValidateTemplateVariables(templateName string, providedVars map[string]string) error {
	template, exists := Templates[templateName]
	if !exists {
		return &errors.ErrTemplateNotFound{TemplateName: templateName}
	}

	for i, varName := range template.Variables {
		varKey := string(rune('1' + i)) // "1", "2", "3", etc
		if _, ok := providedVars[varKey]; !ok {
			return &errors.ErrMissingVariable{
				TemplateName: templateName,
				VariableName: varName,
				VariableKey:  varKey,
			}
		}
	}

	return nil
}

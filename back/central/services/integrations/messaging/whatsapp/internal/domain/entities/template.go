package entities

import (
	"github.com/secamc93/probability/back/central/services/integrations/messaging/whatsapp/internal/domain/errors"
)

// TemplateDefinition define la estructura de una plantilla de WhatsApp
type TemplateDefinition struct {
	Name         string   // Nombre de la plantilla en Meta
	Language     string   // Código de idioma (ej: "es", "en")
	Variables    []string // Lista de variables {{1}}, {{2}}, etc
	HasButtons   bool     // Si la plantilla tiene botones
	ButtonLabels []string // Etiquetas de los botones Quick Reply
	Description  string   // Descripción de la plantilla
}

// Templates contiene el catálogo completo de las 11 plantillas aprobadas
var Templates = map[string]TemplateDefinition{
	"confirmacion_pedido_contraentrega": {
		Name:     "confirmacion_pedido_contraentrega",
		Language: "es",
		Variables: []string{
			"nombre",           // {{1}}
			"tienda",           // {{2}}
			"numero_orden",     // {{3}}
			"direccion",        // {{4}}
			"productos",        // {{5}}
		},
		HasButtons: true,
		ButtonLabels: []string{
			"Confirmar pedido",
			"No confirmar",
		},
		Description: "Plantilla inicial de confirmación de pedido contra entrega",
	},
	"pedido_confirmado": {
		Name:     "pedido_confirmado",
		Language: "es",
		Variables: []string{
			"numero_pedido", // {{1}}
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmación de que el pedido ha sido confirmado exitosamente",
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

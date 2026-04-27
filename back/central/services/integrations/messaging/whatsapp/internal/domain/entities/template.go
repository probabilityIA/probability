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

var Templates = map[string]TemplateDefinition{
	"confirmacion_pedido_contraentrega": {
		Name:     "confirmacion_pedido_contraentrega",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_orden",
			"direccion",
			"ciudad",
			"departamento",
			"productos",
			"metodo_pago",
			"valor_recaudar",
		},
		HasButtons:   true,
		ButtonLabels: []string{"Confirmar pedido", "No confirmar"},
		Description:  "Confirmacion de pedido con direccion desglosada, metodo de pago y valor a recaudar",
	},
	"pedido_confirmado_v2": {
		Name:     "pedido_confirmado_v2",
		Language: "es",
		Variables: []string{
			"nombre",
			"numero_pedido",
			"tienda",
			"direccion",
			"productos",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de que el pedido ha sido confirmado exitosamente",
	},
	"menu_no_confirmacion": {
		Name:     "menu_no_confirmacion",
		Language: "es",
		Variables: []string{
			"numero_pedido",
		},
		HasButtons:   true,
		ButtonLabels: []string{"Presentar novedad", "Cancelar pedido", "Asesor"},
		Description:  "Menu de opciones cuando el usuario no confirma el pedido",
	},
	"tipo_novedad_pedido": {
		Name:         "tipo_novedad_pedido",
		Language:     "es",
		Variables:    []string{},
		HasButtons:   true,
		ButtonLabels: []string{"Cambio de dirección", "Cambio de productos", "Cambio medio de pago"},
		Description:  "Menu de seleccion del tipo de novedad",
	},
	"confirmar_cancelacion_pedido": {
		Name:     "confirmar_cancelacion_pedido",
		Language: "es",
		Variables: []string{
			"numero_pedido",
		},
		HasButtons:   true,
		ButtonLabels: []string{"Sí, cancelar", "No, volver"},
		Description:  "Confirmacion antes de cancelar el pedido",
	},
	"motivo_cancelacion_pedido": {
		Name:         "motivo_cancelacion_pedido",
		Language:     "es",
		Variables:    []string{},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Solicita al usuario que escriba el motivo de cancelacion (texto libre)",
	},
	"pedido_cancelado": {
		Name:     "pedido_cancelado",
		Language: "es",
		Variables: []string{
			"numero_pedido",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de que el pedido ha sido cancelado",
	},
	"novedad_cambio_direccion": {
		Name:         "novedad_cambio_direccion",
		Language:     "es",
		Variables:    []string{},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de recepcion de solicitud de cambio de direccion",
	},
	"novedad_cambio_productos": {
		Name:         "novedad_cambio_productos",
		Language:     "es",
		Variables:    []string{},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de recepcion de solicitud de cambio de productos",
	},
	"novedad_cambio_medio_pago": {
		Name:         "novedad_cambio_medio_pago",
		Language:     "es",
		Variables:    []string{},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de recepcion de solicitud de cambio de medio de pago",
	},
	"handoff_asesor": {
		Name:         "handoff_asesor",
		Language:     "es",
		Variables:    []string{},
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
		Description:  "Mensaje de prueba para verificar conexion de WhatsApp",
	},
	"alerta_servidor": {
		Name:     "alerta_servidor",
		Language: "es",
		Variables: []string{
			"tipo_alerta",
			"descripcion",
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
			"link_rastreo",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Notificacion de guia de envio generada (pago online)",
	},
	"guia_envio_generada_cod": {
		Name:     "guia_envio_generada_cod",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_pedido",
			"numero_guia",
			"transportadora",
			"valor_recaudar",
			"link_rastreo",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Notificacion de guia de envio generada para contra entrega",
	},
	"pedido_en_reparto": {
		Name:     "pedido_en_reparto",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_pedido",
			"numero_guia",
			"transportadora",
			"link_rastreo",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Notificacion al cliente cuando su pedido esta en reparto (pago online)",
	},
	"pedido_en_reparto_cod": {
		Name:     "pedido_en_reparto_cod",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_pedido",
			"numero_guia",
			"transportadora",
			"valor_recaudar",
			"link_rastreo",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Notificacion al cliente cuando su pedido esta en reparto (contra entrega)",
	},
	"pedido_entregado": {
		Name:     "pedido_entregado",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_pedido",
			"direccion",
			"ciudad",
			"departamento",
			"productos",
			"numero_guia",
			"transportadora",
			"link_rastreo",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de entrega exitosa (pago online)",
	},
	"pedido_entregado_cod": {
		Name:     "pedido_entregado_cod",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_pedido",
			"direccion",
			"ciudad",
			"departamento",
			"productos",
			"metodo_pago",
			"numero_guia",
			"transportadora",
			"valor_recaudado",
			"link_rastreo",
		},
		HasButtons:   false,
		ButtonLabels: []string{},
		Description:  "Confirmacion de entrega exitosa para contra entrega",
	},
	"confirmacion_pedido": {
		Name:     "confirmacion_pedido",
		Language: "es",
		Variables: []string{
			"nombre",
			"tienda",
			"numero_orden",
			"direccion",
			"ciudad",
			"departamento",
			"productos",
		},
		HasButtons:   true,
		ButtonLabels: []string{"Confirmar pedido", "No confirmar"},
		Description:  "Confirmacion de pedido para pago online (sin metodo de pago ni valor a recaudar)",
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

func GetTemplateDefinition(templateName string) (TemplateDefinition, bool) {
	template, exists := Templates[templateName]
	return template, exists
}

func ValidateTemplateVariables(templateName string, providedVars map[string]string) error {
	template, exists := Templates[templateName]
	if !exists {
		return &errors.ErrTemplateNotFound{TemplateName: templateName}
	}

	for i, varName := range template.Variables {
		varKey := string(rune('1' + i))
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

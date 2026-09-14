package repository

import (
	"context"
	"fmt"
)

type systemTemplateText struct {
	Header  string
	Body    string
	Footer  string
	Buttons string
}

var systemTemplateTexts = map[string]systemTemplateText{
	"confirmacion_pedido": {
		Body:    "Hola {{1}}, tu pedido en {{2}} ha sido recibido.\n\n\U0001f9fe Orden: {{3}}\n\U0001f4cd Direcci\u00f3n: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n\n\u00bfConfirmas tu pedido?",
		Buttons: "[{\"Text\": \"Confirmar pedido\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"No confirmar\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}]",
	},
	"confirmacion_pedido_contraentrega": {
		Body:    "Hola {{1}}, tu pedido en {{2}} ha sido recibido.\n\n\U0001f9feOrden: {{3}}\n\U0001f4cdDirecci\u00f3n: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n- m\u00e9todo de pago: {{8}} \U0001f69a\n- valor a recaudar: {{9}}\n\n\u00bfConfirmas tu pedido?",
		Buttons: "[{\"Text\": \"Confirmar pedido\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"No confirmar\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}]",
	},
	"confirmacion_pedido_contraentrega_sin_valor": {
		Body:    "Hola {{1}}, tu pedido en {{2}} ha sido recibido.\n\n\U0001f9feOrden: {{3}}\n\U0001f4cdDirecci\u00f3n: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n- m\u00e9todo de pago: {{8}} \U0001f69a\n\n\u00bfConfirmas tu pedido?",
		Buttons: "[{\"Text\": \"Confirmar pedido\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"No confirmar\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}]",
	},
	"confirmar_cancelacion_pedido": {
		Body:    "\u00bfEst\u00e1s seguro de que deseas cancelar el pedido {{1}}? Esta acci\u00f3n no se puede deshacer.",
		Buttons: "[{\"Text\": \"S\u00ed, cancelar\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"No, volver\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}]",
	},
	"guia_envio_generada": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\n\nTu pedido en {{2}} ya tiene gu\u00eda de env\u00edo asignada. \U0001f4e6\U0001f69a\n\n\U0001f4cb N\u00famero de pedido: {{3}}\n\U0001f522 N\u00famero de gu\u00eda: {{4}}\n\U0001f69b Transportadora: {{5}}\n\nPuedes hacer seguimiento de tu env\u00edo aqu\u00ed: {{6}}\n\n\u00a1Pronto lo recibir\u00e1s!",
	},
	"guia_envio_generada_cod": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\n\nTu pedido en {{2}} ya tiene gu\u00eda de env\u00edo asignada. \U0001f4e6\U0001f69a\n\n\U0001f4cb N\u00famero de pedido: {{3}}\n\U0001f522 N\u00famero de gu\u00eda: {{4}}\n\U0001f69b Transportadora: {{5}}\n\U0001f4b2Valor a recaudar: {{6}}\n\nPuedes hacer seguimiento de tu env\u00edo aqu\u00ed: {{7}}\n\n\u00a1Pronto lo recibir\u00e1s!",
	},
	"handoff_asesor": {
		Body: "Te estamos conectando con un asesor. Por favor espera un momento, pronto te atenderemos.",
	},
	"menu_no_confirmacion": {
		Body:    "Entendemos que no deseas confirmar el pedido {{1}}. \u00bfQu\u00e9 te gustar\u00eda hacer?",
		Buttons: "[{\"Text\": \"Presentar novedad\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"Cancelar pedido\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"Asesor\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}]",
	},
	"motivo_cancelacion_pedido": {
		Body: "Por favor, cu\u00e9ntanos el motivo por el cual deseas cancelar tu pedido. Escribe tu respuesta a continuaci\u00f3n:",
	},
	"novedad_cambio_direccion": {
		Body: "Hemos recibido tu solicitud de cambio de direcci\u00f3n. Nuestro equipo la procesar\u00e1 a la brevedad.",
	},
	"novedad_cambio_medio_pago": {
		Body: "Hemos recibido tu solicitud de cambio de medio de pago. Nuestro equipo la procesar\u00e1 a la brevedad.",
	},
	"novedad_cambio_productos": {
		Body: "Hemos recibido tu solicitud de cambio de productos. Nuestro equipo la revisar\u00e1 y te contactaremos pronto.",
	},
	"pedido_cancelado": {
		Body: "Tu pedido {{1}} ha sido cancelado. Si tienes alguna duda, cont\u00e1ctanos.",
	},
	"pedido_confirmado_v2": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\nTu pedido ha sido *confirmado* exitosamente \u2705\n\n\U0001f4cb *Resumen de tu pedido:*\n\U0001f6d2 *Pedido:* {{2}}\n\U0001f3ea *Tienda:* {{3}}\n\U0001f4cd *Direcci\u00f3n de entrega:* {{4}}\n\U0001f4e6 *Productos:* {{5}}\n\nPronto estar\u00e1 en camino \U0001f69a\n\u00a1Gracias por tu compra! \U0001f64f",
	},
	"pedido_en_reparto": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\n\nTu pedido en {{2}} ya est\u00e1 en reparto \U0001f4e6\U0001f69a\n\n\U0001f4cb N\u00famero de pedido: {{3}}\n\U0001f522 N\u00famero de gu\u00eda: {{4}}\n\U0001f69b Transportadora: {{5}}\n\nPuedes hacer seguimiento de tu env\u00edo aqu\u00ed: {{6}}\n\n\u00a1Pronto lo recibir\u00e1s!",
	},
	"pedido_en_reparto_cod": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\n\nTu pedido en {{2}} ya est\u00e1 en reparto \U0001f4e6\U0001f69a\n\n\U0001f4cb N\u00famero de pedido: {{3}}\n\U0001f522 N\u00famero de gu\u00eda: {{4}}\n\U0001f69b Transportadora: {{5}}\n\U0001f4b2Valor a recaudar: {{6}}\n\nPuedes hacer seguimiento de tu env\u00edo aqu\u00ed: {{7}}\n\n\u00a1Pronto lo recibir\u00e1s!",
	},
	"pedido_entregado": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\n\nTu pedido en {{2}} fue entregado\n\n\U0001f9feOrden: {{3}}\n\U0001f4cdDirecci\u00f3n: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n\U0001f522 N\u00famero de gu\u00eda: {{8}}\n\U0001f69b Transportadora: {{9}}\n\nConsulta el detalle de tu env\u00edo aqu\u00ed: {{10}}\n\n\u00a1Gracias por tu compra!",
	},
	"pedido_entregado_cod": {
		Body: "\u00a1Hola {{1}}! \U0001f44b\n\nTu pedido en {{2}} fue entregado\n\n\U0001f9feOrden: {{3}}\n\U0001f4cdDirecci\u00f3n: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n- M\u00e9todo de pago: {{8}} \U0001f69a\n\U0001f522 N\u00famero de gu\u00eda: {{9}}\n\U0001f69b Transportadora: {{10}}\n\U0001f4b2Valor recaudado: {{11}}\n\nConsulta el detalle de tu env\u00edo aqu\u00ed: {{12}}\n\n\u00a1Gracias por tu compra!",
	},
	"tipo_novedad_pedido": {
		Body:    "Selecciona el tipo de novedad que deseas reportar para tu pedido:",
		Buttons: "[{\"Text\": \"Cambio de direcci\u00f3n\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"Cambio de productos\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}, {\"Text\": \"Cambio medio de pago\", \"Type\": \"QUICK_REPLY\", \"URL\": \"\"}]",
	},
	"alerta_servidor": {
		Body: "\u26a0\ufe0f Alerta del servidor - Tipo: {{1}}\n\nDetalle: {{2}}\n\nPor favor revisa el estado del sistema.",
	},
	"recuperacion_codigo": {
		Body:    "Tu c\u00f3digo de verificaci\u00f3n es *{{1}}*. Por tu seguridad, no lo compartas.",
		Footer:  "Vence en 10 minutos.",
		Buttons: "[{\"Text\": \"Copiar c\u00f3digo\", \"Type\": \"URL\", \"URL\": \"https://www.whatsapp.com/otp/code/?otp_type=COPY_CODE&code_expiration_minutes=10&code=otp{{1}}\"}]",
	},
	"reporte_saldo_billetera": {
		Body: "Hola {{1}}, este es tu reporte del d\u00eda en Probability\n\nSaldo disponible en Billetera: {{2}}\n\nRecuerda recargar para que no te tome por sorpresa al momento de generar tus gu\u00edas \U0001f4a1\nFeliz d\u00eda!!",
	},
	"reporte_saldo_billetera_v2": {
		Body: "Hola {{1}}, este es tu reporte del d\u00eda en Probability \U0001f4ca. Tu saldo disponible en Billetera es {{2}} pesos.",
	},
	"resumen_pago_suscripcion": {
		Body: "Hola {{1}}, este es el resumen de tu suscripcion en Probability. Tu periodo de pago vence el {{2}}. Valor del ciclo: {{3}} pesos. Saldo disponible en tu Billetera en este momento: {{4}} pesos.",
	},
	"prueba_conexion": {
		Body: "\u00a1Hola! Este es un mensaje de prueba de Probability. Tu conexi\u00f3n de WhatsApp est\u00e1 funcionando correctamente. \u2705",
	},
	"hello_world": {
		Header: "Hello World",
		Body:   "Welcome and congratulations!! This message demonstrates your ability to send a WhatsApp message notification from the Cloud API, hosted by Meta. Thank you for taking the time to test with us.",
		Footer: "WhatsApp Business Platform sample message",
	},
}

func (r *Repository) migrateSystemTemplateBodies(ctx context.Context) error {
	conn := r.db.Conn(ctx)

	for name, text := range systemTemplateTexts {
		if err := conn.Exec(`
			UPDATE whatsapp_templates
			SET body_text = ?, header_text = ?, footer_text = ?, updated_at = NOW()
			WHERE business_id IS NULL AND name = ? AND deleted_at IS NULL AND btrim(body_text) = ''`,
			text.Body, text.Header, text.Footer, name,
		).Error; err != nil {
			return fmt.Errorf("failed to fill system template body %s: %w", name, err)
		}

		if text.Buttons == "" {
			continue
		}

		if err := conn.Exec(`
			UPDATE whatsapp_templates
			SET buttons = CAST(? AS jsonb), updated_at = NOW()
			WHERE business_id IS NULL AND name = ? AND deleted_at IS NULL
			  AND (buttons IS NULL OR buttons::text IN ('null', '[]'))`,
			text.Buttons, name,
		).Error; err != nil {
			return fmt.Errorf("failed to fill system template buttons %s: %w", name, err)
		}
	}

	return nil
}

package entities

var templateBodies = map[string]string{
	"guia_envio_generada": "¡Hola {{1}}! 👋\n\nTu pedido en {{2}} ya tiene guía de envío asignada. 📦🚚\n\n📋 Número de pedido: {{3}}\n🔢 Número de guía: {{4}}\n🚛 Transportadora: {{5}}\n\nPuedes hacer seguimiento de tu envío aquí: {{6}}\n\n¡Pronto lo recibirás!",

	"guia_envio_generada_cod": "¡Hola {{1}}! 👋\n\nTu pedido en {{2}} ya tiene guía de envío asignada. 📦🚚\n\n📋 Número de pedido: {{3}}\n🔢 Número de guía: {{4}}\n🚛 Transportadora: {{5}}\n💲Valor a recaudar: {{6}}\n\nPuedes hacer seguimiento de tu envío aquí: {{7}}\n\n¡Pronto lo recibirás!",

	"pedido_en_reparto": "¡Hola {{1}}! 👋\n\nTu pedido en {{2}} ya está en reparto 📦🚚\n\n📋 Número de pedido: {{3}}\n🔢 Número de guía: {{4}}\n🚛 Transportadora: {{5}}\n\nPuedes hacer seguimiento de tu envío aquí: {{6}}\n\n¡Pronto lo recibirás!",

	"pedido_en_reparto_cod": "¡Hola {{1}}! 👋\n\nTu pedido en {{2}} ya está en reparto 📦🚚\n\n📋 Número de pedido: {{3}}\n🔢 Número de guía: {{4}}\n🚛 Transportadora: {{5}}\n💲Valor a recaudar: {{6}}\n\nPuedes hacer seguimiento de tu envío aquí: {{7}}\n\n¡Pronto lo recibirás!",

	"pedido_entregado": "¡Hola {{1}}! 👋\n\nTu pedido en {{2}} fue entregado\n\n🧾Orden: {{3}}\n📍Dirección: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n🔢 Número de guía: {{8}}\n🚛 Transportadora: {{9}}\n\nConsulta el detalle de tu envío aquí: {{10}}\n\n¡Gracias por tu compra!",

	"pedido_entregado_cod": "¡Hola {{1}}! 👋\n\nTu pedido en {{2}} fue entregado\n\n🧾Orden: {{3}}\n📍Dirección: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n- Método de pago: {{8}} 🚚\n🔢 Número de guía: {{9}}\n🚛 Transportadora: {{10}}\n💲Valor recaudado: {{11}}\n\nConsulta el detalle de tu envío aquí: {{12}}\n\n¡Gracias por tu compra!",

	"confirmacion_pedido": "Hola {{1}}, tu pedido en {{2}} ha sido recibido.\n\n🧾 Orden: {{3}}\n📍 Dirección: {{4}}, Ciudad: {{5}},\nDepartamento: {{6}}, Colombia\nProductos:\n{{7}}\n\n¿Confirmas tu pedido?",
}

func init() {
	for name, body := range templateBodies {
		if tpl, ok := Templates[name]; ok {
			tpl.Body = body
			Templates[name] = tpl
		}
	}
}

package navigation

type subpage struct {
	Key          string
	Parent       string
	Label        string
	Route        string
	Description  string
	AlsoRequires string
}

var subpages = []subpage{
	{Key: "shipments.cod", Parent: "shipments", Label: "Recaudo contra entrega", Route: "/shipments/cod", Description: "Recaudo de las transportadoras y cortes de pago"},
	{Key: "shipments.quotes", Parent: "shipments", Label: "Cotizaciones", Route: "/shipments/quotes", Description: "Cotizaciones de env\u00edo de las tiendas y el panel"},
	{Key: "delivery.drivers", Parent: "delivery", Label: "Conductores", Route: "/delivery/drivers", Description: "Conductores propios"},
	{Key: "delivery.vehicles", Parent: "delivery", Label: "Veh\u00edculos", Route: "/delivery/vehicles", Description: "Veh\u00edculos propios"},
	{Key: "delivery.geozones", Parent: "delivery", Label: "Geozonas", Route: "/delivery/geozones", Description: "Zonas de cobertura en el mapa"},
	{Key: "inventory.movements", Parent: "inventory", Label: "Movimientos", Route: "/inventory/movements", Description: "Entradas y salidas de inventario"},
	{Key: "invoicing.configs", Parent: "invoicing", Label: "Configuraci\u00f3n de facturaci\u00f3n", Route: "/invoicing/configs", Description: "Proveedor y facturaci\u00f3n autom\u00e1tica por integraci\u00f3n"},
	{Key: "wallet.finanzas", Parent: "wallet", Label: "Finanzas", Route: "/wallet/finanzas", Description: "Resumen financiero y utilidad de env\u00edos"},
	{Key: "integrations.hub", Parent: "integrations", AlsoRequires: "orders", Label: "Tus Integraciones", Route: "/orders", Description: "Tus canales conectados: estado, \u00f3rdenes, productos y sincronizaci\u00f3n"},
	{Key: "integrations.hub.products", Parent: "integrations", AlsoRequires: "orders", Label: "Comparar productos", Route: "/orders", Description: "Qu\u00e9 producto est\u00e1 en cada canal y cu\u00e1l falta por publicar"},
	{Key: "integrations.hub.data", Parent: "integrations", AlsoRequires: "orders", Label: "Actualizar productos", Route: "/orders", Description: "Datos del canal que pueden entrar a Probability"},
	{Key: "integrations.hub.inventory", Parent: "integrations", AlsoRequires: "orders", Label: "Sincronizar inventario", Route: "/orders", Description: "Stock de cada canal frente al de Probability"},
	{Key: "integrations.hub.orders", Parent: "integrations", AlsoRequires: "orders", Label: "Comparar \u00f3rdenes", Route: "/orders", Description: "\u00d3rdenes del canal que no est\u00e1n en Probability"},
}

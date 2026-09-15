package navigation

type subpage struct {
	Key         string
	Parent      string
	Label       string
	Route       string
	Description string
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
}

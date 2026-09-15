package navigation

type highlight struct {
	Target string
	Title  string
	Hint   string
}

type subpage struct {
	Key          string
	Parent       string
	Label        string
	Route        string
	Description  string
	AlsoRequires string
	Highlight    *highlight
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
	{
		Key: "integrations.hub", Parent: "integrations", AlsoRequires: "orders", Label: "Tus Integraciones", Route: "/orders",
		Description: "Tus canales conectados: estado, \u00f3rdenes, productos y sincronizaci\u00f3n",
		Highlight: &highlight{
			Target: `[data-tour="integrations.hub-button"]`,
			Title:  "Aqu\u00ed est\u00e1n tus integraciones",
			Hint:   "Este bot\u00f3n abre Tus Integraciones: tus canales conectados, su estado y la sincronizaci\u00f3n de productos, inventario y \u00f3rdenes.",
		},
	},
	{Key: "integrations.hub.products", Parent: "integrations", AlsoRequires: "orders", Label: "Comparar productos", Route: "/orders", Description: "Qu\u00e9 producto est\u00e1 en cada canal y cu\u00e1l falta por publicar"},
	{Key: "integrations.hub.data", Parent: "integrations", AlsoRequires: "orders", Label: "Actualizar productos", Route: "/orders", Description: "Datos del canal que pueden entrar a Probability"},
	{Key: "integrations.hub.inventory", Parent: "integrations", AlsoRequires: "orders", Label: "Sincronizar inventario", Route: "/orders", Description: "Stock de cada canal frente al de Probability"},
	{Key: "integrations.hub.orders", Parent: "integrations", AlsoRequires: "orders", Label: "Comparar \u00f3rdenes", Route: "/orders", Description: "\u00d3rdenes del canal que no est\u00e1n en Probability"},
	{
		Key: "orders.generate_guide", Parent: "orders", AlsoRequires: "shipments", Label: "Generar gu\u00eda de env\u00edo", Route: "/orders",
		Description: "Bot\u00f3n \u00abVer recomendaci\u00f3n IA\u00bb en cada orden",
		Highlight: &highlight{
			Target: `[data-tour="orders.row.generate-guide"]:not([disabled])`,
			Title:  "Genera la gu\u00eda desde aqu\u00ed",
			Hint:   "Toca este bot\u00f3n en la orden que quieres enviar. Se abre \u00abGenerar Gu\u00eda de Env\u00edo\u00bb con 4 pasos: origen y destino, cotizaci\u00f3n, detalles y pago con tu billetera.",
		},
	},
	{
		Key: "orders.create", Parent: "orders", Label: "Nueva orden", Route: "/orders",
		Description: "Crear una orden a mano",
		Highlight: &highlight{
			Target: `[data-tour="orders.create"]`,
			Title:  "Crea una orden a mano",
			Hint:   "Este bot\u00f3n abre el formulario de la orden: cliente, direcci\u00f3n, productos, pago y si es contra entrega.",
		},
	},
	{
		Key: "orders.bulk_upload", Parent: "orders", Label: "Carga masiva de \u00f3rdenes", Route: "/orders",
		Description: "Subir varias \u00f3rdenes con un archivo",
		Highlight: &highlight{
			Target: `[data-tour="orders.bulk-upload"]`,
			Title:  "Sube varias \u00f3rdenes a la vez",
			Hint:   "Este bot\u00f3n abre la carga masiva: subes un archivo CSV o Excel con tus \u00f3rdenes.",
		},
	},
	{
		Key: "orders.bulk_guides", Parent: "orders", AlsoRequires: "shipments", Label: "Generaci\u00f3n masiva de gu\u00edas", Route: "/orders",
		Description: "Generar gu\u00edas de varias \u00f3rdenes a la vez",
		Highlight: &highlight{
			Target: `[data-tour="orders.bulk-guides"]`,
			Title:  "Genera varias gu\u00edas a la vez",
			Hint:   "Este bot\u00f3n abre la generaci\u00f3n masiva: eliges la bodega de origen, cotizas las \u00f3rdenes y generas sus gu\u00edas juntas.",
		},
	},
}

package moduleregistry

type ModuleCode string

const (
	ModuleIAM           ModuleCode = "iam"
	ModuleOrders        ModuleCode = "orders"
	ModuleShipments     ModuleCode = "shipments"
	ModuleInventory     ModuleCode = "inventory"
	ModuleInvoicing     ModuleCode = "invoicing"
	ModuleDelivery      ModuleCode = "delivery"
	ModuleCustomers     ModuleCode = "customers"
	ModuleStorefront    ModuleCode = "storefront"
	ModuleWallet        ModuleCode = "wallet"
	ModuleAnnouncements ModuleCode = "announcements"
	ModuleTickets       ModuleCode = "tickets"
	ModuleIntegrations  ModuleCode = "integrations"
	ModuleNotifications ModuleCode = "notification_config"
)

var All = []ModuleCode{
	ModuleIAM,
	ModuleOrders,
	ModuleShipments,
	ModuleInventory,
	ModuleInvoicing,
	ModuleDelivery,
	ModuleCustomers,
	ModuleStorefront,
	ModuleWallet,
	ModuleAnnouncements,
	ModuleTickets,
	ModuleIntegrations,
	ModuleNotifications,
}

var DisplayNames = map[ModuleCode]string{
	ModuleIAM:           "Usuarios y Roles",
	ModuleOrders:        "Ordenes",
	ModuleShipments:     "Envios",
	ModuleInventory:     "Inventario",
	ModuleInvoicing:     "Facturacion",
	ModuleDelivery:      "Ultima Milla",
	ModuleCustomers:     "Clientes",
	ModuleStorefront:    "Tienda",
	ModuleWallet:        "Billetera",
	ModuleAnnouncements: "Anuncios",
	ModuleTickets:       "Tickets",
	ModuleIntegrations:  "Integraciones",
	ModuleNotifications: "Notificaciones",
}

// RestrictedByDefault son modulos ocultos para todos los negocios sin
// importar su plan (incluso sin ningun plan asignado). Solo se habilitan
// otorgando un override puntual al negocio.
var RestrictedByDefault = []ModuleCode{
	ModuleStorefront,
	ModuleTickets,
	ModuleDelivery,
}

func IsValid(code string) bool {
	for _, m := range All {
		if string(m) == code {
			return true
		}
	}
	return false
}

func DisplayName(code string) string {
	if name, ok := DisplayNames[ModuleCode(code)]; ok {
		return name
	}
	return code
}

func IsRestrictedByDefault(code string) bool {
	for _, m := range RestrictedByDefault {
		if string(m) == code {
			return true
		}
	}
	return false
}

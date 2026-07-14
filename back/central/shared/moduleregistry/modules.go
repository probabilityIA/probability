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

func IsValid(code string) bool {
	for _, m := range All {
		if string(m) == code {
			return true
		}
	}
	return false
}

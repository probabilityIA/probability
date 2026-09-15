package authz

import (
	"strings"

	"github.com/secamc93/probability/back/central/shared/moduleregistry"
)

const (
	ActionRead   = "read"
	ActionCreate = "create"
	ActionUpdate = "update"
	ActionDelete = "delete"
)

type Resource struct {
	Code              string
	Label             string
	Module            string
	LegacyNames       []string
	SuperAdminOnly    bool
	InheritFromParent string
	ExcludeRoles      []string
}

type NavRule string

const (
	NavRulePermission     NavRule = "permission"
	NavRuleAlways         NavRule = "always"
	NavRuleSuperAdminOnly NavRule = "super_admin_only"
	NavRuleRole           NavRule = "role"
)

type NavEntry struct {
	Key          string
	Label        string
	Route        string
	Section      string
	Description  string
	Rule         NavRule
	Resource     string
	Module       string
	RoleCodes    []string
	ExcludeRoles []string
}

var actionAliases = map[string]string{
	"read":      ActionRead,
	"list":      ActionRead,
	"create":    ActionCreate,
	"update":    ActionUpdate,
	"delete":    ActionDelete,
	"manage":    "manage",
	"approve":   "approve",
	"reject":    "reject",
	"assign":    "assign",
	"schedule":  "schedule",
	"report":    "report",
	"configure": "configure",
	"audit":     "audit",
	"migrate":   "migrate",
}

var Resources = []Resource{
	{Code: "users", Label: "Usuarios", Module: string(moduleregistry.ModuleIAM), LegacyNames: []string{"Usuarios", "Users", "Empleados"}},
	{Code: "roles", Label: "Roles", Module: string(moduleregistry.ModuleIAM), LegacyNames: []string{"Roles", "Roles y Permisos"}},
	{Code: "permissions", Label: "Permisos", Module: string(moduleregistry.ModuleIAM), LegacyNames: []string{"Permisos", "Permissions"}},
	{Code: "resources", Label: "Recursos", Module: string(moduleregistry.ModuleIAM), LegacyNames: []string{"Recursos", "Resources"}, SuperAdminOnly: true},
	{Code: "businesses", Label: "Empresas", Module: string(moduleregistry.ModuleIAM), LegacyNames: []string{"Empresas", "Businesses"}},
	{Code: "orders", Label: "\u00d3rdenes", Module: string(moduleregistry.ModuleOrders), LegacyNames: []string{"Ordenes", "Orders"}},
	{Code: "products", Label: "Productos", Module: string(moduleregistry.ModuleOrders), LegacyNames: []string{"Productos", "Products"}},
	{Code: "shipments", Label: "Env\u00edos", Module: string(moduleregistry.ModuleShipments), LegacyNames: []string{"Envios", "Shipments"}},
	{Code: "invoicing", Label: "Facturaci\u00f3n", Module: string(moduleregistry.ModuleInvoicing), LegacyNames: []string{"Facturacion", "Invoicing"}},
	{Code: "customers", Label: "Clientes", Module: string(moduleregistry.ModuleCustomers), LegacyNames: []string{"Clientes", "Customers"}},
	{Code: "wallet", Label: "Billetera", Module: string(moduleregistry.ModuleWallet), LegacyNames: []string{"Billetera", "Wallet"}},
	{Code: "notifications", Label: "Notificaciones", Module: string(moduleregistry.ModuleNotifications), LegacyNames: []string{"Notificaciones", "Configuraci\u00f3n de Notificaciones"}},
	{Code: "delivery", Label: "\u00daltima milla", Module: string(moduleregistry.ModuleDelivery), LegacyNames: []string{"Ultima Milla", "Delivery"}},
	{Code: "storefront", Label: "Cat\u00e1logo", Module: string(moduleregistry.ModuleStorefront), LegacyNames: []string{"Storefront"}},
	{Code: "integrations", Label: "Integraciones", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones", "Integrations"}},
	{Code: "integrations.platform", Label: "Integraciones de plataforma", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-Platform"}},
	{Code: "integrations.ecommerce", Label: "Integraciones e-commerce", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-E-commerce"}},
	{Code: "integrations.einvoicing", Label: "Facturaci\u00f3n electr\u00f3nica", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-Facturacion-Electronica"}},
	{Code: "integrations.messaging", Label: "Mensajer\u00eda", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-Mensajeria"}},
	{Code: "integrations.payments", Label: "Pagos", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-Pagos"}},
	{Code: "integrations.logistics", Label: "Log\u00edstica", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-Logistica"}},
	{Code: "integrations.types", Label: "Tipos de integraci\u00f3n", Module: string(moduleregistry.ModuleIntegrations), LegacyNames: []string{"Integraciones-Tipos-de-integracion"}, SuperAdminOnly: true},
	{Code: "inventory", Label: "Inventario", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario", "Inventory"}},
	{Code: "warehouses", Label: "Bodegas", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Bodegas", "Warehouses"}},
	{Code: "inventory.stock", Label: "Stock", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Stock"}, InheritFromParent: "inventory"},
	{Code: "inventory.movements", Label: "Movimientos", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Movimientos"}, InheritFromParent: "inventory"},
	{Code: "inventory.traceability", Label: "Trazabilidad", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Trazabilidad"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.kardex", Label: "Kardex", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Kardex"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.operations", Label: "Operaciones", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Operaciones"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.slotting", Label: "Slotting", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Slotting"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.audit", Label: "Auditor\u00eda", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Auditoria"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.lpn", Label: "LPN", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-LPN"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.scan", Label: "Escaneo", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Scan"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
	{Code: "inventory.sync_logs", Label: "Logs de sincronizaci\u00f3n", Module: string(moduleregistry.ModuleInventory), LegacyNames: []string{"Inventario-Sync-Logs"}, InheritFromParent: "inventory", ExcludeRoles: []string{"demo"}},
}

var Navigation = []NavEntry{
	{Key: "home", Label: "Inicio", Route: "/home", Section: "main", Rule: NavRuleAlways, Description: "Resumen del negocio"},
	{Key: "orders", Label: "\u00d3rdenes", Route: "/orders", Section: "operations", Rule: NavRulePermission, Resource: "orders", Description: "Pedidos de todos los canales de venta"},
	{Key: "products", Label: "Productos", Route: "/products", Section: "operations", Rule: NavRulePermission, Resource: "products", Description: "Cat\u00e1logo de productos"},
	{Key: "shipments", Label: "Env\u00edos", Route: "/shipments", Section: "operations", Rule: NavRulePermission, Resource: "shipments", Description: "Gu\u00edas, seguimiento y contra entrega"},
	{Key: "customers", Label: "Clientes", Route: "/customers", Section: "operations", Rule: NavRulePermission, Resource: "customers", Description: "Clientes y su historial"},
	{Key: "delivery", Label: "\u00daltima milla", Route: "/delivery/routes", Section: "operations", Rule: NavRulePermission, Resource: "delivery", Description: "Rutas y repartidores propios"},
	{Key: "inventory", Label: "Inventario", Route: "/inventory", Section: "inventory", Rule: NavRulePermission, Resource: "inventory", Description: "Existencias y movimientos"},
	{Key: "warehouses", Label: "Bodegas", Route: "/warehouses", Section: "inventory", Rule: NavRulePermission, Resource: "warehouses", Description: "Bodegas y ubicaciones"},
	{Key: "invoicing", Label: "Facturaci\u00f3n", Route: "/invoicing/invoices", Section: "finance", Rule: NavRulePermission, Resource: "invoicing", Description: "Facturas electr\u00f3nicas"},
	{Key: "wallet", Label: "Billetera", Route: "/wallet", Section: "finance", Rule: NavRulePermission, Resource: "wallet", Description: "Saldo, recargas y movimientos"},
	{Key: "notifications", Label: "Notificaciones", Route: "/notification-config", Section: "communication", Rule: NavRulePermission, Resource: "notifications", Description: "Mensajes autom\u00e1ticos y conversaciones"},
	{Key: "storefront", Label: "Cat\u00e1logo", Route: "/storefront/catalogo", Section: "sales", Rule: NavRulePermission, Resource: "storefront", Description: "Tienda en l\u00ednea"},
	{Key: "integrations", Label: "Integraciones", Route: "/integrations", Section: "settings", Rule: NavRulePermission, Resource: "integrations", Description: "Canales de venta, transportadoras y facturaci\u00f3n"},
	{Key: "users", Label: "Usuarios", Route: "/users", Section: "iam", Rule: NavRulePermission, Resource: "users", Description: "Equipo del negocio"},
	{Key: "roles", Label: "Roles", Route: "/roles", Section: "iam", Rule: NavRulePermission, Resource: "roles", Description: "Roles y sus permisos"},
	{Key: "permissions", Label: "Permisos", Route: "/permissions", Section: "iam", Rule: NavRulePermission, Resource: "permissions", Description: "Permisos del sistema"},
	{Key: "businesses", Label: "Empresas", Route: "/businesses", Section: "iam", Rule: NavRulePermission, Resource: "businesses", Description: "Datos del negocio"},
	{Key: "resources", Label: "Recursos", Route: "/resources", Section: "iam", Rule: NavRuleSuperAdminOnly, Description: "Recursos protegidos"},
	{Key: "website_config", Label: "Sitio web", Route: "/website-config", Section: "settings", Rule: NavRuleRole, RoleCodes: []string{"administrador"}, Description: "P\u00e1gina web del negocio"},
	{Key: "subscription", Label: "Suscripci\u00f3n", Route: "/subscription", Section: "settings", Rule: NavRuleAlways, ExcludeRoles: []string{"demo"}, Description: "Plan y pagos"},
	{Key: "profile", Label: "Perfil", Route: "/profile", Section: "account", Rule: NavRuleAlways, Description: "Datos de la cuenta"},
	{Key: "tickets", Label: "Tickets", Route: "/tickets", Section: "platform", Rule: NavRuleSuperAdminOnly, Module: string(moduleregistry.ModuleTickets), Description: "Soporte y desarrollo"},
	{Key: "announcements", Label: "Anuncios", Route: "/announcements", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Anuncios a los negocios"},
	{Key: "accounting", Label: "Contabilidad", Route: "/accounting", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Contabilidad de la plataforma"},
	{Key: "commercial", Label: "Comercial", Route: "/commercial", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Prospectos y seguimiento comercial"},
	{Key: "marketing_leads", Label: "Leads", Route: "/marketing-leads", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Contactos del sitio web"},
	{Key: "siigo_referrals", Label: "Referidos Siigo", Route: "/siigo-referrals", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Referidos de Siigo"},
	{Key: "notification_channels", Label: "Canales de notificaci\u00f3n", Route: "/notification-channels", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Canales de env\u00edo de notificaciones"},
	{Key: "notification_event_types", Label: "Eventos de notificaci\u00f3n", Route: "/notification-event-types", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "Tipos de evento que disparan notificaciones"},
	{Key: "shipping_margins", Label: "M\u00e1rgenes de env\u00edo", Route: "/shipping-margins", Section: "platform", Rule: NavRuleSuperAdminOnly, Description: "M\u00e1rgenes por transportadora"},
}

var resourceByCode = func() map[string]Resource {
	m := make(map[string]Resource, len(Resources))
	for _, r := range Resources {
		m[r.Code] = r
	}
	return m
}()

var resourceByLegacyName = func() map[string]Resource {
	m := make(map[string]Resource)
	for _, r := range Resources {
		for _, name := range r.LegacyNames {
			m[normalize(name)] = r
		}
	}
	return m
}()

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

func ResourceByCode(code string) (Resource, bool) {
	r, ok := resourceByCode[code]
	return r, ok
}

func ResourceByLegacyName(name string) (Resource, bool) {
	r, ok := resourceByLegacyName[normalize(name)]
	return r, ok
}

func ActionCode(legacyName string) (string, bool) {
	code, ok := actionAliases[normalize(legacyName)]
	return code, ok
}

func PermissionCode(resourceCode, actionCode string) string {
	return resourceCode + "." + actionCode
}

func RoleCode(roleName string) string {
	replacer := strings.NewReplacer(" ", "_", "-", "_",
		"\u00e1", "a", "\u00e9", "e", "\u00ed", "i", "\u00f3", "o", "\u00fa", "u", "\u00f1", "n")
	return replacer.Replace(normalize(roleName))
}

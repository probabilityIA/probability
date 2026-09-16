package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/infra/primary/handlers/response"
)

var alertRoutes = map[string]response.Destination{
	"orders":        {Key: "orders", Label: "Órdenes", Route: "/orders", Description: "Pedidos de todos los canales de venta"},
	"shipments":     {Key: "shipments", Label: "Envíos", Route: "/shipments", Description: "Guías, seguimiento y contra entrega"},
	"invoicing":     {Key: "invoicing", Label: "Facturación", Route: "/invoicing/invoices", Description: "Facturas electrónicas"},
	"wallet":        {Key: "wallet", Label: "Billetera", Route: "/wallet", Description: "Saldo, recargas y movimientos"},
	"notifications": {Key: "notifications", Label: "Notificaciones", Route: "/notification-config", Description: "Mensajes autom\u00e1ticos y conversaciones"},
}

func FromAlert(alert entities.Alert) response.Alert {
	out := response.Alert{
		ID:        alert.ID,
		EventType: alert.EventType,
		Severity:  alert.Severity,
		Title:     alert.Title,
		Body:      alert.Body,
		Unread:    alert.Unread,
		CreatedAt: alert.CreatedAt,
	}
	if dest, ok := alertRoutes[alert.DestinationKey]; ok {
		if alert.DestinationRoute != "" {
			dest.Route = alert.DestinationRoute
		}
		out.Destination = &dest
	}
	if alert.ReferenceID != "" {
		out.Reference = &response.Reference{Type: alert.ReferenceType, ID: alert.ReferenceID}
	}
	return out
}

func FromAlerts(page *dtos.PaginatedResponse[entities.Alert]) response.Alerts {
	items := make([]response.Alert, 0, len(page.Data))
	for _, alert := range page.Data {
		items = append(items, FromAlert(alert))
	}
	return response.Alerts{Data: items, Total: page.Total, Page: page.Page, PageSize: page.PageSize, TotalPages: page.TotalPages}
}

func FromAlertsUnread(unread *entities.AlertsUnread) response.AlertsUnread {
	out := response.AlertsUnread{Count: unread.Count}
	if unread.Latest != nil {
		latest := FromAlert(*unread.Latest)
		out.Latest = &latest
	}
	return out
}

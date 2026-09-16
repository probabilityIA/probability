package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const (
	toolFindOrder     = "consultar_orden"
	toolFindShipment  = "consultar_guia"
	toolListOrders    = "listar_ordenes"
	toolOrdersSummary = "resumen_ordenes"
	toolListAlerts    = "consultar_alertas"

	dateLayout = "2006-01-02"
)

type dataAccess struct {
	BusinessID *uint
	Tools      []dtos.ToolDefinition
	Blocked    string
}

var shipmentStatusLabels = map[string]string{
	"pending":            "Pendiente",
	"picked_up":          "Recolectado",
	"in_transit":         "En tr\u00e1nsito",
	"out_for_delivery":   "En reparto",
	"delivered":          "Entregado",
	"on_hold":            "Novedad",
	"returned":           "Devuelto",
	"failed":             "Fallido",
	"cancelled":          "Cancelado",
	"needs_verification": "Por verificar",
}

func shipmentStatusLabel(code string) string {
	if label, ok := shipmentStatusLabels[code]; ok {
		return label
	}
	return code
}

func (uc *UseCase) resolveDataAccess(catalog *entities.NavigationCatalog, scope dtos.AccessScope) dataAccess {
	_, canOrders := catalog.Find("orders")
	_, canShipments := catalog.Find("shipments")
	if uc.businessData == nil || (!canOrders && !canShipments) {
		return dataAccess{}
	}

	businessID := businessOf(scope)
	if businessID == nil {
		return dataAccess{Blocked: "no_business"}
	}

	access := dataAccess{BusinessID: businessID}
	if canOrders {
		access.Tools = append(access.Tools, findOrderTool(), listOrdersTool(), ordersSummaryTool())
	}
	access.Tools = append(access.Tools, findShipmentTool())
	if uc.alerts != nil {
		access.Tools = append(access.Tools, listAlertsTool())
	}
	return access
}

func listAlertsTool() dtos.ToolDefinition {
	return dtos.ToolDefinition{
		Name:        toolListAlerts,
		Description: "Lista las \u00faltimas alertas o novedades del negocio que V\u00eda registr\u00f3: cancelaciones, gu\u00edas rechazadas, novedades de transportadora, saldo bajo, facturas fallidas. Dice cu\u00e1ndo pas\u00f3 cada una y a qu\u00e9 orden o gu\u00eda pertenece.",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"limite": map[string]any{"type": "integer", "minimum": 1, "maximum": RecentAlertsForModel},
			},
		},
	}
}

func findOrderTool() dtos.ToolDefinition {
	return dtos.ToolDefinition{
		Name:        toolFindOrder,
		Description: "Busca una orden del negocio por el n\u00famero que ve la persona (por ejemplo DEM-0048, 1042 o el n\u00famero del canal). Devuelve estado, pago, total, contra entrega, productos y la gu\u00eda.",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"numero": map[string]any{"type": "string", "description": "N\u00famero de la orden tal como lo escribi\u00f3 la persona."},
			},
			"required": []string{"numero"},
		},
	}
}

func findShipmentTool() dtos.ToolDefinition {
	return dtos.ToolDefinition{
		Name:        toolFindShipment,
		Description: "Busca un env\u00edo del negocio por n\u00famero de gu\u00eda o de rastreo. Devuelve transportadora, estado, novedad, valor a recaudar, fechas y los \u00faltimos eventos de rastreo.",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"numero_guia": map[string]any{"type": "string", "description": "N\u00famero de gu\u00eda o de rastreo."},
			},
			"required": []string{"numero_guia"},
		},
	}
}

func listOrdersTool() dtos.ToolDefinition {
	return dtos.ToolDefinition{
		Name:        toolListOrders,
		Description: "Lista hasta 10 \u00f3rdenes recientes del negocio con filtros opcionales y dice cu\u00e1ntas hay en total con esos filtros.",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"estado":              map[string]any{"type": "string", "description": "Texto del estado de la orden, por ejemplo pendiente, en camino, entregada, cancelada o novedad."},
				"desde":               map[string]any{"type": "string", "description": "Fecha inicial AAAA-MM-DD."},
				"hasta":               map[string]any{"type": "string", "description": "Fecha final AAAA-MM-DD, incluida."},
				"solo_contra_entrega": map[string]any{"type": "boolean"},
				"sin_guia":            map[string]any{"type": "boolean", "description": "Solo \u00f3rdenes que todav\u00eda no tienen gu\u00eda."},
				"con_novedad":         map[string]any{"type": "boolean", "description": "Solo \u00f3rdenes o env\u00edos con novedad."},
				"limite":              map[string]any{"type": "integer", "minimum": 1, "maximum": MaxListedOrders},
			},
		},
	}
}

func ordersSummaryTool() dtos.ToolDefinition {
	return dtos.ToolDefinition{
		Name:        toolOrdersSummary,
		Description: "Resume las \u00f3rdenes del negocio en un per\u00edodo: cantidad, valor total, cu\u00e1ntas por estado, contra entrega, sin gu\u00eda, con novedad y estados de los env\u00edos. Sin fechas usa los \u00faltimos 7 d\u00edas.",
		Schema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"desde": map[string]any{"type": "string", "description": "Fecha inicial AAAA-MM-DD."},
				"hasta": map[string]any{"type": "string", "description": "Fecha final AAAA-MM-DD, incluida."},
			},
		},
	}
}

func (uc *UseCase) runTools(ctx context.Context, access dataAccess, calls []dtos.ToolCall) []dtos.ToolResult {
	allowed := make(map[string]bool, len(access.Tools))
	for _, tool := range access.Tools {
		allowed[tool.Name] = true
	}

	results := make([]dtos.ToolResult, 0, len(calls))
	for i, call := range calls {
		if i >= MaxToolCallsPerTurn {
			results = append(results, dtos.ToolResult{ToolCallID: call.ID, Content: map[string]any{"error": "Demasiadas consultas en un solo paso."}})
			continue
		}
		if !allowed[call.Name] || access.BusinessID == nil {
			results = append(results, dtos.ToolResult{ToolCallID: call.ID, Content: map[string]any{"error": "Esa consulta no est\u00e1 disponible para este usuario."}})
			continue
		}

		content, err := uc.runTool(ctx, *access.BusinessID, call)
		if err != nil {
			uc.log.Warn(ctx).Err(err).Str("tool", call.Name).Uint("business_id", *access.BusinessID).Msg("[ai.assistant] fallo una consulta de datos")
			content = map[string]any{"error": "No se pudo consultar en este momento."}
		}
		results = append(results, dtos.ToolResult{ToolCallID: call.ID, Content: content})
	}
	return results
}

func (uc *UseCase) runTool(ctx context.Context, businessID uint, call dtos.ToolCall) (map[string]any, error) {
	switch call.Name {
	case toolFindOrder:
		number := stringArg(call.Input, "numero")
		if number == "" {
			return map[string]any{"error": "Falta el n\u00famero de la orden."}, nil
		}
		orders, err := uc.businessData.FindOrders(ctx, businessID, number)
		if err != nil {
			return nil, err
		}
		if len(orders) == 0 {
			return map[string]any{"encontradas": 0, "mensaje": "No hay una orden con ese n\u00famero en este negocio."}, nil
		}
		items := make([]map[string]any, 0, len(orders))
		for _, order := range orders {
			items = append(items, orderToMap(order))
		}
		return map[string]any{"encontradas": len(orders), "ordenes": items}, nil

	case toolFindShipment:
		number := stringArg(call.Input, "numero_guia")
		if number == "" {
			return map[string]any{"error": "Falta el n\u00famero de la gu\u00eda."}, nil
		}
		shipments, err := uc.businessData.FindShipments(ctx, businessID, number)
		if err != nil {
			return nil, err
		}
		if len(shipments) == 0 {
			return map[string]any{"encontradas": 0, "mensaje": "No hay una gu\u00eda con ese n\u00famero en este negocio."}, nil
		}
		items := make([]map[string]any, 0, len(shipments))
		for _, shipment := range shipments {
			items = append(items, shipmentToMap(shipment))
		}
		return map[string]any{"encontradas": len(shipments), "guias": items}, nil

	case toolListOrders:
		query := dtos.OrderQuery{
			Status:       stringArg(call.Input, "estado"),
			CodOnly:      boolArg(call.Input, "solo_contra_entrega"),
			WithoutGuide: boolArg(call.Input, "sin_guia"),
			WithNovelty:  boolArg(call.Input, "con_novedad"),
			Limit:        intArg(call.Input, "limite", MaxListedOrders),
		}
		if query.Limit < 1 || query.Limit > MaxListedOrders {
			query.Limit = MaxListedOrders
		}
		query.From, query.To = dateRangeArgs(call.Input)
		orders, total, err := uc.businessData.ListOrders(ctx, businessID, query)
		if err != nil {
			return nil, err
		}
		items := make([]map[string]any, 0, len(orders))
		for _, order := range orders {
			items = append(items, summaryToMap(order))
		}
		return map[string]any{"total_con_filtros": total, "mostradas": len(items), "ordenes": items}, nil

	case toolOrdersSummary:
		from, to := dateRangeArgs(call.Input)
		now := uc.now().In(colombia)
		if to == nil {
			end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, colombia).AddDate(0, 0, 1)
			to = &end
		}
		if from == nil {
			start := to.AddDate(0, 0, -DefaultSummaryDays)
			from = &start
		}
		overview, err := uc.businessData.SummarizeOrders(ctx, businessID, *from, *to)
		if err != nil {
			return nil, err
		}
		return overviewToMap(overview), nil

	case toolListAlerts:
		limit := intArg(call.Input, "limite", RecentAlertsForModel)
		if limit < 1 || limit > RecentAlertsForModel {
			limit = RecentAlertsForModel
		}
		alerts, err := uc.alerts.RecentAlerts(ctx, businessID, limit)
		if err != nil {
			return nil, err
		}
		if len(alerts) == 0 {
			return map[string]any{"encontradas": 0, "mensaje": "No hay alertas registradas para este negocio."}, nil
		}
		items := make([]map[string]any, 0, len(alerts))
		for _, alert := range alerts {
			items = append(items, alertToMap(alert))
		}
		return map[string]any{"encontradas": len(items), "alertas": items}, nil
	}
	return map[string]any{"error": "Consulta desconocida."}, nil
}

func orderToMap(order entities.OrderInfo) map[string]any {
	products := make([]map[string]any, 0, len(order.Items))
	for _, item := range order.Items {
		product := map[string]any{"nombre": item.Name, "cantidad": item.Quantity}
		if item.Variant != "" {
			product["variante"] = item.Variant
		}
		products = append(products, product)
	}

	out := map[string]any{
		"numero":         order.Number,
		"canal":          order.Platform,
		"fecha":          order.CreatedAt.In(colombia).Format("2006-01-02 15:04"),
		"estado":         order.StatusName,
		"estado_pago":    order.PaymentStatus,
		"pagada":         order.IsPaid,
		"total":          order.Total,
		"moneda":         order.Currency,
		"contra_entrega": order.IsCod,
		"cliente":        order.CustomerName,
		"ciudad":         order.City,
		"confirmada":     order.IsConfirmed,
		"facturada":      order.HasInvoice,
		"productos":      products,
	}
	if order.CodTotal != nil {
		out["valor_contra_entrega"] = *order.CodTotal
	}
	if order.Novelty != "" {
		out["novedad"] = order.Novelty
	}
	if order.IsTest {
		out["es_prueba"] = true
	}
	if order.Shipment != nil {
		out["guia"] = shipmentToMap(*order.Shipment)
	} else {
		out["guia"] = "Sin gu\u00eda"
	}
	return out
}

func shipmentToMap(shipment entities.ShipmentInfo) map[string]any {
	out := map[string]any{
		"numero_guia":    shipment.TrackingNumber,
		"transportadora": shipment.Carrier,
		"estado":         shipmentStatusLabel(shipment.Status),
		"creada":         shipment.CreatedAt.In(colombia).Format("2006-01-02 15:04"),
	}
	if shipment.CarrierStatus != "" {
		out["estado_transportadora"] = shipment.CarrierStatus
	}
	if shipment.CarrierStatusDetail != "" {
		out["detalle"] = shipment.CarrierStatusDetail
	}
	if shipment.DestinationCity != "" {
		out["ciudad_destino"] = shipment.DestinationCity
	}
	if shipment.OrderNumber != "" {
		out["orden"] = shipment.OrderNumber
	}
	if shipment.CodCollectAmount != nil {
		out["valor_a_recaudar"] = *shipment.CodCollectAmount
	}
	if shipment.TotalCost != nil {
		out["costo_guia"] = *shipment.TotalCost
	}
	if shipment.ShippedAt != nil {
		out["despachada"] = shipment.ShippedAt.In(colombia).Format(dateLayout)
	}
	if shipment.DeliveredAt != nil {
		out["entregada"] = shipment.DeliveredAt.In(colombia).Format(dateLayout)
	}
	if shipment.EstimatedDelivery != nil {
		out["entrega_estimada"] = shipment.EstimatedDelivery.In(colombia).Format(dateLayout)
	}
	if shipment.IsTest {
		out["es_prueba"] = true
	}
	if len(shipment.Events) > 0 {
		events := make([]map[string]any, 0, len(shipment.Events))
		for _, event := range shipment.Events {
			events = append(events, map[string]any{"fecha": event.Date, "estado": event.RawStatus, "detalle": event.Description})
		}
		out["ultimos_eventos"] = events
	}
	return out
}

func summaryToMap(order entities.OrderSummary) map[string]any {
	out := map[string]any{
		"numero":         order.Number,
		"fecha":          order.CreatedAt.In(colombia).Format("2006-01-02 15:04"),
		"cliente":        order.CustomerName,
		"total":          order.Total,
		"moneda":         order.Currency,
		"estado":         order.StatusName,
		"contra_entrega": order.IsCod,
	}
	if order.TrackingNumber != "" {
		out["guia"] = order.TrackingNumber
		out["estado_envio"] = shipmentStatusLabel(order.ShipmentStatus)
	} else {
		out["guia"] = "Sin gu\u00eda"
	}
	if order.IsTest {
		out["es_prueba"] = true
	}
	return out
}

func overviewToMap(overview *entities.OrdersOverview) map[string]any {
	byStatus := make([]map[string]any, 0, len(overview.ByStatus))
	for _, s := range overview.ByStatus {
		byStatus = append(byStatus, map[string]any{"estado": s.Name, "cantidad": s.Count})
	}
	shipments := make([]map[string]any, 0, len(overview.ShipmentStatus))
	for _, s := range overview.ShipmentStatus {
		shipments = append(shipments, map[string]any{"estado": shipmentStatusLabel(s.Name), "cantidad": s.Count})
	}
	return map[string]any{
		"desde":             overview.From.In(colombia).Format(dateLayout),
		"hasta":             overview.To.In(colombia).AddDate(0, 0, -1).Format(dateLayout),
		"ordenes":           overview.Orders,
		"valor_total":       overview.Amount,
		"contra_entrega":    overview.CashOnDelivery,
		"sin_guia":          overview.WithoutGuide,
		"con_novedad":       overview.WithNovelty,
		"por_estado":        byStatus,
		"envios_por_estado": shipments,
	}
}

func stringArg(input map[string]any, key string) string {
	value, ok := input[key]
	if !ok || value == nil {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(value))
}

func boolArg(input map[string]any, key string) bool {
	switch v := input[key].(type) {
	case bool:
		return v
	case string:
		return strings.EqualFold(v, "true")
	}
	return false
}

func intArg(input map[string]any, key string, fallback int) int {
	switch v := input[key].(type) {
	case float64:
		return int(v)
	case int:
		return v
	case int64:
		return int(v)
	}
	return fallback
}

func dateRangeArgs(input map[string]any) (*time.Time, *time.Time) {
	var from, to *time.Time
	if raw := stringArg(input, "desde"); raw != "" {
		if parsed, err := time.ParseInLocation(dateLayout, raw, colombia); err == nil {
			from = &parsed
		}
	}
	if raw := stringArg(input, "hasta"); raw != "" {
		if parsed, err := time.ParseInLocation(dateLayout, raw, colombia); err == nil {
			end := parsed.AddDate(0, 0, 1)
			to = &end
		}
	}
	return from, to
}

func alertToMap(alert entities.Alert) map[string]any {
	item := map[string]any{
		"titulo":  alert.Title,
		"detalle": alert.Body,
		"fecha":   alert.CreatedAt.In(colombia).Format("2006-01-02 15:04"),
		"tipo":    alert.EventType,
	}
	if alert.ReferenceID != "" {
		item["referencia"] = map[string]any{"tipo": alert.ReferenceType, "numero": alert.ReferenceID}
	}
	if alert.DestinationKey != "" {
		item["modulo"] = alert.DestinationKey
	}
	return item
}

package entities

// OrderStatus define los posibles estados de una orden en Probability
type OrderStatus string

const (
	// ════════════════════════════════════════════
	// Estados iniciales
	// ════════════════════════════════════════════

	// OrderStatusPending - Orden recibida, pendiente de procesamiento
	OrderStatusPending OrderStatus = "pending"

	// OrderStatusOnHold - Orden en espera
	OrderStatusOnHold OrderStatus = "on_hold"

	// ════════════════════════════════════════════
	// Fase Almacén (Fulfillment)
	// ════════════════════════════════════════════

	// OrderStatusPicking - Seleccionando productos del inventario
	OrderStatusPicking OrderStatus = "picking"

	// OrderStatusPacking - Empacando el pedido
	OrderStatusPacking OrderStatus = "packing"

	// OrderStatusReadyToShip - Listo para despacho
	OrderStatusReadyToShip OrderStatus = "ready_to_ship"

	// OrderStatusInventoryIssue - Novedad de inventario (sin stock, producto dañado)
	OrderStatusInventoryIssue OrderStatus = "inventory_issue"

	// ════════════════════════════════════════════
	// Fase Asignación y Recogida
	// ════════════════════════════════════════════

	// OrderStatusAssignedToDriver - Asignado a piloto/conductor
	OrderStatusAssignedToDriver OrderStatus = "assigned_to_driver"

	// OrderStatusPickedUp - Recogido por el piloto
	OrderStatusPickedUp OrderStatus = "picked_up"

	// ════════════════════════════════════════════
	// Fase Tránsito / Última Milla
	// ════════════════════════════════════════════

	// OrderStatusInTransit - En camino al destino
	OrderStatusInTransit OrderStatus = "in_transit"

	// OrderStatusOutForDelivery - En reparto final (última milla)
	OrderStatusOutForDelivery OrderStatus = "out_for_delivery"

	// ════════════════════════════════════════════
	// Resultado de Entrega
	// ════════════════════════════════════════════

	// OrderStatusDelivered - Entregada al cliente
	OrderStatusDelivered OrderStatus = "delivered"

	// OrderStatusDeliveryNovelty - Novedad de entrega (dirección incorrecta, cliente ausente)
	OrderStatusDeliveryNovelty OrderStatus = "delivery_novelty"

	// OrderStatusDeliveryFailed - Entrega fallida (marcado manualmente por usuario)
	OrderStatusDeliveryFailed OrderStatus = "delivery_failed"

	// OrderStatusRejected - Rechazado por el cliente
	OrderStatusRejected OrderStatus = "rejected"

	// ════════════════════════════════════════════
	// Devoluciones
	// ════════════════════════════════════════════

	// OrderStatusReturnInTransit - Devolución en camino al almacén
	OrderStatusReturnInTransit OrderStatus = "return_in_transit"

	// OrderStatusReturned - Devuelto al almacén
	OrderStatusReturned OrderStatus = "returned"

	// ════════════════════════════════════════════
	// Estados finales / financieros
	// ════════════════════════════════════════════

	// OrderStatusCompleted - Orden completada exitosamente
	OrderStatusCompleted OrderStatus = "completed"

	// OrderStatusCancelled - Orden cancelada
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderStatusRefunded - Orden reembolsada
	OrderStatusRefunded OrderStatus = "refunded"

	// OrderStatusFailed - Fallo del sistema durante procesamiento
	OrderStatusFailed OrderStatus = "failed"

	// ════════════════════════════════════════════
	// Deprecados (backward compat - no usar en código nuevo)
	// ════════════════════════════════════════════

	// OrderStatusProcessing - DEPRECADO: usar OrderStatusPicking
	OrderStatusProcessing OrderStatus = "processing"

	// OrderStatusShipped - DEPRECADO: usar OrderStatusInTransit
	OrderStatusShipped OrderStatus = "shipped"
)

// validStatuses contiene todos los estados válidos del sistema
var validStatuses = map[OrderStatus]bool{
	OrderStatusPending:         true,
	OrderStatusOnHold:          true,
	OrderStatusPicking:         true,
	OrderStatusPacking:         true,
	OrderStatusReadyToShip:     true,
	OrderStatusInventoryIssue:  true,
	OrderStatusAssignedToDriver: true,
	OrderStatusPickedUp:        true,
	OrderStatusInTransit:       true,
	OrderStatusOutForDelivery:  true,
	OrderStatusDelivered:       true,
	OrderStatusDeliveryNovelty: true,
	OrderStatusDeliveryFailed:  true,
	OrderStatusRejected:        true,
	OrderStatusReturnInTransit: true,
	OrderStatusReturned:        true,
	OrderStatusCompleted:       true,
	OrderStatusCancelled:       true,
	OrderStatusRefunded:        true,
	OrderStatusFailed:          true,
	// Deprecados - aún válidos para órdenes históricas
	OrderStatusProcessing: true,
	OrderStatusShipped:    true,
}

// terminalStatuses son estados sin transiciones salientes
var terminalStatuses = map[OrderStatus]bool{
	OrderStatusCancelled: true,
	OrderStatusRefunded:  true,
}

// validTransitions define las transiciones permitidas entre estados
var validTransitions = map[OrderStatus][]OrderStatus{
	OrderStatusPending: {
		OrderStatusPicking,
		OrderStatusOnHold,
	},
	OrderStatusPicking: {
		OrderStatusPacking,
		OrderStatusInventoryIssue,
		OrderStatusOnHold,
	},
	OrderStatusPacking: {
		OrderStatusReadyToShip,
		OrderStatusOnHold,
	},
	OrderStatusReadyToShip: {
		OrderStatusAssignedToDriver,
		OrderStatusOnHold,
	},
	OrderStatusAssignedToDriver: {
		OrderStatusPickedUp,
	},
	OrderStatusPickedUp: {
		OrderStatusInTransit,
	},
	OrderStatusInTransit: {
		OrderStatusOutForDelivery,
	},
	OrderStatusOutForDelivery: {
		OrderStatusDelivered,
		OrderStatusDeliveryNovelty,
		OrderStatusRejected,
		OrderStatusDeliveryFailed,
	},
	OrderStatusDelivered: {
		OrderStatusCompleted,
		OrderStatusRefunded,
		OrderStatusReturnInTransit,
	},
	OrderStatusDeliveryNovelty: {
		OrderStatusAssignedToDriver,
		OrderStatusOutForDelivery,
		OrderStatusDeliveryFailed,
		OrderStatusReturnInTransit,
	},
	OrderStatusDeliveryFailed: {
		OrderStatusReturnInTransit,
	},
	OrderStatusRejected: {
		OrderStatusReturnInTransit,
	},
	OrderStatusReturnInTransit: {
		OrderStatusReturned,
	},
	OrderStatusReturned: {
		OrderStatusRefunded,
	},
	OrderStatusInventoryIssue: {
		OrderStatusPicking,
	},
	OrderStatusOnHold: {
		OrderStatusPending,
		OrderStatusPicking,
	},
	OrderStatusCompleted: {
		OrderStatusRefunded,
	},
	OrderStatusFailed: {},
}

// IsValid verifica si el estado es válido
func (s OrderStatus) IsValid() bool {
	return validStatuses[s]
}

// IsTerminal verifica si el estado es terminal (sin transiciones salientes)
func (s OrderStatus) IsTerminal() bool {
	return terminalStatuses[s]
}

// String retorna la representación en string del estado
func (s OrderStatus) String() string {
	return string(s)
}

// CanTransitionTo verifica si se puede transicionar al estado objetivo
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	// Cancelled es accesible desde cualquier estado no terminal
	if target == OrderStatusCancelled {
		return !s.IsTerminal()
	}

	allowedTargets, exists := validTransitions[s]
	if !exists {
		return false
	}

	for _, allowed := range allowedTargets {
		if allowed == target {
			return true
		}
	}
	return false
}

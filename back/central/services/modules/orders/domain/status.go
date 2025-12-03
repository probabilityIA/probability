package domain

// OrderStatus define los posibles estados de una orden en Probability
type OrderStatus string

const (
	// OrderStatusPending - Orden recibida, pendiente de procesamiento
	OrderStatusPending OrderStatus = "pending"

	// OrderStatusProcessing - Orden en proceso de preparación
	OrderStatusProcessing OrderStatus = "processing"

	// OrderStatusCompleted - Orden completada exitosamente
	OrderStatusCompleted OrderStatus = "completed"

	// OrderStatusCancelled - Orden cancelada
	OrderStatusCancelled OrderStatus = "cancelled"

	// OrderStatusFailed - Orden fallida
	OrderStatusFailed OrderStatus = "failed"

	// OrderStatusRefunded - Orden reembolsada
	OrderStatusRefunded OrderStatus = "refunded"

	// OrderStatusOnHold - Orden en espera
	OrderStatusOnHold OrderStatus = "on_hold"

	// OrderStatusShipped - Orden enviada
	OrderStatusShipped OrderStatus = "shipped"

	// OrderStatusDelivered - Orden entregada
	OrderStatusDelivered OrderStatus = "delivered"
)

// IsValid verifica si el estado es válido
func (s OrderStatus) IsValid() bool {
	switch s {
	case OrderStatusPending, OrderStatusProcessing, OrderStatusCompleted,
		OrderStatusCancelled, OrderStatusFailed, OrderStatusRefunded,
		OrderStatusOnHold, OrderStatusShipped, OrderStatusDelivered:
		return true
	}
	return false
}

// String retorna la representación en string del estado
func (s OrderStatus) String() string {
	return string(s)
}

// CanTransitionTo verifica si se puede transicionar al estado objetivo
func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	// Definir las transiciones válidas
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending: {
			OrderStatusProcessing,
			OrderStatusCancelled,
			OrderStatusOnHold,
		},
		OrderStatusProcessing: {
			OrderStatusCompleted,
			OrderStatusCancelled,
			OrderStatusOnHold,
			OrderStatusShipped,
		},
		OrderStatusOnHold: {
			OrderStatusPending,
			OrderStatusProcessing,
			OrderStatusCancelled,
		},
		OrderStatusShipped: {
			OrderStatusDelivered,
			OrderStatusFailed,
		},
		OrderStatusDelivered: {
			OrderStatusRefunded,
		},
		OrderStatusCompleted: {
			OrderStatusRefunded,
		},
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

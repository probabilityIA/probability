package request

// OrderEvent representa el evento de orden recibido de RabbitMQ
type OrderEvent struct {
	EventType     string `json:"event_type"`      // "order.created", "order.updated", "order.paid"
	OrderID       string `json:"order_id"`        // UUID de la orden
	BusinessID    uint   `json:"business_id"`     // ID del negocio
	IntegrationID uint   `json:"integration_id"`  // ID de la integración
	TotalAmount   float64 `json:"total_amount"`   // Monto total
	Currency      string `json:"currency"`        // Moneda
	IsPaid        bool   `json:"is_paid"`         // Si está pagada
	Status        string `json:"status"`          // Estado de la orden
	Timestamp     string `json:"timestamp"`       // Timestamp del evento
}

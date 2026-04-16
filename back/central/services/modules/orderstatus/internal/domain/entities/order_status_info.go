package entities

// OrderStatusInfo contiene información básica del estado de orden
// PURO - Sin tags JSON
type OrderStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
	Priority    int
	IsActive    bool
}

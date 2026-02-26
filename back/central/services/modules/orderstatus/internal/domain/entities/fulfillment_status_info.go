package entities

// FulfillmentStatusInfo contiene información de un estado de fulfillment
type FulfillmentStatusInfo struct {
	ID          uint
	Code        string
	Name        string
	Description string
	Category    string
	Color       string
}

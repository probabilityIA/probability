package domain

type ProbabilityShipmentStatus string

const (
	StatusPending        ProbabilityShipmentStatus = "pending"
	StatusPickedUp       ProbabilityShipmentStatus = "picked_up"
	StatusInTransit      ProbabilityShipmentStatus = "in_transit"
	StatusOutForDelivery ProbabilityShipmentStatus = "out_for_delivery"
	StatusDelivered      ProbabilityShipmentStatus = "delivered"
	StatusOnHold         ProbabilityShipmentStatus = "on_hold"
	StatusFailed         ProbabilityShipmentStatus = "failed"
	StatusReturned       ProbabilityShipmentStatus = "returned"
	StatusCancelled      ProbabilityShipmentStatus = "cancelled"
)

func (s ProbabilityShipmentStatus) String() string {
	return string(s)
}

func (s ProbabilityShipmentStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusPickedUp, StatusInTransit, StatusOutForDelivery,
		StatusDelivered, StatusOnHold, StatusFailed, StatusReturned, StatusCancelled:
		return true
	}
	return false
}

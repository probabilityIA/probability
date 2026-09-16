package entities

type OrderStatus string

const (
	OrderStatusPending OrderStatus = "pending"

	OrderStatusOnHold OrderStatus = "on_hold"

	OrderStatusPicking OrderStatus = "picking"

	OrderStatusPacking OrderStatus = "packing"

	OrderStatusReadyToShip OrderStatus = "ready_to_ship"

	OrderStatusInventoryIssue OrderStatus = "inventory_issue"

	OrderStatusAssignedToDriver OrderStatus = "assigned_to_driver"

	OrderStatusPickedUp OrderStatus = "picked_up"

	OrderStatusInTransit OrderStatus = "in_transit"

	OrderStatusOutForDelivery OrderStatus = "out_for_delivery"

	OrderStatusDelivered OrderStatus = "delivered"

	OrderStatusDeliveryNovelty OrderStatus = "delivery_novelty"

	OrderStatusDeliveryFailed OrderStatus = "delivery_failed"

	OrderStatusRejected OrderStatus = "rejected"

	OrderStatusReturnInTransit OrderStatus = "return_in_transit"

	OrderStatusReturned OrderStatus = "returned"

	OrderStatusCompleted OrderStatus = "completed"

	OrderStatusCancelled OrderStatus = "cancelled"

	OrderStatusCancelRequested OrderStatus = "cancel_requested"

	OrderStatusRefunded OrderStatus = "refunded"

	OrderStatusFailed OrderStatus = "failed"

	OrderStatusProcessing OrderStatus = "processing"

	OrderStatusShipped OrderStatus = "shipped"
)

var validStatuses = map[OrderStatus]bool{
	OrderStatusPending:          true,
	OrderStatusOnHold:           true,
	OrderStatusPicking:          true,
	OrderStatusPacking:          true,
	OrderStatusReadyToShip:      true,
	OrderStatusInventoryIssue:   true,
	OrderStatusAssignedToDriver: true,
	OrderStatusPickedUp:         true,
	OrderStatusInTransit:        true,
	OrderStatusOutForDelivery:   true,
	OrderStatusDelivered:        true,
	OrderStatusDeliveryNovelty:  true,
	OrderStatusDeliveryFailed:   true,
	OrderStatusRejected:         true,
	OrderStatusReturnInTransit:  true,
	OrderStatusReturned:         true,
	OrderStatusCompleted:        true,
	OrderStatusCancelled:        true,
	OrderStatusCancelRequested:  true,
	OrderStatusRefunded:         true,
	OrderStatusFailed:           true,
	OrderStatusProcessing:       true,
	OrderStatusShipped:          true,
}

var terminalStatuses = map[OrderStatus]bool{
	OrderStatusCancelled: true,
	OrderStatusRefunded:  true,
}

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
	OrderStatusCancelRequested: {
		OrderStatusOnHold,
		OrderStatusPicking,
		OrderStatusPacking,
		OrderStatusReadyToShip,
		OrderStatusAssignedToDriver,
		OrderStatusPickedUp,
		OrderStatusInTransit,
		OrderStatusOutForDelivery,
	},
}

func (s OrderStatus) IsValid() bool {
	return validStatuses[s]
}

func (s OrderStatus) IsTerminal() bool {
	return terminalStatuses[s]
}

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) CanTransitionTo(target OrderStatus) bool {
	if target == OrderStatusCancelled {
		return !s.IsTerminal()
	}
	if target == OrderStatusCancelRequested {
		return !s.IsTerminal() && s != OrderStatusCancelRequested
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

var channelOverrides = map[OrderStatus]bool{
	OrderStatusCancelled: true,
	OrderStatusRefunded:  true,
}

func (s OrderStatus) IsChannelOverride() bool {
	return channelOverrides[s]
}

package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrderStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		// Active statuses
		{"pending is valid", OrderStatusPending, true},
		{"picking is valid", OrderStatusPicking, true},
		{"packing is valid", OrderStatusPacking, true},
		{"ready_to_ship is valid", OrderStatusReadyToShip, true},
		{"assigned_to_driver is valid", OrderStatusAssignedToDriver, true},
		{"picked_up is valid", OrderStatusPickedUp, true},
		{"in_transit is valid", OrderStatusInTransit, true},
		{"out_for_delivery is valid", OrderStatusOutForDelivery, true},
		{"delivered is valid", OrderStatusDelivered, true},
		{"delivery_novelty is valid", OrderStatusDeliveryNovelty, true},
		{"delivery_failed is valid", OrderStatusDeliveryFailed, true},
		{"rejected is valid", OrderStatusRejected, true},
		{"return_in_transit is valid", OrderStatusReturnInTransit, true},
		{"returned is valid", OrderStatusReturned, true},
		{"inventory_issue is valid", OrderStatusInventoryIssue, true},
		{"on_hold is valid", OrderStatusOnHold, true},
		{"completed is valid", OrderStatusCompleted, true},
		{"cancelled is valid", OrderStatusCancelled, true},
		{"refunded is valid", OrderStatusRefunded, true},
		{"failed is valid", OrderStatusFailed, true},
		// Deprecated but valid
		{"processing is valid (deprecated)", OrderStatusProcessing, true},
		{"shipped is valid (deprecated)", OrderStatusShipped, true},
		// Invalid
		{"empty is invalid", OrderStatus(""), false},
		{"random string is invalid", OrderStatus("nonexistent"), false},
		{"typo is invalid", OrderStatus("pendng"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsValid())
		})
	}
}

func TestOrderStatus_IsTerminal(t *testing.T) {
	tests := []struct {
		name   string
		status OrderStatus
		want   bool
	}{
		{"cancelled is terminal", OrderStatusCancelled, true},
		{"refunded is terminal", OrderStatusRefunded, true},
		{"pending is not terminal", OrderStatusPending, false},
		{"delivered is not terminal", OrderStatusDelivered, false},
		{"completed is not terminal", OrderStatusCompleted, false},
		{"failed is not terminal", OrderStatusFailed, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.status.IsTerminal())
		})
	}
}

func TestOrderStatus_String(t *testing.T) {
	assert.Equal(t, "pending", OrderStatusPending.String())
	assert.Equal(t, "cancelled", OrderStatusCancelled.String())
}

func TestOrderStatus_CanTransitionTo_HappyPath(t *testing.T) {
	tests := []struct {
		name   string
		from   OrderStatus
		to     OrderStatus
		expect bool
	}{
		// Pending transitions
		{"pending -> picking", OrderStatusPending, OrderStatusPicking, true},
		{"pending -> on_hold", OrderStatusPending, OrderStatusOnHold, true},
		{"pending -> packing NOT allowed", OrderStatusPending, OrderStatusPacking, false},

		// Warehouse flow
		{"picking -> packing", OrderStatusPicking, OrderStatusPacking, true},
		{"picking -> inventory_issue", OrderStatusPicking, OrderStatusInventoryIssue, true},
		{"picking -> on_hold", OrderStatusPicking, OrderStatusOnHold, true},
		{"packing -> ready_to_ship", OrderStatusPacking, OrderStatusReadyToShip, true},
		{"ready_to_ship -> assigned_to_driver", OrderStatusReadyToShip, OrderStatusAssignedToDriver, true},

		// Driver flow
		{"assigned_to_driver -> picked_up", OrderStatusAssignedToDriver, OrderStatusPickedUp, true},
		{"picked_up -> in_transit", OrderStatusPickedUp, OrderStatusInTransit, true},
		{"in_transit -> out_for_delivery", OrderStatusInTransit, OrderStatusOutForDelivery, true},

		// Delivery outcomes
		{"out_for_delivery -> delivered", OrderStatusOutForDelivery, OrderStatusDelivered, true},
		{"out_for_delivery -> delivery_novelty", OrderStatusOutForDelivery, OrderStatusDeliveryNovelty, true},
		{"out_for_delivery -> rejected", OrderStatusOutForDelivery, OrderStatusRejected, true},
		{"out_for_delivery -> delivery_failed", OrderStatusOutForDelivery, OrderStatusDeliveryFailed, true},

		// Post-delivery
		{"delivered -> completed", OrderStatusDelivered, OrderStatusCompleted, true},
		{"delivered -> refunded", OrderStatusDelivered, OrderStatusRefunded, true},
		{"delivered -> return_in_transit", OrderStatusDelivered, OrderStatusReturnInTransit, true},

		// Returns
		{"delivery_novelty -> assigned_to_driver (retry)", OrderStatusDeliveryNovelty, OrderStatusAssignedToDriver, true},
		{"delivery_novelty -> out_for_delivery", OrderStatusDeliveryNovelty, OrderStatusOutForDelivery, true},
		{"delivery_novelty -> delivery_failed", OrderStatusDeliveryNovelty, OrderStatusDeliveryFailed, true},
		{"delivery_novelty -> return_in_transit", OrderStatusDeliveryNovelty, OrderStatusReturnInTransit, true},
		{"delivery_failed -> return_in_transit", OrderStatusDeliveryFailed, OrderStatusReturnInTransit, true},
		{"rejected -> return_in_transit", OrderStatusRejected, OrderStatusReturnInTransit, true},
		{"return_in_transit -> returned", OrderStatusReturnInTransit, OrderStatusReturned, true},
		{"returned -> refunded", OrderStatusReturned, OrderStatusRefunded, true},

		// Inventory issue recovery
		{"inventory_issue -> picking", OrderStatusInventoryIssue, OrderStatusPicking, true},

		// On hold recovery
		{"on_hold -> pending", OrderStatusOnHold, OrderStatusPending, true},
		{"on_hold -> picking", OrderStatusOnHold, OrderStatusPicking, true},

		// Completed -> refunded
		{"completed -> refunded", OrderStatusCompleted, OrderStatusRefunded, true},

		// Failed has no transitions
		{"failed -> pending NOT allowed", OrderStatusFailed, OrderStatusPending, false},
		{"failed -> picking NOT allowed", OrderStatusFailed, OrderStatusPicking, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expect, tt.from.CanTransitionTo(tt.to))
		})
	}
}

func TestOrderStatus_CanTransitionTo_CancelledFromAnyNonTerminal(t *testing.T) {
	nonTerminalStatuses := []OrderStatus{
		OrderStatusPending, OrderStatusPicking, OrderStatusPacking,
		OrderStatusReadyToShip, OrderStatusAssignedToDriver, OrderStatusPickedUp,
		OrderStatusInTransit, OrderStatusOutForDelivery, OrderStatusDelivered,
		OrderStatusDeliveryNovelty, OrderStatusDeliveryFailed, OrderStatusRejected,
		OrderStatusReturnInTransit, OrderStatusReturned, OrderStatusInventoryIssue,
		OrderStatusOnHold, OrderStatusCompleted, OrderStatusFailed,
	}

	for _, status := range nonTerminalStatuses {
		t.Run("cancelled from "+status.String(), func(t *testing.T) {
			assert.True(t, status.CanTransitionTo(OrderStatusCancelled),
				"%s should be able to transition to cancelled", status)
		})
	}
}

func TestOrderStatus_CanTransitionTo_TerminalCannotTransition(t *testing.T) {
	terminalStatuses := []OrderStatus{OrderStatusCancelled, OrderStatusRefunded}
	targets := []OrderStatus{
		OrderStatusPending, OrderStatusPicking, OrderStatusCancelled, OrderStatusRefunded,
	}

	for _, terminal := range terminalStatuses {
		for _, target := range targets {
			t.Run(terminal.String()+" -> "+target.String(), func(t *testing.T) {
				assert.False(t, terminal.CanTransitionTo(target),
					"terminal status %s should NOT transition to %s", terminal, target)
			})
		}
	}
}

func TestOrderStatus_CanTransitionTo_InvalidSourceStatus(t *testing.T) {
	// A status not in validTransitions (e.g. deprecated ones without explicit transitions)
	assert.False(t, OrderStatusProcessing.CanTransitionTo(OrderStatusPicking))
	assert.False(t, OrderStatusShipped.CanTransitionTo(OrderStatusDelivered))
	// But cancelled is still accessible from deprecated non-terminal statuses
	assert.True(t, OrderStatusProcessing.CanTransitionTo(OrderStatusCancelled))
	assert.True(t, OrderStatusShipped.CanTransitionTo(OrderStatusCancelled))
}

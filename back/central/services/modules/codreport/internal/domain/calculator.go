package domain

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/codreport/internal/domain/entities"
)

const (
	CodStateCollected      = "collected"
	CodStatePendingPayment = "pending_payment"
	CodStateInProgress     = "in_progress"
	CodStatePending        = "pending"
	CodStateNotCollectable = "not_collectable"
)

var inProgressStatuses = map[string]bool{
	"picked_up":        true,
	"in_transit":       true,
	"out_for_delivery": true,
	"on_hold":          true,
}

func CollectionState(shipmentStatus string) string {
	switch {
	case shipmentStatus == "delivered":
		return CodStateCollected
	case inProgressStatuses[shipmentStatus]:
		return CodStateInProgress
	case shipmentStatus == "pending":
		return CodStatePending
	default:
		return CodStateNotCollectable
	}
}

func IsCollected(shipmentStatus string) bool {
	return CollectionState(shipmentStatus) == CodStateCollected
}

func PaymentState(shipmentStatus string, paid bool) string {
	if paid {
		return CodStateCollected
	}
	if shipmentStatus == "delivered" {
		return CodStatePendingPayment
	}
	return CollectionState(shipmentStatus)
}

func ApplyDiscount(collected, pct float64) (discount, net float64) {
	if pct < 0 {
		pct = 0
	}
	discount = collected * pct / 100
	net = collected - discount
	return discount, net
}

func EnrichCarrierAggregates(aggs []entities.CarrierAggregate, discountMap map[string]float64) {
	for i := range aggs {
		pct := discountMap[aggs[i].Carrier]
		discount, net := ApplyDiscount(aggs[i].TotalCollected, pct)
		aggs[i].DiscountPct = pct
		aggs[i].TotalDiscount = discount
		aggs[i].TotalNet = net
	}
}

func SumAggregates(aggs []entities.CarrierAggregate) (orders int, collected, discount, net float64) {
	for i := range aggs {
		orders += aggs[i].OrdersCount
		collected += aggs[i].TotalCollected
		discount += aggs[i].TotalDiscount
		net += aggs[i].TotalNet
	}
	return orders, collected, discount, net
}

func WeekBounds(t time.Time) (start, end time.Time) {
	t = t.UTC()
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	start = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).AddDate(0, 0, -(weekday - 1))
	end = start.AddDate(0, 0, 6)
	return start, end
}

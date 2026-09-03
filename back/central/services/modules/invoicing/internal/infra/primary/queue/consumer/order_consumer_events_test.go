package consumer

import "testing"

func TestShouldProcessEvent(t *testing.T) {
	c := &OrderConsumer{}

	relevantes := []string{"order.created", "order.paid", "order.guide_generated"}
	for _, ev := range relevantes {
		if !c.shouldProcessEvent(ev) {
			t.Errorf("%q deberia disparar facturacion", ev)
		}
	}

	ignorados := []string{"order.updated", "order.cancelled", "shipment.guide_generated", ""}
	for _, ev := range ignorados {
		if c.shouldProcessEvent(ev) {
			t.Errorf("%q no deberia disparar facturacion", ev)
		}
	}
}

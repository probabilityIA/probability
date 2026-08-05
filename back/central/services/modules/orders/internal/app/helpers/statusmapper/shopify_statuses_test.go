package statusmapper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapShopifyFinancialStatusToPaymentStatus(t *testing.T) {
	cases := map[string]string{
		"pending":            "pending",
		"authorized":         "authorized",
		"paid":               "paid",
		"partially_paid":     "partially_paid",
		"refunded":           "refunded",
		"partially_refunded": "partially_refunded",
		"voided":             "voided",
		"unpaid":             "unpaid",
	}

	for entrada, esperado := range cases {
		t.Run(entrada, func(t *testing.T) {
			assert.Equal(t, esperado, MapShopifyFinancialStatusToPaymentStatus(entrada))
		})
	}
}

func TestMapShopifyFinancialStatusToPaymentStatus_DesconocidoCaeEnPending(t *testing.T) {
	desconocidos := []string{"", "expired", "PAID", "Paid", "algo-nuevo-de-shopify", "  paid  "}

	for _, entrada := range desconocidos {
		t.Run(entrada, func(t *testing.T) {
			assert.Equal(t, "pending", MapShopifyFinancialStatusToPaymentStatus(entrada))
		})
	}
}

func TestMapShopifyFulfillmentStatusToFulfillmentStatus(t *testing.T) {
	cases := map[string]string{
		"unfulfilled": "unfulfilled",
		"partial":     "partial",
		"fulfilled":   "fulfilled",
		"shipped":     "shipped",
	}

	for entrada, esperado := range cases {
		t.Run(entrada, func(t *testing.T) {
			valor := entrada
			assert.Equal(t, esperado, MapShopifyFulfillmentStatusToFulfillmentStatus(&valor))
		})
	}
}

func TestMapShopifyFulfillmentStatusToFulfillmentStatus_NilCaeEnUnfulfilled(t *testing.T) {
	assert.Equal(t, "unfulfilled", MapShopifyFulfillmentStatusToFulfillmentStatus(nil))
}

func TestMapShopifyFulfillmentStatusToFulfillmentStatus_VacioCaeEnUnfulfilled(t *testing.T) {
	vacio := ""

	assert.Equal(t, "unfulfilled", MapShopifyFulfillmentStatusToFulfillmentStatus(&vacio))
}

func TestMapShopifyFulfillmentStatusToFulfillmentStatus_DesconocidoCaeEnUnfulfilled(t *testing.T) {
	desconocidos := []string{"restocked", "FULFILLED", "Fulfilled", "  fulfilled  ", "algo-nuevo"}

	for _, entrada := range desconocidos {
		t.Run(entrada, func(t *testing.T) {
			valor := entrada
			assert.Equal(t, "unfulfilled", MapShopifyFulfillmentStatusToFulfillmentStatus(&valor))
		})
	}
}

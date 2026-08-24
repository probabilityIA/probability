package canonical

import (
	"encoding/json"
	"testing"
)

func decodificar(t *testing.T, crudo []byte) map[string]interface{} {
	t.Helper()
	if crudo == nil {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(crudo, &m); err != nil {
		t.Fatalf("la seccion no es JSON valido: %v", err)
	}
	return m
}

func TestExtractSectionsTiendanube(t *testing.T) {
	raw := []byte(`{
		"id": 2051703380,
		"status": "open",
		"payment_status": "paid",
		"gateway": "offline",
		"gateway_name": "A convenir",
		"paid_at": "2026-08-22T16:11:19+0000",
		"total": "55000.00",
		"subtotal": "55000.00",
		"currency": "COP",
		"discount": "0.00",
		"shipping_address": {"city": "Bogota"},
		"shipping_option": "A convenir",
		"shipping_cost_customer": "0.00",
		"shipping_status": "unpacked",
		"fulfillments": [],
		"contact_email": "carlos@example.test",
		"note": null
	}`)

	secciones := ExtractSections(raw, TiendanubeSections)

	pago := decodificar(t, secciones.Payment)
	if pago["payment_status"] != "paid" || pago["gateway_name"] != "A convenir" {
		t.Fatalf("la seccion de pago no trae lo esperado: %v", pago)
	}
	if _, sobra := pago["contact_email"]; sobra {
		t.Fatal("la seccion de pago no debe arrastrar llaves ajenas")
	}

	envio := decodificar(t, secciones.Shipping)
	if envio["shipping_option"] != "A convenir" {
		t.Fatalf("la seccion de envio no trae lo esperado: %v", envio)
	}

	finanzas := decodificar(t, secciones.Financial)
	if finanzas["total"] != "55000.00" || finanzas["currency"] != "COP" {
		t.Fatalf("la seccion financiera no trae lo esperado: %v", finanzas)
	}

	fulfillment := decodificar(t, secciones.Fulfillment)
	if fulfillment["shipping_status"] != "unpacked" || fulfillment["status"] != "open" {
		t.Fatalf("la seccion de fulfillment no trae lo esperado: %v", fulfillment)
	}
	if _, presente := fulfillment["note"]; presente {
		t.Fatal("las llaves nulas no deben guardarse")
	}
}

func TestExtractSectionsJumpsellerVieneAnidada(t *testing.T) {
	raw := []byte(`{"order": {
		"status": "Paid",
		"shipment_status": "in_transit",
		"payment_method_name": "Transferencia",
		"payment_method_type": "manual",
		"shipping_method_name": "Envio estandar",
		"currency": "CLP",
		"total": 19990
	}}`)

	secciones := ExtractSections(raw, JumpsellerSections)

	pago := decodificar(t, secciones.Payment)
	if pago["payment_method_name"] != "Transferencia" {
		t.Fatalf("no se leyo la orden anidada: %v", pago)
	}
	fulfillment := decodificar(t, secciones.Fulfillment)
	if fulfillment["shipment_status"] != "in_transit" {
		t.Fatalf("la seccion de fulfillment no trae lo esperado: %v", fulfillment)
	}
}

func TestExtractSectionsRaizAusente(t *testing.T) {
	secciones := ExtractSections([]byte(`{"status": "Paid"}`), JumpsellerSections)

	if secciones.Payment != nil || secciones.Shipping != nil || secciones.Financial != nil || secciones.Fulfillment != nil {
		t.Fatal("sin la llave raiz esperada no se inventan secciones")
	}
}

func TestExtractSectionsEntradasInvalidas(t *testing.T) {
	casos := map[string][]byte{
		"nulo":       nil,
		"vacio":      {},
		"no es json": []byte("<html>"),
		"es lista":   []byte(`[{"status":"paid"}]`),
	}

	for nombre, raw := range casos {
		t.Run(nombre, func(t *testing.T) {
			secciones := ExtractSections(raw, TiendanubeSections)
			if secciones.Payment != nil || secciones.Shipping != nil || secciones.Financial != nil || secciones.Fulfillment != nil {
				t.Fatal("un payload ilegible no debe producir secciones")
			}
		})
	}
}

func TestExtractSectionsSinCoincidenciasDevuelveNil(t *testing.T) {
	secciones := ExtractSections([]byte(`{"otra_cosa": 1}`), WooCommerceSections)

	if secciones.Payment != nil || secciones.Financial != nil {
		t.Fatal("si ninguna llave coincide la seccion queda nula, no un objeto vacio")
	}
}

func TestExtractSectionsWooYMeli(t *testing.T) {
	woo := ExtractSections([]byte(`{
		"status": "processing",
		"payment_method": "cod",
		"payment_method_title": "Pago contra entrega",
		"date_paid": null,
		"shipping_total": "8000",
		"currency": "COP",
		"total": "154000"
	}`), WooCommerceSections)

	pago := decodificar(t, woo.Payment)
	if pago["payment_method"] != "cod" {
		t.Fatalf("woo: seccion de pago inesperada: %v", pago)
	}
	if _, presente := pago["date_paid"]; presente {
		t.Fatal("woo: una llave nula no debe guardarse")
	}

	meli := ExtractSections([]byte(`{
		"status": "paid",
		"status_detail": null,
		"fulfilled": true,
		"currency_id": "COP",
		"total_amount": 120000,
		"paid_amount": 120000,
		"shipping_cost": 0,
		"payments": [{"status": "approved"}]
	}`), MeliSections)

	fulfillment := decodificar(t, meli.Fulfillment)
	if fulfillment["fulfilled"] != true || fulfillment["status"] != "paid" {
		t.Fatalf("meli: seccion de fulfillment inesperada: %v", fulfillment)
	}
	pagoMeli := decodificar(t, meli.Payment)
	if pagoMeli["payments"] == nil {
		t.Fatalf("meli: falta la lista de pagos: %v", pagoMeli)
	}
}

func TestExtractSectionsVTEX(t *testing.T) {
	raw := []byte(`{
		"orderId": "1234-01",
		"status": "invoiced",
		"statusDescription": "Faturado",
		"value": 154000,
		"totalFreight": 8000,
		"totalDiscount": 0,
		"lastChange": "2026-08-22T16:00:00",
		"shippingData": {"selectedSla": "Expresso"},
		"paymentData": {"transactions": []},
		"packageAttachment": null
	}`)

	secciones := ExtractSections(raw, VTEXSections)

	finanzas := decodificar(t, secciones.Financial)
	if finanzas["value"] != float64(154000) || finanzas["totalFreight"] != float64(8000) {
		t.Fatalf("vtex: seccion financiera inesperada: %v", finanzas)
	}

	envio := decodificar(t, secciones.Shipping)
	if envio["shippingData"] == nil {
		t.Fatalf("vtex: falta shippingData: %v", envio)
	}
	if _, presente := envio["packageAttachment"]; presente {
		t.Fatal("vtex: una llave nula no debe guardarse")
	}

	fulfillment := decodificar(t, secciones.Fulfillment)
	if fulfillment["status"] != "invoiced" || fulfillment["statusDescription"] != "Faturado" {
		t.Fatalf("vtex: seccion de fulfillment inesperada: %v", fulfillment)
	}

	pago := decodificar(t, secciones.Payment)
	if pago["paymentData"] == nil {
		t.Fatalf("vtex: falta paymentData: %v", pago)
	}
}

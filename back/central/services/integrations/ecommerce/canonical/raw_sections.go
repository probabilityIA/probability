package canonical

import "encoding/json"

// SectionKeys describe que llaves del JSON original del canal componen cada
// seccion de la orden. Root sirve para los canales que envuelven la orden en un
// objeto, como Jumpseller con "order".
type SectionKeys struct {
	Root        string
	Financial   []string
	Shipping    []string
	Payment     []string
	Fulfillment []string
}

// Sections son las secciones ya serializadas, listas para el DTO.
type Sections struct {
	Financial   []byte
	Shipping    []byte
	Payment     []byte
	Fulfillment []byte
}

// ExtractSections arma las secciones de la orden a partir del JSON original del
// canal. Es lo que permite consultar el pago o el envio tal como los reporto la
// tienda sin volver a leer el payload entero. Si el JSON no se puede interpretar
// devuelve secciones vacias: nunca es motivo para descartar la orden.
func ExtractSections(raw []byte, keys SectionKeys) Sections {
	if len(raw) == 0 {
		return Sections{}
	}

	var order map[string]interface{}
	if err := json.Unmarshal(raw, &order); err != nil {
		return Sections{}
	}

	if keys.Root != "" {
		anidado, ok := order[keys.Root].(map[string]interface{})
		if !ok {
			return Sections{}
		}
		order = anidado
	}

	return Sections{
		Financial:   subconjunto(order, keys.Financial),
		Shipping:    subconjunto(order, keys.Shipping),
		Payment:     subconjunto(order, keys.Payment),
		Fulfillment: subconjunto(order, keys.Fulfillment),
	}
}

func subconjunto(order map[string]interface{}, llaves []string) []byte {
	if len(llaves) == 0 {
		return nil
	}

	seccion := make(map[string]interface{}, len(llaves))
	for _, llave := range llaves {
		if valor, existe := order[llave]; existe && valor != nil {
			seccion[llave] = valor
		}
	}
	if len(seccion) == 0 {
		return nil
	}

	crudo, err := json.Marshal(seccion)
	if err != nil {
		return nil
	}
	return crudo
}

// Las llaves de cada perfil salen del JSON real que ya guardamos en
// order_channel_metadata, no de la documentacion.

var TiendanubeSections = SectionKeys{
	Financial: []string{
		"currency", "subtotal", "subtotal_without_taxes", "discount", "promotional_discount",
		"discount_coupon", "discount_gateway", "coupon", "total", "total_usd", "total_paid",
	},
	Shipping: []string{
		"shipping", "shipping_address", "shipping_option", "shipping_option_code",
		"shipping_option_reference", "shipping_suboption", "shipping_carrier_name",
		"shipping_cost_customer", "shipping_cost_owner", "shipping_min_days", "shipping_max_days",
		"shipping_pickup_type", "shipping_pickup_details", "shipping_store_branch_name",
		"has_shippable_products", "weight",
	},
	Payment: []string{
		"payment_status", "payment_details", "payment_count", "gateway", "gateway_id",
		"gateway_name", "gateway_link", "paid_at", "total_paid", "billing_address",
		"billing_name", "billing_document_type", "billing_phone",
	},
	Fulfillment: []string{
		"status", "shipping_status", "fulfillments", "shipped_at", "shipping_tracking_number",
		"shipping_tracking_url", "closed_at", "completed_at", "cancelled_at", "cancel_reason",
	},
}

var JumpsellerSections = SectionKeys{
	Root:      "order",
	Financial: []string{"currency", "subtotal", "discount", "tax", "total", "shipping_discount"},
	Shipping: []string{
		"shipping", "shipping_address", "shipping_method_id", "shipping_method_name",
		"shipping_required", "shipping_tax",
	},
	Payment:     []string{"payment_information", "payment_method_name", "payment_method_type", "billing_address"},
	Fulfillment: []string{"status", "shipment_status"},
}

var MeliSections = SectionKeys{
	Financial:   []string{"currency_id", "total_amount", "paid_amount", "coupon_amount", "coupon_id", "taxes"},
	Shipping:    []string{"shipping", "shipping_cost", "shipment_detail"},
	Payment:     []string{"payments", "paid_amount", "taxes"},
	Fulfillment: []string{"status", "status_detail", "fulfilled", "cancel_detail", "date_closed", "tags"},
}

var WooCommerceSections = SectionKeys{
	Financial: []string{
		"currency", "currency_symbol", "total", "total_tax", "discount_total", "discount_tax",
		"cart_tax", "prices_include_tax", "coupon_lines", "fee_lines", "tax_lines",
	},
	Shipping:    []string{"shipping", "shipping_lines", "shipping_total", "shipping_tax", "customer_note"},
	Payment:     []string{"payment_method", "payment_method_title", "transaction_id", "date_paid", "date_paid_gmt", "needs_payment", "billing"},
	Fulfillment: []string{"status", "date_completed", "date_completed_gmt", "needs_processing", "refunds"},
}

// VTEX es la excepcion: en produccion todavia no hay una orden real con su JSON
// guardado, asi que estas llaves salen de los tags del response del cliente
// (vtex_order_response.go). Conviene contrastarlas con el primer pedido real.
var VTEXSections = SectionKeys{
	Financial: []string{
		"value", "totalItems", "totalDiscount", "totalFreight", "totals",
		"ratesAndBenefitsData", "storePreferencesData",
	},
	Shipping:    []string{"shippingData", "packageAttachment", "sellers"},
	Payment:     []string{"paymentData", "clientProfileData"},
	Fulfillment: []string{"status", "statusDescription", "lastChange", "packageAttachment", "marketplaceOrderId"},
}

package mapper

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
)

func metaValue(meta []domain.WooCommerceMetaData, key string) string {
	for _, m := range meta {
		if m.Key != key {
			continue
		}
		switch v := m.Value.(type) {
		case string:
			return v
		case float64:
			return strconv.FormatInt(int64(v), 10)
		case int:
			return strconv.Itoa(v)
		case int64:
			return strconv.FormatInt(v, 10)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

// MapWooOrderToProbability convierte una orden WooCommerce a el DTO canónico de Probability.
// quoteRepo se usa para recuperar, por quote_id/rate_index, la tarifa exacta cotizada por
// Probability (flete, seguro minimo, comision y margen COD) cuando el metodo de envio
// vino del cotizador propio. Puede ser nil (metodo de envio ajeno o sin datos de quote).
func MapWooOrderToProbability(ctx context.Context, order *domain.WooCommerceOrder, rawJSON []byte, quoteRepo domain.IProductRepository, integrationID uint) *canonical.ProbabilityOrderDTO {
	now := time.Now()
	totalAmount := parseFloat(order.Total)
	totalTax := parseFloat(order.TotalTax)
	discount := parseFloat(order.DiscountTotal)
	shippingCost := parseFloat(order.ShippingTotal)
	freeShipping := shippingCost <= 0 && len(order.ShippingLines) > 0
	subtotal := totalAmount - totalTax - shippingCost + discount

	// Customer name from billing
	customerName := strings.TrimSpace(fmt.Sprintf("%s %s", order.Billing.FirstName, order.Billing.LastName))

	// Notes
	var notes *string
	if order.CustomerNote != "" {
		notes = &order.CustomerNote
	}

	// Coupon
	var coupon *string
	if len(order.CouponLines) > 0 {
		codes := make([]string, len(order.CouponLines))
		for i, cl := range order.CouponLines {
			codes[i] = cl.Code
		}
		joined := strings.Join(codes, ", ")
		coupon = &joined
	}

	// Map status
	status := MapWooStatus(order.Status)

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "woocommerce",
		Platform:        "woocommerce",
		ExternalID:      fmt.Sprintf("%d", order.ID),
		OrderNumber:     order.Number,
		Subtotal:        subtotal,
		Tax:             totalTax,
		Discount:        discount,
		ShippingCost:    shippingCost,
		FreeShipping:    freeShipping,
		// TotalAmount es solo el valor de productos, igual que en las ordenes
		// manuales (ver .claude/bitacora): el total con envio/comision se arma
		// sumando ShippingCost/CodTotal, no se guarda un total distinto aqui.
		TotalAmount:     subtotal,
		Currency:        order.Currency,
		CustomerName:    customerName,
		CustomerEmail:   order.Billing.Email,
		CustomerPhone:   order.Billing.Phone,
		Status:          status,
		OriginalStatus:  order.Status,
		Notes:           notes,
		Coupon:          coupon,
		OccurredAt:      order.DateCreated,
		ImportedAt:      now,
	}

	// Order items
	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.LineItems))
	for _, item := range order.LineItems {
		productID := fmt.Sprintf("%d", item.ProductID)
		unitPrice := item.Price
		totalPrice := parseFloat(item.Total)
		tax := parseFloat(item.TotalTax)
		sku := item.SKU

		var variantID *string
		if item.VariationID > 0 {
			vid := fmt.Sprintf("%d", item.VariationID)
			variantID = &vid
		}

		var imageURL *string
		if item.ImageURL != "" {
			imageURL = &item.ImageURL
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   sku,
			ProductName:  item.Name,
			ProductTitle: item.Name,
			VariantID:    variantID,
			Quantity:     item.Quantity,
			UnitPrice:    unitPrice,
			TotalPrice:   totalPrice,
			Currency:     order.Currency,
			Tax:          tax,
			ImageURL:     imageURL,
		})
	}

	// Addresses
	dto.Addresses = make([]canonical.ProbabilityAddressDTO, 0, 2)

	// Billing address
	dto.Addresses = append(dto.Addresses, canonical.ProbabilityAddressDTO{
		Type:       "billing",
		FirstName:  order.Billing.FirstName,
		LastName:   order.Billing.LastName,
		Company:    order.Billing.Company,
		Phone:      order.Billing.Phone,
		Street:     order.Billing.Address1,
		Street2:    order.Billing.Address2,
		City:       order.Billing.City,
		State:      order.Billing.State,
		Country:    order.Billing.Country,
		PostalCode: order.Billing.Postcode,
	})

	// Shipping address
	// WooCommerce deja shipping vacio cuando el comprador no pide una direccion
	// de envio distinta a la de facturacion: en ese caso se envia a la de billing.
	shippingAddress := canonical.ProbabilityAddressDTO{
		Type:       "shipping",
		FirstName:  order.Shipping.FirstName,
		LastName:   order.Shipping.LastName,
		Company:    order.Shipping.Company,
		Phone:      order.Shipping.Phone,
		Street:     order.Shipping.Address1,
		Street2:    order.Shipping.Address2,
		City:       order.Shipping.City,
		State:      order.Shipping.State,
		Country:    order.Shipping.Country,
		PostalCode: order.Shipping.Postcode,
	}
	if strings.TrimSpace(shippingAddress.Street) == "" && strings.TrimSpace(shippingAddress.City) == "" {
		shippingAddress.FirstName = order.Billing.FirstName
		shippingAddress.LastName = order.Billing.LastName
		shippingAddress.Company = order.Billing.Company
		shippingAddress.Phone = order.Billing.Phone
		shippingAddress.Street = order.Billing.Address1
		shippingAddress.Street2 = order.Billing.Address2
		shippingAddress.City = order.Billing.City
		shippingAddress.State = order.Billing.State
		shippingAddress.Country = order.Billing.Country
		shippingAddress.PostalCode = order.Billing.Postcode
	}
	dto.Addresses = append(dto.Addresses, shippingAddress)

	// Payment
	isCOD := false
	if order.PaymentMethod != "" {
		paymentStatus := "pending"
		var paidAt *time.Time
		if order.DatePaid != nil {
			paidAt = order.DatePaid
			paymentStatus = "completed"
		}

		gateway := order.PaymentMethod
		paymentMethodID := mapWooPaymentMethod(order.PaymentMethod)
		dto.Payments = append(dto.Payments, canonical.ProbabilityPaymentDTO{
			PaymentMethodID: paymentMethodID,
			Amount:          totalAmount,
			Currency:        order.Currency,
			Status:          paymentStatus,
			PaidAt:          paidAt,
			Gateway:         &gateway,
		})

		isCOD = paymentMethodID == paymentMethodCOD
	}

	// Shipments from shipping lines.
	// shippingCost (guia) empieza en el total crudo de Woo como fallback; si la linea trae
	// quote_id/rate_index de Probability, se corrige al valor real cotizado (flete + seguro
	// minimo, SIN comision de contra entrega), evitando que la comision quede contada dos
	// veces en shipping_cost/cod_total (ver .claude/rules/guias-contra-entrega.md).
	correctedShippingCost := shippingCost
	shippingCostResolved := false
	shippingLineDetails := make([]map[string]interface{}, 0, len(order.ShippingLines))
	for _, sl := range order.ShippingLines {
		carrier := sl.MethodTitle
		carrierCode := sl.MethodID
		shCost := parseFloat(sl.Total)

		dto.Shipments = append(dto.Shipments, canonical.ProbabilityShipmentDTO{
			Carrier:      &carrier,
			CarrierCode:  &carrierCode,
			Status:       mapShipmentStatus(order.Status),
			ShippingCost: &shCost,
		})

		source := sl.MethodID
		code := ""
		if quoteIDStr := metaValue(sl.MetaData, "quote_id"); quoteIDStr != "" {
			source = "probability"
			rateIndexStr := metaValue(sl.MetaData, "rate_index")
			code = "pq-" + quoteIDStr + "-" + rateIndexStr

			if quoteRepo != nil {
				quoteID64, errQ := strconv.ParseUint(quoteIDStr, 10, 64)
				rateIndex, errR := strconv.Atoi(rateIndexStr)
				if errQ == nil && errR == nil {
					if rate, err := quoteRepo.GetShippingQuoteRate(ctx, uint(quoteID64), rateIndex); err == nil && rate != nil {
						// guia = flete + seguro minimo + seguro extra + margen COD
						// (solo si la orden es contra entrega y la tarifa lo soporta),
						// ver .claude/rules/guias-contra-entrega.md.
						correctedShippingCost = rate.Flete + rate.MinimumInsurance + rate.ExtraInsurance
						if isCOD && rate.COD {
							correctedShippingCost += rate.CODProbabilityMargin
						}
						shippingCostResolved = true
					}
				}
			}
		}

		shippingLineDetails = append(shippingLineDetails, map[string]interface{}{
			"title":  sl.MethodTitle,
			"price":  sl.Total,
			"source": source,
			"code":   code,
		})
	}

	if shippingCostResolved {
		dto.ShippingCost = correctedShippingCost
	}

	if isCOD {
		codIncludesShipping := true
		if quoteRepo != nil {
			if v, err := quoteRepo.GetIntegrationCodIncludesShipping(ctx, integrationID); err == nil {
				codIncludesShipping = v
			}
		}

		if codIncludesShipping {
			// La tienda ya le cobro el envio (con comision incluida) al cliente en
			// el checkout: se recauda el total tal cual, igual que en ordenes
			// manuales donde shipping_cost es el valor completo cobrado.
			codTotal := totalAmount - totalTax
			dto.CodTotal = &codTotal
		} else {
			// Las guias de contra entrega se generan en Probability aparte: el
			// total del canal trae solo productos, se le suma el flete real de
			// la guia (sin comision, que se agrega despues al generarla).
			codTotal := subtotal + dto.ShippingCost
			dto.CodTotal = &codTotal
		}
	}

	if len(shippingLineDetails) > 0 {
		if sd, err := json.Marshal(map[string]interface{}{"shipping_lines": shippingLineDetails}); err == nil {
			dto.ShippingDetails = sd
		}
	}

	// Channel metadata with raw data
	if rawJSON != nil {
		secciones := canonical.ExtractSections(rawJSON, canonical.WooCommerceSections)
		dto.FinancialDetails = secciones.Financial
		dto.ShippingDetails = secciones.Shipping
		dto.PaymentDetails = secciones.Payment
		dto.FulfillmentDetails = secciones.Fulfillment

		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "woocommerce",
			RawData:       rawJSON,
			Version:       "v3",
			ReceivedAt:    now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	dto.Invoiceable = strings.EqualFold(order.Currency, "COP")

	return dto
}

// IDs del catalogo seed payment_methods (migration/shared).
const (
	paymentMethodCreditCard   uint = 1
	paymentMethodPaypal       uint = 3
	paymentMethodBankTransfer uint = 4
	paymentMethodCash         uint = 5
	paymentMethodCOD          uint = 6
	paymentMethodMercadoPago  uint = 7
	paymentMethodStripe       uint = 8
)

// mapWooPaymentMethod mapea el slug del gateway de WooCommerce al catalogo payment_methods.
func mapWooPaymentMethod(method string) uint {
	m := strings.ToLower(method)
	switch {
	case m == "cod" || strings.Contains(m, "contra"):
		return paymentMethodCOD
	case m == "bacs" || strings.Contains(m, "transfer"):
		return paymentMethodBankTransfer
	case m == "cheque" || strings.Contains(m, "cash"):
		return paymentMethodCash
	case strings.Contains(m, "paypal") || m == "ppcp-gateway":
		return paymentMethodPaypal
	case strings.Contains(m, "stripe"):
		return paymentMethodStripe
	case strings.Contains(m, "mercado"):
		return paymentMethodMercadoPago
	default:
		return paymentMethodCreditCard
	}
}

// mapWooStatus mapea el estado de WooCommerce al estado canónico de Probability.
func MapWooStatus(wooStatus string) string {
	switch wooStatus {
	case "pending", "checkout-draft", "processing", "addi-approved":
		return "pending"
	case "on-hold":
		return "on_hold"
	case "completed":
		return "completed"
	case "cancelled", "trash":
		return "cancelled"
	case "refunded":
		return "refunded"
	case "failed":
		return "failed"
	default:
		return "pending"
	}
}

// mapShipmentStatus mapea el estado de la orden WooCommerce a un estado de envío.
func mapShipmentStatus(wooStatus string) string {
	switch wooStatus {
	case "completed":
		return "delivered"
	case "processing":
		return "pending"
	case "on-hold":
		return "pending"
	case "cancelled", "refunded", "failed":
		return "cancelled"
	default:
		return "pending"
	}
}

func parseFloat(s string) float64 {
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

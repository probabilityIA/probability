package mapper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/canonical"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

const meliAPIBaseURL = "https://api.mercadolibre.com"

func MapMeliOrderToProbability(order *domain.MeliOrder, shippingDetail *domain.MeliShippingDetail, rawJSON []byte, shipmentRaw []byte) *canonical.ProbabilityOrderDTO {
	now := time.Now()

	shippingCost := 0.0
	if shippingDetail != nil && shippingDetail.ShippingOption != nil {
		shippingCost = shippingDetail.ShippingOption.Cost
	}

	subtotal := order.TotalAmount + order.CouponAmount

	customerName := strings.TrimSpace(fmt.Sprintf("%s %s", order.Buyer.FirstName, order.Buyer.LastName))
	if customerName == "" {
		customerName = order.Buyer.Nickname
	}
	if customerName == "" && shippingDetail != nil && shippingDetail.ReceiverAddress != nil {
		customerName = shippingDetail.ReceiverAddress.ReceiverName
	}

	customerPhone := ""
	if order.Buyer.Phone.Number != "" {
		if order.Buyer.Phone.AreaCode != "" {
			customerPhone = order.Buyer.Phone.AreaCode + order.Buyer.Phone.Number
		} else {
			customerPhone = order.Buyer.Phone.Number
		}
	}
	if customerPhone == "" && shippingDetail != nil && shippingDetail.ReceiverAddress != nil {
		customerPhone = sanitizePhone(shippingDetail.ReceiverAddress.ReceiverPhone)
	}

	customerDNI := ""
	if order.Buyer.BillingInfo != nil {
		customerDNI = order.Buyer.BillingInfo.DocNumber
	}

	var coupon *string
	if order.CouponID != nil && *order.CouponID != "" {
		coupon = order.CouponID
	}

	status := mapMeliOrderStatus(order.Status)

	packID := ""
	if order.PackID != nil && *order.PackID > 0 {
		packID = fmt.Sprintf("%d", *order.PackID)
	}

	dto := &canonical.ProbabilityOrderDTO{
		IntegrationType: "mercado_libre",
		Platform:        "mercadolibre",
		ExternalID:      fmt.Sprintf("%d", order.ID),
		ChannelPackID:   packID,
		OrderNumber:     fmt.Sprintf("%d", order.ID),
		Subtotal:        subtotal,
		Tax:             0,
		Discount:        order.CouponAmount,
		ShippingCost:    shippingCost,
		TotalAmount:     order.TotalAmount,
		Currency:        order.CurrencyID,
		CustomerName:    customerName,
		CustomerEmail:   order.Buyer.Email,
		CustomerPhone:   customerPhone,
		CustomerDNI:     customerDNI,
		Status:          status,
		OriginalStatus:  order.Status,
		Coupon:          coupon,
		OrderTypeName:   deriveOrderType(shippingDetail),
		OccurredAt:      order.DateCreated,
		ImportedAt:      now,
	}

	dto.OrderItems = make([]canonical.ProbabilityOrderItemDTO, 0, len(order.OrderItems))
	for _, item := range order.OrderItems {
		productID := item.Item.ID

		var variantID *string
		if item.Item.VariationID != nil {
			vid := fmt.Sprintf("%d", *item.Item.VariationID)
			variantID = &vid
		}

		sku := ""
		if item.Item.SellerSKU != nil {
			sku = *item.Item.SellerSKU
		}

		dto.OrderItems = append(dto.OrderItems, canonical.ProbabilityOrderItemDTO{
			ProductID:    &productID,
			ProductSKU:   sku,
			ProductName:  item.Item.Title,
			ProductTitle: item.Item.Title,
			VariantID:    variantID,
			Quantity:     item.Quantity,
			UnitPrice:    item.UnitPrice,
			TotalPrice:   item.UnitPrice * float64(item.Quantity),
			Currency:     item.Currency,
		})
	}

	if shippingDetail != nil && shippingDetail.ReceiverAddress != nil {
		addr := shippingDetail.ReceiverAddress
		street := addr.StreetName
		if addr.StreetNumber != "" {
			street = addr.StreetName + " " + addr.StreetNumber
		}

		city := addr.City.Name
		if strings.EqualFold(addr.State.Name, "Bogota D.C.") || strings.EqualFold(addr.State.Name, "Bogotá D.C.") {
			if addr.City.Name != "" {
				street = strings.TrimSpace(street + " - " + addr.City.Name)
			}
			city = addr.State.Name
		}

		address := canonical.ProbabilityAddressDTO{
			Type:       "shipping",
			FirstName:  firstNonEmptyStr(addr.ReceiverName, customerName),
			Phone:      customerPhone,
			Street:     street,
			City:       city,
			State:      addr.State.Name,
			Country:    addr.Country.Name,
			PostalCode: addr.ZipCode,
			Latitude:   addr.Latitude,
			Longitude:  addr.Longitude,
		}
		if addr.Comment != "" {
			instructions := addr.Comment
			address.Instructions = &instructions
		}
		dto.Addresses = append(dto.Addresses, address)
	}

	dto.Payments = make([]canonical.ProbabilityPaymentDTO, 0, len(order.Payments))
	for _, p := range order.Payments {
		gateway := p.PaymentMethodID
		paymentStatus := mapMeliPaymentStatus(p.Status)

		dto.Payments = append(dto.Payments, canonical.ProbabilityPaymentDTO{
			Amount:   p.TransactionAmount,
			Currency: p.CurrencyID,
			Status:   paymentStatus,
			PaidAt:   p.DateApproved,
			Gateway:  &gateway,
		})
	}

	if shippingDetail != nil {
		carrier := "mercadoenvios"
		shStatus := mapMeliShippingStatus(shippingDetail.Status)

		shipment := canonical.ProbabilityShipmentDTO{
			Carrier:      &carrier,
			Status:       shStatus,
			ShippingCost: &shippingCost,
		}

		if shippingDetail.ID > 0 {
			guideID := fmt.Sprintf("%d", shippingDetail.ID)
			shipment.GuideID = &guideID
			guideURL := fmt.Sprintf("%s/shipment_labels?shipment_ids=%d&response_type=pdf", meliAPIBaseURL, shippingDetail.ID)
			shipment.GuideURL = &guideURL
		}

		if shippingDetail.TrackingNumber != "" {
			tn := shippingDetail.TrackingNumber
			shipment.TrackingNumber = &tn
		}

		if shippingDetail.ShippingOption != nil {
			optName := shippingDetail.ShippingOption.Name
			shipment.CarrierCode = &optName
			if shippingDetail.ShippingOption.EstimatedDeliveryTime != nil {
				shipment.EstimatedDelivery = shippingDetail.ShippingOption.EstimatedDeliveryTime.Date
			}
		}

		dto.Shipments = append(dto.Shipments, shipment)
	}

	if rawJSON != nil {
		secciones := canonical.ExtractSections(rawJSON, canonical.MeliSections)
		dto.FinancialDetails = secciones.Financial
		dto.ShippingDetails = secciones.Shipping
		dto.PaymentDetails = secciones.Payment
		dto.FulfillmentDetails = secciones.Fulfillment

		dto.ChannelMetadata = &canonical.ProbabilityChannelMetadataDTO{
			ChannelSource: "mercadolibre",
			RawData:       withPackOrderIDs(withShipmentRaw(rawJSON, shipmentRaw), order.PackOrderIDs),
			Version:       "v2",
			ReceivedAt:    now,
			IsLatest:      true,
			SyncStatus:    "synced",
		}
	}

	return dto
}

func deriveOrderType(shippingDetail *domain.MeliShippingDetail) string {
	if shippingDetail == nil {
		return "marketplace"
	}
	switch shippingDetail.LogisticType {
	case "fulfillment":
		return "full"
	case "self_service":
		return "flex"
	case "cross_docking":
		return "cross_docking"
	case "drop_off", "xd_drop_off":
		return "mercado_envios"
	}
	switch shippingDetail.LogisticMode {
	case "me1":
		return "me1"
	case "me2":
		return "mercado_envios"
	}
	return "marketplace"
}

func mapMeliOrderStatus(meliStatus string) string {
	switch meliStatus {
	case "confirmed":
		return "pending"
	case "payment_required":
		return "pending"
	case "payment_in_process":
		return "pending"
	case "paid":
		return "paid"
	case "partially_paid":
		return "paid"
	case "cancelled":
		return "cancelled"
	case "invalid":
		return "cancelled"
	default:
		return meliStatus
	}
}

func mapMeliPaymentStatus(meliStatus string) string {
	switch meliStatus {
	case "approved":
		return "paid"
	case "pending", "in_process", "in_mediation":
		return "pending"
	case "rejected":
		return "failed"
	case "refunded", "charged_back":
		return "refunded"
	case "cancelled":
		return "cancelled"
	default:
		return meliStatus
	}
}

func mapMeliShippingStatus(meliStatus string) string {
	switch meliStatus {
	case "pending", "handling", "ready_to_ship":
		return "pending"
	case "shipped":
		return "shipped"
	case "delivered":
		return "delivered"
	case "not_delivered":
		return "failed"
	case "cancelled":
		return "cancelled"
	default:
		return "pending"
	}
}

// withShipmentRaw adjunta la respuesta cruda de /shipments/{id} al JSON de la
// orden bajo la clave shipment_detail. MercadoLibre no manda la direccion del
// comprador en la orden (solo shipping.id), asi que sin esto el JSON guardado
// queda sin el destino.
func withShipmentRaw(orderRaw []byte, shipmentRaw []byte) []byte {
	if len(shipmentRaw) == 0 {
		return orderRaw
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(orderRaw, &merged); err != nil {
		return orderRaw
	}
	merged["shipment_detail"] = json.RawMessage(shipmentRaw)

	out, err := json.Marshal(merged)
	if err != nil {
		return orderRaw
	}
	return out
}

// sanitizePhone descarta los telefonos enmascarados que manda MercadoLibre
// (llegan como "XXXXXXX" cuando el comprador tiene los datos protegidos).
func sanitizePhone(phone string) string {
	digits := 0
	for _, r := range phone {
		if r >= '0' && r <= '9' {
			digits++
		}
	}
	if digits < 7 {
		return ""
	}
	return phone
}

func firstNonEmptyStr(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

// withPackOrderIDs deja en el raw los ids reales de las ordenes del carrito.
// La orden consolidada se guarda con el pack_id como numero, que no es
// consultable en la API de MeLi; sin estos ids no hay como volver a pedirlas.
func withPackOrderIDs(orderRaw []byte, packOrderIDs []int64) []byte {
	if len(packOrderIDs) == 0 {
		return orderRaw
	}

	ids, err := json.Marshal(packOrderIDs)
	if err != nil {
		return orderRaw
	}

	var merged map[string]json.RawMessage
	if err := json.Unmarshal(orderRaw, &merged); err != nil {
		return orderRaw
	}
	merged["pack_order_ids"] = json.RawMessage(ids)

	out, err := json.Marshal(merged)
	if err != nil {
		return orderRaw
	}
	return out
}

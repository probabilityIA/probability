package repository

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
)

func (r *Repository) GetGuidePDFContext(ctx context.Context, shipmentID uint) (*domain.GuidePDFContext, error) {
	var row struct {
		ID                 uint
		TrackingNumber     *string
		Carrier            *string
		Weight             *float64
		Height             *float64
		Width              *float64
		Length             *float64
		CreatedAt          *time.Time
		EstimatedDelivery  *time.Time
		DestinationAddress string
		DestinationCity   string
		DestinationState   string
		DestinationSuburb string
		CodTotal           *float64
		CodCarrierFee      *float64
		OrderNumber        *string
		CustomerName       *string
		CustomerPhone      *string
		CustomerDNI        *string
		CustomerEmail      *string
		TotalAmount        *float64
		Currency           *string
		BusinessName       *string
		BusinessAddress    *string
		WName              *string
		WCompany           *string
		WFirst             *string
		WLast              *string
		WAddress           *string
		WStreet            *string
		WCity              *string
		WState             *string
		WPhone             *string
		WPostal            *string
		GuideURL           *string
		Metadata           map[string]interface{}
	}

	var items []struct {
		SKU         *string
		ProductName *string
		Quantity    *int
		UnitPrice   *float64
	}

	err := r.db.Conn(ctx).Raw(`
		SELECT
			s.id,
			s.tracking_number,
			s.carrier,
			s.weight, s.height, s.width, s.length,
			s.created_at,
			s.estimated_delivery,
			s.destination_address, s.destination_city, s.destination_state, s.destination_suburb,
			o.cod_total,
			s.cod_carrier_fee,
			o.order_number,
			o.customer_name,
			o.customer_phone,
			o.customer_dni,
			o.customer_email,
			o.total_amount,
			o.currency,
			b.name AS business_name,
			b.address AS business_address,
			COALESCE(w.name, oa.alias) AS w_name,
			COALESCE(w.company, oa.company) AS w_company,
			COALESCE(w.first_name, oa.first_name) AS w_first,
			COALESCE(w.last_name, oa.last_name) AS w_last,
			w.address AS w_address,
			COALESCE(w.street, oa.street) AS w_street,
			COALESCE(w.city, oa.city) AS w_city,
			COALESCE(w.state, oa.state) AS w_state,
			COALESCE(w.phone, oa.phone) AS w_phone,
			COALESCE(w.zip_code, oa.postal_code) AS w_postal,
			s.guide_url,
			s.metadata
		FROM shipments s
		LEFT JOIN orders o ON o.id = s.order_id
		LEFT JOIN business b ON b.id = o.business_id
		LEFT JOIN warehouses w ON w.id = s.warehouse_id
		LEFT JOIN origin_address oa ON oa.business_id = o.business_id AND oa.is_default = true AND oa.deleted_at IS NULL
		WHERE s.id = ? AND s.deleted_at IS NULL
	`, shipmentID).Scan(&row).Error

	if err == nil && row.ID != 0 {
		itemsErr := r.db.Conn(ctx).Raw(`
			SELECT
				p.sku,
				p.name AS product_name,
				oi.quantity,
				oi.unit_price
			FROM order_items oi
			LEFT JOIN products p ON p.id = oi.product_id
			WHERE oi.order_id = ? AND oi.deleted_at IS NULL
			ORDER BY oi.created_at ASC
			LIMIT 5
		`, row.OrderNumber).Scan(&items).Error
		if itemsErr != nil {
			items = []struct {
				SKU         *string
				ProductName *string
				Quantity    *int
				UnitPrice   *float64
			}{}
		}
	}
	if err != nil {
		return nil, err
	}
	if row.ID == 0 {
		return nil, nil
	}

	val := func(p *string) string {
		if p == nil {
			return ""
		}
		return strings.TrimSpace(*p)
	}
	valF := func(p *float64) float64 {
		if p == nil {
			return 0
		}
		return *p
	}
	contact := strings.TrimSpace(val(row.WFirst) + " " + val(row.WLast))
	wAddr := val(row.WAddress)
	if wAddr == "" {
		wAddr = val(row.WStreet)
	}

	metaStr := func(key string) string {
		if row.Metadata == nil {
			return ""
		}
		if v, ok := row.Metadata[key]; ok {
			if s, ok := v.(string); ok {
				return strings.TrimSpace(s)
			}
		}
		return ""
	}

	orderItems := make([]domain.OrderItemContext, 0)
	for _, item := range items {
		sku := ""
		if item.SKU != nil {
			sku = strings.TrimSpace(*item.SKU)
		}
		pName := ""
		if item.ProductName != nil {
			pName = strings.TrimSpace(*item.ProductName)
		}
		qty := 1
		if item.Quantity != nil && *item.Quantity > 0 {
			qty = *item.Quantity
		}
		price := 0.0
		if item.UnitPrice != nil {
			price = *item.UnitPrice
		}
		orderItems = append(orderItems, domain.OrderItemContext{
			SKU:         sku,
			ProductName: pName,
			Quantity:    qty,
			UnitPrice:   price,
		})
	}

	result := &domain.GuidePDFContext{
		ShipmentID:         row.ID,
		TrackingNumber:     val(row.TrackingNumber),
		Carrier:            val(row.Carrier),
		Weight:             valF(row.Weight),
		Height:             valF(row.Height),
		Width:              valF(row.Width),
		Length:             valF(row.Length),
		CreatedAt:          row.CreatedAt,
		EstimatedDelivery:  row.EstimatedDelivery,
		OrderNumber:        val(row.OrderNumber),
		CustomerName:       val(row.CustomerName),
		CustomerPhone:      val(row.CustomerPhone),
		CustomerDNI:        val(row.CustomerDNI),
		CustomerEmail:      val(row.CustomerEmail),
		DestinationAddress: row.DestinationAddress,
		DestinationCity:    row.DestinationCity,
		DestinationState:   row.DestinationState,
		DestinationSuburb:  row.DestinationSuburb,
		DeclaredValue:      valF(row.TotalAmount),
		Currency:           val(row.Currency),
		CodTotal:           valF(row.CodTotal),
		CodCarrierFee:      valF(row.CodCarrierFee),
		BusinessName:       val(row.BusinessName),
		BusinessAddress:    val(row.BusinessAddress),
		WarehouseName:      val(row.WName),
		WarehouseCompany:   val(row.WCompany),
		WarehouseContact:   contact,
		WarehouseAddress:   wAddr,
		WarehouseCity:      val(row.WCity),
		WarehouseState:     val(row.WState),
		WarehousePhone:     val(row.WPhone),
		WarehousePostal:    firstNonEmpty(metaStr("postal_origen"), val(row.WPostal)),
		Origen:             metaStr("origen"),
		AsCode:             metaStr("as_code"),
		Paq:                metaStr("paq"),
		Unidad:             metaStr("unidad"),
		Destino:            metaStr("destino"),
		ZonaHub:            metaStr("zona_hub"),
		EquipoReparto:      metaStr("equipo_reparto"),
		Ref:                metaStr("ref"),
		Guia:               metaStr("guia"),
		Observaciones:      metaStr("observaciones"),
		OrderItems:         orderItems,
	}

	return result, nil
}

package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

const (
	orderStatusSourceCarrier = "carrier"
	orderStatusCarrierLabel  = "Transportadora"
)

func (r *Repository) GetOrderBusinessID(ctx context.Context, orderUUID string) (uint, error) {
	var result struct {
		BusinessID uint `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select("business_id").
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return 0, err
	}
	if result.BusinessID == 0 {
		return 0, fmt.Errorf("orden %s no encontrada o sin business asociado", orderUUID)
	}

	return result.BusinessID, nil
}

func (r *Repository) GetShipmentBusinessIDByTracking(ctx context.Context, trackingNumber string) (uint, error) {
	var result struct {
		BusinessID uint `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("shipments s").
		Select("o.business_id").
		Joins("INNER JOIN orders o ON s.order_id = o.id").
		Where("s.tracking_number = ?", trackingNumber).
		Where("s.deleted_at IS NULL").
		Where("o.deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("envío con tracking %s no encontrado", trackingNumber)
		}
		return 0, err
	}
	if result.BusinessID == 0 {
		return 0, fmt.Errorf("envío con tracking %s no encontrado o sin business asociado", trackingNumber)
	}

	return result.BusinessID, nil
}

func (r *Repository) GetShipmentBusinessIDByID(ctx context.Context, shipmentID uint) (uint, error) {
	var result struct {
		BusinessID uint `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("shipments s").
		Select("o.business_id").
		Joins("INNER JOIN orders o ON s.order_id = o.id").
		Where("s.id = ?", shipmentID).
		Where("s.deleted_at IS NULL").
		Where("o.deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("envío con ID %d no encontrado", shipmentID)
		}
		return 0, err
	}
	if result.BusinessID == 0 {
		return 0, fmt.Errorf("envío con ID %d no encontrado o sin business asociado", shipmentID)
	}

	return result.BusinessID, nil
}

func (r *Repository) GetOrderIntegrationID(ctx context.Context, orderUUID string) (uint, error) {
	var result struct {
		IntegrationID uint `gorm:"column:integration_id"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select("integration_id").
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return 0, err
	}

	return result.IntegrationID, nil
}

func (r *Repository) GetOrderCodTotal(ctx context.Context, orderUUID string) (*float64, error) {
	var result struct {
		CodTotal *float64 `gorm:"column:cod_total"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select("cod_total").
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return nil, err
	}

	return result.CodTotal, nil
}

func (r *Repository) GetOrderCodBasis(ctx context.Context, orderUUID string) (*domain.OrderCodBasis, error) {
	var result struct {
		TotalAmount         float64  `gorm:"column:total_amount"`
		CodTotal            *float64 `gorm:"column:cod_total"`
		CodIncludesShipping bool     `gorm:"column:cod_includes_shipping"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select("total_amount, cod_total, cod_includes_shipping").
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return nil, err
	}

	codTotal := 0.0
	if result.CodTotal != nil {
		codTotal = *result.CodTotal
	}

	return &domain.OrderCodBasis{
		TotalAmount:         result.TotalAmount,
		CodTotal:            codTotal,
		CodIncludesShipping: result.CodIncludesShipping,
	}, nil
}

func (r *Repository) GetOrderExternalGuide(ctx context.Context, orderUUID string) (*domain.OrderExternalGuide, error) {
	var result struct {
		Platform       string  `gorm:"column:platform"`
		TrackingNumber *string `gorm:"column:tracking_number"`
		GuideID        *string `gorm:"column:guide_id"`
		GuideLink      *string `gorm:"column:guide_link"`
		Carrier        *string `gorm:"column:carrier"`
		CustomerName   string  `gorm:"column:customer_name"`
		ShippingStreet string  `gorm:"column:shipping_street"`
		ShippingCity   string  `gorm:"column:shipping_city"`
		ShippingState  string  `gorm:"column:shipping_state"`
		ShippingCost   float64 `gorm:"column:shipping_cost"`
		IsTest         bool    `gorm:"column:is_test"`
		Status         string  `gorm:"column:status"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select(`platform, tracking_number, guide_id, guide_link, carrier, customer_name,
			shipping_street, shipping_city, shipping_state, COALESCE(shipping_cost,0) AS shipping_cost,
			COALESCE(is_test,false) AS is_test, status`).
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("orden %s no encontrada", orderUUID)
		}
		return nil, err
	}

	return &domain.OrderExternalGuide{
		Platform:       result.Platform,
		TrackingNumber: strings.TrimSpace(derefStr(result.TrackingNumber)),
		GuideID:        strings.TrimSpace(derefStr(result.GuideID)),
		GuideURL:       strings.TrimSpace(derefStr(result.GuideLink)),
		Carrier:        strings.TrimSpace(derefStr(result.Carrier)),
		CustomerName:   result.CustomerName,
		ShippingStreet: result.ShippingStreet,
		ShippingCity:   result.ShippingCity,
		ShippingState:  result.ShippingState,
		ShippingCost:   result.ShippingCost,
		IsTest:         result.IsTest,
		Status:         result.Status,
	}, nil
}

func derefStr(v *string) string {
	if v == nil {
		return ""
	}
	return *v
}

func (r *Repository) UpdateOrderGuideLink(ctx context.Context, orderID string, guideLink string, trackingNumber string, carrier string, shippingCost float64) error {
	updates := map[string]interface{}{}
	if guideLink != "" {
		updates["guide_link"] = guideLink
	}
	if trackingNumber != "" {
		updates["tracking_number"] = trackingNumber
	}
	if carrier != "" {
		updates["carrier"] = carrier
	}
	if shippingCost > 0 {
		updates["shipping_cost"] = shippingCost
		updates["cod_total"] = gorm.Expr(
			"CASE WHEN cod_total > 0 AND cod_includes_shipping = false THEN total_amount + ? ELSE cod_total END",
			shippingCost,
		)
	}
	if len(updates) == 0 {
		return nil
	}

	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Where("deleted_at IS NULL").
		Updates(updates).Error
}

func (r *Repository) UpdateOrderStatusByOrderID(ctx context.Context, orderID string, status string) error {
	if orderID == "" || status == "" {
		return nil
	}

	var current struct {
		Status string `gorm:"column:status"`
	}
	_ = r.db.Conn(ctx).
		Table("orders").
		Select("status").
		Where("id = ? AND deleted_at IS NULL", orderID).
		Limit(1).
		Scan(&current).Error

	updates := map[string]any{
		"status":            status,
		"status_source":     orderStatusSourceCarrier,
		"status_changed_by": orderStatusCarrierLabel,
		"status_changed_at": time.Now(),
	}

	var statusID struct{ ID uint }
	err := r.db.Conn(ctx).
		Table("order_statuses").
		Select("id").
		Where("code = ? AND deleted_at IS NULL", status).
		Limit(1).
		Scan(&statusID).Error
	if err == nil && statusID.ID > 0 {
		updates["status_id"] = statusID.ID
	}

	if err := r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ?", orderID).
		Where("deleted_at IS NULL").
		Updates(updates).Error; err != nil {
		return err
	}

	if current.Status == status {
		return nil
	}

	return r.db.Conn(ctx).Exec(`
		INSERT INTO order_history (created_at, updated_at, order_id, previous_status, new_status, changed_by_name, source)
		VALUES (NOW(), NOW(), ?, ?, ?, ?, ?)`,
		orderID, current.Status, status, orderStatusCarrierLabel, orderStatusSourceCarrier).Error
}

func (r *Repository) ClearOrderGuideData(ctx context.Context, orderID string) error {
	if orderID == "" {
		return nil
	}
	return r.db.Conn(ctx).
		Model(&models.Order{}).
		Where("id = ? AND deleted_at IS NULL", orderID).
		Updates(map[string]any{
			"tracking_number": nil,
			"tracking_link":   nil,
			"guide_link":      nil,
			"guide_id":        nil,
			"carrier":         nil,
			"status":          "pending",
		}).Error
}

func (r *Repository) EnsureAllBusinessesActive(ctx context.Context) error {
	return r.db.Conn(ctx).
		Table("business").
		Where("deleted_at IS NULL").
		Update("is_active", true).Error
}

func (r *Repository) GetOrderRecipient(ctx context.Context, orderUUID string) (*domain.OrderRecipient, error) {
	var result struct {
		CustomerName         string  `gorm:"column:customer_name"`
		CustomerFirstName    string  `gorm:"column:customer_first_name"`
		CustomerLastName     string  `gorm:"column:customer_last_name"`
		CustomerEmail        string  `gorm:"column:customer_email"`
		CustomerPhone        string  `gorm:"column:customer_phone"`
		ShippingStreet       string  `gorm:"column:shipping_street"`
		ShippingCity         string  `gorm:"column:shipping_city"`
		ShippingState        string  `gorm:"column:shipping_state"`
		ShippingNeighborhood string  `gorm:"column:shipping_neighborhood"`
		BusinessID           *uint   `gorm:"column:business_id"`
	}

	err := r.db.Conn(ctx).
		Table("orders").
		Select(`COALESCE(customer_name,'') AS customer_name,
			COALESCE(customer_first_name,'') AS customer_first_name,
			COALESCE(customer_last_name,'') AS customer_last_name,
			COALESCE(customer_email,'') AS customer_email,
			COALESCE(customer_phone,'') AS customer_phone,
			COALESCE(shipping_street,'') AS shipping_street,
			COALESCE(shipping_city,'') AS shipping_city,
			COALESCE(shipping_state,'') AS shipping_state,
			COALESCE(shipping_neighborhood,'') AS shipping_neighborhood,
			business_id`).
		Where("id = ?", orderUUID).
		Where("deleted_at IS NULL").
		Limit(1).
		Scan(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if result.CustomerName == "" && result.ShippingStreet == "" && result.BusinessID == nil {
		return nil, nil
	}

	first := strings.TrimSpace(result.CustomerFirstName)
	last := strings.TrimSpace(result.CustomerLastName)
	if first == "" && last == "" {
		parts := strings.Fields(strings.TrimSpace(result.CustomerName))
		if len(parts) > 0 {
			first = parts[0]
		}
		if len(parts) > 1 {
			last = strings.Join(parts[1:], " ")
		}
	}

	return &domain.OrderRecipient{
		BusinessID:   result.BusinessID,
		FullName:     strings.TrimSpace(result.CustomerName),
		FirstName:    first,
		LastName:     last,
		Email:        strings.TrimSpace(result.CustomerEmail),
		Phone:        strings.TrimSpace(result.CustomerPhone),
		Address:      strings.TrimSpace(result.ShippingStreet),
		City:         strings.TrimSpace(result.ShippingCity),
		State:        strings.TrimSpace(result.ShippingState),
		Neighborhood: strings.TrimSpace(result.ShippingNeighborhood),
	}, nil
}

func (r *Repository) GetUserDisplayName(ctx context.Context, userID uint) string {
	if userID == 0 {
		return ""
	}
	var row struct {
		Name  string `gorm:"column:name"`
		Email string `gorm:"column:email"`
	}
	err := r.db.Conn(ctx).
		Table(`"user"`).
		Select("COALESCE(name,'') AS name, COALESCE(email,'') AS email").
		Where("id = ?", userID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return ""
	}
	if name := strings.TrimSpace(row.Name); name != "" {
		return name
	}
	return strings.TrimSpace(row.Email)
}

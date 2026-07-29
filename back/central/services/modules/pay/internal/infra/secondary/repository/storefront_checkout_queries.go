package repository

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/pay/internal/domain/entities"
	"gorm.io/gorm"
)

type publicCheckoutRow struct {
	ID              string
	BusinessID      uint
	IntegrationID   uint
	Slug            string
	Reference       string
	Status          string
	Amount          float64
	Currency        string
	Items           json.RawMessage
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	CustomerDni     string
	ShippingAddress json.RawMessage
}

func (r *Repository) GetStorefrontCheckoutByReference(ctx context.Context, reference string) (*entities.StorefrontCheckoutSnapshot, error) {
	var row publicCheckoutRow
	err := r.db.Conn(ctx).
		Table("public_checkouts").
		Select("id, business_id, integration_id, slug, reference, status, amount, currency, items, customer_name, customer_email, customer_phone, customer_dni, shipping_address").
		Where("reference = ?", reference).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if row.ID == "" {
		return nil, nil
	}
	return &entities.StorefrontCheckoutSnapshot{
		ID:                  row.ID,
		BusinessID:          row.BusinessID,
		IntegrationID:       row.IntegrationID,
		Slug:                row.Slug,
		Reference:           row.Reference,
		Status:              row.Status,
		Amount:              row.Amount,
		Currency:            row.Currency,
		ItemsJSON:           row.Items,
		CustomerName:        row.CustomerName,
		CustomerEmail:       row.CustomerEmail,
		CustomerPhone:       row.CustomerPhone,
		CustomerDni:         row.CustomerDni,
		ShippingAddressJSON: row.ShippingAddress,
	}, nil
}

func (r *Repository) MarkStorefrontCheckoutStatus(ctx context.Context, reference, status string) error {
	updates := map[string]interface{}{"status": status, "updated_at": time.Now()}
	if status == "paid" {
		updates["paid_at"] = time.Now()
	}
	return r.db.Conn(ctx).Table("public_checkouts").Where("reference = ?", reference).Updates(updates).Error
}

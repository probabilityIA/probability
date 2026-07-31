package repository

import (
	"context"
	"encoding/json"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (r *Repository) GetTiendaIntegrationID(ctx context.Context, businessID uint) (uint, error) {
	var result struct {
		ID uint
	}
	err := r.db.Conn(ctx).
		Table("integrations").
		Select("id").
		Where("business_id = ? AND integration_type_id = ? AND is_active = true AND deleted_at IS NULL", businessID, 31).
		Limit(1).
		Scan(&result).Error
	if err != nil {
		return 0, err
	}
	if result.ID == 0 {
		return 0, domainerrors.ErrBusinessNotFound
	}
	return result.ID, nil
}

func checkoutToModel(c *entities.PublicCheckout) (*models.PublicCheckout, error) {
	itemsJSON, err := json.Marshal(c.Items)
	if err != nil {
		return nil, err
	}
	var addressJSON []byte
	if c.ShippingAddress != nil {
		addressJSON, err = json.Marshal(c.ShippingAddress)
		if err != nil {
			return nil, err
		}
	}
	return &models.PublicCheckout{
		BusinessID:      c.BusinessID,
		IntegrationID:   c.IntegrationID,
		Slug:            c.Slug,
		Reference:       c.Reference,
		Status:          c.Status,
		Amount:          c.Amount,
		Currency:        c.Currency,
		Items:           datatypes.JSON(itemsJSON),
		CustomerName:    c.CustomerName,
		CustomerEmail:   c.CustomerEmail,
		CustomerPhone:   c.CustomerPhone,
		CustomerDni:     c.CustomerDni,
		ShippingAddress: datatypes.JSON(addressJSON),
		BoldOrderID:     c.BoldOrderID,
	}, nil
}

func checkoutToEntity(m *models.PublicCheckout) *entities.PublicCheckout {
	var items []entities.CheckoutItem
	_ = json.Unmarshal(m.Items, &items)
	var address *entities.CheckoutAddress
	if len(m.ShippingAddress) > 0 {
		address = &entities.CheckoutAddress{}
		_ = json.Unmarshal(m.ShippingAddress, address)
	}
	return &entities.PublicCheckout{
		ID:              m.ID.String(),
		BusinessID:      m.BusinessID,
		IntegrationID:   m.IntegrationID,
		Slug:            m.Slug,
		Reference:       m.Reference,
		Status:          m.Status,
		Amount:          m.Amount,
		Currency:        m.Currency,
		Items:           items,
		CustomerName:    m.CustomerName,
		CustomerEmail:   m.CustomerEmail,
		CustomerPhone:   m.CustomerPhone,
		CustomerDni:     m.CustomerDni,
		ShippingAddress: address,
		BoldOrderID:     m.BoldOrderID,
		CreatedAt:       m.CreatedAt,
		PaidAt:          m.PaidAt,
	}
}

func (r *Repository) CreateCheckout(ctx context.Context, checkout *entities.PublicCheckout) error {
	model, err := checkoutToModel(checkout)
	if err != nil {
		return err
	}
	if err := r.db.Conn(ctx).Create(model).Error; err != nil {
		return err
	}
	checkout.ID = model.ID.String()
	checkout.CreatedAt = model.CreatedAt
	return nil
}

func (r *Repository) GetCheckoutByReference(ctx context.Context, reference string) (*entities.PublicCheckout, error) {
	var model models.PublicCheckout
	err := r.db.Conn(ctx).Where("reference = ?", reference).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrCheckoutNotFound
		}
		return nil, err
	}
	return checkoutToEntity(&model), nil
}

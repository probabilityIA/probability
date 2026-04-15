package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/customers/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

func (r *Repository) UpsertCustomerAddress(ctx context.Context, address *entities.CustomerAddress) error {
	var existing models.CustomerAddress
	err := r.db.Conn(ctx).
		Where("customer_id = ? AND business_id = ? AND street = ? AND city = ? AND state = ? AND country = ? AND postal_code = ?",
			address.CustomerID, address.BusinessID, address.Street, address.City, address.State, address.Country, address.PostalCode).
		First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		model := mapCustomerAddressFromEntity(address)
		return r.db.Conn(ctx).Create(model).Error
	}
	if err != nil {
		return err
	}

	updates := map[string]any{
		"times_used":   existing.TimesUsed + 1,
		"last_used_at": address.LastUsedAt,
	}
	if address.Latitude != nil {
		updates["latitude"] = address.Latitude
	}
	if address.Longitude != nil {
		updates["longitude"] = address.Longitude
	}
	return r.db.Conn(ctx).Model(&existing).Updates(updates).Error
}

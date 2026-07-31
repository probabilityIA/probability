package repository

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	"gorm.io/gorm"
)

func (r *Repository) UpdateClientProfile(ctx context.Context, businessID, userID uint, name, phone, dni string) error {
	updates := map[string]interface{}{}
	if name != "" {
		updates["name"] = name
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if dni != "" {
		updates["dni"] = dni
	}
	if len(updates) == 0 {
		return nil
	}

	err := r.db.Conn(ctx).
		Table("client").
		Where("business_id = ? AND user_id = ? AND deleted_at IS NULL", businessID, userID).
		Updates(updates).Error
	if err != nil {
		return err
	}

	userUpdates := map[string]interface{}{}
	if name != "" {
		userUpdates["name"] = name
	}
	if phone != "" {
		userUpdates["phone"] = phone
	}
	if len(userUpdates) == 0 {
		return nil
	}
	return r.db.Conn(ctx).
		Table("user").
		Where("id = ? AND deleted_at IS NULL", userID).
		Updates(userUpdates).Error
}

func (r *Repository) UpdateUserAvatar(ctx context.Context, userID uint, avatarURL string) error {
	return r.db.Conn(ctx).
		Table("user").
		Where("id = ? AND deleted_at IS NULL", userID).
		Update("avatar_url", avatarURL).Error
}

func (r *Repository) GetUserAvatar(ctx context.Context, userID uint) (string, error) {
	var user struct {
		AvatarURL string
	}
	err := r.db.Conn(ctx).
		Table("user").
		Select("avatar_url").
		Where("id = ? AND deleted_at IS NULL", userID).
		Limit(1).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return user.AvatarURL, nil
}

func (r *Repository) ListCustomerAddresses(ctx context.Context, businessID, customerID uint) ([]entities.CustomerAddress, error) {
	var rows []struct {
		ID         uint
		Street     string
		City       string
		State      string
		Country    string
		PostalCode string
		IsPrimary  bool
	}
	err := r.db.Conn(ctx).
		Table("customer_address").
		Select("id, street, city, state, country, postal_code, is_primary").
		Where("customer_id = ? AND business_id = ? AND deleted_at IS NULL", customerID, businessID).
		Order("is_primary DESC, last_used_at DESC").
		Find(&rows).Error
	if err != nil {
		return nil, err
	}
	addresses := make([]entities.CustomerAddress, 0, len(rows))
	for _, row := range rows {
		addresses = append(addresses, entities.CustomerAddress{
			ID:         row.ID,
			Street:     row.Street,
			City:       row.City,
			State:      row.State,
			Country:    row.Country,
			PostalCode: row.PostalCode,
			IsPrimary:  row.IsPrimary,
		})
	}
	return addresses, nil
}

func (r *Repository) CreateCustomerAddress(ctx context.Context, businessID, customerID uint, addr *entities.CustomerAddress) error {
	if addr.IsPrimary {
		if err := r.clearPrimaryAddress(ctx, businessID, customerID); err != nil {
			return err
		}
	}
	return r.db.Conn(ctx).
		Table("customer_address").
		Create(map[string]interface{}{
			"customer_id":  customerID,
			"business_id":  businessID,
			"street":       addr.Street,
			"city":         addr.City,
			"state":        addr.State,
			"country":      addr.Country,
			"postal_code":  addr.PostalCode,
			"is_primary":   addr.IsPrimary,
			"times_used":   1,
			"last_used_at": time.Now(),
			"created_at":   time.Now(),
			"updated_at":   time.Now(),
		}).Error
}

func (r *Repository) SetPrimaryAddress(ctx context.Context, businessID, customerID, addressID uint) error {
	if err := r.clearPrimaryAddress(ctx, businessID, customerID); err != nil {
		return err
	}
	return r.db.Conn(ctx).
		Table("customer_address").
		Where("id = ? AND customer_id = ? AND business_id = ? AND deleted_at IS NULL", addressID, customerID, businessID).
		Update("is_primary", true).Error
}

func (r *Repository) DeleteCustomerAddress(ctx context.Context, businessID, customerID, addressID uint) error {
	return r.db.Conn(ctx).
		Table("customer_address").
		Where("id = ? AND customer_id = ? AND business_id = ? AND deleted_at IS NULL", addressID, customerID, businessID).
		Update("deleted_at", time.Now()).Error
}

func (r *Repository) clearPrimaryAddress(ctx context.Context, businessID, customerID uint) error {
	return r.db.Conn(ctx).
		Table("customer_address").
		Where("customer_id = ? AND business_id = ? AND is_primary = true AND deleted_at IS NULL", customerID, businessID).
		Update("is_primary", false).Error
}

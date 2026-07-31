package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/publicsite/internal/domain/entities"
	"gorm.io/gorm"
)

func (r *Repository) GetClientSession(ctx context.Context, businessID uint, userID uint) (*entities.TiendaSession, error) {
	var client struct {
		ID    uint
		Name  string
		Email *string
		Phone string
		Dni   *string
	}
	err := r.db.Conn(ctx).
		Table("client").
		Select("id, name, email, phone, dni").
		Where("business_id = ? AND user_id = ? AND deleted_at IS NULL", businessID, userID).
		Limit(1).
		First(&client).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	session := &entities.TiendaSession{
		Name:       client.Name,
		Phone:      client.Phone,
		CustomerID: client.ID,
	}
	if client.Email != nil {
		session.Email = *client.Email
	}
	if client.Dni != nil {
		session.Dni = *client.Dni
	}

	var addr struct {
		Street     string
		City       string
		State      string
		Country    string
		PostalCode string
	}
	err = r.db.Conn(ctx).
		Table("customer_address").
		Select("street, city, state, country, postal_code").
		Where("customer_id = ? AND business_id = ? AND deleted_at IS NULL", client.ID, businessID).
		Order("is_primary DESC, last_used_at DESC").
		Limit(1).
		First(&addr).Error
	if err == nil {
		session.Address = &entities.CheckoutAddress{
			Street:     addr.Street,
			City:       addr.City,
			State:      addr.State,
			Country:    addr.Country,
			PostalCode: addr.PostalCode,
		}
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return session, nil
}

func (r *Repository) GetUserBasic(ctx context.Context, userID uint) (*entities.TiendaSession, error) {
	var user struct {
		Name  string
		Email string
		Phone string
	}
	err := r.db.Conn(ctx).
		Table("users").
		Select("name, email, phone").
		Where("id = ? AND deleted_at IS NULL", userID).
		Limit(1).
		First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &entities.TiendaSession{
		Name:  user.Name,
		Email: user.Email,
		Phone: user.Phone,
	}, nil
}

func (r *Repository) UserBelongsToBusiness(ctx context.Context, businessID uint, userID uint) (bool, error) {
	var result struct {
		ID uint
	}
	err := r.db.Conn(ctx).
		Table("business_staff").
		Select("id").
		Where("business_id = ? AND user_id = ? AND deleted_at IS NULL", businessID, userID).
		Limit(1).
		First(&result).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

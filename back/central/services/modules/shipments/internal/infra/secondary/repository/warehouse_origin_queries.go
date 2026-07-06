package repository

import (
	"context"
	"strings"

	"errors"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"gorm.io/gorm"
)

func (r *Repository) GetDefaultWarehouseOrigin(ctx context.Context, businessID uint) (*domain.OriginAddress, error) {
	var wh struct {
		ID           uint
		Name         string
		Address      string
		Street       string
		City         string
		State        string
		CityDaneCode string
		Suburb       string
		PostalCode   string
		ZipCode      string
		Company      string
		FirstName    string
		LastName     string
		Email        string
		ContactEmail string
		ContactName  string
		Phone        string
	}

	err := r.db.Conn(ctx).
		Table("warehouses").
		Select("id, name, address, street, city, state, city_dane_code, suburb, postal_code, zip_code, company, first_name, last_name, email, contact_email, contact_name, phone").
		Where("business_id = ? AND is_active = ? AND deleted_at IS NULL", businessID, true).
		Order("is_default DESC, id ASC").
		Limit(1).
		Scan(&wh).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	if wh.ID == 0 || wh.CityDaneCode == "" {
		return nil, nil
	}

	street := firstNonEmpty(wh.Street, wh.Address)
	company := firstNonEmpty(wh.Company, wh.Name)
	email := firstNonEmpty(wh.Email, wh.ContactEmail)
	postal := firstNonEmpty(wh.PostalCode, wh.ZipCode)

	firstName := wh.FirstName
	lastName := wh.LastName
	if firstName == "" && lastName == "" && wh.ContactName != "" {
		parts := strings.Fields(wh.ContactName)
		if len(parts) > 0 {
			firstName = parts[0]
		}
		if len(parts) > 1 {
			lastName = strings.Join(parts[1:], " ")
		}
	}

	return &domain.OriginAddress{
		ID:           wh.ID,
		BusinessID:   businessID,
		Alias:        wh.Name,
		Company:      company,
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		Phone:        wh.Phone,
		Street:       street,
		Suburb:       wh.Suburb,
		CityDaneCode: wh.CityDaneCode,
		City:         wh.City,
		State:        wh.State,
		PostalCode:   postal,
		IsDefault:    true,
	}, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

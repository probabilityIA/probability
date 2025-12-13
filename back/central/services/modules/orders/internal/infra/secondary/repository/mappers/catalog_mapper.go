package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/orders/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToDBProduct convierte un producto de dominio a modelo de base de datos
func ToDBProduct(p *domain.Product) *models.Product {
	if p == nil {
		return nil
	}
	return &models.Product{
		ID:         p.ID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		DeletedAt:  p.DeletedAt,
		BusinessID: p.BusinessID,
		SKU:        p.SKU,
		Name:       p.Name,
		ExternalID: p.ExternalID,
	}
}

// ToDomainProduct convierte un producto de base de datos a dominio
func ToDomainProduct(p *models.Product) *domain.Product {
	if p == nil {
		return nil
	}
	return &domain.Product{
		ID:         p.ID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		DeletedAt:  p.DeletedAt,
		BusinessID: p.BusinessID,
		SKU:        p.SKU,
		Name:       p.Name,
		ExternalID: p.ExternalID,
	}
}

// ToDBClient convierte un cliente de dominio a modelo de base de datos
func ToDBClient(c *domain.Client) *models.Client {
	if c == nil {
		return nil
	}
	return &models.Client{
		Model: gorm.Model{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		BusinessID: c.BusinessID,
		Name:       c.Name,
		Email:      c.Email,
		Phone:      c.Phone,
		Dni:        c.Dni,
	}
}

// ToDomainClient convierte un cliente de base de datos a dominio
func ToDomainClient(c *models.Client) *domain.Client {
	if c == nil {
		return nil
	}
	var deletedAt *time.Time
	if c.DeletedAt.Valid {
		deletedAt = &c.DeletedAt.Time
	}
	return &domain.Client{
		ID:         c.ID,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
		DeletedAt:  deletedAt,
		BusinessID: c.BusinessID,
		Name:       c.Name,
		Email:      c.Email,
		Phone:      c.Phone,
		Dni:        c.Dni,
	}
}

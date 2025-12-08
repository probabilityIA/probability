package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/products/internal/domain"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

// ToDBProduct convierte un producto de dominio a modelo de base de datos
func ToDBProduct(p *domain.Product) *models.Product {
	if p == nil {
		return nil
	}
	dbProduct := &models.Product{
		Model: gorm.Model{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
			DeletedAt: gorm.DeletedAt{},
		},
		BusinessID: p.BusinessID,
		SKU:        p.SKU,
		Name:       p.Name,
		ExternalID: p.ExternalID,
	}
	if p.DeletedAt != nil {
		dbProduct.DeletedAt = gorm.DeletedAt{Time: *p.DeletedAt, Valid: true}
	}
	return dbProduct
}

// ToDomainProduct convierte un producto de base de datos a dominio
func ToDomainProduct(p *models.Product) *domain.Product {
	if p == nil {
		return nil
	}
	var deletedAt *time.Time
	if p.DeletedAt.Valid {
		deletedAt = &p.DeletedAt.Time
	}
	return &domain.Product{
		ID:         p.ID,
		CreatedAt:  p.CreatedAt,
		UpdatedAt:  p.UpdatedAt,
		DeletedAt:  deletedAt,
		BusinessID: p.BusinessID,
		SKU:        p.SKU,
		Name:       p.Name,
		ExternalID: p.ExternalID,
	}
}


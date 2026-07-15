package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/db"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ProductReadRepository struct {
	db  db.IDatabase
	log log.ILogger
}

func NewProductRepository(database db.IDatabase, logger log.ILogger) ports.IProductReadRepository {
	return &ProductReadRepository{
		db:  database,
		log: logger.WithModule("siigo.product_repository"),
	}
}

func (r *ProductReadRepository) ListProductsByBusiness(ctx context.Context, businessID uint) ([]dtos.ProductForSync, error) {
	var rows []struct {
		ID   string
		SKU  string
		Name string
	}

	err := r.db.Conn(ctx).
		Table("products").
		Select("id, sku, name").
		Where("business_id = ? AND deleted_at IS NULL AND is_active = ?", businessID, true).
		Order("created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	products := make([]dtos.ProductForSync, 0, len(rows))
	for _, row := range rows {
		products = append(products, dtos.ProductForSync{
			ID:   row.ID,
			SKU:  row.SKU,
			Name: row.Name,
		})
	}
	return products, nil
}

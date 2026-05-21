package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/pricing/internal/domain/errors"
	"github.com/secamc93/probability/back/migration/shared/models"
	"gorm.io/gorm"
)

type catalogRow struct {
	ProductID      string   `gorm:"column:id"`
	ProductName    string   `gorm:"column:name"`
	ProductSKU     string   `gorm:"column:sku"`
	ImageURL       string   `gorm:"column:image_url"`
	FamilyImageURL string   `gorm:"column:family_image_url"`
	Currency       string   `gorm:"column:currency"`
	BasePrice      float64  `gorm:"column:base_price"`
	CustomPrice    *float64 `gorm:"column:custom_price"`
}

func (r *Repository) ListCatalogPrices(ctx context.Context, params dtos.ListCatalogPricesParams) ([]entities.CatalogPriceRow, int64, error) {
	joinCond := "cpp.product_id = p.id AND cpp.business_id = p.business_id AND cpp.deleted_at IS NULL"
	var joinArg any
	if params.Target.ClientGroupID != nil {
		joinCond += " AND cpp.client_group_id = ?"
		joinArg = *params.Target.ClientGroupID
	} else {
		joinCond += " AND cpp.client_id = ?"
		joinArg = *params.Target.ClientID
	}

	base := r.db.Conn(ctx).
		Table("products p").
		Where("p.business_id = ? AND p.deleted_at IS NULL AND p.status = ?", params.Target.BusinessID, "active")
	if params.Search != "" {
		like := "%" + params.Search + "%"
		base = base.Where("p.name ILIKE ? OR p.sku ILIKE ?", like, like)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []catalogRow
	if err := base.
		Joins("LEFT JOIN custom_product_price cpp ON "+joinCond, joinArg).
		Joins("LEFT JOIN product_families pf ON p.family_id = pf.id").
		Select("p.id, p.name, p.sku, p.image_url, COALESCE(pf.image_url, '') AS family_image_url, p.currency, p.price AS base_price, cpp.price AS custom_price").
		Order("p.name ASC").
		Offset(params.Offset()).Limit(params.PageSize).
		Scan(&rows).Error; err != nil {
		return nil, 0, err
	}

	result := make([]entities.CatalogPriceRow, len(rows))
	for i, row := range rows {
		result[i] = entities.CatalogPriceRow{
			ProductID:      row.ProductID,
			ProductName:    row.ProductName,
			ProductSKU:     row.ProductSKU,
			ImageURL:       row.ImageURL,
			FamilyImageURL: row.FamilyImageURL,
			Currency:       row.Currency,
			BasePrice:      row.BasePrice,
			CustomPrice:    row.CustomPrice,
		}
	}
	return result, total, nil
}

func (r *Repository) SaveCatalogPrices(ctx context.Context, dto dtos.SaveCatalogPricesDTO) error {
	if len(dto.Items) == 0 {
		return nil
	}

	return r.db.Conn(ctx).Transaction(func(tx *gorm.DB) error {
		productIDs := make([]string, len(dto.Items))
		for i, item := range dto.Items {
			productIDs[i] = item.ProductID
		}

		var validIDs []string
		if err := tx.Model(&models.Product{}).
			Where("business_id = ? AND id IN ?", dto.Target.BusinessID, productIDs).
			Pluck("id", &validIDs).Error; err != nil {
			return err
		}
		valid := make(map[string]bool, len(validIDs))
		for _, id := range validIDs {
			valid[id] = true
		}

		for _, item := range dto.Items {
			if !valid[item.ProductID] {
				continue
			}

			scope := tx.Unscoped().Model(&models.CustomProductPrice{}).
				Where("business_id = ? AND product_id = ?", dto.Target.BusinessID, item.ProductID)
			if dto.Target.ClientGroupID != nil {
				scope = scope.Where("client_group_id = ? AND client_id IS NULL", *dto.Target.ClientGroupID)
			} else {
				scope = scope.Where("client_id = ? AND client_group_id IS NULL", *dto.Target.ClientID)
			}

			if item.Price == nil {
				if err := scope.Delete(&models.CustomProductPrice{}).Error; err != nil {
					return err
				}
				continue
			}

			var existing models.CustomProductPrice
			err := scope.Session(&gorm.Session{}).First(&existing).Error
			if err == gorm.ErrRecordNotFound {
				if err := tx.Create(&models.CustomProductPrice{
					BusinessID:    dto.Target.BusinessID,
					ProductID:     item.ProductID,
					ClientGroupID: dto.Target.ClientGroupID,
					ClientID:      dto.Target.ClientID,
					Price:         *item.Price,
					IsActive:      true,
				}).Error; err != nil {
					return err
				}
				continue
			}
			if err != nil {
				return err
			}
			if err := tx.Unscoped().Model(&models.CustomProductPrice{}).
				Where("id = ?", existing.ID).
				Updates(map[string]any{
					"price":      *item.Price,
					"is_active":  true,
					"deleted_at": nil,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *Repository) GetEffectivePrice(ctx context.Context, params dtos.EffectivePriceParams) (*entities.EffectivePrice, error) {
	var product models.Product
	err := r.db.Conn(ctx).Select("id, price").
		Where("id = ? AND business_id = ? AND deleted_at IS NULL", params.ProductID, params.BusinessID).
		First(&product).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainerrors.ErrProductNotFound
		}
		return nil, err
	}

	groupID, err := r.GetClientGroupID(ctx, params.BusinessID, params.ClientID)
	if err != nil {
		return nil, err
	}

	var groupPrice *float64
	var groupName string
	if groupID != nil {
		var group models.ClientGroup
		if e := r.db.Conn(ctx).Select("name").Where("id = ?", *groupID).First(&group).Error; e == nil {
			groupName = group.Name
		}

		var row models.CustomProductPrice
		e := r.db.Conn(ctx).
			Where("business_id = ? AND product_id = ? AND client_group_id = ? AND is_active = true", params.BusinessID, params.ProductID, *groupID).
			First(&row).Error
		if e == nil {
			price := row.Price
			groupPrice = &price
		} else if e != gorm.ErrRecordNotFound {
			return nil, e
		}
	}

	var clientPrice *float64
	var clientRow models.CustomProductPrice
	e := r.db.Conn(ctx).
		Where("business_id = ? AND product_id = ? AND client_id = ? AND is_active = true", params.BusinessID, params.ProductID, params.ClientID).
		First(&clientRow).Error
	if e == nil {
		price := clientRow.Price
		clientPrice = &price
	} else if e != gorm.ErrRecordNotFound {
		return nil, e
	}

	result := domain.ResolveEffectivePrice(params.ProductID, product.Price, groupPrice, clientPrice, groupID)
	result.GroupName = groupName
	return &result, nil
}

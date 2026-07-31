package repository

import (
	"context"

	"gorm.io/gorm"
)

type priceRow struct {
	ProductID string
	Price     float64
}

func (r *Repository) GetCustomerPrices(ctx context.Context, businessID uint, clientID uint) (map[string]float64, error) {
	prices := make(map[string]float64)

	var groupID *uint
	var member struct {
		ClientGroupID uint
	}
	err := r.db.Conn(ctx).
		Table("client_group_member").
		Select("client_group_id").
		Where("business_id = ? AND client_id = ? AND deleted_at IS NULL", businessID, clientID).
		Limit(1).
		First(&member).Error
	if err == nil {
		groupID = &member.ClientGroupID
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if groupID != nil {
		var rows []priceRow
		err = r.db.Conn(ctx).
			Table("custom_product_price").
			Select("product_id, price").
			Where("business_id = ? AND client_group_id = ? AND is_active = true AND deleted_at IS NULL", businessID, *groupID).
			Find(&rows).Error
		if err != nil {
			return nil, err
		}
		for _, row := range rows {
			prices[row.ProductID] = row.Price
		}
	}

	var clientRows []priceRow
	err = r.db.Conn(ctx).
		Table("custom_product_price").
		Select("product_id, price").
		Where("business_id = ? AND client_id = ? AND is_active = true AND deleted_at IS NULL", businessID, clientID).
		Find(&clientRows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range clientRows {
		prices[row.ProductID] = row.Price
	}

	return prices, nil
}

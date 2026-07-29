package repository

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func (r *Repository) GetProductMatchRules(ctx context.Context, integrationID uint) []productmatch.Rule {
	if integrationID == 0 {
		return productmatch.DefaultRules()
	}

	var row struct {
		ProductMatchRules        []byte
		DefaultProductMatchRules []byte
		ProductMatchOptions      []byte
	}
	err := r.db.Conn(ctx).
		Table("integrations AS i").
		Select("i.product_match_rules, t.default_product_match_rules, t.product_match_options").
		Joins("LEFT JOIN integration_types t ON t.id = i.integration_type_id").
		Where("i.id = ? AND i.deleted_at IS NULL", integrationID).
		Limit(1).
		Scan(&row).Error
	if err != nil {
		return productmatch.DefaultRules()
	}

	rules := productmatch.Resolve(row.ProductMatchRules, row.DefaultProductMatchRules)
	return productmatch.AllowedByOptions(rules, productmatch.ParseOptions(row.ProductMatchOptions))
}

func orderItemChannelValue(item dtos.ProbabilityOrderItemDTO, field string) string {
	switch field {
	case productmatch.FieldSKU:
		return item.ProductSKU
	case productmatch.FieldBarcode:
		if item.ExternalBarcode != nil {
			return *item.ExternalBarcode
		}
	case productmatch.FieldExternalID:
		if item.ProductID != nil {
			return *item.ProductID
		}
	case productmatch.FieldVariantID:
		if item.VariantID != nil {
			return *item.VariantID
		}
	case productmatch.FieldName:
		if item.ProductName != "" {
			return item.ProductName
		}
		return item.ProductTitle
	}
	return ""
}

var mappingColumnByChannelField = map[string]string{
	productmatch.FieldSKU:        "external_sku",
	productmatch.FieldBarcode:    "external_barcode",
	productmatch.FieldExternalID: "external_product_id",
	productmatch.FieldVariantID:  "external_variant_id",
}

var productColumnByProbabilityField = map[string]string{
	productmatch.FieldSKU:        "sku",
	productmatch.FieldBarcode:    "barcode",
	productmatch.FieldExternalID: "external_id",
	productmatch.FieldName:       "name",
}

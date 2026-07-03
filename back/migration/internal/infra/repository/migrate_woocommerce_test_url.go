package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	woocommerceTypeCode    = "woocommerce"
	woocommerceBaseURLTest = "http://back-testing:9096"
)

func (r *Repository) migrateWooCommerceTestURL(ctx context.Context) error {
	result := r.db.Conn(ctx).
		Model(&models.IntegrationType{}).
		Where("code = ? AND (base_url_test IS NULL OR base_url_test = '')", woocommerceTypeCode).
		Update("base_url_test", woocommerceBaseURLTest)

	if result.Error != nil {
		return fmt.Errorf("migrateWooCommerceTestURL: %w", result.Error)
	}
	return nil
}

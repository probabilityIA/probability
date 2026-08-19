package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

const (
	tiendanubeTypeCode    = "tiendanube"
	tiendanubeBaseURL     = "https://api.tiendanube.com/v1"
	tiendanubeBaseURLTest = "http://back-testing:9102/v1"
)

func (r *Repository) migrateTiendanubeURLs(ctx context.Context) error {
	result := r.db.Conn(ctx).
		Model(&models.IntegrationType{}).
		Where("code = ? AND (base_url IS NULL OR base_url = '')", tiendanubeTypeCode).
		Update("base_url", tiendanubeBaseURL)
	if result.Error != nil {
		return fmt.Errorf("migrateTiendanubeURLs base_url: %w", result.Error)
	}

	result = r.db.Conn(ctx).
		Model(&models.IntegrationType{}).
		Where("code = ? AND (base_url_test IS NULL OR base_url_test = '')", tiendanubeTypeCode).
		Update("base_url_test", tiendanubeBaseURLTest)
	if result.Error != nil {
		return fmt.Errorf("migrateTiendanubeURLs base_url_test: %w", result.Error)
	}

	return nil
}

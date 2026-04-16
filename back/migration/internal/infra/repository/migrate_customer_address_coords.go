package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateCustomerAddressCoords(ctx context.Context) error {
	db := r.db.Conn(ctx)
	if err := db.AutoMigrate(&models.CustomerAddress{}); err != nil {
		return fmt.Errorf("failed to auto-migrate customer_addresses: %w", err)
	}

	db.Exec(`
		UPDATE customer_address ca
		SET latitude = o.shipping_lat, longitude = o.shipping_lng
		FROM (
			SELECT DISTINCT ON (customer_id, business_id, shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code)
				customer_id, business_id, shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code,
				shipping_lat, shipping_lng
			FROM orders
			WHERE deleted_at IS NULL
				AND customer_id IS NOT NULL
				AND shipping_lat IS NOT NULL
				AND shipping_lng IS NOT NULL
			ORDER BY customer_id, business_id, shipping_street, shipping_city, shipping_state, shipping_country, shipping_postal_code, created_at DESC
		) o
		WHERE ca.customer_id = o.customer_id
			AND ca.business_id = o.business_id
			AND ca.street = o.shipping_street
			AND ca.city = o.shipping_city
			AND ca.state = o.shipping_state
			AND ca.country = o.shipping_country
			AND ca.postal_code = o.shipping_postal_code
			AND ca.latitude IS NULL
	`)

	return nil
}

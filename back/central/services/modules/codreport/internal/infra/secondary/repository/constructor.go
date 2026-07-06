package repository

import "github.com/secamc93/probability/back/central/shared/db"

type Repository struct {
	db db.IDatabase
}

func New(database db.IDatabase) *Repository {
	return &Repository{db: database}
}

const pendingStatuses = "'pending','picked_up','in_transit','out_for_delivery','on_hold'"

const latestShipmentJoin = `
JOIN LATERAL (
	SELECT sh.id, sh.carrier, sh.status, sh.delivered_at, sh.updated_at, sh.shipping_cost, sh.cod_carrier_fee, sh.cod_probability_margin, sh.guide_id, sh.guide_url, sh.probability_guide_url
	FROM shipments sh
	WHERE sh.order_id = o.id AND sh.deleted_at IS NULL
	ORDER BY sh.created_at DESC
	LIMIT 1
) s ON true`

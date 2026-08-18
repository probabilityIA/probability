package repository

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

// migrateOrderChannelPack agrega orders.channel_pack_id: el id del carrito del
// canal (pack de MercadoLibre) cuando la orden agrupa varias ordenes del canal.
// Sin esto no hay forma de distinguir en la UI una orden normal de una que
// consolida un carrito.
func (r *Repository) migrateOrderChannelPack(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.Order{}); err != nil {
		return fmt.Errorf("automigrate orders (channel_pack_id): %w", err)
	}

	const backfill = `
		UPDATE orders o
		SET channel_pack_id = m.raw_data->>'pack_id'
		FROM order_channel_metadata m
		WHERE m.order_id = o.id
		  AND m.deleted_at IS NULL
		  AND o.platform = 'mercadolibre'
		  AND COALESCE(o.channel_pack_id, '') = ''
		  AND m.raw_data->>'pack_id' IS NOT NULL
		  AND m.raw_data->>'pack_id' = o.external_id`

	if err := r.db.Conn(ctx).Exec(backfill).Error; err != nil {
		return fmt.Errorf("backfill channel_pack_id: %w", err)
	}

	return nil
}

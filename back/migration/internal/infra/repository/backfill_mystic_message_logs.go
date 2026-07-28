package repository

import (
	"context"
	"fmt"
)

const mysticBusinessID = 36

func (r *Repository) backfillMysticMessageLogs(ctx context.Context) error {
	db := r.db.Conn(ctx)

	res := db.Exec(`
INSERT INTO whatsapp_message_logs (
    id, conversation_id, direction, message_id, template_name, content, status, created_at
)
SELECT
    gen_random_uuid(),
    wc.id,
    'outbound',
    wc.last_message_id,
    wc.last_template_id,
    CASE
        WHEN wc.last_template_id LIKE 'guia_envio_generada%' THEN
            '[recuperado] Guia de envio generada'
            || ' - Pedido: ' || COALESCE(wc.metadata->>'numero_pedido', wc.order_number)
            || ' - Guia: ' || COALESCE(wc.metadata->>'numero_guia', '-')
            || ' - Transportadora: ' || COALESCE(wc.metadata->>'transportadora', '-')
            || ' - Rastreo: ' || COALESCE(wc.metadata->>'link_rastreo', '-')
        ELSE
            '[recuperado] ' || wc.last_template_id
    END,
    'sent',
    COALESCE(wc.created_at, wc.updated_at)
FROM whatsapp_conversations wc
WHERE wc.business_id = ?
  AND COALESCE(wc.last_message_id, '') <> ''
  AND COALESCE(wc.last_template_id, '') <> ''
  AND NOT EXISTS (
      SELECT 1 FROM whatsapp_message_logs ml
      WHERE ml.message_id = wc.last_message_id
  )
`, mysticBusinessID)

	if res.Error != nil {
		return fmt.Errorf("backfill mystic message logs: %w", res.Error)
	}

	return nil
}

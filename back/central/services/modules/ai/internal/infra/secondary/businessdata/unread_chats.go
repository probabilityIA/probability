package businessdata

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const unreadChatsSQL = `
WITH conv AS (
	SELECT c.id, c.business_id, regexp_replace(c.phone_number, '[^0-9]', '', 'g') AS phone_key
	FROM whatsapp_conversations c
	WHERE c.created_at > NOW() - INTERVAL '60 days'
),
entrantes AS (
	SELECT conv.business_id, conv.phone_key, ml.created_at
	FROM conv
	JOIN whatsapp_message_logs ml ON ml.conversation_id = conv.id
	WHERE ml.direction = 'inbound' AND ml.created_at > NOW() - INTERVAL '60 days'
)
SELECT e.business_id, COUNT(DISTINCT e.phone_key) AS total, MIN(e.created_at) AS oldest_at
FROM entrantes e
LEFT JOIN whatsapp_conversation_reads rd ON rd.business_id = e.business_id AND rd.phone_key = e.phone_key
WHERE e.created_at > COALESCE(rd.last_read_at, 'epoch'::timestamptz)
GROUP BY e.business_id
ORDER BY e.business_id`

func (r *Reader) CountUnreadWhatsAppChats(ctx context.Context) ([]entities.UnreadChats, error) {
	var rows []struct {
		BusinessID uint
		Total      int64
		OldestAt   time.Time
	}
	if err := r.db.Conn(ctx).Raw(unreadChatsSQL).Scan(&rows).Error; err != nil {
		return nil, err
	}
	items := make([]entities.UnreadChats, 0, len(rows))
	for _, row := range rows {
		items = append(items, entities.UnreadChats{BusinessID: row.BusinessID, Count: row.Total, OldestAt: row.OldestAt})
	}
	return items, nil
}

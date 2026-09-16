package businessdata

import (
	"context"
	"strings"
)

const customerNameByPhoneSQL = `
SELECT cl.name
FROM client cl
WHERE cl.business_id = ?
  AND cl.deleted_at IS NULL
  AND btrim(COALESCE(cl.name, '')) <> ''
  AND right(regexp_replace(cl.phone, '[^0-9]', '', 'g'), 10) = right(?, 10)
ORDER BY cl.updated_at DESC
LIMIT 1`

func (r *Reader) FindCustomerNameByPhone(ctx context.Context, businessID uint, phone string) (string, error) {
	digits := digitsOnly(phone)
	if len(digits) < 7 {
		return "", nil
	}
	var name string
	if err := r.db.Conn(ctx).Raw(customerNameByPhoneSQL, businessID, digits).Scan(&name).Error; err != nil {
		return "", err
	}
	return strings.TrimSpace(name), nil
}

func digitsOnly(raw string) string {
	var b strings.Builder
	for _, r := range raw {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

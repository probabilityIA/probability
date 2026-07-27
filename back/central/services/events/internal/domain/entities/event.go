package entities

import (
	"time"
)

type Event struct {
	ID            string
	Type          string
	Category      string
	BusinessID    uint
	IntegrationID uint
	Timestamp     time.Time
	Data          map[string]interface{}
	Metadata      map[string]interface{}
}

const PaymentMethodIDCOD uint = 6

func (e Event) CodAmountIsFinal() bool {
	includes, ok := e.Data["cod_includes_shipping"].(bool)
	return ok && includes
}

func (e Event) IsCOD() bool {
	if flag, ok := e.Data["is_cod"].(bool); ok && flag {
		return true
	}

	val, ok := e.Data["payment_method_id"]
	if !ok || val == nil {
		return false
	}
	switch v := val.(type) {
	case uint:
		return v == PaymentMethodIDCOD
	case int:
		return uint(v) == PaymentMethodIDCOD
	case int64:
		return uint(v) == PaymentMethodIDCOD
	case float64:
		return uint(v) == PaymentMethodIDCOD
	default:
		return false
	}
}

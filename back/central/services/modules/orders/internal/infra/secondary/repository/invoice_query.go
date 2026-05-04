package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
	"github.com/secamc93/probability/back/central/services/modules/orders/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/db"
	"gorm.io/datatypes"
)

type InvoiceQuery struct {
	db db.IDatabase
}

func NewInvoiceQuery(database db.IDatabase) *InvoiceQuery {
	return &InvoiceQuery{db: database}
}

func (iq *InvoiceQuery) GetInvoiceByOrderID(ctx context.Context, orderID string) (*dtos.InvoiceData, error) {
	conn := iq.db.Conn(ctx)
	var invoice models.Invoice

	result := conn.
		Where("order_id = ? AND status = ?", orderID, "issued").
		Order("issued_at DESC").
		Limit(1).
		First(&invoice)

	if result.Error != nil {
		return nil, nil
	}

	retentionAmount := extractRetentionFromProviderResponse(invoice.ProviderResponse)

	return &dtos.InvoiceData{
		ID:              invoice.ID,
		InvoiceNumber:   invoice.InvoiceNumber,
		Status:          invoice.Status,
		IssuedAt:        invoice.IssuedAt,
		RetentionAmount: retentionAmount,
	}, nil
}

func extractRetentionFromProviderResponse(providerResponse datatypes.JSON) float64 {
	if len(providerResponse) == 0 {
		return 0
	}

	var data map[string]interface{}
	if err := json.Unmarshal(providerResponse, &data); err != nil {
		return 0
	}

	if retention, ok := data["totalWithholdingTax"]; ok {
		switch v := retention.(type) {
		case float64:
			return v
		case string:
			var result float64
			fmt.Sscanf(v, "%f", &result)
			return result
		}
	}

	return 0
}

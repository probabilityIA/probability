package dian

import (
	"context"
	"fmt"
	"math"

	factusinv "github.com/secamc93/probability/back/central/services/integrations/invoicing/factus"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/ports"
)

type FactusEmitter struct {
	emitter *factusinv.PlatformEmitter
}

func NewFactusEmitter(emitter *factusinv.PlatformEmitter) ports.IDianEmitter {
	return &FactusEmitter{emitter: emitter}
}

func (f *FactusEmitter) Emit(ctx context.Context, integrationID uint, inv *entities.Invoice, config map[string]interface{}) (*entities.DianEmitResult, error) {
	var ivaRate *float64
	for _, line := range inv.TaxDetail {
		if line.TaxCode == "IVA" {
			rate := line.RatePercent
			ivaRate = &rate
			break
		}
	}
	if ivaRate == nil {
		zero := 0.0
		ivaRate = &zero
	}

	items := make([]factusinv.PlatformInvoiceItem, len(inv.Items))
	for i, item := range inv.Items {
		qty := int(math.Round(item.Quantity))
		if qty < 1 {
			qty = 1
		}
		items[i] = factusinv.PlatformInvoiceItem{
			Code:      fmt.Sprintf("%s-%d", inv.Number, i+1),
			Name:      item.Description,
			Quantity:  qty,
			UnitPrice: item.UnitPrice,
			TaxRate:   ivaRate,
		}
	}

	result, err := f.emitter.Emit(ctx, factusinv.PlatformInvoiceRequest{
		IntegrationID:   integrationID,
		ReferenceCode:   inv.Number,
		CustomerName:    inv.BusinessName,
		CustomerEmail:   inv.EmailTo,
		CustomerPhone:   inv.CustomerPhone,
		CustomerDNI:     inv.CustomerDocument,
		CustomerAddress: inv.CustomerAddress,
		Items:           items,
		Subtotal:        inv.Subtotal,
		Tax:             inv.TaxTotal,
		Total:           inv.Total,
		Config:          config,
	})
	if err != nil {
		return nil, err
	}
	return &entities.DianEmitResult{
		Number:     result.Number,
		CUFE:       result.CUFE,
		QRCode:     result.QRCode,
		ExternalID: result.ExternalID,
		IssuedAt:   result.IssuedAt,
	}, nil
}

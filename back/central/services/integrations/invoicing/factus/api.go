package factus

import (
	"context"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	factuscore "github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/core"
)

type PlatformInvoiceItem struct {
	Code      string
	Name      string
	Quantity  int
	UnitPrice float64
	TaxRate   *float64
}

type PlatformInvoiceRequest struct {
	IntegrationID   uint
	ReferenceCode   string
	CustomerName    string
	CustomerEmail   string
	CustomerPhone   string
	CustomerDNI     string
	CustomerAddress string
	Items           []PlatformInvoiceItem
	Subtotal        float64
	Tax             float64
	Total           float64
	Config          map[string]interface{}
}

type PlatformInvoiceResult struct {
	Number     string
	CUFE       string
	QRCode     string
	ExternalID string
	IssuedAt   string
}

type PlatformEmitter struct {
	core *factuscore.FactusCore
}

func NewPlatformEmitter(core *factuscore.FactusCore) *PlatformEmitter {
	return &PlatformEmitter{core: core}
}

func (e *PlatformEmitter) Emit(ctx context.Context, req PlatformInvoiceRequest) (*PlatformInvoiceResult, error) {
	items := make([]dtos.ItemData, len(req.Items))
	for i, item := range req.Items {
		items[i] = dtos.ItemData{
			SKU:        item.Code,
			Name:       item.Name,
			Quantity:   item.Quantity,
			UnitPrice:  item.UnitPrice,
			TotalPrice: float64(item.Quantity) * item.UnitPrice,
			TaxRate:    item.TaxRate,
		}
	}
	result, err := e.core.ProcessInvoice(ctx, &dtos.ProcessInvoiceRequest{
		Operation:     "create",
		IntegrationID: req.IntegrationID,
		Customer: dtos.CustomerData{
			Name:    req.CustomerName,
			Email:   req.CustomerEmail,
			Phone:   req.CustomerPhone,
			DNI:     req.CustomerDNI,
			Address: req.CustomerAddress,
		},
		Items:    items,
		Subtotal: req.Subtotal,
		Tax:      req.Tax,
		Total:    req.Total,
		Currency: "COP",
		OrderID:  req.ReferenceCode,
		Config:   req.Config,
	})
	if err != nil {
		return nil, err
	}
	return &PlatformInvoiceResult{
		Number:     result.InvoiceNumber,
		CUFE:       result.CUFE,
		QRCode:     result.QRCode,
		ExternalID: result.ExternalID,
		IssuedAt:   result.IssuedAt,
	}, nil
}

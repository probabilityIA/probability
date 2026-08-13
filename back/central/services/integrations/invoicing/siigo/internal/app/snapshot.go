package app

import (
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func siigoSnapshot(p dtos.ProductItem) productmatch.ChannelSnapshot {
	snapshot := productmatch.ChannelSnapshot{
		Name:        p.Name,
		Description: p.Description,
		Barcode:     p.Barcode,
		Brand:       p.Brand,
	}
	if p.Price > 0 {
		price := p.Price
		snapshot.Price = &price
	}
	return snapshot
}

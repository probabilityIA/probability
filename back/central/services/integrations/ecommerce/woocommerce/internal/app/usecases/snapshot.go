package usecases

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/woocommerce/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func wooSnapshot(p domain.WooProduct) productmatch.ChannelSnapshot {
	snapshot := productmatch.ChannelSnapshot{
		Name:        p.Name,
		Description: p.Description,
		Barcode:     p.Barcode,
		Brand:       p.Brand,
		Category:    p.Category,
		ImageURL:    p.ImageURL,
		Weight:      p.Weight,
		Length:      p.Length,
		Width:       p.Width,
		Height:      p.Height,
	}
	if len(p.VariantAttributes) > 0 && p.ParentID != "" {
		snapshot.FamilyName = nombreDeFamiliaWoo(p)
		snapshot.Color = p.VariantAttributes["color"]
		snapshot.Size = firstNonEmptyWoo(p.VariantAttributes["talla"], p.VariantAttributes["size"], p.VariantAttributes["tamaño"])
	}
	if p.Price > 0 {
		price := p.Price
		snapshot.Price = &price
	}
	return snapshot
}

func nombreDeFamiliaWoo(p domain.WooProduct) string {
	if nombre := strings.TrimSpace(p.ParentName); nombre != "" {
		return nombre
	}
	return strings.TrimSpace(p.Name)
}

func firstNonEmptyWoo(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}

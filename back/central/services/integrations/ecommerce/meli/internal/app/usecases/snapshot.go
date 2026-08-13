package usecases

import (
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
	"github.com/secamc93/probability/back/central/shared/productmatch"
)

func meliSnapshot(p domain.MeliProduct) productmatch.ChannelSnapshot {
	snapshot := productmatch.ChannelSnapshot{
		Name:     p.Name,
		Barcode:  p.Barcode,
		Brand:    p.Brand,
		Category: p.Category,
		Color:    p.Color,
		Size:     p.Size,
		ImageURL: p.ImageURL,
	}
	if p.VariationID != "" {
		snapshot.FamilyName = nombreDeFamilia(p)
	}
	if p.Price > 0 {
		price := p.Price
		snapshot.Price = &price
	}
	return snapshot
}

func atributosDeVariante(p domain.MeliProduct) map[string]string {
	if p.VariationID == "" {
		return nil
	}
	atributos := make(map[string]string, 2)
	if strings.TrimSpace(p.Color) != "" {
		atributos["color"] = strings.TrimSpace(p.Color)
	}
	if strings.TrimSpace(p.Size) != "" {
		atributos["talla"] = strings.TrimSpace(p.Size)
	}
	if len(atributos) == 0 {
		return nil
	}
	return atributos
}

func nombreDeFamilia(p domain.MeliProduct) string {
	if nombre := strings.TrimSpace(p.ParentName); nombre != "" {
		return nombre
	}
	return strings.TrimSpace(p.Name)
}

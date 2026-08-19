package response

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/tiendanube/internal/domain"
)

var preferredLangs = []string{"es", "es_mx", "es_ar", "pt", "pt_br", "en"}

type Localized map[string]string

func (l Localized) First() string {
	if len(l) == 0 {
		return ""
	}
	for _, lang := range preferredLangs {
		if v, ok := l[lang]; ok && strings.TrimSpace(v) != "" {
			return v
		}
	}
	for _, v := range l {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func (l *Localized) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" || trimmed == "" {
		*l = Localized{}
		return nil
	}
	if strings.HasPrefix(trimmed, "\"") {
		var plain string
		if err := json.Unmarshal(data, &plain); err != nil {
			return err
		}
		*l = Localized{"es": plain}
		return nil
	}
	var m map[string]string
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*l = Localized(m)
	return nil
}

type Numeric float64

func (n *Numeric) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" || trimmed == "" || trimmed == `""` {
		*n = 0
		return nil
	}
	if strings.HasPrefix(trimmed, "\"") {
		var raw string
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		raw = strings.TrimSpace(raw)
		if raw == "" {
			*n = 0
			return nil
		}
		parsed, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			*n = 0
			return nil
		}
		*n = Numeric(parsed)
		return nil
	}
	var parsed float64
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	*n = Numeric(parsed)
	return nil
}

type Stock struct {
	Value    int
	Infinite bool
}

func (s *Stock) UnmarshalJSON(data []byte) error {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "null" || trimmed == "" || trimmed == `""` {
		s.Value = 0
		s.Infinite = true
		return nil
	}
	if strings.HasPrefix(trimmed, "\"") {
		var raw string
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		raw = strings.TrimSpace(raw)
		if raw == "" {
			s.Value = 0
			s.Infinite = true
			return nil
		}
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			s.Value = 0
			return nil
		}
		s.Value = parsed
		return nil
	}
	var parsed int
	if err := json.Unmarshal(data, &parsed); err != nil {
		return err
	}
	s.Value = parsed
	return nil
}

type Variant struct {
	ID              int64   `json:"id"`
	ProductID       int64   `json:"product_id"`
	SKU             string  `json:"sku"`
	Barcode         string  `json:"barcode"`
	Price           Numeric `json:"price"`
	Stock           Stock   `json:"stock"`
	StockManagement bool    `json:"stock_management"`
	Weight          Numeric `json:"weight"`
	Depth           Numeric `json:"depth"`
	Width           Numeric `json:"width"`
	Height          Numeric `json:"height"`
}

type Image struct {
	ID       int64  `json:"id"`
	Src      string `json:"src"`
	Position int    `json:"position"`
}

type Product struct {
	ID          int64     `json:"id"`
	Name        Localized `json:"name"`
	Description Localized `json:"description"`
	Published   bool      `json:"published"`
	Variants    []Variant `json:"variants"`
	Images      []Image   `json:"images"`
}

func (p Product) ToDomain() domain.TiendanubeProduct {
	variants := make([]domain.TiendanubeVariant, 0, len(p.Variants))
	for _, v := range p.Variants {
		variants = append(variants, domain.TiendanubeVariant{
			ID:              v.ID,
			ProductID:       p.ID,
			SKU:             strings.TrimSpace(v.SKU),
			Barcode:         strings.TrimSpace(v.Barcode),
			Price:           float64(v.Price),
			Stock:           v.Stock.Value,
			StockManagement: v.StockManagement,
			Weight:          float64(v.Weight),
			Depth:           float64(v.Depth),
			Width:           float64(v.Width),
			Height:          float64(v.Height),
		})
	}

	image := ""
	for _, img := range p.Images {
		if img.Src != "" {
			image = img.Src
			break
		}
	}

	return domain.TiendanubeProduct{
		ID:          p.ID,
		Name:        p.Name.First(),
		Description: p.Description.First(),
		ImageURL:    image,
		Published:   p.Published,
		Variants:    variants,
	}
}

type Store struct {
	ID           int64     `json:"id"`
	Name         Localized `json:"name"`
	URL          string    `json:"url"`
	Country      string    `json:"country"`
	Currency     string    `json:"currency"`
	MainLanguage string    `json:"main_language"`
}

package productmatch

import "strings"

type ChannelSnapshot struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Barcode     string   `json:"barcode,omitempty"`
	Brand       string   `json:"brand,omitempty"`
	Category    string   `json:"category,omitempty"`
	Color       string   `json:"color,omitempty"`
	Size        string   `json:"size,omitempty"`
	FamilyName  string   `json:"family_name,omitempty"`
	ImageURL    string   `json:"image_url,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	CostPrice   *float64 `json:"cost_price,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
	Length      *float64 `json:"length,omitempty"`
	Width       *float64 `json:"width,omitempty"`
	Height      *float64 `json:"height,omitempty"`
}

func (s ChannelSnapshot) IsEmpty() bool {
	textos := []string{s.Name, s.Description, s.Barcode, s.Brand, s.Category, s.Color, s.Size, s.FamilyName, s.ImageURL}
	for _, t := range textos {
		if strings.TrimSpace(t) != "" {
			return false
		}
	}
	numeros := []*float64{s.Price, s.CostPrice, s.Weight, s.Length, s.Width, s.Height}
	for _, n := range numeros {
		if n != nil {
			return false
		}
	}
	return true
}

type SnapshotEntry struct {
	ProbabilityID    string
	ChannelVariantID string
	Snapshot         ChannelSnapshot
}

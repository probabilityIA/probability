package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	SourceChannel = "channel"
	SourceManual  = "manual"
	SourceImport  = "import"
)

const (
	FieldName        = "name"
	FieldDescription = "description"
	FieldBarcode     = "barcode"
	FieldBrand       = "brand"
	FieldCategory    = "category"
	FieldPrice       = "price"
	FieldCostPrice   = "cost_price"
	FieldImageURL    = "image_url"
	FieldWeight      = "weight"
	FieldDimensions  = "dimensions"
	FieldVariant     = "variant_label"
)

const (
	MaxOldValueLen     = 1000
	MaxChangesPerField = 5
)

type FieldOrigin struct {
	Source        string
	IntegrationID *uint
	UserID        *uint
}

func ChannelOrigin(integrationID uint) FieldOrigin {
	if integrationID == 0 {
		return FieldOrigin{Source: SourceChannel}
	}
	id := integrationID
	return FieldOrigin{Source: SourceChannel, IntegrationID: &id}
}

func ManualOrigin(userID uint) FieldOrigin {
	if userID == 0 {
		return FieldOrigin{Source: SourceManual}
	}
	id := userID
	return FieldOrigin{Source: SourceManual, UserID: &id}
}

type FieldChange struct {
	ProductID  string
	BusinessID uint
	Field      string
	OldValue   string
	Truncated  bool
	Origin     FieldOrigin
	BatchID    string
	ChangedAt  time.Time
}

type FieldOriginView struct {
	Field           string     `json:"field"`
	Source          string     `json:"source"`
	IntegrationID   *uint      `json:"integration_id,omitempty"`
	IntegrationName string     `json:"integration_name,omitempty"`
	UserID          *uint      `json:"user_id,omitempty"`
	UserName        string     `json:"user_name,omitempty"`
	ChangedAt       time.Time  `json:"changed_at"`
	PreviousValue   string     `json:"previous_value,omitempty"`
	PreviousAt      *time.Time `json:"previous_at,omitempty"`
}

func clip(value string) (string, bool) {
	if len(value) <= MaxOldValueLen {
		return value, false
	}
	return value[:MaxOldValueLen], true
}

func textOf(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func numberOf(value *float64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatFloat(*value, 'f', -1, 64)
}

func dimensionsOf(p *Product) string {
	if p.Length == nil && p.Width == nil && p.Height == nil {
		return ""
	}
	return fmt.Sprintf("%s x %s x %s %s", numberOf(p.Length), numberOf(p.Width), numberOf(p.Height), p.DimensionUnit)
}

func snapshot(p *Product) map[string]string {
	return map[string]string{
		FieldName:        p.Name,
		FieldDescription: p.Description,
		FieldBarcode:     textOf(p.Barcode),
		FieldBrand:       p.Brand,
		FieldCategory:    p.Category,
		FieldPrice:       strconv.FormatFloat(p.Price, 'f', -1, 64),
		FieldCostPrice:   numberOf(p.CostPrice),
		FieldImageURL:    p.ImageURL,
		FieldWeight:      numberOf(p.Weight),
		FieldDimensions:  dimensionsOf(p),
		FieldVariant:     p.VariantLabel,
	}
}

func DiffProduct(before, after *Product, origin FieldOrigin, batchID string, at time.Time) []FieldChange {
	if before == nil || after == nil {
		return nil
	}

	previo := snapshot(before)
	nuevo := snapshot(after)

	changes := make([]FieldChange, 0, len(previo))
	for field, viejo := range previo {
		if strings.TrimSpace(viejo) == strings.TrimSpace(nuevo[field]) {
			continue
		}
		valor, recortado := clip(viejo)
		changes = append(changes, FieldChange{
			ProductID:  after.ID,
			BusinessID: after.BusinessID,
			Field:      field,
			OldValue:   valor,
			Truncated:  recortado,
			Origin:     origin,
			BatchID:    batchID,
			ChangedAt:  at,
		})
	}
	return changes
}

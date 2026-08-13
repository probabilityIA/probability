package domain

import (
	"context"
	"time"
)

const (
	DataFieldName        = "name"
	DataFieldDescription = "description"
	DataFieldBarcode     = "barcode"
	DataFieldBrand       = "brand"
	DataFieldCategory    = "category"
	DataFieldImageURL    = "image_url"
	DataFieldWeight      = "weight"
	DataFieldVariant     = "variant_label"
)

const (
	ModeFillEmpty = "fill_empty"
	ModeOverwrite = "overwrite"
)

type DataField struct {
	Key   string
	Label string
	Note  string
}

func DataFields() []DataField {
	return []DataField{
		{Key: DataFieldDescription, Label: "Descripcion"},
		{Key: DataFieldCategory, Label: "Categoria"},
		{Key: DataFieldBrand, Label: "Marca"},
		{Key: DataFieldImageURL, Label: "Imagen", Note: "se guarda la URL del canal, no se copia el archivo"},
		{Key: DataFieldBarcode, Label: "Codigo de barras", Note: "se ignora cuando el canal repite el SKU"},
		{Key: DataFieldVariant, Label: "Variante", Note: "color y talla del canal"},
		{Key: DataFieldWeight, Label: "Peso"},
		{Key: DataFieldName, Label: "Nombre"},
	}
}

func IsDataField(key string) bool {
	for _, f := range DataFields() {
		if f.Key == key {
			return true
		}
	}
	return false
}

type DataSummaryCell struct {
	IntegrationID   uint
	IntegrationName string
	ChannelCode     string
	CanFill         int
	CanOverwrite    int
}

type DataSummaryRow struct {
	Field string
	Label string
	Note  string
	Cells []DataSummaryCell
}

type DataSummary struct {
	Rows       []DataSummaryRow
	Snapshotal *time.Time
}

type DataConflict struct {
	Source          string
	IntegrationID   *uint
	IntegrationName string
	Count           int
	LastChangeAt    *time.Time
}

type DataSample struct {
	ProductID string
	SKU       string
	Name      string
	Current   string
	Incoming  string
}

type DataPreviewQuery struct {
	BusinessID    uint
	IntegrationID uint
	Field         string
	Mode          string
	SKUs          []string
}

type DataPreview struct {
	Field        string
	Mode         string
	WouldFill    int
	WouldReplace int
	Conflicts    []DataConflict
	Samples      []DataSample
}

type IDataDiffRepository interface {
	DataSummary(ctx context.Context, businessID uint) (*DataSummary, error)
	DataPreview(ctx context.Context, query DataPreviewQuery) (*DataPreview, error)
}

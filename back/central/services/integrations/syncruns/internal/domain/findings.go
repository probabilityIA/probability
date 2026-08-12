package domain

import "context"

const (
	SeverityError = "error"
	SeverityWarn  = "warn"
	SeverityInfo  = "info"
)

// Codigos de hallazgo. Los primeros salen del comparativo de cada canal; los
// ultimos solo aparecen al cruzar todos los canales del negocio.
const (
	FindingChannelNoSKU  = "channel_no_sku"
	FindingSKUChanged    = "sku_changed"
	FindingSKUTypo       = "sku_typo"
	FindingNotAssociated = "not_associated"
	FindingNotPublished  = "not_published"
	FindingSoldNotOwned  = "sold_not_owned"
	FindingImbalance     = "channel_imbalance"
)

// ChannelSummary es el estado de un canal segun su ultimo comparativo.
type ChannelSummary struct {
	IntegrationID   uint
	IntegrationName string
	ChannelCode     string
	Matched         int
	NotAssociated   int
	OnlyInChannel   int
	ChannelNoSKU    int
	SKUChanged      int
	SKUTypo         int
	ComparedAt      string
}

// CrossChannelCounts son los conteos que solo se pueden calcular mirando todos
// los canales a la vez.
type CrossChannelCounts struct {
	NotPublished int
	SoldNotOwned int
	Imbalance    int
}

// Finding es un hallazgo accionable para el negocio.
type Finding struct {
	Code     string
	Severity string
	Title    string
	Detail   string
	Count    int
	Channels []string
}

type FindingsReport struct {
	BusinessID uint
	Channels   []ChannelSummary
	Cross      CrossChannelCounts
	Findings   []Finding
	Total      int
}

type IFindingsRepository interface {
	ChannelSummaries(ctx context.Context, businessID uint) ([]ChannelSummary, error)
	CrossChannel(ctx context.Context, businessID uint) (CrossChannelCounts, error)
	FindingItems(ctx context.Context, query FindingItemsQuery) (*FindingItemsPage, error)
}

// FindingItem es una fila del detalle de un hallazgo. Channels lleva los canales
// involucrados: un mismo SKU puede tener el problema en mas de uno.
type FindingItem struct {
	SKU      string
	Name     string
	Detail   string
	Channels []string
}

type FindingItemsQuery struct {
	BusinessID uint
	Code       string
	Search     string
	Page       int
	PageSize   int
	All        bool
}

func (q FindingItemsQuery) Offset() int {
	if q.Page < 1 {
		return 0
	}
	return (q.Page - 1) * q.PageSize
}

func (q *FindingItemsQuery) Normalize() {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 || q.PageSize > 200 {
		q.PageSize = 50
	}
}

type FindingItemsPage struct {
	Items      []FindingItem
	Total      int64
	Page       int
	PageSize   int
	TotalPages int
}

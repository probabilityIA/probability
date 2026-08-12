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
}

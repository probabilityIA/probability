package entities

import "time"

const (
	TemplateStatusDraft    = "draft"
	TemplateStatusPending  = "pending"
	TemplateStatusApproved = "approved"
	TemplateStatusRejected = "rejected"
	TemplateStatusPaused   = "paused"
	TemplateStatusDisabled = "disabled"
	TemplateStatusFailed   = "failed"
)

const (
	TemplateCategoryMarketing = "MARKETING"
	TemplateCategoryUtility   = "UTILITY"
)

type TemplateVariable struct {
	Position int
	Source   string
	Label    string
	Fallback string
}

type TemplateButton struct {
	Type string
	Text string
	URL  string
}

const (
	TemplateOriginSystem   = "system"
	TemplateOriginInternal = "internal"
	TemplateOriginBusiness = "business"
)

const (
	TemplateScopeOrderEvent = "order_event"
	TemplateScopeScheduled  = "scheduled"
	TemplateScopeInternal   = "internal"
)

func IsAllowedScope(scope string) bool {
	return scope == TemplateScopeOrderEvent || scope == TemplateScopeScheduled
}

type WhatsappTemplate struct {
	ID            uint
	BusinessID    *uint
	Origin        string
	Scope         string
	IntegrationID *uint

	Name     string
	Language string
	Category string

	BodyText   string
	HeaderText string
	FooterText string

	Variables  []TemplateVariable
	Buttons    []TemplateButton
	Components []map[string]any

	WABAID         string
	MetaTemplateID string

	Status         string
	RejectedReason string
	SubmittedAt    *time.Time
	ReviewedAt     *time.Time
	LastSyncedAt   *time.Time

	CreatedByID *uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *WhatsappTemplate) IsUsable() bool {
	return t.Status == TemplateStatusApproved
}

func (t *WhatsappTemplate) IsEditable() bool {
	return t.Origin == TemplateOriginBusiness
}

func (t *WhatsappTemplate) BelongsTo(businessID uint) bool {
	if t.Origin != TemplateOriginBusiness {
		return true
	}
	return t.BusinessID != nil && *t.BusinessID == businessID
}

var allowedVariableSources = map[string]string{
	"customer.first_name":    "Nombre del cliente",
	"customer.full_name":     "Nombre completo del cliente",
	"customer.days_inactive": "Dias sin comprar",
	"customer.last_product":  "Ultimo producto comprado",
	"customer.total_orders":  "Cantidad de pedidos",
	"business.name":          "Nombre de la tienda",
}

func IsAllowedVariableSource(source string) bool {
	_, ok := allowedVariableSources[source]
	return ok
}

func VariableSourceLabel(source string) string {
	return allowedVariableSources[source]
}

func AllowedVariableSources() map[string]string {
	out := make(map[string]string, len(allowedVariableSources))
	for key, value := range allowedVariableSources {
		out[key] = value
	}
	return out
}

package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/subscriptions/internal/domain/entities"
)

type SubscriptionTypeResponse struct {
	ID                   uint      `json:"id"`
	Name                 string    `json:"name"`
	Code                 string    `json:"code"`
	Description          string    `json:"description"`
	Price                float64   `json:"price"`
	BillingPeriod        string    `json:"billing_period"`
	Active               bool      `json:"active"`
	ModuleCodes          []string  `json:"module_codes"`
	MaxEcommerceChannels int       `json:"max_ecommerce_channels"`
	BusinessID           *uint     `json:"business_id,omitempty"`
	IncludedShipments    *int      `json:"included_shipments,omitempty"`
	ShipmentOveragePrice *float64  `json:"shipment_overage_price,omitempty"`
	IncludedInvoices     *int      `json:"included_invoices,omitempty"`
	InvoiceOveragePrice  *float64  `json:"invoice_overage_price,omitempty"`
	IncludedOrders       *int      `json:"included_orders,omitempty"`
	OrderOveragePrice    *float64  `json:"order_overage_price,omitempty"`
	CreatedAt            time.Time `json:"created_at"`
	UpdatedAt            time.Time `json:"updated_at"`
}

func FromSubscriptionType(t *entities.SubscriptionType) SubscriptionTypeResponse {
	return SubscriptionTypeResponse{
		ID:                   t.ID,
		Name:                 t.Name,
		Code:                 t.Code,
		Description:          t.Description,
		Price:                t.Price,
		BillingPeriod:        t.BillingPeriod,
		Active:               t.Active,
		ModuleCodes:          t.ModuleCodes,
		MaxEcommerceChannels: t.MaxEcommerceChannels,
		BusinessID:           t.BusinessID,
		IncludedShipments:    t.IncludedShipments,
		ShipmentOveragePrice: t.ShipmentOveragePrice,
		IncludedInvoices:     t.IncludedInvoices,
		InvoiceOveragePrice:  t.InvoiceOveragePrice,
		IncludedOrders:       t.IncludedOrders,
		OrderOveragePrice:    t.OrderOveragePrice,
		CreatedAt:            t.CreatedAt,
		UpdatedAt:            t.UpdatedAt,
	}
}

func FromSubscriptionTypes(types []entities.SubscriptionType) []SubscriptionTypeResponse {
	result := make([]SubscriptionTypeResponse, len(types))
	for i, t := range types {
		result[i] = FromSubscriptionType(&t)
	}
	return result
}

type SubscriptionResponse struct {
	ID                         uint                      `json:"id"`
	BusinessID                 uint                      `json:"business_id"`
	SubscriptionTypeID         uint                      `json:"subscription_type_id"`
	SubscriptionTypeName       string                    `json:"subscription_type_name"`
	Months                     int                       `json:"months"`
	Amount                     float64                   `json:"amount"`
	StartDate                  time.Time                 `json:"start_date"`
	EndDate                    time.Time                 `json:"end_date"`
	Status                     string                    `json:"status"`
	PaymentMethod              string                    `json:"payment_method"`
	PaymentReference           *string                   `json:"payment_reference,omitempty"`
	Notes                      *string                   `json:"notes,omitempty"`
	CreatedAt                  time.Time                 `json:"created_at"`
	SubscriptionType           *SubscriptionTypeResponse `json:"subscription_type,omitempty"`
	BusinessSubscriptionStatus string                    `json:"business_subscription_status,omitempty"`
}

func FromSubscription(s *entities.BusinessSubscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                   s.ID,
		BusinessID:           s.BusinessID,
		SubscriptionTypeID:   s.SubscriptionTypeID,
		SubscriptionTypeName: s.SubscriptionTypeName,
		Months:               s.Months,
		Amount:               s.Amount,
		StartDate:            s.StartDate,
		EndDate:              s.EndDate,
		Status:               s.Status,
		PaymentMethod:        s.PaymentMethod,
		PaymentReference:     s.PaymentReference,
		Notes:                s.Notes,
		CreatedAt:            s.CreatedAt,
	}
}

func FromSubscriptionWithType(s *entities.BusinessSubscription, t *entities.SubscriptionType) SubscriptionResponse {
	resp := FromSubscription(s)
	if t != nil {
		typeResp := FromSubscriptionType(t)
		resp.SubscriptionType = &typeResp
	}
	return resp
}

func FromSubscriptions(subs []entities.BusinessSubscription) []SubscriptionResponse {
	result := make([]SubscriptionResponse, len(subs))
	for i, s := range subs {
		result[i] = FromSubscription(&s)
	}
	return result
}

type OverrideResponse struct {
	ID              uint      `json:"id"`
	BusinessID      uint      `json:"business_id"`
	ModuleCode      string    `json:"module_code"`
	GrantedByUserID uint      `json:"granted_by_user_id"`
	Notes           *string   `json:"notes,omitempty"`
	ExpiresAt       *string   `json:"expires_at,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

func FromOverride(o *entities.BusinessModuleOverride) OverrideResponse {
	resp := OverrideResponse{
		ID:              o.ID,
		BusinessID:      o.BusinessID,
		ModuleCode:      o.ModuleCode,
		GrantedByUserID: o.GrantedByUserID,
		Notes:           o.Notes,
		CreatedAt:       o.CreatedAt,
	}
	if o.ExpiresAt != nil {
		formatted := o.ExpiresAt.Format("2006-01-02")
		resp.ExpiresAt = &formatted
	}
	return resp
}

func FromOverrides(overrides []entities.BusinessModuleOverride) []OverrideResponse {
	result := make([]OverrideResponse, len(overrides))
	for i, o := range overrides {
		result[i] = FromOverride(&o)
	}
	return result
}

type AuditLogResponse struct {
	ID          uint      `json:"id"`
	BusinessID  uint      `json:"business_id"`
	ActorUserID *uint     `json:"actor_user_id,omitempty"`
	ActorLabel  string    `json:"actor_label"`
	Action      string    `json:"action"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func FromAuditLog(l *entities.SubscriptionAuditLog) AuditLogResponse {
	return AuditLogResponse{
		ID:          l.ID,
		BusinessID:  l.BusinessID,
		ActorUserID: l.ActorUserID,
		ActorLabel:  l.ActorLabel,
		Action:      l.Action,
		Description: l.Description,
		CreatedAt:   l.CreatedAt,
	}
}

func FromAuditLogs(logs []entities.SubscriptionAuditLog) []AuditLogResponse {
	result := make([]AuditLogResponse, len(logs))
	for i, l := range logs {
		result[i] = FromAuditLog(&l)
	}
	return result
}

type AdminBusinessRowResponse struct {
	ID                uint       `json:"id"`
	Name              string     `json:"name"`
	Code              string     `json:"code"`
	PlanName          string     `json:"plan_name,omitempty"`
	Status            string     `json:"status"`
	CycleStartDate    *time.Time `json:"cycle_start_date,omitempty"`
	CycleEndDate      *time.Time `json:"cycle_end_date,omitempty"`
	LastPaymentAmount *float64   `json:"last_payment_amount,omitempty"`
	LastPaymentDate   *time.Time `json:"last_payment_date,omitempty"`
	ForecastedPayment *float64   `json:"forecasted_payment,omitempty"`
	CutoffDay         *int       `json:"cutoff_day,omitempty"`
}

func FromAdminBusinessRow(r *entities.AdminBusinessRow) AdminBusinessRowResponse {
	return AdminBusinessRowResponse{
		ID:                r.ID,
		Name:              r.Name,
		Code:              r.Code,
		PlanName:          r.PlanName,
		Status:            r.Status,
		CycleStartDate:    r.CycleStartDate,
		CycleEndDate:      r.CycleEndDate,
		LastPaymentAmount: r.LastPaymentAmount,
		LastPaymentDate:   r.LastPaymentDate,
		ForecastedPayment: r.ForecastedPayment,
		CutoffDay:         r.CutoffDay,
	}
}

func FromAdminBusinessRows(rows []entities.AdminBusinessRow) []AdminBusinessRowResponse {
	result := make([]AdminBusinessRowResponse, len(rows))
	for i, r := range rows {
		result[i] = FromAdminBusinessRow(&r)
	}
	return result
}

type AdminKPIsResponse struct {
	ActiveCount             int64   `json:"active_count"`
	ExpiringSoonCount       int64   `json:"expiring_soon_count"`
	ExpiredOrSuspendedCount int64   `json:"expired_or_suspended_count"`
	MRR                     float64 `json:"mrr"`
}

type SubscriptionUsageResponse struct {
	PlanName             string    `json:"plan_name"`
	PlanPrice            float64   `json:"plan_price"`
	BillingPeriod        string    `json:"billing_period"`
	ModuleCodes          []string  `json:"module_codes"`
	MaxEcommerceChannels int       `json:"max_ecommerce_channels"`
	CycleStartDate       time.Time `json:"cycle_start_date"`
	CycleEndDate         time.Time `json:"cycle_end_date"`

	IncludedShipments    *int     `json:"included_shipments,omitempty"`
	ShipmentOveragePrice *float64 `json:"shipment_overage_price,omitempty"`
	ShipmentsUsed        int64    `json:"shipments_used"`

	IncludedInvoices    *int     `json:"included_invoices,omitempty"`
	InvoiceOveragePrice *float64 `json:"invoice_overage_price,omitempty"`
	InvoicesUsed        int64    `json:"invoices_used"`

	IncludedOrders    *int     `json:"included_orders,omitempty"`
	OrderOveragePrice *float64 `json:"order_overage_price,omitempty"`
	OrdersUsed        int64    `json:"orders_used"`

	ForecastedPayment *float64 `json:"forecasted_payment,omitempty"`
}

func FromSubscriptionUsage(u *entities.SubscriptionUsage) SubscriptionUsageResponse {
	return SubscriptionUsageResponse{
		PlanName:             u.PlanName,
		PlanPrice:            u.PlanPrice,
		BillingPeriod:        u.BillingPeriod,
		ModuleCodes:          u.ModuleCodes,
		MaxEcommerceChannels: u.MaxEcommerceChannels,
		CycleStartDate:       u.CycleStartDate,
		CycleEndDate:         u.CycleEndDate,
		IncludedShipments:    u.IncludedShipments,
		ShipmentOveragePrice: u.ShipmentOveragePrice,
		ShipmentsUsed:        u.ShipmentsUsed,
		IncludedInvoices:     u.IncludedInvoices,
		InvoiceOveragePrice:  u.InvoiceOveragePrice,
		InvoicesUsed:         u.InvoicesUsed,
		IncludedOrders:       u.IncludedOrders,
		OrderOveragePrice:    u.OrderOveragePrice,
		OrdersUsed:           u.OrdersUsed,
		ForecastedPayment:    u.ForecastedPayment,
	}
}

func FromAdminKPIs(k entities.AdminKPIs) AdminKPIsResponse {
	return AdminKPIsResponse{
		ActiveCount:             k.ActiveCount,
		ExpiringSoonCount:       k.ExpiringSoonCount,
		ExpiredOrSuspendedCount: k.ExpiredOrSuspendedCount,
		MRR:                     k.MRR,
	}
}

package entities

import "time"

type SubscriptionAuditLog struct {
	ID          uint
	BusinessID  uint
	ActorUserID *uint
	ActorLabel  string
	Action      string
	Description string
	CreatedAt   time.Time
}

const (
	AuditActionPaymentRegistered       = "payment_registered"
	AuditActionPaymentReverted         = "payment_reverted"
	AuditActionDatesEdited             = "dates_edited"
	AuditActionCourtesyExtended        = "courtesy_extended"
	AuditActionOverrideGranted         = "override_granted"
	AuditActionOverrideRevoked         = "override_revoked"
	AuditActionSubscriptionSuspended   = "subscription_suspended"
	AuditActionSubscriptionReactivated = "subscription_reactivated"
	AuditActionTrialDowngraded         = "trial_downgraded_to_free"
	AuditActionOverageSettled          = "overage_settled"
	AuditActionOverageDuePaid          = "overage_due_paid"
)

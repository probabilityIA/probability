package ports

import (
	"context"
	"time"
)

type ITemplateOverageBillingRepository interface {
	GetActiveTemplatePlanLimits(ctx context.Context, businessID uint) (cycleStart, cycleEnd time.Time, includedTemplates *int, overagePrice *float64, found bool, err error)
	CountTemplatesInCycle(ctx context.Context, businessID uint, cycleStart, cycleEnd time.Time) (int64, error)
	DebitWalletForTemplateOverage(ctx context.Context, businessID uint, amount float64, templateID uint) error
}

package ports

import (
	"context"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
)

type ITemplateRepository interface {
	CreateTemplate(ctx context.Context, template *entities.WhatsappTemplate) error
	UpdateTemplate(ctx context.Context, template *entities.WhatsappTemplate) error
	GetTemplateByID(ctx context.Context, id uint) (*entities.WhatsappTemplate, error)
	GetTemplateByName(ctx context.Context, businessID uint, name, language string) (*entities.WhatsappTemplate, error)
	ListTemplates(ctx context.Context, businessID uint, scope, status string, page, pageSize int) ([]entities.WhatsappTemplate, int64, error)
	DeleteTemplate(ctx context.Context, id uint) error
	UpdateTemplateStatusByMeta(ctx context.Context, wabaID, name, language, status, reason string) error
}

type IScheduledRuleRepository interface {
	CreateRule(ctx context.Context, rule *entities.ScheduledRule) error
	UpdateRule(ctx context.Context, rule *entities.ScheduledRule) error
	GetRuleByID(ctx context.Context, id uint) (*entities.ScheduledRule, error)
	ListRules(ctx context.Context, businessID uint, page, pageSize int) ([]entities.ScheduledRule, int64, error)
	DeleteRule(ctx context.Context, id uint) error
	ListDueRules(ctx context.Context, now time.Time, limit int) ([]entities.ScheduledRule, error)
	MarkRuleRan(ctx context.Context, ruleID uint, lastRunAt, nextRunAt time.Time) error
}

type IScheduledRunRepository interface {
	CreateRun(ctx context.Context, run *entities.ScheduledRun) error
	FinishRun(ctx context.Context, run *entities.ScheduledRun) error
	ListRuns(ctx context.Context, ruleID uint, page, pageSize int) ([]entities.ScheduledRun, int64, error)
}

type IScheduledSendRepository interface {
	ReserveSend(ctx context.Context, send *entities.ScheduledSend) (bool, error)
	MarkSendResult(ctx context.Context, sendID uint, status, messageID, errorMessage string) error
	CountSendsToday(ctx context.Context, ruleID uint, sendDate string) (int64, error)
}

type ISegmentQuerier interface {
	FindInactiveCustomers(ctx context.Context, businessID uint, params entities.SegmentParams, requiresOptIn bool, cooldownDays uint, ruleID uint, limit int) ([]entities.SegmentCandidate, error)
}

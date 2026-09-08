package scheduled

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type ISendPublisher interface {
	PublishScheduledSend(ctx context.Context, message dtos.ScheduledSendMessage) error
}

type IUseCase interface {
	CreateRule(ctx context.Context, dto dtos.CreateScheduledRuleDTO) (*entities.ScheduledRule, error)
	UpdateRule(ctx context.Context, dto dtos.UpdateScheduledRuleDTO) (*entities.ScheduledRule, error)
	GetRule(ctx context.Context, id, businessID uint) (*entities.ScheduledRule, error)
	ListRules(ctx context.Context, businessID uint, page, pageSize int) ([]entities.ScheduledRule, int64, error)
	DeleteRule(ctx context.Context, id, businessID uint) error
	ListRuns(ctx context.Context, ruleID, businessID uint, page, pageSize int) ([]entities.ScheduledRun, int64, error)
	RunDueRules(ctx context.Context) error
	RunRuleNow(ctx context.Context, id, businessID uint) (*entities.ScheduledRun, error)
	MarkSendResult(ctx context.Context, sendID uint, status, messageID, errorMessage string) error
}

type useCase struct {
	rules     ports.IScheduledRuleRepository
	runs      ports.IScheduledRunRepository
	sends     ports.IScheduledSendRepository
	segments  ports.ISegmentQuerier
	templates ports.ITemplateRepository
	publisher ISendPublisher
	logger    log.ILogger
}

func New(
	rules ports.IScheduledRuleRepository,
	runs ports.IScheduledRunRepository,
	sends ports.IScheduledSendRepository,
	segments ports.ISegmentQuerier,
	templates ports.ITemplateRepository,
	publisher ISendPublisher,
	logger log.ILogger,
) IUseCase {
	return &useCase{
		rules:     rules,
		runs:      runs,
		sends:     sends,
		segments:  segments,
		templates: templates,
		publisher: publisher,
		logger:    logger.WithModule("scheduled_notifications_usecase"),
	}
}

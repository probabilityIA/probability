package runner

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/secamc93/probability/back/central/shared/rabbitmq"
)

type progressPublisher struct {
	queue rabbitmq.IQueue
	log   log.ILogger
}

func NewProgressPublisher(queue rabbitmq.IQueue, logger log.ILogger) ports.IProgressPublisher {
	return &progressPublisher{queue: queue, log: logger.WithModule("notification_backfill.progress")}
}

func (p *progressPublisher) PublishProgress(ctx context.Context, job *entities.JobState) {
	if p.queue == nil {
		return
	}

	bizID := uint(0)
	if job.BusinessID != nil {
		bizID = *job.BusinessID
	}

	envelope := rabbitmq.EventEnvelope{
		Type:       "backfill.progress",
		Category:   "notification_backfill",
		BusinessID: bizID,
		Data: map[string]interface{}{
			"job_id":         job.ID,
			"event_code":     job.EventCode,
			"status":         string(job.Status),
			"total_eligible": job.TotalEligible,
			"sent":           job.Sent,
			"skipped":        job.Skipped,
			"failed":         job.Failed,
			"dry_run":        job.DryRun,
			"error_message":  job.ErrorMessage,
		},
	}

	if err := rabbitmq.PublishEvent(ctx, p.queue, envelope); err != nil {
		p.log.Warn(ctx).Err(err).Str("job_id", job.ID).Msg("Failed to publish backfill progress event")
	}
}

package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
)

const (
	defaultThrottlePerSec = 10
	minThrottleInterval   = 20 * time.Millisecond
)

func (uc *useCase) Run(ctx context.Context, req dtos.RunRequest, createdBy uint) (*entities.JobState, error) {
	selector, ok := uc.registry.Get(req.EventCode)
	if !ok {
		return nil, fmt.Errorf("event_code no soportado: %s", req.EventCode)
	}

	filter := dtos.BackfillFilter{
		EventCode:  req.EventCode,
		BusinessID: req.BusinessID,
		Days:       req.Days,
		Limit:      req.Limit,
	}

	candidates, err := selector.Preview(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("preview failed: %w", err)
	}

	job := &entities.JobState{
		ID:            uuid.New().String(),
		EventCode:     req.EventCode,
		BusinessID:    req.BusinessID,
		Status:        entities.JobStatusRunning,
		TotalEligible: len(candidates),
		StartedAt:     time.Now(),
		CreatedBy:     createdBy,
	}
	uc.store.Create(job)
	uc.progressPublisher.PublishProgress(context.Background(), job)

	go uc.runJob(job.ID, selector, candidates)

	return job, nil
}

func (uc *useCase) GetJob(jobID string) (*entities.JobState, bool) {
	return uc.store.Get(jobID)
}

func (uc *useCase) runJob(jobID string, selector interface{ Dispatch(context.Context, entities.Candidate) error }, candidates []entities.Candidate) {
	bgCtx := context.Background()
	interval := time.Second / time.Duration(defaultThrottlePerSec)
	if interval < minThrottleInterval {
		interval = minThrottleInterval
	}

	for i, c := range candidates {
		if err := selector.Dispatch(bgCtx, c); err != nil {
			uc.log.Warn(bgCtx).
				Err(err).
				Str("job_id", jobID).
				Str("order_id", c.OrderID).
				Msg("Backfill dispatch failed")
			uc.store.Update(jobID, func(j *entities.JobState) {
				j.Failed++
			})
		} else {
			uc.store.Update(jobID, func(j *entities.JobState) {
				j.Sent++
			})
		}

		if snapshot, ok := uc.store.Get(jobID); ok {
			uc.progressPublisher.PublishProgress(bgCtx, snapshot)
		}

		if i < len(candidates)-1 {
			time.Sleep(interval)
		}
	}

	now := time.Now()
	uc.store.Update(jobID, func(j *entities.JobState) {
		j.Status = entities.JobStatusCompleted
		j.FinishedAt = &now
	})
	if snapshot, ok := uc.store.Get(jobID); ok {
		uc.progressPublisher.PublishProgress(bgCtx, snapshot)
	}
}

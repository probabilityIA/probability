package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
)

type IEligibilitySelector interface {
	EventCode() string
	EventName() string
	Channel() string
	Preview(ctx context.Context, filter dtos.BackfillFilter) ([]entities.Candidate, error)
	Dispatch(ctx context.Context, candidate entities.Candidate) error
}

type ISelectorRegistry interface {
	Register(selector IEligibilitySelector)
	Get(eventCode string) (IEligibilitySelector, bool)
	List() []IEligibilitySelector
}

type IJobStore interface {
	Create(job *entities.JobState)
	Get(jobID string) (*entities.JobState, bool)
	Update(jobID string, mutator func(*entities.JobState))
	List(businessID *uint) []*entities.JobState
}

type IProgressPublisher interface {
	PublishProgress(ctx context.Context, job *entities.JobState)
}

type IBusinessNameResolver interface {
	ResolveNames(ctx context.Context, ids []uint) (map[uint]string, error)
}

type IUseCase interface {
	ListEvents(ctx context.Context) []dtos.RegisteredEventResponse
	Preview(ctx context.Context, filter dtos.BackfillFilter) (*dtos.PreviewResponse, error)
	Run(ctx context.Context, req dtos.RunRequest, createdBy uint) (*entities.JobState, error)
	GetJob(jobID string) (*entities.JobState, bool)
}

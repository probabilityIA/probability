package mocks

import (
	"context"
	"sync"

	"github.com/rs/zerolog"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
	"github.com/secamc93/probability/back/central/shared/log"
)

type SelectorMock struct {
	CodeFn     func() string
	NameFn     func() string
	ChannelFn  func() string
	PreviewFn  func(ctx context.Context, filter dtos.BackfillFilter) ([]entities.Candidate, error)
	DispatchFn func(ctx context.Context, candidate entities.Candidate) error

	mu             sync.Mutex
	dispatched     []entities.Candidate
	previewFilters []dtos.BackfillFilter
}

var _ ports.IEligibilitySelector = (*SelectorMock)(nil)

func (m *SelectorMock) EventCode() string {
	if m.CodeFn != nil {
		return m.CodeFn()
	}
	return "guide_created"
}

func (m *SelectorMock) EventName() string {
	if m.NameFn != nil {
		return m.NameFn()
	}
	return "Guia creada"
}

func (m *SelectorMock) Channel() string {
	if m.ChannelFn != nil {
		return m.ChannelFn()
	}
	return "whatsapp"
}

func (m *SelectorMock) Preview(ctx context.Context, filter dtos.BackfillFilter) ([]entities.Candidate, error) {
	m.mu.Lock()
	m.previewFilters = append(m.previewFilters, filter)
	m.mu.Unlock()
	if m.PreviewFn != nil {
		return m.PreviewFn(ctx, filter)
	}
	return nil, nil
}

func (m *SelectorMock) Dispatch(ctx context.Context, candidate entities.Candidate) error {
	m.mu.Lock()
	m.dispatched = append(m.dispatched, candidate)
	m.mu.Unlock()
	if m.DispatchFn != nil {
		return m.DispatchFn(ctx, candidate)
	}
	return nil
}

func (m *SelectorMock) Dispatched() []entities.Candidate {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]entities.Candidate, len(m.dispatched))
	copy(out, m.dispatched)
	return out
}

func (m *SelectorMock) PreviewFilters() []dtos.BackfillFilter {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]dtos.BackfillFilter, len(m.previewFilters))
	copy(out, m.previewFilters)
	return out
}

type RegistryMock struct {
	GetFn  func(eventCode string) (ports.IEligibilitySelector, bool)
	ListFn func() []ports.IEligibilitySelector

	Registered []ports.IEligibilitySelector
}

var _ ports.ISelectorRegistry = (*RegistryMock)(nil)

func (m *RegistryMock) Register(selector ports.IEligibilitySelector) {
	m.Registered = append(m.Registered, selector)
}

func (m *RegistryMock) Get(eventCode string) (ports.IEligibilitySelector, bool) {
	if m.GetFn != nil {
		return m.GetFn(eventCode)
	}
	return &SelectorMock{}, true
}

func (m *RegistryMock) List() []ports.IEligibilitySelector {
	if m.ListFn != nil {
		return m.ListFn()
	}
	return nil
}

type JobStoreMock struct {
	mu   sync.Mutex
	jobs map[string]*entities.JobState
}

var _ ports.IJobStore = (*JobStoreMock)(nil)

func NewJobStore() *JobStoreMock {
	return &JobStoreMock{jobs: map[string]*entities.JobState{}}
}

func (m *JobStoreMock) Create(job *entities.JobState) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.jobs == nil {
		m.jobs = map[string]*entities.JobState{}
	}
	copy := *job
	m.jobs[job.ID] = &copy
}

func (m *JobStoreMock) Get(jobID string) (*entities.JobState, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	j, ok := m.jobs[jobID]
	if !ok {
		return nil, false
	}
	copy := *j
	return &copy, true
}

func (m *JobStoreMock) Update(jobID string, mutator func(*entities.JobState)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if j, ok := m.jobs[jobID]; ok {
		mutator(j)
	}
}

func (m *JobStoreMock) List(businessID *uint) []*entities.JobState {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]*entities.JobState, 0, len(m.jobs))
	for _, j := range m.jobs {
		if businessID != nil && (j.BusinessID == nil || *j.BusinessID != *businessID) {
			continue
		}
		copy := *j
		out = append(out, &copy)
	}
	return out
}

type ProgressPublisherMock struct {
	mu        sync.Mutex
	snapshots []entities.JobState
	done      chan struct{}
	once      sync.Once
}

var _ ports.IProgressPublisher = (*ProgressPublisherMock)(nil)

func NewProgressPublisher() *ProgressPublisherMock {
	return &ProgressPublisherMock{done: make(chan struct{})}
}

func (m *ProgressPublisherMock) PublishProgress(ctx context.Context, job *entities.JobState) {
	m.mu.Lock()
	m.snapshots = append(m.snapshots, *job)
	m.mu.Unlock()
	if job.Status == entities.JobStatusCompleted || job.Status == entities.JobStatusFailed {
		m.once.Do(func() { close(m.done) })
	}
}

func (m *ProgressPublisherMock) Snapshots() []entities.JobState {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]entities.JobState, len(m.snapshots))
	copy(out, m.snapshots)
	return out
}

func (m *ProgressPublisherMock) Done() <-chan struct{} { return m.done }

type BusinessNameResolverMock struct {
	ResolveNamesFn func(ctx context.Context, ids []uint) (map[uint]string, error)
	ReceivedIDs    [][]uint
}

var _ ports.IBusinessNameResolver = (*BusinessNameResolverMock)(nil)

func (m *BusinessNameResolverMock) ResolveNames(ctx context.Context, ids []uint) (map[uint]string, error) {
	m.ReceivedIDs = append(m.ReceivedIDs, ids)
	if m.ResolveNamesFn != nil {
		return m.ResolveNamesFn(ctx, ids)
	}
	return map[uint]string{}, nil
}

type SilentLogger struct{}

func NewSilentLogger() log.ILogger { return &SilentLogger{} }

func (l *SilentLogger) nop() zerolog.Logger { return zerolog.Nop() }

func (l *SilentLogger) Info(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Info() }
func (l *SilentLogger) Error(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Error() }
func (l *SilentLogger) Warn(ctx ...context.Context) *zerolog.Event  { n := l.nop(); return n.Warn() }
func (l *SilentLogger) Debug(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Debug() }
func (l *SilentLogger) Fatal(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Fatal() }
func (l *SilentLogger) Panic(ctx ...context.Context) *zerolog.Event { n := l.nop(); return n.Panic() }
func (l *SilentLogger) With() zerolog.Context                       { n := l.nop(); return n.With() }
func (l *SilentLogger) WithService(service string) log.ILogger      { return l }
func (l *SilentLogger) WithModule(module string) log.ILogger        { return l }
func (l *SilentLogger) WithBusinessID(businessID uint) log.ILogger  { return l }

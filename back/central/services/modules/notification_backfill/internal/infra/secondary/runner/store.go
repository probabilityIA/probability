package runner

import (
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
)

type store struct {
	mu   sync.RWMutex
	jobs map[string]*entities.JobState
}

func NewStore() ports.IJobStore {
	return &store{jobs: make(map[string]*entities.JobState)}
}

func (s *store) Create(job *entities.JobState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[job.ID] = job
}

func (s *store) Get(jobID string) (*entities.JobState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[jobID]
	if !ok {
		return nil, false
	}
	copy := *j
	return &copy, true
}

func (s *store) Update(jobID string, mutator func(*entities.JobState)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[jobID]
	if !ok {
		return
	}
	mutator(j)
}

func (s *store) List(businessID *uint) []*entities.JobState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*entities.JobState, 0, len(s.jobs))
	for _, j := range s.jobs {
		if businessID != nil && (j.BusinessID == nil || *j.BusinessID != *businessID) {
			continue
		}
		copy := *j
		out = append(out, &copy)
	}
	return out
}

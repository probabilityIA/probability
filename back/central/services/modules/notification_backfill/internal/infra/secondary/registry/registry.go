package registry

import (
	"sort"
	"sync"

	"github.com/secamc93/probability/back/central/services/modules/notification_backfill/internal/domain/ports"
)

type registry struct {
	mu        sync.RWMutex
	selectors map[string]ports.IEligibilitySelector
}

func New() ports.ISelectorRegistry {
	return &registry{selectors: make(map[string]ports.IEligibilitySelector)}
}

func (r *registry) Register(selector ports.IEligibilitySelector) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.selectors[selector.EventCode()] = selector
}

func (r *registry) Get(eventCode string) (ports.IEligibilitySelector, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	s, ok := r.selectors[eventCode]
	return s, ok
}

func (r *registry) List() []ports.IEligibilitySelector {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]ports.IEligibilitySelector, 0, len(r.selectors))
	for _, s := range r.selectors {
		out = append(out, s)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].EventCode() < out[j].EventCode()
	})
	return out
}

package consumer

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
)

const (
	documentIndexCurrentDayTTL = 2 * time.Minute
	documentIndexPastDayTTL    = 6 * time.Hour
	documentIndexMaxDays       = 24
)

type indexedDocument struct {
	DocumentNumber string
	DocumentDate   string
	Comment        string
}

type dayIndex struct {
	docs      []indexedDocument
	fetchedAt time.Time
}

type documentIndex struct {
	mu    sync.Mutex
	days  map[string]*dayIndex
	locks map[string]*sync.Mutex
}

func newDocumentIndex() *documentIndex {
	return &documentIndex{
		days:  make(map[string]*dayIndex),
		locks: make(map[string]*sync.Mutex),
	}
}

func dayIndexKey(integrationID uint, day string) string {
	return strconv.FormatUint(uint64(integrationID), 10) + "|" + day
}

func (d *documentIndex) ttlFor(day string, now time.Time) time.Duration {
	if day == now.Format("2006-01-02") {
		return documentIndexCurrentDayTTL
	}
	return documentIndexPastDayTTL
}

func (d *documentIndex) cached(key, day string, now time.Time) (*dayIndex, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	entry, ok := d.days[key]
	if !ok {
		return nil, false
	}
	if now.Sub(entry.fetchedAt) > d.ttlFor(day, now) {
		return nil, false
	}
	return entry, true
}

func (d *documentIndex) store(key string, entry *dayIndex) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.days[key] = entry
	if len(d.days) <= documentIndexMaxDays {
		return
	}
	oldestKey := ""
	var oldestAt time.Time
	for k, v := range d.days {
		if oldestKey == "" || v.fetchedAt.Before(oldestAt) {
			oldestKey = k
			oldestAt = v.fetchedAt
		}
	}
	delete(d.days, oldestKey)
}

func (d *documentIndex) lockFor(key string) *sync.Mutex {
	d.mu.Lock()
	defer d.mu.Unlock()
	lock, ok := d.locks[key]
	if !ok {
		lock = &sync.Mutex{}
		d.locks[key] = lock
	}
	return lock
}

type dayFetcher func(ctx context.Context, day string) ([]ports.ListedDocument, error)

func (d *documentIndex) documentsForDay(
	ctx context.Context,
	integrationID uint,
	day string,
	now time.Time,
	fetch dayFetcher,
) ([]indexedDocument, bool, error) {
	key := dayIndexKey(integrationID, day)

	if entry, ok := d.cached(key, day, now); ok {
		return entry.docs, true, nil
	}

	lock := d.lockFor(key)
	lock.Lock()
	defer lock.Unlock()

	if entry, ok := d.cached(key, day, now); ok {
		return entry.docs, true, nil
	}

	raw, err := fetch(ctx, day)
	if err != nil {
		return nil, false, err
	}

	docs := make([]indexedDocument, 0, len(raw))
	for i := range raw {
		docs = append(docs, indexedDocument{
			DocumentNumber: raw[i].DocumentNumber,
			DocumentDate:   raw[i].DocumentDate,
			Comment:        raw[i].Comment,
		})
	}

	d.store(key, &dayIndex{docs: docs, fetchedAt: time.Now()})

	return docs, false, nil
}

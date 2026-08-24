package consumer

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/softpymes/internal/domain/ports"
)

func TestDocumentIndexCachesDayAcrossLookups(t *testing.T) {
	idx := newDocumentIndex()
	calls := 0
	fetch := func(ctx context.Context, day string) ([]ports.ListedDocument, error) {
		calls++
		return []ports.ListedDocument{
			{DocumentNumber: "0000001", Comment: "order:abc"},
		}, nil
	}

	now := time.Now()
	day := now.Format("2006-01-02")

	docs, fromCache, err := idx.documentsForDay(context.Background(), 32, day, now, fetch)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fromCache {
		t.Fatal("first lookup must not come from cache")
	}
	if len(docs) != 1 {
		t.Fatalf("expected 1 doc, got %d", len(docs))
	}

	for i := range 50 {
		docs, fromCache, err = idx.documentsForDay(context.Background(), 32, day, now, fetch)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !fromCache {
			t.Fatalf("lookup %d must come from cache", i)
		}
	}

	if calls != 1 {
		t.Fatalf("expected exactly 1 provider fetch for 51 lookups, got %d", calls)
	}
}

func TestDocumentIndexSeparatesIntegrationsAndDays(t *testing.T) {
	idx := newDocumentIndex()
	calls := 0
	fetch := func(ctx context.Context, day string) ([]ports.ListedDocument, error) {
		calls++
		return nil, nil
	}

	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")

	_, _, _ = idx.documentsForDay(context.Background(), 32, today, now, fetch)
	_, _, _ = idx.documentsForDay(context.Background(), 32, yesterday, now, fetch)
	_, _, _ = idx.documentsForDay(context.Background(), 99, today, now, fetch)

	if calls != 3 {
		t.Fatalf("expected 3 fetches for 3 distinct keys, got %d", calls)
	}
}

func TestDocumentIndexExpiresPastDayAfterTTL(t *testing.T) {
	idx := newDocumentIndex()
	calls := 0
	fetch := func(ctx context.Context, day string) ([]ports.ListedDocument, error) {
		calls++
		return nil, nil
	}

	now := time.Now()
	pastDay := now.AddDate(0, 0, -3).Format("2006-01-02")

	_, _, _ = idx.documentsForDay(context.Background(), 32, pastDay, now, fetch)

	later := now.Add(documentIndexPastDayTTL + time.Minute)
	_, fromCache, _ := idx.documentsForDay(context.Background(), 32, pastDay, later, fetch)

	if fromCache {
		t.Fatal("entry must expire after past-day TTL")
	}
	if calls != 2 {
		t.Fatalf("expected refetch after expiry, got %d calls", calls)
	}
}

func TestDocumentIndexDoesNotCacheFailedFetch(t *testing.T) {
	idx := newDocumentIndex()
	calls := 0
	fetch := func(ctx context.Context, day string) ([]ports.ListedDocument, error) {
		calls++
		return nil, errors.New("softpymes timeout")
	}

	now := time.Now()
	day := now.Format("2006-01-02")

	if _, _, err := idx.documentsForDay(context.Background(), 32, day, now, fetch); err == nil {
		t.Fatal("expected error to propagate")
	}
	if _, _, err := idx.documentsForDay(context.Background(), 32, day, now, fetch); err == nil {
		t.Fatal("expected error to propagate on retry")
	}
	if calls != 2 {
		t.Fatalf("a failed fetch must not be cached, got %d calls", calls)
	}
}

package chatretention

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/notification_config/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
)

type repoMock struct {
	cutoff time.Time
	result entities.ChatPurgeResult
	err    error
}

func (m *repoMock) PurgeChatHistory(_ context.Context, cutoff time.Time) (entities.ChatPurgeResult, error) {
	m.cutoff = cutoff
	return m.result, m.err
}

func TestRunUsesOneYearCutoff(t *testing.T) {
	repo := &repoMock{result: entities.ChatPurgeResult{Messages: 3}}
	now := time.Date(2027, 3, 1, 12, 0, 0, 0, time.UTC)
	uc := &useCase{repo: repo, logger: log.New(), now: func() time.Time { return now }}

	result, err := uc.Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if want := now.AddDate(0, 0, -365); !repo.cutoff.Equal(want) {
		t.Fatalf("cutoff %v, want %v", repo.cutoff, want)
	}
	if result.Messages != 3 {
		t.Fatalf("unexpected result %+v", result)
	}
}

func TestRunPropagatesError(t *testing.T) {
	uc := &useCase{repo: &repoMock{err: errors.New("db down")}, logger: log.New(), now: time.Now}
	if _, err := uc.Run(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

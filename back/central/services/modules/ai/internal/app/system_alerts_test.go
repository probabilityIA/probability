package app

import (
	"context"
	"testing"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
	"github.com/secamc93/probability/back/central/shared/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPublishSystemAlertsAvisaChatsSinLeer(t *testing.T) {
	now := time.Date(2026, 9, 16, 15, 0, 0, 0, colombia)
	reader := &readerFake{unread: []entities.UnreadChats{{BusinessID: 26, Count: 13, OldestAt: now.Add(-3 * time.Hour)}, {BusinessID: 30, Count: 0}}}
	alerts := &alertsFake{}
	uc := New(&recommendationFake{}, nil, nil, nil, nil, nil, reader, alerts, log.New()).(*UseCase)
	uc.now = func() time.Time { return now }

	created, err := uc.PublishSystemAlerts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, created)
	require.Len(t, alerts.saved, 1)
	assert.Equal(t, uint(26), alerts.saved[0].BusinessID)
	assert.Equal(t, "Tienes 13 chats de WhatsApp sin leer. El más antiguo espera desde hace 3 h.", alerts.saved[0].Body)
	assert.Equal(t, "/notification-config?tab=conversations", alerts.saved[0].DestinationRoute)
	assert.Equal(t, "13", alerts.saved[0].ReferenceID)
}

func TestPublishSystemAlertsNoRepiteSiNoCrecio(t *testing.T) {
	now := time.Date(2026, 9, 16, 15, 0, 0, 0, colombia)
	reader := &readerFake{unread: []entities.UnreadChats{{BusinessID: 26, Count: 5}}}
	alerts := &alertsFake{last: map[uint]*entities.Alert{26: {ReferenceID: "5", CreatedAt: now.Add(-time.Hour)}}}
	uc := New(&recommendationFake{}, nil, nil, nil, nil, nil, reader, alerts, log.New()).(*UseCase)
	uc.now = func() time.Time { return now }

	created, err := uc.PublishSystemAlerts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 0, created)

	reader.unread[0].Count = 7
	created, err = uc.PublishSystemAlerts(context.Background())
	require.NoError(t, err)
	assert.Equal(t, 1, created)
}

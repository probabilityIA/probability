package app

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/secamc93/probability/back/central/services/modules/ai/internal/domain/entities"
)

const (
	SystemAlertInterval    = 15 * time.Minute
	UnreadChatsRepeatAfter = 4 * time.Hour
	eventTypeUnreadChats   = "whatsapp.unread_chats"
)

func (uc *UseCase) PublishSystemAlerts(ctx context.Context) (int, error) {
	if uc.alerts == nil || uc.businessData == nil {
		return 0, nil
	}
	unread, err := uc.businessData.CountUnreadWhatsAppChats(ctx)
	if err != nil {
		return 0, err
	}
	now := uc.now()
	created := 0
	for _, item := range unread {
		if item.Count <= 0 {
			continue
		}
		last, err := uc.alerts.LastAlertOfType(ctx, item.BusinessID, eventTypeUnreadChats)
		if err != nil {
			return created, err
		}
		if last != nil {
			lastCount, _ := strconv.ParseInt(last.ReferenceID, 10, 64)
			if now.Sub(last.CreatedAt) < UnreadChatsRepeatAfter && item.Count <= lastCount {
				continue
			}
		}
		alert := unreadChatsAlert(item, now)
		ok, err := uc.alerts.SaveAlert(ctx, alert)
		if err != nil {
			return created, err
		}
		if ok {
			created++
		}
	}
	return created, nil
}

func unreadChatsAlert(item entities.UnreadChats, now time.Time) entities.Alert {
	body := fmt.Sprintf("Tienes %d chats de WhatsApp sin leer", item.Count)
	if item.Count == 1 {
		body = "Tienes 1 chat de WhatsApp sin leer"
	}
	if !item.OldestAt.IsZero() {
		body += fmt.Sprintf(". El más antiguo espera desde hace %s", humanSince(now.Sub(item.OldestAt)))
	}
	body += "."
	severity := entities.AlertSeverityInfo
	if item.Count >= 10 {
		severity = entities.AlertSeverityWarning
	}
	return entities.Alert{
		ID:               uuid.NewString(),
		BusinessID:       item.BusinessID,
		EventID:          fmt.Sprintf("%s:%d:%s", eventTypeUnreadChats, item.BusinessID, now.UTC().Format("2006010215")),
		EventType:        eventTypeUnreadChats,
		Severity:         severity,
		Title:            "Chats de WhatsApp sin leer",
		Body:             body,
		DestinationKey:   "notifications",
		DestinationRoute: "/notification-config?tab=conversations",
		ReferenceType:    "unread_chats",
		ReferenceID:      strconv.FormatInt(item.Count, 10),
		CreatedAt:        now,
	}
}

func humanSince(d time.Duration) string {
	minutes := int(d.Minutes())
	switch {
	case minutes < 60:
		if minutes < 1 {
			minutes = 1
		}
		return fmt.Sprintf("%d min", minutes)
	case minutes < 24*60:
		return fmt.Sprintf("%d h", minutes/60)
	default:
		return fmt.Sprintf("%d días", minutes/(24*60))
	}
}

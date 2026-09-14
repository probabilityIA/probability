package app

import (
	"context"
	"fmt"
	"strings"
)

func (uc *useCase) MarkConversationRead(ctx context.Context, conversationID string, businessID uint, userID *uint) error {
	if strings.TrimSpace(conversationID) == "" {
		return fmt.Errorf("conversation ID requerido")
	}
	if businessID == 0 {
		return fmt.Errorf("business_id requerido")
	}
	return uc.messageAuditQuerier.MarkConversationRead(ctx, conversationID, businessID, userID)
}

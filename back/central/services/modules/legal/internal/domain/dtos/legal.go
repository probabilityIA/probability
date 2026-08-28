package dtos

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
)

type PendingDocumentsResult struct {
	RequiresAcceptance bool
	Documents          []entities.LegalDocument
}

type AcceptDocumentsInput struct {
	UserID      uint
	BusinessID  *uint
	DocumentIDs []uint
	IPAddress   string
	UserAgent   string
}

type AcceptDocumentsResult struct {
	AcceptedAt  time.Time
	DocumentIDs []uint
}

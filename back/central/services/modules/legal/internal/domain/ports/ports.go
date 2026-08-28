package ports

import (
	"context"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
)

type IRepository interface {
	GetActiveDocuments(ctx context.Context) ([]entities.LegalDocument, error)
	GetDocumentsByIDs(ctx context.Context, ids []uint) ([]entities.LegalDocument, error)
	GetAcceptedDocumentIDs(ctx context.Context, userID uint) (map[uint]bool, error)
	SaveAcceptances(ctx context.Context, acceptances []entities.LegalAcceptance) error
	IsPlatformUser(ctx context.Context, userID uint) (bool, error)
}

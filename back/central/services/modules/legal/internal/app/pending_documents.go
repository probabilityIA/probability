package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
	legalerrors "github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/errors"
)

func (uc *useCase) GetPendingDocuments(ctx context.Context, userID uint) (*dtos.PendingDocumentsResult, error) {
	if userID == 0 {
		return nil, legalerrors.ErrUserRequired
	}

	esPlataforma, err := uc.repo.IsPlatformUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error al verificar el scope del usuario: %w", err)
	}
	if esPlataforma {
		return &dtos.PendingDocumentsResult{RequiresAcceptance: false, Documents: []entities.LegalDocument{}}, nil
	}

	activos, err := uc.repo.GetActiveDocuments(ctx)
	if err != nil {
		return nil, fmt.Errorf("error al consultar los documentos legales vigentes: %w", err)
	}
	if len(activos) == 0 {
		return &dtos.PendingDocumentsResult{RequiresAcceptance: false, Documents: []entities.LegalDocument{}}, nil
	}

	aceptados, err := uc.repo.GetAcceptedDocumentIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("error al consultar las aceptaciones del usuario: %w", err)
	}

	pendientes := make([]entities.LegalDocument, 0, len(activos))
	for _, doc := range activos {
		if !aceptados[doc.ID] {
			pendientes = append(pendientes, doc)
		}
	}

	return &dtos.PendingDocumentsResult{
		RequiresAcceptance: len(pendientes) > 0,
		Documents:          pendientes,
	}, nil
}

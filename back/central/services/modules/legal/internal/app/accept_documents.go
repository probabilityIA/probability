package app

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
	legalerrors "github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/errors"
)

func (uc *useCase) AcceptDocuments(ctx context.Context, input dtos.AcceptDocumentsInput) (*dtos.AcceptDocumentsResult, error) {
	if input.UserID == 0 {
		return nil, legalerrors.ErrUserRequired
	}
	if len(input.DocumentIDs) == 0 {
		return nil, legalerrors.ErrNoDocumentsSelected
	}

	documentos, err := uc.repo.GetDocumentsByIDs(ctx, input.DocumentIDs)
	if err != nil {
		return nil, fmt.Errorf("error al consultar los documentos legales: %w", err)
	}
	if len(documentos) != len(input.DocumentIDs) {
		return nil, legalerrors.ErrDocumentNotAvailable
	}

	aceptadoEn := time.Now()
	aceptaciones := make([]entities.LegalAcceptance, 0, len(documentos))
	ids := make([]uint, 0, len(documentos))

	for _, doc := range documentos {
		if !doc.IsActive {
			return nil, legalerrors.ErrDocumentNotAvailable
		}
		aceptaciones = append(aceptaciones, entities.LegalAcceptance{
			UserID:          input.UserID,
			LegalDocumentID: doc.ID,
			BusinessID:      input.BusinessID,
			DocumentCode:    doc.Code,
			DocumentVersion: doc.Version,
			DocumentSHA256:  doc.SHA256,
			AcceptedAt:      aceptadoEn,
			IPAddress:       input.IPAddress,
			UserAgent:       input.UserAgent,
			Method:          entities.AcceptanceMethodWebModal,
		})
		ids = append(ids, doc.ID)
	}

	if err := uc.repo.SaveAcceptances(ctx, aceptaciones); err != nil {
		return nil, fmt.Errorf("error al registrar la aceptacion: %w", err)
	}

	uc.log.Info(ctx).
		Uint("user_id", input.UserID).
		Int("documentos", len(ids)).
		Str("ip", input.IPAddress).
		Msg("Aceptacion de documentos legales registrada")

	return &dtos.AcceptDocumentsResult{AcceptedAt: aceptadoEn, DocumentIDs: ids}, nil
}

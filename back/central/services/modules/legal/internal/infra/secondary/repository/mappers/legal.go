package mappers

import (
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
	"github.com/secamc93/probability/back/migration/shared/models"
)

func ToDomainDocuments(filas []models.LegalDocument) []entities.LegalDocument {
	docs := make([]entities.LegalDocument, 0, len(filas))
	for _, f := range filas {
		docs = append(docs, entities.LegalDocument{
			ID:            f.ID,
			Code:          f.Code,
			Version:       f.Version,
			Title:         f.Title,
			FileURL:       f.FileURL,
			SHA256:        f.SHA256,
			ContentHTML:   f.ContentHTML,
			EffectiveFrom: f.EffectiveFrom,
			IsActive:      f.IsActive,
		})
	}
	return docs
}

func ToModelAcceptances(aceptaciones []entities.LegalAcceptance) []models.LegalAcceptance {
	filas := make([]models.LegalAcceptance, 0, len(aceptaciones))
	for _, a := range aceptaciones {
		filas = append(filas, models.LegalAcceptance{
			UserID:          a.UserID,
			LegalDocumentID: a.LegalDocumentID,
			BusinessID:      a.BusinessID,
			DocumentCode:    a.DocumentCode,
			DocumentVersion: a.DocumentVersion,
			DocumentSHA256:  a.DocumentSHA256,
			AcceptedAt:      a.AcceptedAt,
			IPAddress:       a.IPAddress,
			UserAgent:       a.UserAgent,
			Method:          a.Method,
		})
	}
	return filas
}

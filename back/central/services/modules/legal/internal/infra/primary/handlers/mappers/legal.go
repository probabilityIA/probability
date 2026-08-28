package mappers

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/legal/internal/domain/entities"
	"github.com/secamc93/probability/back/central/services/modules/legal/internal/infra/primary/handlers/response"
)

func ToResponseDocuments(docs []entities.LegalDocument) []response.LegalDocument {
	salida := make([]response.LegalDocument, 0, len(docs))
	for _, d := range docs {
		salida = append(salida, response.LegalDocument{
			ID:            d.ID,
			Code:          d.Code,
			Version:       d.Version,
			Title:         d.Title,
			FileURL:       d.FileURL,
			SHA256:        d.SHA256,
			EffectiveFrom: d.EffectiveFrom.Format(time.RFC3339),
		})
	}
	return salida
}

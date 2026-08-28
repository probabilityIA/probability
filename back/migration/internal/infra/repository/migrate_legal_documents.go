package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func (r *Repository) migrateLegalDocuments(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.LegalDocument{}, &models.LegalAcceptance{}); err != nil {
		return fmt.Errorf("failed to auto-migrate legal tables: %w", err)
	}
	return r.seedLegalDocumentsV1(ctx)
}

func (r *Repository) seedLegalDocumentsV1(ctx context.Context) error {
	vigencia := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)

	documentos := []models.LegalDocument{
		{
			Code:          "terms_of_service",
			Version:       "1.0",
			Title:         "Terminos y Condiciones de Uso",
			FileURL:       "/legal/terminos-y-condiciones-v1.0.pdf",
			SHA256:        "25639ec867d7261189d61f262a80bdbedd36628024c0d8959d8cfbb06730c5ab",
			EffectiveFrom: vigencia,
			IsActive:      true,
			RequiresStaff: true,
		},
		{
			Code:          "privacy_policy",
			Version:       "1.0",
			Title:         "Politica de Tratamiento de Datos Personales",
			FileURL:       "/legal/politica-datos-personales-v1.0.pdf",
			SHA256:        "f97b9ff947034a6b8b234bda5a039285b833809ec4f9ec0ccb7ec0e8b05eb7ba",
			EffectiveFrom: vigencia,
			IsActive:      true,
			RequiresStaff: true,
		},
	}

	for _, doc := range documentos {
		res := r.db.Conn(ctx).Exec(`
INSERT INTO legal_documents (code, version, title, file_url, sha256, effective_from, is_active, requires_staff, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, TRUE, TRUE, NOW(), NOW())
ON CONFLICT (code, version) DO NOTHING
`, doc.Code, doc.Version, doc.Title, doc.FileURL, doc.SHA256, doc.EffectiveFrom)

		if res.Error != nil {
			return fmt.Errorf("seed legal document %s v%s: %w", doc.Code, doc.Version, res.Error)
		}
		if res.RowsAffected > 0 {
			fmt.Printf("seed legal document: %s v%s creado\n", doc.Code, doc.Version)
		}
	}

	return nil
}

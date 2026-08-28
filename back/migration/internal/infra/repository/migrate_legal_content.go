package repository

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"fmt"

	"github.com/secamc93/probability/back/migration/shared/models"
)

//go:embed data/legal_terms_of_service_v1.0.html
var legalTermsV1HTML string

//go:embed data/legal_privacy_policy_v1.0.html
var legalPrivacyV1HTML string

type legalContent struct {
	code    string
	version string
	html    string
}

func (r *Repository) migrateLegalContent(ctx context.Context) error {
	if err := r.db.Conn(ctx).AutoMigrate(&models.LegalDocument{}); err != nil {
		return fmt.Errorf("failed to add content_html to legal_documents: %w", err)
	}

	documentos := []legalContent{
		{code: "terms_of_service", version: "1.0", html: legalTermsV1HTML},
		{code: "privacy_policy", version: "1.0", html: legalPrivacyV1HTML},
	}

	for _, doc := range documentos {
		suma := sha256.Sum256([]byte(doc.html))
		err := r.db.Conn(ctx).
			Model(&models.LegalDocument{}).
			Where("code = ? AND version = ?", doc.code, doc.version).
			Updates(map[string]any{
				"content_html": doc.html,
				"sha256":       hex.EncodeToString(suma[:]),
			}).Error
		if err != nil {
			return fmt.Errorf("failed to seed content for %s v%s: %w", doc.code, doc.version, err)
		}
	}

	return nil
}

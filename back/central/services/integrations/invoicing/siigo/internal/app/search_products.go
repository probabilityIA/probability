package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/siigo/internal/domain/dtos"
)

func (uc *invoicingUseCase) SearchProducts(ctx context.Context, integrationID uint, termino string, limite int) ([]dtos.CatalogItem, error) {
	if integrationID == 0 {
		return nil, fmt.Errorf("integration_id es requerido")
	}

	creds, err := uc.resolveWebhookCredentials(ctx, fmt.Sprintf("%d", integrationID))
	if err != nil {
		return nil, err
	}

	items, err := uc.siigoClient.SearchProducts(ctx, creds, termino, limite)
	if err != nil {
		return nil, err
	}

	out := make([]dtos.CatalogItem, 0, len(items))
	for _, i := range items {
		out = append(out, dtos.CatalogItem{Code: i.Code, Name: i.Name, Detail: i.Type})
	}
	return out, nil
}

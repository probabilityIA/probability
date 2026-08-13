package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/meli/internal/domain"
)

const loteMultiget = 20

func (c *MeliClient) GetItemsStock(ctx context.Context, accessToken string, itemIDs []string) ([]domain.ChannelStock, error) {
	unicos := make([]string, 0, len(itemIDs))
	vistos := make(map[string]bool, len(itemIDs))
	for _, id := range itemIDs {
		limpio := strings.TrimSpace(id)
		if limpio == "" || vistos[limpio] {
			continue
		}
		vistos[limpio] = true
		unicos = append(unicos, limpio)
	}
	if len(unicos) == 0 {
		return []domain.ChannelStock{}, nil
	}

	stocks := make([]domain.ChannelStock, 0, len(unicos))
	for inicio := 0; inicio < len(unicos); inicio += loteMultiget {
		fin := inicio + loteMultiget
		if fin > len(unicos) {
			fin = len(unicos)
		}

		detalles, err := c.multigetItems(ctx, accessToken, unicos[inicio:fin])
		if err != nil {
			return nil, err
		}

		for _, detalle := range detalles {
			if len(detalle.Variations) == 0 {
				stocks = append(stocks, domain.ChannelStock{
					ExternalItemID:    detalle.ID,
					AvailableQuantity: detalle.AvailableQuantity,
					Found:             true,
				})
				continue
			}
			for _, variacion := range detalle.Variations {
				stocks = append(stocks, domain.ChannelStock{
					ExternalItemID:    detalle.ID,
					ExternalVariantID: fmt.Sprintf("%d", variacion.ID),
					AvailableQuantity: variacion.AvailableQuantity,
					Found:             true,
				})
			}
		}
	}

	return stocks, nil
}

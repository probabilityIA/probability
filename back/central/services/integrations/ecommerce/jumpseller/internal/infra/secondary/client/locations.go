package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/ecommerce/jumpseller/internal/infra/secondary/client/response"
)

func (c *JumpsellerClient) GetLocations(ctx context.Context, cred domain.Credential) ([]domain.Location, error) {
	raw, err := c.do(ctx, cred, http.MethodGet, "/locations.json", nil, nil)
	if err != nil {
		return nil, err
	}

	var envelopes []response.LocationEnvelope
	if err := json.Unmarshal(raw, &envelopes); err != nil {
		return nil, fmt.Errorf("jumpseller client: parsing locations: %w", err)
	}

	locations := make([]domain.Location, 0, len(envelopes))
	for _, envelope := range envelopes {
		locations = append(locations, envelope.ToDomain())
	}
	return locations, nil
}

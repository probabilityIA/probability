package app

import (
	"context"
	"fmt"
	"strings"

	"github.com/secamc93/probability/back/central/services/integrations/transport/shipit/internal/domain"
)

func (uc *useCase) Quote(ctx context.Context, baseURL string, creds domain.Credentials, req domain.GuideRequest) (*domain.QuoteResponse, error) {
	uc.log.Info(ctx).
		Str("destination_city", req.Destination.City).
		Int("packages", len(req.Packages)).
		Msg("Quoting Shipit shipment")

	originID, _, err := uc.resolveCommune(ctx, baseURL, creds, req.Origin)
	if err != nil {
		return nil, fmt.Errorf("comuna de origen: %w", err)
	}
	destinyID, _, err := uc.resolveCommune(ctx, baseURL, creds, req.Destination)
	if err != nil {
		return nil, fmt.Errorf("comuna de destino: %w", err)
	}

	sizes := buildSizes(req.Packages)
	ratesReq := domain.RatesRequest{
		Parcel: domain.Parcel{
			Length:        sizes.Length,
			Width:         sizes.Width,
			Height:        sizes.Height,
			Weight:        sizes.Weight,
			OriginID:      originID,
			DestinyID:     destinyID,
			TypeOfDestiny: "domicilio",
			Algorithm:     1,
		},
	}

	resp, err := uc.client.Rates(baseURL, creds, ratesReq)
	if err != nil {
		return nil, err
	}

	return mapRatesResponse(resp), nil
}

func (uc *useCase) Generate(ctx context.Context, baseURL string, creds domain.Credentials, req domain.GuideRequest, sandbox bool) (*domain.GenerateResponse, error) {
	uc.log.Info(ctx).
		Str("reference", req.MyShipmentReference).
		Bool("sandbox", sandbox).
		Msg("Generating Shipit shipment")

	communeID, communeName, err := uc.resolveCommune(ctx, baseURL, creds, req.Destination)
	if err != nil {
		return nil, fmt.Errorf("comuna de destino: %w", err)
	}

	body := buildShipmentBody(req, communeID, communeName, sandbox)
	resp, err := uc.client.CreateShipment(baseURL, creds, domain.ShipmentRequest{Shipment: body})
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, domain.ErrEmptyResponse
	}

	return mapGenerateResponse(req, resp), nil
}

func (uc *useCase) Track(ctx context.Context, baseURL string, creds domain.Credentials, trackingNumber string) (*domain.TrackOutput, error) {
	uc.log.Info(ctx).Str("tracking_number", trackingNumber).Msg("Tracking Shipit shipment")

	resp, err := uc.client.Track(baseURL, creds, trackingNumber)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, domain.ErrEmptyResponse
	}

	return mapTrackingResponse(resp), nil
}

func (uc *useCase) TestAuthentication(ctx context.Context, baseURL string, creds domain.Credentials) error {
	return uc.client.TestAuthentication(baseURL, creds)
}

func (uc *useCase) resolveCommune(ctx context.Context, baseURL string, creds domain.Credentials, addr domain.Address) (uint, string, error) {
	if addr.CommuneID != 0 {
		name := addr.CommuneName
		if name == "" {
			name = addr.City
		}
		return addr.CommuneID, name, nil
	}

	candidates := []string{addr.CommuneName, addr.City, addr.Suburb, addr.State}

	communes, err := uc.getCommunes(ctx, baseURL, creds)
	if err != nil {
		return 0, "", err
	}

	for _, candidate := range candidates {
		key := normalizeCommuneName(candidate)
		if key == "" {
			continue
		}
		if id, ok := communes[key]; ok {
			return id, strings.ToUpper(candidate), nil
		}
	}

	return 0, "", fmt.Errorf("%w (busque: %s)", domain.ErrCommuneNotFound, strings.Join(nonEmpty(candidates), ", "))
}

func (uc *useCase) getCommunes(ctx context.Context, baseURL string, creds domain.Credentials) (map[string]uint, error) {
	uc.communesMu.Lock()
	cached, ok := uc.communesCache[baseURL]
	uc.communesMu.Unlock()
	if ok {
		return cached, nil
	}

	list, err := uc.client.GetCommunes(baseURL, creds)
	if err != nil {
		return nil, err
	}

	index := make(map[string]uint, len(list))
	for _, commune := range list {
		index[normalizeCommuneName(commune.Name)] = commune.ID
	}

	uc.communesMu.Lock()
	uc.communesCache[baseURL] = index
	uc.communesMu.Unlock()

	uc.log.Info(ctx).Int("communes", len(index)).Msg("Shipit communes cached")
	return index, nil
}

func normalizeCommuneName(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	replacer := strings.NewReplacer(
		"\u00e1", "a", "\u00e9", "e", "\u00ed", "i", "\u00f3", "o", "\u00fa", "u",
		"\u00f1", "n", "\u00fc", "u",
	)
	return replacer.Replace(name)
}

func nonEmpty(values []string) []string {
	result := make([]string, 0, len(values))
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			result = append(result, v)
		}
	}
	return result
}

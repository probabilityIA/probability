package usecaseenvioclick

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/domain"
	"github.com/secamc93/probability/back/central/services/modules/shipments/internal/infra/secondary/envioclick"
	"github.com/secamc93/probability/back/central/shared/log"
)

type UseCaseEnvioClick struct {
	logger           log.ILogger
	envioclickClient *envioclick.Client
	repo             domain.IRepository
}

func New(logger log.ILogger, client *envioclick.Client, repo domain.IRepository) *UseCaseEnvioClick {
	return &UseCaseEnvioClick{
		logger:           logger,
		envioclickClient: client,
		repo:             repo,
	}
}

func (uc *UseCaseEnvioClick) QuoteShipment(ctx context.Context, req domain.EnvioClickQuoteRequest) (*domain.EnvioClickQuoteResponse, error) {
	uc.logger.Info().Msg("Quoting shipment with EnvioClick")
	// Here we could add logic to validate or enrich the request if needed
	return uc.envioclickClient.Quote(req)
}

func (uc *UseCaseEnvioClick) GenerateGuide(ctx context.Context, req domain.EnvioClickQuoteRequest) (*domain.EnvioClickGenerateResponse, error) {
	uc.logger.Info().Msg("Generating guide with EnvioClick")
	// Here we could add logic to save the guide info to the database (update the order)
	// For now we just pass through to the client
	resp, err := uc.envioclickClient.Generate(req)
	if err != nil {
		return nil, err
	}

	// Deduct Wallet Balance
	orderID := req.OrderUUID
	if orderID == "" {
		orderID = req.ExternalOrderID // Fallback to external ID if UUID not provided (though relation might fail)
	}

	if req.TotalCost > 0 && orderID != "" {
		// Attempt to deduct balance
		businessID, err := uc.repo.GetOrderBusinessID(ctx, orderID)
		if err != nil {
			uc.logger.Error().Err(err).Msg("Failed to find order business ID for wallet deduction")
			return nil, fmt.Errorf("no se pudo validar la orden para descuento de billetera")
		}
		if businessID != nil {
			reference := fmt.Sprintf("Guía: %s, Tracking: %s", req.MyShipmentReference, resp.Data.TrackingNumber)
			if err := uc.repo.DeductWalletBalance(ctx, *businessID, req.TotalCost, reference); err != nil {
				uc.logger.Error().Err(err).Msg("Failed to deduct wallet balance")
				return nil, fmt.Errorf("no se pudo deducir el saldo de la billetera: %w", err)
			}
		}
	}

	shipment := &domain.Shipment{
		OrderID:            &orderID,
		TrackingNumber:     &resp.Data.TrackingNumber,
		GuideURL:           &resp.Data.LabelURL,
		Status:             "pending",
		ClientName:         req.Destination.FirstName + " " + req.Destination.LastName,
		DestinationAddress: req.Destination.Address + ", " + req.Destination.Suburb + ", " + req.Destination.DaneCode,
		TotalCost:          &req.TotalCost,
		// Note: Carrier info is not available in GenerateResponse, purely IDRate based.
		// We might need to fetch it or pass it in request if needed.
	}

	if err := uc.repo.CreateShipment(ctx, shipment); err != nil {
		uc.logger.Error().Err(err).Msg("Failed to persist shipment after generating guide")
		return nil, fmt.Errorf("failed to persist shipment locally: %w", err)
	}

	return resp, nil
}

func (uc *UseCaseEnvioClick) TrackShipment(ctx context.Context, trackingNumber string) (*domain.EnvioClickTrackingResponse, error) {
	uc.logger.Info().Str("tracking_number", trackingNumber).Msg("Tracking shipment with EnvioClick")
	return uc.envioclickClient.Track(trackingNumber)
}

func (uc *UseCaseEnvioClick) CancelShipment(ctx context.Context, idShipment string) (*domain.EnvioClickCancelResponse, error) {
	uc.logger.Info().Str("id_shipment", idShipment).Msg("Canceling shipment with EnvioClick")
	return uc.envioclickClient.Cancel(idShipment)
}

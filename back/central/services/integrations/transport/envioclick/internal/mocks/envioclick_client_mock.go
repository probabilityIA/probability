package mocks

import (
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
)

type EnvioClickClientMock struct {
	QuoteFn              func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error)
	GenerateFn           func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error)
	TrackFn              func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error)
	TrackByOrdersBatchFn func(baseURL, apiKey string, orders []int64) (*domain.TrackingResponse, error)
	CancelFn             func(baseURL, apiKey string, idShipment string) (*domain.CancelResponse, error)
	CancelBatchFn        func(baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error)
}

func (m *EnvioClickClientMock) Quote(baseURL, apiKey string, req domain.QuoteRequest, _ *domain.SyncMeta) (*domain.QuoteResponse, error) {
	if m.QuoteFn != nil {
		return m.QuoteFn(baseURL, apiKey, req)
	}
	return &domain.QuoteResponse{}, nil
}

func (m *EnvioClickClientMock) Generate(baseURL, apiKey string, req domain.QuoteRequest, _ *domain.SyncMeta) (*domain.GenerateResponse, error) {
	if m.GenerateFn != nil {
		return m.GenerateFn(baseURL, apiKey, req)
	}
	return &domain.GenerateResponse{}, nil
}

func (m *EnvioClickClientMock) Track(baseURL, apiKey string, trackingNumber string, _ *domain.SyncMeta) (*domain.TrackingResponse, error) {
	if m.TrackFn != nil {
		return m.TrackFn(baseURL, apiKey, trackingNumber)
	}
	return &domain.TrackingResponse{}, nil
}

func (m *EnvioClickClientMock) TrackByOrdersBatch(baseURL, apiKey string, orders []int64, _ *domain.SyncMeta) (*domain.TrackingResponse, error) {
	if m.TrackByOrdersBatchFn != nil {
		return m.TrackByOrdersBatchFn(baseURL, apiKey, orders)
	}
	return &domain.TrackingResponse{}, nil
}

func (m *EnvioClickClientMock) Cancel(baseURL, apiKey string, idShipment string, _ *domain.SyncMeta) (*domain.CancelResponse, error) {
	if m.CancelFn != nil {
		return m.CancelFn(baseURL, apiKey, idShipment)
	}
	return &domain.CancelResponse{}, nil
}

func (m *EnvioClickClientMock) CancelBatch(baseURL, apiKey string, req domain.CancelBatchRequest, _ *domain.SyncMeta) (*domain.CancelBatchResponse, error) {
	if m.CancelBatchFn != nil {
		return m.CancelBatchFn(baseURL, apiKey, req)
	}
	return &domain.CancelBatchResponse{}, nil
}

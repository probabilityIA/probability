package app

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/domain"
	"github.com/secamc93/probability/back/central/services/integrations/transport/envioclick/internal/mocks"
)

const (
	testBaseURL        = "https://api.envioclick.com"
	testAPIKey         = "test-api-key"
	testTrackingNumber = "EC123456789CO"
)

func TestMain(m *testing.M) {
	cancelVerifyDelay = 0
	os.Exit(m.Run())
}

func trackSequence(statuses ...string) (func(string, string, string) (*domain.TrackingResponse, error), *int) {
	calls := 0
	fn := func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
		idx := calls
		calls++
		if idx >= len(statuses) {
			idx = len(statuses) - 1
		}
		return &domain.TrackingResponse{
			Data: domain.TrackingData{Status: statuses[idx]},
		}, nil
	}
	return fn, &calls
}

func TestCancel_Success_WithIdOrder(t *testing.T) {
	ctx := context.Background()

	trackFn, trackCalls := trackSequence("Pendiente de Recoleccion", "Cancelado")
	client := &mocks.EnvioClickClientMock{
		TrackFn: trackFn,
		CancelBatchFn: func(baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error) {
			return &domain.CancelBatchResponse{
				Status: "success",
				Data: domain.CancelBatchData{
					OnlyCancelOrders: []int64{req.IDOrders[0]},
				},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, 42, nil)

	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if result.Status != "success" {
		t.Errorf("se esperaba status 'success', se obtuvo: %q", result.Status)
	}
	if !strings.Contains(result.Message, "exitosamente") {
		t.Errorf("se esperaba mensaje de exito, se obtuvo: %q", result.Message)
	}
	if *trackCalls < 2 {
		t.Errorf("se esperaba un Track de verificacion despues de cancelar, hubo %d llamadas", *trackCalls)
	}
}

func TestCancel_Success_CaseInsensitive(t *testing.T) {
	ctx := context.Background()

	trackFn, _ := trackSequence("pendiente de recoleccion", "CANCELADO")
	client := &mocks.EnvioClickClientMock{
		TrackFn: trackFn,
		CancelBatchFn: func(baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error) {
			return &domain.CancelBatchResponse{
				Status: "success",
				Data:   domain.CancelBatchData{},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, 99, nil)

	if err != nil {
		t.Fatalf("se esperaba nil error con status en minusculas, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if result.Status != "success" {
		t.Errorf("se esperaba status 'success', se obtuvo: %q", result.Status)
	}
}

func TestCancel_Error_CarrierAcceptsButShipmentStaysActive(t *testing.T) {
	ctx := context.Background()

	trackFn, trackCalls := trackSequence("Pendiente de Recoleccion")
	client := &mocks.EnvioClickClientMock{
		TrackFn: trackFn,
		CancelBatchFn: func(baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error) {
			return &domain.CancelBatchResponse{
				Status: "OK",
				Data:   domain.CancelBatchData{},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, 4613038, nil)

	if err == nil {
		t.Fatal("se esperaba error cuando la guia sigue activa tras aceptar la cancelacion, se obtuvo nil")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil para no marcar el envio como cancelado, se obtuvo: %+v", result)
	}
	if !strings.Contains(err.Error(), "sigue activo") {
		t.Errorf("el error deberia indicar que el envio sigue activo, se obtuvo: %q", err.Error())
	}
	if *trackCalls != 1+cancelVerifyAttempts {
		t.Errorf("se esperaban %d Tracks (1 previo + %d de verificacion), hubo %d", 1+cancelVerifyAttempts, cancelVerifyAttempts, *trackCalls)
	}
}

func TestCancel_Success_FallbackSingular(t *testing.T) {
	ctx := context.Background()

	singularCalled := false
	client := &mocks.EnvioClickClientMock{
		TrackFn: func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
			return &domain.TrackingResponse{
				Data: domain.TrackingData{
					Status: "Pendiente",
				},
			}, nil
		},
		CancelFn: func(baseURL, apiKey string, idShipment string) (*domain.CancelResponse, error) {
			singularCalled = true
			if idShipment != testTrackingNumber {
				return nil, errors.New("tracking number incorrecto")
			}
			return &domain.CancelResponse{
				Status:  "success",
				Message: "Cancelado",
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, 0, nil)

	if err != nil {
		t.Fatalf("se esperaba nil error con idOrder=0, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if !singularCalled {
		t.Error("se esperaba que se llamara al Cancel singular cuando idOrder=0")
	}
}

func TestCancel_Error_NotCancellable(t *testing.T) {
	ctx := context.Background()

	client := &mocks.EnvioClickClientMock{
		TrackFn: func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
			return &domain.TrackingResponse{
				Data: domain.TrackingData{
					Status: "En transito",
				},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, 42, nil)

	if err == nil {
		t.Fatal("se esperaba error para status no cancelable, se obtuvo nil")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil cuando no es cancelable, se obtuvo: %+v", result)
	}
	if !strings.Contains(err.Error(), "En transito") {
		t.Errorf("el error deberia mencionar el status actual, se obtuvo: %q", err.Error())
	}
}

func TestCancel_Error_TrackFails(t *testing.T) {
	ctx := context.Background()

	trackErr := errors.New("timeout al conectar con envioclick")
	client := &mocks.EnvioClickClientMock{
		TrackFn: func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
			return nil, trackErr
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, 42, nil)

	if err == nil {
		t.Fatal("se esperaba error cuando Track falla, se obtuvo nil")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil cuando Track falla, se obtuvo: %+v", result)
	}
	if !errors.Is(err, trackErr) {
		t.Errorf("se esperaba que el error envolviera el error original, se obtuvo: %v", err)
	}
}

func TestCancel_Error_BatchReturnsNotValid(t *testing.T) {
	ctx := context.Background()

	const orderID int64 = 777

	client := &mocks.EnvioClickClientMock{
		TrackFn: func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
			return &domain.TrackingResponse{
				Data: domain.TrackingData{
					Status: "Pendiente de Recoleccion",
				},
			}, nil
		},
		CancelBatchFn: func(baseURL, apiKey string, req domain.CancelBatchRequest) (*domain.CancelBatchResponse, error) {
			return &domain.CancelBatchResponse{
				Status: "error",
				Data: domain.CancelBatchData{
					NotValidOrders: []int64{orderID},
				},
			}, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Cancel(ctx, testBaseURL, testAPIKey, testTrackingNumber, orderID, nil)

	if err == nil {
		t.Fatal("se esperaba error cuando la orden es not_valid, se obtuvo nil")
	}
	if result != nil {
		t.Errorf("se esperaba resultado nil para no marcar el envio como cancelado, se obtuvo: %+v", result)
	}
	if !strings.Contains(err.Error(), "no valido") {
		t.Errorf("se esperaba mensaje sobre orden no valida, se obtuvo: %q", err.Error())
	}
}

func TestQuote_Success(t *testing.T) {
	ctx := context.Background()

	expectedResponse := &domain.QuoteResponse{
		Status: "success",
		Data: domain.QuoteData{
			Rates: []domain.Rate{
				{IDRate: 1, Product: "Express", Flete: 15000.0},
			},
		},
	}

	client := &mocks.EnvioClickClientMock{
		QuoteFn: func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.QuoteResponse, error) {
			return expectedResponse, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	req := domain.QuoteRequest{
		Description: "Paquete de prueba",
		Packages: []domain.Package{
			{Weight: 1.0, Height: 10, Width: 10, Length: 10},
		},
	}

	result, err := uc.Quote(ctx, testBaseURL, testAPIKey, req, nil)

	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if len(result.Data.Rates) != 1 {
		t.Errorf("se esperaba 1 tarifa, se obtuvo: %d", len(result.Data.Rates))
	}
	if result.Data.Rates[0].IDRate != 1 {
		t.Errorf("se esperaba IDRate=1, se obtuvo: %d", result.Data.Rates[0].IDRate)
	}
}

func TestGenerate_Success(t *testing.T) {
	ctx := context.Background()

	expectedResponse := &domain.GenerateResponse{
		Status: "success",
		Data: domain.GenerateData{
			TrackingNumber: "EC987654321CO",
			LabelURL:       "https://labels.envioclick.com/EC987654321CO.pdf",
			Carrier:        "Servientrega",
			IDOrder:        55,
		},
	}

	client := &mocks.EnvioClickClientMock{
		GenerateFn: func(baseURL, apiKey string, req domain.QuoteRequest) (*domain.GenerateResponse, error) {
			return expectedResponse, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	req := domain.QuoteRequest{
		IDRate:      10,
		Description: "Pedido 1234",
	}

	result, err := uc.Generate(ctx, testBaseURL, testAPIKey, req, nil)

	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if result.Data.TrackingNumber != "EC987654321CO" {
		t.Errorf("se esperaba TrackingNumber 'EC987654321CO', se obtuvo: %q", result.Data.TrackingNumber)
	}
	if result.Data.IDOrder != 55 {
		t.Errorf("se esperaba IDOrder=55, se obtuvo: %d", result.Data.IDOrder)
	}
}

func TestTrack_Success(t *testing.T) {
	ctx := context.Background()

	expectedResponse := &domain.TrackingResponse{
		Status: "success",
		Data: domain.TrackingData{
			TrackingNumber: testTrackingNumber,
			Carrier:        "Coordinadora",
			Status:         "Entregado",
			StatusDetail:   "Entregado al destinatario",
		},
	}

	client := &mocks.EnvioClickClientMock{
		TrackFn: func(baseURL, apiKey string, trackingNumber string) (*domain.TrackingResponse, error) {
			if trackingNumber != testTrackingNumber {
				return nil, errors.New("tracking number no encontrado")
			}
			return expectedResponse, nil
		},
	}

	uc := New(client, &mocks.LoggerMock{})

	result, err := uc.Track(ctx, testBaseURL, testAPIKey, testTrackingNumber, nil)

	if err != nil {
		t.Fatalf("se esperaba nil error, se obtuvo: %v", err)
	}
	if result == nil {
		t.Fatal("se esperaba resultado no nil")
	}
	if result.Data.Status != "Entregado" {
		t.Errorf("se esperaba status 'Entregado', se obtuvo: %q", result.Data.Status)
	}
	if result.Data.Carrier != "Coordinadora" {
		t.Errorf("se esperaba carrier 'Coordinadora', se obtuvo: %q", result.Data.Carrier)
	}
}

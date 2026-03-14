package usecases

import (
	"testing"

	"github.com/secamc93/probability/back/testing/integrations/envioclick/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/logx"
)

func TestMockSimulator_ConsistentGuideGeneration(t *testing.T) {
	// Setup
	logger := logx.NewLogger()
	repo := domain.NewShipmentRepository()
	s3Provider := nil
	urlBase := ""
	simulator := NewAPISimulator(repo, logger, s3Provider, urlBase)

	// Test data
	quoteReq := domain.QuoteRequest{
		ContentValue: 100000,
		Packages: []domain.Package{
			{
				Weight: 1.0,
				Height: 10,
				Width:  10,
				Length: 10,
			},
		},
		Origin: domain.Address{
			DaneCode: "11001000", // Bogota
		},
		Destination: domain.Address{
			DaneCode: "05001000", // Medellin
		},
	}

	// Step 1: Generate quote
	quoteResp, err := simulator.HandleQuote(quoteReq)
	require.NoError(t, err)
	require.NotNil(t, quoteResp)
	require.Equal(t, "success", quoteResp.Status)
	require.Greater(t, len(quoteResp.Data.Rates), 0)

	// Get the first rate
	firstRate := quoteResp.Data.Rates[0]
	rateID := firstRate.IDRate
	originalCarrier := firstRate.Carrier
	originalProduct := firstRate.Product

	// Step 2: Generate shipment with first rate
	generateReq1 := domain.QuoteRequest{
		IDRate:              rateID,
		ContentValue:        100000,
		MyShipmentReference: "test-ref-1",
		Packages:            quoteReq.Packages,
		Origin:              quoteReq.Origin,
		Destination:         quoteReq.Destination,
	}

	resp1, err := simulator.HandleGenerate(generateReq1)
	require.NoError(t, err)
	require.NotNil(t, resp1)
	require.Equal(t, "success", resp1.Status)
	require.NotEmpty(t, resp1.Data.TrackingNumber)
	require.NotEmpty(t, resp1.Data.LabelURL)

	trackingNumber1 := resp1.Data.TrackingNumber
	labelURL1 := resp1.Data.LabelURL

	// Verify the stored shipment has the correct carrier
	storedShipment1, exists := repo.GetByTracking(trackingNumber1)
	require.True(t, exists, "Shipment should exist in repository")
	require.Equal(t, originalCarrier, storedShipment1.Carrier, "Carrier should match the quote")
	require.Equal(t, originalProduct, storedShipment1.Product, "Product should match the quote")

	// Step 3: Generate shipment AGAIN with the same rate
	generateReq2 := domain.QuoteRequest{
		IDRate:              rateID,
		ContentValue:        100000,
		MyShipmentReference: "test-ref-2",
		Packages:            quoteReq.Packages,
		Origin:              quoteReq.Origin,
		Destination:         quoteReq.Destination,
	}

	resp2, err := simulator.HandleGenerate(generateReq2)
	require.NoError(t, err)
	require.NotNil(t, resp2)
	require.Equal(t, "success", resp2.Status)

	trackingNumber2 := resp2.Data.TrackingNumber
	labelURL2 := resp2.Data.LabelURL

	// Verify consistency: same rate should generate same tracking number
	assert.Equal(t, trackingNumber1, trackingNumber2, "Same rate should generate same tracking number")
	assert.Equal(t, labelURL1, labelURL2, "Same rate should generate same label URL")

	// Verify the stored shipment still has correct carrier
	storedShipment2, exists := repo.GetByTracking(trackingNumber2)
	require.True(t, exists, "Shipment should exist in repository")
	require.Equal(t, originalCarrier, storedShipment2.Carrier, "Carrier should still match the quote")
	require.Equal(t, originalProduct, storedShipment2.Product, "Product should still match the quote")

	// Step 4: Track the shipment and verify data
	trackResp, err := simulator.HandleTrack(trackingNumber1)
	require.NoError(t, err)
	require.NotNil(t, trackResp)
	require.Equal(t, "success", trackResp.Status)
	require.Equal(t, trackingNumber1, trackResp.Data.TrackingNumber)
	require.Equal(t, originalCarrier, trackResp.Data.Carrier, "Tracked shipment should have correct carrier")
	require.Greater(t, len(trackResp.Data.Events), 0, "Tracking should have events")
}

func TestMockSimulator_DifferentRatesDifferentCarriers(t *testing.T) {
	// Setup
	logger := logx.NewLogger()
	repo := domain.NewShipmentRepository()
	simulator := NewAPISimulator(repo, logger, nil, "")

	// Test data
	quoteReq := domain.QuoteRequest{
		ContentValue: 100000,
		Packages: []domain.Package{
			{
				Weight: 1.0,
				Height: 10,
				Width:  10,
				Length: 10,
			},
		},
		Origin: domain.Address{
			DaneCode: "11001000", // Bogota
		},
		Destination: domain.Address{
			DaneCode: "05001000", // Medellin
		},
	}

	// Generate quote with multiple rates
	quoteResp, err := simulator.HandleQuote(quoteReq)
	require.NoError(t, err)
	require.Greater(t, len(quoteResp.Data.Rates), 1, "Should have multiple rates")

	// Get different rates
	rate1 := quoteResp.Data.Rates[0]
	rate2 := quoteResp.Data.Rates[1]

	// Generate shipment with first rate
	resp1, err := simulator.HandleGenerate(domain.QuoteRequest{
		IDRate:       rate1.IDRate,
		ContentValue: 100000,
		Packages:     quoteReq.Packages,
		Origin:       quoteReq.Origin,
		Destination:  quoteReq.Destination,
	})
	require.NoError(t, err)

	// Generate shipment with second rate
	resp2, err := simulator.HandleGenerate(domain.QuoteRequest{
		IDRate:       rate2.IDRate,
		ContentValue: 100000,
		Packages:     quoteReq.Packages,
		Origin:       quoteReq.Origin,
		Destination:  quoteReq.Destination,
	})
	require.NoError(t, err)

	// Get stored shipments
	ship1, _ := repo.GetByTracking(resp1.Data.TrackingNumber)
	ship2, _ := repo.GetByTracking(resp2.Data.TrackingNumber)

	// Verify different rates should have different tracking numbers
	assert.NotEqual(t, resp1.Data.TrackingNumber, resp2.Data.TrackingNumber, "Different rates should have different tracking numbers")

	// Verify that if rates are from different carriers, shipments use those carriers
	if rate1.IDCarrier != rate2.IDCarrier {
		assert.NotEqual(t, ship1.Carrier, ship2.Carrier, "Different rate carriers should result in different shipment carriers")
	}
}

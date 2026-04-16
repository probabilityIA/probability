package domain

import (
	"testing"
	"time"
)

func TestGetOrdersParams_ToQueryString_Defaults(t *testing.T) {
	// Arrange: parametros vacios
	p := &GetOrdersParams{}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["status"] != "any" {
		t.Errorf("status por defecto incorrecto: got %q, want %q", result["status"], "any")
	}
	if result["limit"] != "250" {
		t.Errorf("limit por defecto incorrecto: got %q, want %q", result["limit"], "250")
	}
}

func TestGetOrdersParams_ToQueryString_WithStatus(t *testing.T) {
	// Arrange
	p := &GetOrdersParams{
		Status: "open",
		Limit:  50,
	}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["status"] != "open" {
		t.Errorf("status incorrecto: got %q, want %q", result["status"], "open")
	}
	if result["limit"] != "50" {
		t.Errorf("limit incorrecto: got %q, want %q", result["limit"], "50")
	}
}

func TestGetOrdersParams_ToQueryString_LimitCappedAt250(t *testing.T) {
	// Arrange: el limit no puede superar 250
	p := &GetOrdersParams{
		Limit: 500,
	}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["limit"] != "250" {
		t.Errorf("limit deberia estar limitado a 250, got %q", result["limit"])
	}
}

func TestGetOrdersParams_ToQueryString_WithDates(t *testing.T) {
	// Arrange
	minTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	maxTime := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	p := &GetOrdersParams{
		CreatedAtMin: &minTime,
		CreatedAtMax: &maxTime,
	}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["created_at_min"] != minTime.Format(time.RFC3339) {
		t.Errorf("created_at_min incorrecto: got %q", result["created_at_min"])
	}
	if result["created_at_max"] != maxTime.Format(time.RFC3339) {
		t.Errorf("created_at_max incorrecto: got %q", result["created_at_max"])
	}
}

func TestGetOrdersParams_ToQueryString_WithFields(t *testing.T) {
	// Arrange
	p := &GetOrdersParams{
		Fields: []string{"id", "name", "financial_status"},
	}

	// Act
	result := p.ToQueryString()

	// Assert
	expected := "id,name,financial_status"
	if result["fields"] != expected {
		t.Errorf("fields incorrecto: got %q, want %q", result["fields"], expected)
	}
}

func TestGetOrdersParams_ToQueryString_WithFinancialAndFulfillmentStatus(t *testing.T) {
	// Arrange
	p := &GetOrdersParams{
		FinancialStatus:   "paid",
		FulfillmentStatus: "shipped",
	}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["financial_status"] != "paid" {
		t.Errorf("financial_status incorrecto: got %q, want %q", result["financial_status"], "paid")
	}
	if result["fulfillment_status"] != "shipped" {
		t.Errorf("fulfillment_status incorrecto: got %q, want %q", result["fulfillment_status"], "shipped")
	}
}

func TestGetOrdersParams_ToQueryString_WithSinceID(t *testing.T) {
	// Arrange
	sinceID := int64(12345678)
	p := &GetOrdersParams{
		SinceID: &sinceID,
	}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["since_id"] != "12345678" {
		t.Errorf("since_id incorrecto: got %q, want %q", result["since_id"], "12345678")
	}
}

func TestGetOrdersParams_ToQueryString_WithUpdatedAtRange(t *testing.T) {
	// Arrange
	updatedMin := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)
	updatedMax := time.Date(2025, 6, 30, 23, 59, 59, 0, time.UTC)

	p := &GetOrdersParams{
		UpdatedAtMin: &updatedMin,
		UpdatedAtMax: &updatedMax,
	}

	// Act
	result := p.ToQueryString()

	// Assert
	if result["updated_at_min"] != updatedMin.Format(time.RFC3339) {
		t.Errorf("updated_at_min incorrecto: got %q", result["updated_at_min"])
	}
	if result["updated_at_max"] != updatedMax.Format(time.RFC3339) {
		t.Errorf("updated_at_max incorrecto: got %q", result["updated_at_max"])
	}
}

func TestGetOrdersParams_ToQueryString_NoFields_KeyAbsent(t *testing.T) {
	// Arrange: sin fields, la clave no debe aparecer en el resultado
	p := &GetOrdersParams{}

	// Act
	result := p.ToQueryString()

	// Assert
	if _, exists := result["fields"]; exists {
		t.Error("la clave 'fields' no deberia estar presente cuando Fields esta vacio")
	}
}

func TestGetOrdersParams_ToQueryString_NoFinancialStatus_KeyAbsent(t *testing.T) {
	// Arrange
	p := &GetOrdersParams{}

	// Act
	result := p.ToQueryString()

	// Assert
	if _, exists := result["financial_status"]; exists {
		t.Error("la clave 'financial_status' no deberia estar presente cuando no se especifica")
	}
}

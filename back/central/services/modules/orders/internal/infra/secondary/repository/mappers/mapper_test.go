package mappers

import (
	"testing"

	"github.com/secamc93/probability/back/migration/shared/models"
)

func TestToDomainOrder_BusinessNil_NoPanicYNombreVacio(t *testing.T) {
	order := &models.Order{
		ExternalID:    "6881022115887",
		OrderNumber:   "WEB-1018",
		IntegrationID: 230,
		Business:      nil,
	}

	var result = ToDomainOrder(order, "")

	if result == nil {
		t.Fatal("se esperaba una orden de dominio, se obtuvo nil")
	}
	if result.BusinessName != "" {
		t.Errorf("se esperaba BusinessName vacio, se obtuvo: %q", result.BusinessName)
	}
	if result.OrderNumber != "WEB-1018" {
		t.Errorf("se esperaba OrderNumber WEB-1018, se obtuvo: %q", result.OrderNumber)
	}
}

func TestToDomainOrder_ConBusiness_MapeaNombre(t *testing.T) {
	order := &models.Order{
		ExternalID:  "EXT-1",
		OrderNumber: "ORD-1",
		Business:    &models.Business{Name: "Gamers Paradise"},
	}

	result := ToDomainOrder(order, "")

	if result == nil {
		t.Fatal("se esperaba una orden de dominio, se obtuvo nil")
	}
	if result.BusinessName != "Gamers Paradise" {
		t.Errorf("se esperaba BusinessName 'Gamers Paradise', se obtuvo: %q", result.BusinessName)
	}
}

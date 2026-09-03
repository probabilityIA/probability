package app

import (
	"testing"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
	invoicingErrors "github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/errors"
)

func requireGuideValidator(t *testing.T, config *entities.FilterConfig) *RequireGuideValidator {
	t.Helper()
	for _, v := range CreateValidators(config, "") {
		if rg, ok := v.(*RequireGuideValidator); ok {
			return rg
		}
	}
	return nil
}

func TestRequireGuideApagadoNoCreaValidador(t *testing.T) {
	if requireGuideValidator(t, &entities.FilterConfig{}) != nil {
		t.Fatal("sin require_guide no debe existir el validador")
	}
	falso := false
	if requireGuideValidator(t, &entities.FilterConfig{RequireGuide: &falso}) != nil {
		t.Fatal("con require_guide en false no debe existir el validador")
	}
}

func TestRequireGuideBloqueaOrdenSinGuia(t *testing.T) {
	activo := true
	v := requireGuideValidator(t, &entities.FilterConfig{RequireGuide: &activo})
	if v == nil {
		t.Fatal("con require_guide en true debe existir el validador")
	}

	if err := v.Validate(&dtos.OrderData{HasGuide: false}); err != invoicingErrors.ErrOrderWithoutGuide {
		t.Fatalf("se esperaba ErrOrderWithoutGuide, llego %v", err)
	}
	if err := v.Validate(&dtos.OrderData{HasGuide: true}); err != nil {
		t.Fatalf("con guia generada la orden debe pasar, llego %v", err)
	}
}

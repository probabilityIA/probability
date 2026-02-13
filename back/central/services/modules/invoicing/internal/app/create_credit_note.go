package app

import (
	"context"
	"fmt"

	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/invoicing/internal/domain/entities"
)

// CreateCreditNote crea una nota de crédito para una factura
// DEPRECATED: Esta funcionalidad fue migrada a integrations/invoicing/softpymes/
// TODO: Re-implementar usando integrationCore si es necesario
func (uc *useCase) CreateCreditNote(ctx context.Context, dto *dtos.CreateCreditNoteDTO) (*entities.CreditNote, error) {
	return nil, fmt.Errorf("CreateCreditNote is deprecated and was moved to softpymes integration")
}

// GetCreditNote obtiene una nota de crédito por ID
func (uc *useCase) GetCreditNote(ctx context.Context, id uint) (*entities.CreditNote, error) {
	return uc.repo.GetCreditNoteByID(ctx, id)
}

// ListCreditNotes lista notas de crédito con filtros
func (uc *useCase) ListCreditNotes(ctx context.Context, filters map[string]interface{}) ([]*entities.CreditNote, error) {
	return uc.repo.ListCreditNotes(ctx, filters)
}

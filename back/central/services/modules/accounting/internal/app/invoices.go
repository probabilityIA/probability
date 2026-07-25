package app

import (
	"context"
	"strings"
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
	domainerrors "github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/errors"
)

func (uc *UseCase) buildInvoice(ctx context.Context, conceptID uint, taxIDs []uint, items []dtos.InvoiceItemDTO) (*entities.Invoice, error) {
	if len(items) == 0 {
		return nil, domainerrors.ErrInvoiceNoItems
	}
	concept, err := uc.repo.GetConceptByID(ctx, conceptID)
	if err != nil {
		return nil, err
	}
	if concept.Kind != entities.KindIncome {
		return nil, domainerrors.ErrInvoiceConceptKind
	}
	inv := &entities.Invoice{ConceptID: concept.ID}
	for _, item := range items {
		desc := strings.TrimSpace(item.Description)
		if desc == "" || item.Quantity <= 0 || item.UnitPrice <= 0 {
			return nil, domainerrors.ErrInvoiceInvalidItem
		}
		amount := round2(item.Quantity * item.UnitPrice)
		inv.Items = append(inv.Items, entities.InvoiceItem{
			ServiceID:   item.ServiceID,
			UnspscCode:  strings.TrimSpace(item.UnspscCode),
			Description: desc,
			Quantity:    item.Quantity,
			UnitPrice:   item.UnitPrice,
			Amount:      amount,
		})
		inv.Subtotal += amount
	}
	inv.Subtotal = round2(inv.Subtotal)
	inv.TaxDetail = []entities.TaxLine{}
	inv.WithholdingDetail = []entities.TaxLine{}
	if len(taxIDs) > 0 {
		taxes, err := uc.repo.GetTaxesByIDs(ctx, taxIDs)
		if err != nil {
			return nil, err
		}
		for _, tax := range taxes {
			if !tax.IsActive || tax.Kind == entities.TaxKindOther {
				continue
			}
			amount := round2(inv.Subtotal * tax.RatePercent / 100)
			line := entities.TaxLine{
				TaxCode:     tax.Code,
				TaxName:     tax.Name,
				RatePercent: tax.RatePercent,
				Amount:      amount,
			}
			if tax.Kind == entities.TaxKindWithholding {
				inv.WithholdingDetail = append(inv.WithholdingDetail, line)
				inv.WithholdingTotal += amount
				continue
			}
			inv.TaxDetail = append(inv.TaxDetail, line)
			inv.TaxTotal += amount
		}
	}
	inv.TaxTotal = round2(inv.TaxTotal)
	inv.WithholdingTotal = round2(inv.WithholdingTotal)
	inv.Total = round2(inv.Subtotal + inv.TaxTotal)
	return inv, nil
}

func (uc *UseCase) CreateInvoice(ctx context.Context, dto dtos.CreateInvoiceDTO) (*entities.Invoice, error) {
	if _, err := uc.repo.GetBusinessName(ctx, dto.BusinessID); err != nil {
		return nil, err
	}
	inv, err := uc.buildInvoice(ctx, dto.ConceptID, dto.TaxIDs, dto.Items)
	if err != nil {
		return nil, err
	}
	inv.BusinessID = dto.BusinessID
	inv.IssueDate = dto.IssueDate
	inv.DueDate = dto.DueDate
	inv.Notes = strings.TrimSpace(dto.Notes)
	inv.EmailTo = strings.TrimSpace(dto.EmailTo)
	inv.CustomerDocument = strings.TrimSpace(dto.CustomerDocument)
	inv.CustomerPhone = strings.TrimSpace(dto.CustomerPhone)
	inv.CustomerAddress = strings.TrimSpace(dto.CustomerAddress)
	inv.IsTest = dto.IsTest
	inv.Status = entities.InvoiceStatusDraft
	created, err := uc.repo.CreateInvoice(ctx, inv)
	if err != nil {
		return nil, err
	}
	return uc.repo.GetInvoiceByID(ctx, created.ID)
}

func (uc *UseCase) UpdateInvoice(ctx context.Context, dto dtos.UpdateInvoiceDTO) (*entities.Invoice, error) {
	existing, err := uc.repo.GetInvoiceByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	if existing.Status != entities.InvoiceStatusDraft {
		return nil, domainerrors.ErrInvoiceNotDraft
	}
	inv, err := uc.buildInvoice(ctx, dto.ConceptID, dto.TaxIDs, dto.Items)
	if err != nil {
		return nil, err
	}
	inv.ID = existing.ID
	inv.BusinessID = existing.BusinessID
	inv.IssueDate = dto.IssueDate
	inv.DueDate = dto.DueDate
	inv.Notes = strings.TrimSpace(dto.Notes)
	inv.EmailTo = strings.TrimSpace(dto.EmailTo)
	inv.CustomerDocument = strings.TrimSpace(dto.CustomerDocument)
	inv.CustomerPhone = strings.TrimSpace(dto.CustomerPhone)
	inv.CustomerAddress = strings.TrimSpace(dto.CustomerAddress)
	inv.IsTest = dto.IsTest
	inv.Status = existing.Status
	if _, err := uc.repo.UpdateInvoice(ctx, inv); err != nil {
		return nil, err
	}
	return uc.repo.GetInvoiceByID(ctx, inv.ID)
}

func (uc *UseCase) GetInvoice(ctx context.Context, id uint) (*entities.Invoice, error) {
	return uc.repo.GetInvoiceByID(ctx, id)
}

func (uc *UseCase) ListInvoices(ctx context.Context, params dtos.ListInvoicesParams) ([]entities.Invoice, int64, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 {
		params.PageSize = 10
	}
	if params.PageSize > 100 {
		params.PageSize = 100
	}
	return uc.repo.ListInvoices(ctx, params)
}

func (uc *UseCase) DeleteInvoice(ctx context.Context, id uint) error {
	inv, err := uc.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return err
	}
	if inv.Status != entities.InvoiceStatusDraft && inv.Status != entities.InvoiceStatusCancelled {
		return domainerrors.ErrInvoiceNotDraft
	}
	return uc.repo.DeleteInvoice(ctx, id)
}

func (uc *UseCase) SendInvoice(ctx context.Context, dto dtos.SendInvoiceDTO) (*entities.Invoice, error) {
	inv, err := uc.repo.GetInvoiceByID(ctx, dto.ID)
	if err != nil {
		return nil, err
	}
	if inv.Status == entities.InvoiceStatusCancelled {
		return nil, domainerrors.ErrInvoiceCancelled
	}
	to := strings.TrimSpace(dto.EmailTo)
	if to == "" {
		to = inv.EmailTo
	}
	if to == "" {
		to, err = uc.repo.GetBusinessDefaultEmail(ctx, inv.BusinessID)
		if err != nil {
			return nil, err
		}
	}
	if to == "" {
		return nil, domainerrors.ErrInvoiceNoEmail
	}
	subject := "Factura " + inv.Number + " - Probability"
	if err := uc.email.SendHTML(ctx, to, subject, renderInvoiceHTML(inv)); err != nil {
		return nil, err
	}
	now := time.Now()
	status := inv.Status
	if status == entities.InvoiceStatusDraft {
		status = entities.InvoiceStatusSent
	}
	if err := uc.repo.SetInvoiceStatus(ctx, inv.ID, status, &now, inv.PaidAt, to); err != nil {
		return nil, err
	}
	uc.log.Info(ctx).Uint("invoice_id", inv.ID).Str("to", to).Msg("Accounting invoice sent by email")
	return uc.repo.GetInvoiceByID(ctx, inv.ID)
}

func (uc *UseCase) MarkInvoicePaid(ctx context.Context, id uint) (*entities.Invoice, error) {
	inv, err := uc.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	switch inv.Status {
	case entities.InvoiceStatusPaid:
		return nil, domainerrors.ErrInvoiceAlreadyPaid
	case entities.InvoiceStatusCancelled:
		return nil, domainerrors.ErrInvoiceCancelled
	case entities.InvoiceStatusDraft, entities.InvoiceStatusSent:
	default:
		return nil, domainerrors.ErrInvoiceNotPayable
	}
	now := time.Now()
	if inv.IsTest {
		if err := uc.repo.SetInvoiceStatus(ctx, inv.ID, entities.InvoiceStatusPaid, inv.SentAt, &now, inv.EmailTo); err != nil {
			return nil, err
		}
		uc.log.Info(ctx).Uint("invoice_id", inv.ID).Msg("Test invoice marked paid: no ledger entry created")
		return uc.repo.GetInvoiceByID(ctx, inv.ID)
	}
	businessID := inv.BusinessID
	entry := &entities.Entry{
		ConceptID:   inv.ConceptID,
		BusinessID:  &businessID,
		EntryDate:   now,
		Amount:      inv.Subtotal,
		Kind:        entities.KindIncome,
		SourceType:  entities.SourceInvoice,
		SourceID:    inv.Number,
		Description: "Factura " + inv.Number,
		IsAutomatic: true,
		TaxTotal:    inv.TaxTotal,
		TaxDetail:   inv.TaxDetail,
	}
	if _, err := uc.repo.CreateEntry(ctx, entry); err != nil {
		return nil, err
	}
	if err := uc.repo.SetInvoiceStatus(ctx, inv.ID, entities.InvoiceStatusPaid, inv.SentAt, &now, inv.EmailTo); err != nil {
		return nil, err
	}
	return uc.repo.GetInvoiceByID(ctx, inv.ID)
}

func (uc *UseCase) CancelInvoice(ctx context.Context, id uint) (*entities.Invoice, error) {
	inv, err := uc.repo.GetInvoiceByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if inv.Status == entities.InvoiceStatusCancelled {
		return nil, domainerrors.ErrInvoiceCancelled
	}
	if inv.Status == entities.InvoiceStatusPaid {
		if err := uc.repo.DeleteEntryBySource(ctx, entities.SourceInvoice, inv.Number); err != nil {
			return nil, err
		}
	}
	if err := uc.repo.SetInvoiceStatus(ctx, inv.ID, entities.InvoiceStatusCancelled, inv.SentAt, nil, inv.EmailTo); err != nil {
		return nil, err
	}
	return uc.repo.GetInvoiceByID(ctx, inv.ID)
}

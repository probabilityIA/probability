package mappers

import (
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/response"
)

// BillsToListResult mapea la respuesta de GET /v1/bills al tipo de dominio
func BillsToListResult(apiResp *response.Bills) *dtos.ListBillsResult {
	bills := make([]dtos.Bill, 0, len(apiResp.Data.Data))

	for _, b := range apiResp.Data.Data {
		creditNotes := make([]dtos.BillNote, 0, len(b.CreditNotes))
		for _, n := range b.CreditNotes {
			creditNotes = append(creditNotes, dtos.BillNote{
				ID:     n.ID,
				Number: n.Number,
			})
		}

		debitNotes := make([]dtos.BillNote, 0, len(b.DebitNotes))
		for _, n := range b.DebitNotes {
			debitNotes = append(debitNotes, dtos.BillNote{
				ID:     n.ID,
				Number: n.Number,
			})
		}

		bills = append(bills, dtos.Bill{
			ID: b.ID,
			Document: dtos.BillDocument{
				Code: b.Document.Code,
				Name: b.Document.Name,
			},
			Number:                    b.Number,
			APIClientName:             b.APIClientName,
			ReferenceCode:             b.ReferenceCode,
			Identification:            b.Identification,
			GraphicRepresentationName: b.GraphicRepresentationName,
			Company:                   b.Company,
			TradeName:                 b.TradeName,
			Names:                     b.Names,
			Email:                     b.Email,
			Total:                     b.Total,
			Status:                    b.Status,
			Errors:                    b.Errors,
			SendEmail:                 b.SendEmail == 1,
			HasClaim:                  b.HasClaim == 1,
			IsNegotiableInstrument:    b.IsNegotiableInstrument == 1,
			PaymentForm: dtos.BillPaymentForm{
				Code: b.PaymentForm.Code,
				Name: b.PaymentForm.Name,
			},
			CreatedAt:   b.CreatedAt,
			CreditNotes: creditNotes,
			DebitNotes:  debitNotes,
		})
	}

	p := apiResp.Data.Pagination
	return &dtos.ListBillsResult{
		Bills: bills,
		Pagination: dtos.BillsPagination{
			Total:       p.Total,
			PerPage:     p.PerPage,
			CurrentPage: p.CurrentPage,
			LastPage:    p.LastPage,
			From:        p.From,
			To:          p.To,
		},
	}
}

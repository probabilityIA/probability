package mappers

import (
	"fmt"

	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/domain/dtos"
	"github.com/secamc93/probability/back/central/services/integrations/invoicing/factus/internal/infra/secondary/client/response"
)

// GetBillToDetail mapea la respuesta de GET /v1/bills/show/:number al tipo de dominio
func GetBillToDetail(apiResp *response.GetBillDetail) *dtos.BillDetail {
	b := apiResp.Data.Bill
	c := apiResp.Data.Customer

	items := make([]dtos.BillDetailItem, 0, len(apiResp.Data.Items))
	for _, it := range apiResp.Data.Items {
		qty := 0.0
		fmt.Sscanf(it.Quantity, "%f", &qty)
		items = append(items, dtos.BillDetailItem{
			CodeReference: it.CodeReference,
			Name:          it.Name,
			Quantity:      qty,
			Price:         it.Price,
			DiscountRate:  it.DiscountRate,
			Discount:      it.Discount,
			TaxRate:       it.TaxRate,
			TaxAmount:     it.TaxAmount,
			Total:         it.Total,
		})
	}

	creditNotes := make([]dtos.BillNote, 0, len(b.CreditNotes))
	for _, n := range b.CreditNotes {
		creditNotes = append(creditNotes, dtos.BillNote{ID: n.ID, Number: n.Number})
	}

	debitNotes := make([]dtos.BillNote, 0, len(b.DebitNotes))
	for _, n := range b.DebitNotes {
		debitNotes = append(debitNotes, dtos.BillNote{ID: n.ID, Number: n.Number})
	}

	customer := dtos.BillDetailCustomer{
		Identification: c.Identification,
		DV:             c.DV,
		Names:          c.Names,
		Company:        c.Company,
		TradeName:      c.TradeName,
		Address:        c.Address,
	}
	if c.Email != nil {
		customer.Email = *c.Email
	}
	if c.Phone != nil {
		customer.Phone = *c.Phone
	}

	return &dtos.BillDetail{
		ID:            b.ID,
		Number:        b.Number,
		ReferenceCode: b.ReferenceCode,
		CUFE:          b.CUFE,
		QRCode:        b.QR,
		QRImage:       b.QRImage,
		Status:        b.Status,
		Total:         b.Total,
		TaxAmount:     b.TaxAmount,
		GrossValue:    b.GrossValue,
		Discount:      b.Discount,
		Validated:     b.Validated,
		CreatedAt:     b.CreatedAt,
		Document: dtos.BillDocument{
			Code: b.Document.Code,
			Name: b.Document.Name,
		},
		PaymentForm: dtos.BillPaymentForm{
			Code: b.PaymentForm.Code,
			Name: b.PaymentForm.Name,
		},
		Customer:    customer,
		Items:       items,
		CreditNotes: creditNotes,
		DebitNotes:  debitNotes,
	}
}

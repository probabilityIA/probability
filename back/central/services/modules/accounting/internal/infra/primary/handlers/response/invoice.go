package response

import (
	"time"

	"github.com/secamc93/probability/back/central/services/modules/accounting/internal/domain/entities"
)

type InvoiceItemResponse struct {
	ID          uint    `json:"id"`
	Description string  `json:"description"`
	Quantity    float64 `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Amount      float64 `json:"amount"`
}

type InvoiceResponse struct {
	ID           uint                  `json:"id"`
	Number       string                `json:"number"`
	BusinessID   uint                  `json:"business_id"`
	BusinessName string                `json:"business_name"`
	ConceptID    uint                  `json:"concept_id"`
	ConceptName  string                `json:"concept_name"`
	IssueDate    string                `json:"issue_date"`
	DueDate      *string               `json:"due_date"`
	Status       string                `json:"status"`
	Notes        string                `json:"notes"`
	Subtotal     float64               `json:"subtotal"`
	TaxTotal     float64               `json:"tax_total"`
	Total        float64               `json:"total"`
	TaxDetail    []TaxLineResponse     `json:"tax_detail"`
	EmailTo      string                `json:"email_to"`
	SentAt       *time.Time            `json:"sent_at"`
	PaidAt       *time.Time            `json:"paid_at"`
	Items        []InvoiceItemResponse `json:"items"`
	CreatedAt    time.Time             `json:"created_at"`

	CustomerDocument string     `json:"customer_document"`
	CustomerPhone    string     `json:"customer_phone"`
	CustomerAddress  string     `json:"customer_address"`
	DianStatus       string     `json:"dian_status"`
	Cufe             string     `json:"cufe"`
	DianNumber       string     `json:"dian_number"`
	DianQR           string     `json:"dian_qr"`
	DianEmittedAt    *time.Time `json:"dian_emitted_at"`
	IsTest           bool       `json:"is_test"`

	WithholdingTotal  float64           `json:"withholding_total"`
	WithholdingDetail []TaxLineResponse `json:"withholding_detail"`
	NetReceivable     float64           `json:"net_receivable"`
}

func FromInvoice(e *entities.Invoice) InvoiceResponse {
	var dueDate *string
	if e.DueDate != nil {
		d := e.DueDate.Format("2006-01-02")
		dueDate = &d
	}
	taxes := make([]TaxLineResponse, len(e.TaxDetail))
	for i, t := range e.TaxDetail {
		taxes[i] = TaxLineResponse{TaxCode: t.TaxCode, TaxName: t.TaxName, RatePercent: t.RatePercent, Amount: t.Amount}
	}
	withholdings := make([]TaxLineResponse, len(e.WithholdingDetail))
	for i, t := range e.WithholdingDetail {
		withholdings[i] = TaxLineResponse{TaxCode: t.TaxCode, TaxName: t.TaxName, RatePercent: t.RatePercent, Amount: t.Amount}
	}
	items := make([]InvoiceItemResponse, len(e.Items))
	for i, item := range e.Items {
		items[i] = InvoiceItemResponse{ID: item.ID, Description: item.Description, Quantity: item.Quantity, UnitPrice: item.UnitPrice, Amount: item.Amount}
	}
	return InvoiceResponse{
		ID:           e.ID,
		Number:       e.Number,
		BusinessID:   e.BusinessID,
		BusinessName: e.BusinessName,
		ConceptID:    e.ConceptID,
		ConceptName:  e.ConceptName,
		IssueDate:    e.IssueDate.Format("2006-01-02"),
		DueDate:      dueDate,
		Status:       e.Status,
		Notes:        e.Notes,
		Subtotal:     e.Subtotal,
		TaxTotal:     e.TaxTotal,
		Total:        e.Total,
		TaxDetail:    taxes,
		EmailTo:      e.EmailTo,
		SentAt:       e.SentAt,
		PaidAt:       e.PaidAt,
		Items:        items,
		CreatedAt:    e.CreatedAt,

		CustomerDocument: e.CustomerDocument,
		CustomerPhone:    e.CustomerPhone,
		CustomerAddress:  e.CustomerAddress,
		DianStatus:       e.DianStatus,
		Cufe:             e.Cufe,
		DianNumber:       e.DianNumber,
		DianQR:           e.DianQR,
		DianEmittedAt:    e.DianEmittedAt,
		IsTest:           e.IsTest,

		WithholdingTotal:  e.WithholdingTotal,
		WithholdingDetail: withholdings,
		NetReceivable:     e.Total - e.WithholdingTotal,
	}
}

type InvoicesListResponse struct {
	Success    bool              `json:"success"`
	Data       []InvoiceResponse `json:"data"`
	Total      int64             `json:"total"`
	Page       int               `json:"page"`
	PageSize   int               `json:"page_size"`
	TotalPages int               `json:"total_pages"`
}

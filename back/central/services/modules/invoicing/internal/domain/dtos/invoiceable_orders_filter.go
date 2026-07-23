package dtos

import "time"

type InvoiceableOrdersFilter struct {
	BusinessID          uint
	Page                int
	PageSize            int
	StartDate           *time.Time
	EndDate             *time.Time
	OrderNumber         string
	CustomerName        string
	CustomerEmail       string
	PaymentStatusID     uint
	FulfillmentStatusID uint
	SortBy              string
	SortOrder           string
}

const (
	InvoiceableOrdersMaxPageSize     = 1000
	InvoiceableOrdersDefaultPageSize = 20
)

var invoiceableSortColumns = map[string]string{
	"created_at":   "created_at",
	"occurred_at":  "occurred_at",
	"order_number": "order_number",
	"total_amount": "total_amount",
	"customer_name": "customer_name",
}

func (f *InvoiceableOrdersFilter) Sanitize() {
	if f.Page < 1 {
		f.Page = 1
	}
	if f.PageSize < 1 {
		f.PageSize = InvoiceableOrdersDefaultPageSize
	}
	if f.PageSize > InvoiceableOrdersMaxPageSize {
		f.PageSize = InvoiceableOrdersMaxPageSize
	}
	if _, ok := invoiceableSortColumns[f.SortBy]; !ok {
		f.SortBy = "created_at"
	}
	if f.SortOrder != "asc" && f.SortOrder != "ASC" {
		f.SortOrder = "DESC"
	} else {
		f.SortOrder = "ASC"
	}
}

func (f *InvoiceableOrdersFilter) SortColumn() string {
	if col, ok := invoiceableSortColumns[f.SortBy]; ok {
		return col
	}
	return "created_at"
}

package dtos

// CreateJournalDTO contiene los datos para crear un comprobante contable (journal)
type CreateJournalDTO struct {
	OrderID         string
	IsManual        bool
	CreatedByUserID *uint
}

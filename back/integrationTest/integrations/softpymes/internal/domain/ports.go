package domain

// IAPIClient define la interfaz para el cliente HTTP que simula SoftPymes API
type IAPIClient interface {
	HandleAuth(apiKey, apiSecret, referer string) (string, error)
	HandleCreateInvoice(token string, invoiceData map[string]interface{}) (*Invoice, error)
	HandleCreateCreditNote(token string, creditNoteData map[string]interface{}) (*CreditNote, error)
	HandleListDocuments(token string, filters map[string]interface{}) ([]Invoice, error)
}

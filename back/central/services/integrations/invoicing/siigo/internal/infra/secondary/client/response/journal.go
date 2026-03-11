package response

// CreateJournalResponse representa la respuesta de Siigo al crear un comprobante contable
type CreateJournalResponse struct {
	ID       string              `json:"id"`
	Document JournalDocumentRef  `json:"document"`
	Date     string              `json:"date"`
	Number   int                 `json:"number"`
	Name     string              `json:"name"` // Ej: "CC-10-20"
	Items    []JournalItemRef    `json:"items,omitempty"`
	Total    float64             `json:"total"`
	Errors   []SiigoError        `json:"Errors,omitempty"`
	Metadata JournalMetadata     `json:"metadata,omitempty"`
}

// JournalDocumentRef referencia al tipo de documento del journal
type JournalDocumentRef struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// JournalItemRef item en la respuesta del journal
type JournalItemRef struct {
	Account  map[string]interface{} `json:"account,omitempty"`
	Movement string                 `json:"movement,omitempty"`
	Value    float64                `json:"value,omitempty"`
}

// JournalMetadata metadatos del journal
type JournalMetadata struct {
	CreatedAt string `json:"created,omitempty"`
}

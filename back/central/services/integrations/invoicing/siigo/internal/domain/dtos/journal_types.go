package dtos

// JournalItemData datos de un item para un comprobante contable en Siigo
type JournalItemData struct {
	SKU             string
	Name            string
	Quantity        float64
	TotalPrice      float64
	CustomerDNI     string
	AccountCode     string
	Movement        string
	WarehouseID     int
	CostCenter      int
	TaxID           int
	TaxName         string
	TaxType         string
	TaxPercentage   float64
}

// CreateJournalRequest datos tipados para crear un comprobante contable en Siigo
type CreateJournalRequest struct {
	Items        []JournalItemData
	Currency     string
	Observations string
	Date         string
	OrderID      string
	Credentials  Credentials
	Config       map[string]interface{}
}

// CreateJournalResult resultado de crear un comprobante contable en Siigo
type CreateJournalResult struct {
	JournalName string
	JournalID   string
	Number      int
	Date        string
	ProviderInfo map[string]interface{}
	AuditData   *AuditData
}

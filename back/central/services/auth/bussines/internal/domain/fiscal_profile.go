package domain

type FiscalProfile struct {
	ID                uint    `json:"id"`
	BusinessID        uint    `json:"business_id"`
	DocumentType      string  `json:"document_type"`
	DocumentNumber    string  `json:"document_number"`
	DV                string  `json:"dv"`
	PersonType        string  `json:"person_type"`
	TaxRegime         string  `json:"tax_regime"`
	MunicipalityID    string  `json:"municipality_id"`
	Address           string  `json:"address"`
	Phone             string  `json:"phone"`
	BillingEmail      string  `json:"billing_email"`
	AppliesIVA        bool    `json:"applies_iva"`
	IVARate           float64 `json:"iva_rate"`
	AppliesReteFuente bool    `json:"applies_retefuente"`
	ReteFuenteRate    float64 `json:"retefuente_rate"`
	AppliesReteICA    bool    `json:"applies_reteica"`
	ReteICARate       float64 `json:"reteica_rate"`
	Notes             string  `json:"notes"`
}

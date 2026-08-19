package dtos

type CatalogItem struct {
	ID      int    `json:"id"`
	Code    string `json:"code,omitempty"`
	Name    string `json:"name"`
	Detail  string `json:"detail,omitempty"`
	Percent string `json:"percent,omitempty"`
}

type Catalogs struct {
	DocumentTypesFV []CatalogItem `json:"document_types_fv"`
	DocumentTypesNC []CatalogItem `json:"document_types_nc"`
	DocumentTypesRC []CatalogItem `json:"document_types_rc"`
	DocumentTypesCC []CatalogItem `json:"document_types_cc"`
	PaymentTypesFV  []CatalogItem `json:"payment_types_fv"`
	PaymentTypesRC  []CatalogItem `json:"payment_types_rc"`
	Sellers         []CatalogItem `json:"sellers"`
	Taxes           []CatalogItem `json:"taxes"`
	CostCenters     []CatalogItem `json:"cost_centers"`
	Warehouses      []CatalogItem `json:"warehouses"`
	Errors          []string      `json:"errors,omitempty"`
}

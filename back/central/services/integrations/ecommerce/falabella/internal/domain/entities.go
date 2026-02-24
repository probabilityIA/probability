package domain

// Integration representa los datos de una integración de Falabella
// tal como se obtienen del core de integraciones.
type Integration struct {
	ID              uint
	BusinessID      *uint
	Name            string
	StoreID         string
	IntegrationType int
	Config          map[string]interface{}
}

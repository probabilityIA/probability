package domain

type ShipmentLabel struct {
	Content     []byte
	ContentType string
	Filename    string
}

type MeliLabelRef struct {
	BusinessID     uint
	IntegrationID  uint
	MeliShipmentID int64
}

package dtos

type RegisterDeviceDTO struct {
	UserID     uint
	BusinessID *uint
	Token      string
	Platform   string
	AppVersion string
	DeviceName string
}

type PushEventDTO struct {
	EventType     string
	EventCategory string
	BusinessID    uint
	OrderNumber   string
	OrderID       string
	Carrier       string
	TrackingNo    string
	StatusCode    string
	StatusName    string
	CustomerName  string
	CodTotal      float64
	Balance       float64
}

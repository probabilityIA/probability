package request

type RegisterDevice struct {
	Token      string `json:"token" binding:"required"`
	Platform   string `json:"platform" binding:"required"`
	AppVersion string `json:"app_version"`
	DeviceName string `json:"device_name"`
	BusinessID *uint  `json:"business_id"`
}

type UnregisterDevice struct {
	Token string `json:"token" binding:"required"`
}

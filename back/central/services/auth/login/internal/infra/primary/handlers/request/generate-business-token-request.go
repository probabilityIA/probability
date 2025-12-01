package request

// GenerateBusinessTokenRequest representa la solicitud para generar un token de business
type GenerateBusinessTokenRequest struct {
	BusinessID uint `json:"business_id"` // 0 para super admin, >0 para usuarios normales
}

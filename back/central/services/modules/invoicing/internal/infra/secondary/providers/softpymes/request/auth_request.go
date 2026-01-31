package request

// AuthRequest representa la solicitud de autenticación a Softpymes
type AuthRequest struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_secret"`
}

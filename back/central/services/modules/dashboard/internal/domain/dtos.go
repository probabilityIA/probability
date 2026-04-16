package domain

// DashboardStatsResponse es la respuesta del endpoint de estadísticas
type DashboardStatsResponse struct {
	Success bool           `json:"success"`
	Message string         `json:"message"`
	Data    DashboardStats `json:"data"`
}
